"""Publish or verify an immutable release and wait for GitHub Latest visibility."""

import argparse
import fnmatch
import json
from pathlib import Path
import re
import subprocess
import time
from urllib.parse import quote

from go_module import ROOT_MODULE
from release import STABLE_VERSION, release_payload

REPOSITORY = ROOT_MODULE.removeprefix("github.com/")

API_TIMEOUT_SECONDS = 20
LATEST_ATTEMPTS = 6
LATEST_RETRY_SECONDS = 5


class GitHubError(ValueError):
    """Keep HTTP failures distinct from transport errors and malformed responses."""

    def __init__(self, status: int, detail: str):
        self.status = status
        super().__init__(f"GitHub HTTP {status}: {detail}")


def github_api(path: str, *, payload: dict | None = None) -> dict:
    """Use gh credentials without putting tokens in arguments or diagnostics."""
    command = ["gh", "api", "--include", "--method", "GET" if payload is None else "POST", path]
    if payload is not None:
        command.extend(["--input", "-"])
    try:
        result = subprocess.run(
            command, input=None if payload is None else json.dumps(payload),
            capture_output=True, text=True, check=False, timeout=API_TIMEOUT_SECONDS,
        )
    except (OSError, subprocess.SubprocessError) as error:
        raise ValueError(f"GitHub request failed for {path}: {error}") from error
    headers, separator, body = result.stdout.partition("\n\n")
    status = re.match(r"HTTP/\S+\s+(\d{3})\b", headers)
    if not separator or status is None:
        raise ValueError(f"Malformed GitHub response for {path}: {result.stderr.strip()}")
    code = int(status[1])
    if not 200 <= code < 300:
        raise GitHubError(code, f"request failed for {path}")
    if result.returncode:
        raise ValueError(f"GitHub request failed for {path}: {result.stderr.strip()}")
    try:
        data = json.loads(body)
        if not isinstance(data, dict):
            raise ValueError("expected a JSON object")
        return data
    except ValueError as error:
        raise ValueError(f"Malformed GitHub response for {path}: {error}") from error


def stable_version(tag: str) -> tuple[int, int, int]:
    match = STABLE_VERSION.fullmatch(tag)
    if match is None:
        raise ValueError(f"Expected a stable Latest tag, received {tag!r}")
    return tuple(int(part) for part in match.groups())


def wait_for_latest(repository: str, tag: str) -> bool:
    """Bound propagation waits without promoting an older release over Latest."""
    current = None
    for attempt in range(LATEST_ATTEMPTS):
        try:
            current = github_api(f"repos/{repository}/releases/latest").get("tag_name")
            if not isinstance(current, str):
                raise ValueError("Latest response is missing its tag_name")
            if current == tag:
                return True
            if stable_version(current) > stable_version(tag):
                return False
        except GitHubError as error:
            if error.status not in (404, 429) and error.status < 500:
                raise
            current = None
        if attempt + 1 < LATEST_ATTEMPTS:
            time.sleep(LATEST_RETRY_SECONDS)
    if current is None:
        raise ValueError("GitHub Latest remains unavailable after bounded retries")
    raise ValueError(f"GitHub Latest did not converge to {tag} after bounded retries; received {current}")


def publish(repository: str, payload: dict, commit: str) -> bool:
    """Reuse matching published content; never edit notes or retry uncertain writes."""
    tag = payload["tag_name"]
    prefix = f"repos/{repository}"
    tag_ref = quote(f"refs/tags/{tag}", safe="")
    if github_api(f"{prefix}/commits/{tag_ref}").get("sha") != commit:
        raise ValueError(f"Release tag {tag} no longer matches checked commit {commit}")
    try:
        published = github_api(f"{prefix}/releases/tags/{tag}")
    except GitHubError as error:
        if error.status != 404:
            raise
        creation = dict(payload)
        if payload["make_latest"] == "true":
            try:
                latest = github_api(f"{prefix}/releases/latest")
            except GitHubError as latest_error:
                if latest_error.status != 404:
                    raise
            else:
                if stable_version(latest.get("tag_name", "")) > stable_version(tag):
                    creation["make_latest"] = "false"
        published = github_api(f"{prefix}/releases", payload=creation)
    verify_content(published, payload)
    if payload["make_latest"] == "false":
        return False
    return wait_for_latest(repository, tag)


def verify_content(published: dict, payload: dict) -> None:
    """Refuse to overwrite an existing Release with different maintained content."""
    for field in ("tag_name", "name", "body", "draft", "prerelease"):
        if published.get(field) != payload[field]:
            raise ValueError(f"Published release {payload['tag_name']} does not match maintained {field}")


def require_latest(repository: str, payload: dict) -> None:
    """Allow Pages to build and deploy only the current stable root documentation."""
    tag = payload["tag_name"]
    if STABLE_VERSION.fullmatch(tag) is None:
        raise ValueError("Pages requires a stable root release tag")
    verify_content(github_api(f"repos/{repository}/releases/tags/{tag}"), payload)
    if not wait_for_latest(repository, tag):
        raise ValueError(f"Pages requires the current Latest root release; {tag} is not Latest")


def check_pages(repository: str, tag: str) -> None:
    """Catch repository deployment settings before publishing a root release."""
    prefix = f"repos/{repository}"
    if github_api(f"{prefix}/pages").get("build_type") != "workflow":
        raise ValueError("GitHub Pages must use the workflow build type")
    environment = github_api(f"{prefix}/environments/github-pages")
    if "deployment_branch_policy" not in environment:
        raise ValueError("Missing github-pages deployment policy metadata")
    policy = environment["deployment_branch_policy"]
    if policy is None:
        return
    if not isinstance(policy, dict) or not policy.get("custom_branch_policies"):
        raise ValueError("github-pages requires a tag-compatible deployment policy")
    rules = github_api(f"{prefix}/environments/github-pages/deployment-branch-policies?per_page=100")
    entries = rules.get("branch_policies") or []
    if rules.get("total_count") != len(entries):
        raise ValueError("Cannot verify the complete github-pages deployment policy")
    if not any(
        item.get("type") == "tag" and fnmatch.fnmatchcase(tag, item.get("name", ""))
        for item in entries
    ):
        raise ValueError(f"github-pages environment does not allow tag {tag}; add a v* tag rule")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    parser.add_argument("--repository", default=REPOSITORY)
    parser.add_argument("--module", choices=("root", "tools"), required=True)
    parser.add_argument("--tag", required=True)
    parser.add_argument("--commit")
    checks = parser.add_mutually_exclusive_group()
    checks.add_argument("--require-latest", action="store_true")
    checks.add_argument("--check-pages", action="store_true")
    args = parser.parse_args()
    try:
        readonly = args.require_latest or args.check_pages
        if not readonly and not re.fullmatch(r"[0-9a-f]{40}", args.commit or ""):
            raise ValueError("A full checked source commit SHA is required")
        payload = release_payload(args.root, args.module, args.tag)
        if args.require_latest:
            require_latest(args.repository, payload)
            return
        if args.check_pages:
            if args.module != "root":
                raise ValueError("Pages settings apply only to root releases")
            check_pages(args.repository, args.tag)
            return
        latest = publish(args.repository, payload, args.commit)
        print(f"latest={str(latest).lower()}")
    except (ValueError, OSError) as error:
        parser.error(str(error))


if __name__ == "__main__":
    main()
