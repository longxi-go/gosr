"""Validate stable Go module tags and their maintained release notes."""

import argparse
import json
from pathlib import Path
import re
import subprocess

from go_module import ROOT_MODULE, read_module

STABLE_VERSION = re.compile(r"v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)")


def required_root_version(root: Path) -> str:
    """Validate the declared dependency; publication is checked separately."""
    requirements = read_module(root / "tools/go.mod").get("Require") or []
    versions = [item["Version"] for item in requirements if item["Path"] == ROOT_MODULE]
    version = versions[0] if len(versions) == 1 else ""
    match = STABLE_VERSION.fullmatch(version)
    if match is None or int(match[1]) >= 2 or version == "v0.0.0":
        raise ValueError("tools release requires a stable gosr dependency version")
    return version


def check_root_release(version: str, repository: str) -> None:
    """Require a completed stable root release before releasing its adapters."""
    message = f"tools release requires a published stable root release {version} in {repository}"
    try:
        response = subprocess.run(
            [
                "gh", "release", "view", version, "--repo", repository,
                "--json", "tagName,isDraft,isPrerelease",
            ],
            check=True,
            capture_output=True,
            text=True,
            timeout=30,
        )
        data = json.loads(response.stdout)
    except (OSError, subprocess.SubprocessError, ValueError) as error:
        raise ValueError(message) from error
    if (
        not isinstance(data, dict)
        or data.get("tagName") != version
        or data.get("isDraft") is not False
        or data.get("isPrerelease") is not False
    ):
        raise ValueError(message)


def validate_tag(root: Path, module: str, tag: str) -> str:
    """Validate the repository tag, module identity and notes before release."""
    if module == "root":
        prefix = ""
        module_file = root / "go.mod"
        title = "gosr"
    elif module == "tools":
        prefix = "tools/"
        module_file = root / "tools/go.mod"
        title = "gosr tools"
    else:
        raise ValueError(f"Unknown release module: {module}")

    if not tag.startswith(prefix):
        raise ValueError(f"{module} tags must begin with {prefix or 'v'}")
    version = tag[len(prefix) :]
    match = STABLE_VERSION.fullmatch(version)
    if match is None:
        raise ValueError("Release tags must use stable vX.Y.Z versions without leading zeroes")
    if int(match[1]) >= 2:
        raise ValueError("v2+ requires a versioned module path and a release-policy update")

    module_data = read_module(module_file)
    expected_module = ROOT_MODULE if module == "root" else f"{ROOT_MODULE}/tools"
    declared_module = (module_data.get("Module") or {}).get("Path")
    if not declared_module:
        raise ValueError(f"Missing module directive in {module_file}")
    if declared_module != expected_module:
        raise ValueError(f"Expected module {expected_module}, found {declared_module}")

    if module == "tools":
        required_root_version(root)

    notes = root / "docs/releases" / f"{tag}.md"
    if not notes.is_file():
        raise ValueError(f"Missing release notes: {notes}")
    lines = notes.read_text(encoding="utf-8").strip().splitlines()
    expected_title = f"# {title} {version}"
    if not lines or lines[0] != expected_title or not "\n".join(lines[1:]).strip():
        raise ValueError(f"Release notes must start with {expected_title!r} and have a body")
    return version[1:]


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    parser.add_argument("--module", choices=("root", "tools"), required=True)
    parser.add_argument("--tag", required=True)
    parser.add_argument("--github-output", type=Path)
    parser.add_argument("--check-root-release", action="store_true")
    parser.add_argument("--repository", default=ROOT_MODULE.removeprefix("github.com/"))
    args = parser.parse_args()
    try:
        version = validate_tag(args.root, args.module, args.tag)
        metadata = {"module": args.module}
        metadata["title"] = f"gosr{' tools' if args.module == 'tools' else ''} v{version}"
        if args.module == "tools":
            dependency = required_root_version(args.root)
            metadata["root_version"] = dependency
            if args.check_root_release:
                check_root_release(dependency, args.repository)
        elif args.check_root_release:
            raise ValueError("--check-root-release applies only to the tools module")
        if args.github_output:
            with args.github_output.open("a", encoding="utf-8") as output:
                output.writelines(f"{key}={value}\n" for key, value in metadata.items())
        print(version)
    except (ValueError, OSError) as error:
        parser.error(str(error))


if __name__ == "__main__":
    main()
