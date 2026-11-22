import copy
import json
import os
from pathlib import Path
import subprocess
import sys
import unittest
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import publish_release

REPOSITORY = "longxi-go/gosr"
TAG = "v0.1.0"
COMMIT = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
PREFIX = f"repos/{REPOSITORY}"


class PublishReleaseTests(unittest.TestCase):
    def setUp(self):
        self.payload = {
            "tag_name": TAG,
            "name": "gosr v0.1.0",
            "body": "# gosr v0.1.0\n\nConfiguration environments.\n",
            "draft": False,
            "prerelease": False,
            "make_latest": "true",
        }
        self.published = {
            **self.payload,
            "id": 454553171,
            "html_url": f"https://github.com/{REPOSITORY}/releases/tag/{TAG}",
            "target_commitish": "main",
            "created_at": "2026-11-22T15:37:48Z",
            "published_at": "2026-11-22T15:43:06Z",
            "assets": [],
        }
        self.published.pop("make_latest")

    def test_existing_matching_release_is_reused_without_writes(self):
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, self.published, self.published,
        ]) as api:
            self.assertTrue(publish_release.publish(REPOSITORY, self.payload, COMMIT))
        self.assertEqual([
            ((f"{PREFIX}/commits/refs%2Ftags%2F{TAG}",), {}),
            ((f"{PREFIX}/releases/tags/{TAG}",), {}),
            ((f"{PREFIX}/releases/latest",), {}),
        ], api.call_args_list)

    def test_missing_release_is_created_once_despite_delayed_latest(self):
        missing = publish_release.GitHubError(404, "Not Found")
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, missing, missing, self.published, missing, self.published,
        ]) as api, patch.object(publish_release.time, "sleep") as sleep:
            self.assertTrue(publish_release.publish(REPOSITORY, self.payload, COMMIT))
        writes = [call for call in api.call_args_list if call.kwargs.get("payload")]
        self.assertEqual(1, len(writes))
        self.assertEqual(f"{PREFIX}/releases", writes[0].args[0])
        self.assertEqual(self.payload, writes[0].kwargs["payload"])
        self.assertEqual(1, sleep.call_count)

    def test_conflicting_existing_release_fails_without_mutation(self):
        for field, changed in (
            ("body", "Different notes"), ("name", "Different name"),
            ("tag_name", "v0.2.0"), ("draft", True), ("prerelease", True),
        ):
            conflict = copy.deepcopy(self.published)
            conflict[field] = changed
            with self.subTest(field=field), patch.object(
                publish_release, "github_api", side_effect=[{"sha": COMMIT}, conflict]
            ) as api:
                with self.assertRaisesRegex(ValueError, "does not match"):
                    publish_release.publish(REPOSITORY, self.payload, COMMIT)
                self.assertEqual(2, api.call_count)

    def test_moved_tag_is_rejected_before_release_access(self):
        with patch.object(publish_release, "github_api", return_value={"sha": "f" * 40}) as api:
            with self.assertRaisesRegex(ValueError, "commit"):
                publish_release.publish(REPOSITORY, self.payload, COMMIT)
            self.assertEqual(1, api.call_count)

    def test_tools_creation_does_not_select_latest_or_read_root_latest(self):
        payload = {**self.payload, "tag_name": "tools/v0.1.0", "name": "gosr tools v0.1.0", "make_latest": "false"}
        published = {**self.published, **payload}
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, publish_release.GitHubError(404, "Not Found"), published,
        ]) as api:
            self.assertFalse(publish_release.publish(REPOSITORY, payload, COMMIT))
        self.assertEqual(3, api.call_count)
        self.assertEqual("false", api.call_args.kwargs["payload"]["make_latest"])
        self.assertIn("tools%2Fv0.1.0", api.call_args_list[0].args[0])

    def test_existing_tools_release_is_reused_without_changing_root_latest(self):
        payload = {**self.payload, "tag_name": "tools/v0.1.0", "name": "gosr tools v0.1.0", "make_latest": "false"}
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, {**self.published, **payload},
        ]) as api:
            self.assertFalse(publish_release.publish(REPOSITORY, payload, COMMIT))
        self.assertEqual(2, api.call_count)

    def test_new_older_root_release_does_not_replace_newer_latest(self):
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, publish_release.GitHubError(404, "Not Found"),
            {"tag_name": "v0.2.0"}, self.published, {"tag_name": "v0.2.0"},
        ]) as api:
            self.assertFalse(publish_release.publish(REPOSITORY, self.payload, COMMIT))
        writes = [call for call in api.call_args_list if call.kwargs.get("payload")]
        self.assertEqual(1, len(writes))
        self.assertEqual("false", writes[0].kwargs["payload"]["make_latest"])

    def test_require_latest_is_read_only_and_rejects_older_or_tools_release(self):
        with patch.object(publish_release, "github_api", side_effect=[self.published, self.published]) as api:
            publish_release.require_latest(REPOSITORY, self.payload)
        self.assertEqual(2, api.call_count)
        with patch.object(publish_release, "github_api", side_effect=[self.published, {"tag_name": "v0.2.0"}]):
            with self.assertRaisesRegex(ValueError, "current Latest"):
                publish_release.require_latest(REPOSITORY, self.payload)
        with self.assertRaisesRegex(ValueError, "root"):
            publish_release.require_latest(REPOSITORY, {**self.payload, "tag_name": "tools/v0.1.0"})

    def test_check_pages_accepts_only_action_build_and_matching_tag_policy(self):
        environment = {"deployment_branch_policy": {"custom_branch_policies": True}}
        for kind, allowed in (("branch", False), ("tag", True)):
            with self.subTest(kind=kind), patch.object(publish_release, "github_api", side_effect=[
                {"build_type": "workflow"}, environment,
                {"total_count": 1, "branch_policies": [{"type": kind, "name": "v*"}]},
            ]):
                if allowed:
                    publish_release.check_pages(REPOSITORY, TAG)
                else:
                    with self.assertRaisesRegex(ValueError, "does not allow"):
                        publish_release.check_pages(REPOSITORY, TAG)

    def test_check_pages_rejects_missing_or_incomplete_policy_metadata(self):
        for responses in (
            [{"build_type": "legacy"}],
            [{"build_type": "workflow"}, {}],
            [{"build_type": "workflow"}, {"deployment_branch_policy": {"custom_branch_policies": True}},
             {"total_count": 2, "branch_policies": [{"type": "tag", "name": "v*"}]}],
        ):
            with self.subTest(responses=responses), patch.object(publish_release, "github_api", side_effect=responses):
                with self.assertRaises(ValueError):
                    publish_release.check_pages(REPOSITORY, TAG)

    def test_previous_latest_is_retried_until_new_release_is_visible(self):
        with patch.object(publish_release, "github_api", side_effect=[
            {"tag_name": "v0.0.9"}, self.published,
        ]), patch.object(publish_release.time, "sleep") as sleep:
            self.assertTrue(publish_release.wait_for_latest(REPOSITORY, TAG))
        self.assertEqual(1, sleep.call_count)

    def test_superseded_release_does_not_deploy_or_promote_itself(self):
        with patch.object(publish_release, "github_api", return_value={"tag_name": "v0.2.0"}) as api:
            self.assertFalse(publish_release.wait_for_latest(REPOSITORY, TAG))
        self.assertEqual(1, api.call_count)

    def test_latest_404_retries_are_bounded_and_failure_is_not_success(self):
        with patch.object(
            publish_release, "github_api", side_effect=publish_release.GitHubError(404, "Not Found")
        ) as api, patch.object(publish_release.time, "sleep") as sleep:
            with self.assertRaisesRegex(ValueError, "Latest.*unavailable"):
                publish_release.wait_for_latest(REPOSITORY, TAG)
        self.assertGreater(api.call_count, 1)
        self.assertLessEqual(api.call_count, 10)
        self.assertEqual(api.call_count - 1, sleep.call_count)

    def test_server_failure_is_retried_but_permission_errors_are_not(self):
        for status in (429, 503):
            with self.subTest(status=status), patch.object(
                publish_release, "github_api", side_effect=[
                    publish_release.GitHubError(status, "Unavailable"), self.published,
                ]
            ), patch.object(publish_release.time, "sleep"):
                self.assertTrue(publish_release.wait_for_latest(REPOSITORY, TAG))
        for status in (401, 403):
            with self.subTest(status=status), patch.object(
                publish_release, "github_api", side_effect=publish_release.GitHubError(status, "Denied")
            ) as api, patch.object(publish_release.time, "sleep") as sleep:
                with self.assertRaises(publish_release.GitHubError):
                    publish_release.wait_for_latest(REPOSITORY, TAG)
                self.assertEqual(1, api.call_count)
                sleep.assert_not_called()

    def test_release_lookup_permission_error_never_creates_release(self):
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, publish_release.GitHubError(403, "Denied"),
        ]) as api:
            with self.assertRaises(publish_release.GitHubError):
                publish_release.publish(REPOSITORY, self.payload, COMMIT)
        self.assertEqual(2, api.call_count)

    def test_create_failure_is_not_retried_or_reported_as_success(self):
        with patch.object(publish_release, "github_api", side_effect=[
            {"sha": COMMIT}, publish_release.GitHubError(404, "Not Found"),
            publish_release.GitHubError(404, "No Latest"),
            publish_release.GitHubError(502, "Uncertain write result"),
        ]) as api:
            with self.assertRaises(publish_release.GitHubError):
                publish_release.publish(REPOSITORY, self.payload, COMMIT)
        self.assertEqual(4, api.call_count)

    def test_stale_older_latest_after_bounded_wait_fails_instead_of_skipping_pages(self):
        with patch.object(publish_release, "github_api", return_value={"tag_name": "v0.0.9"}) as api, patch.object(
            publish_release.time, "sleep"
        ):
            with self.assertRaisesRegex(ValueError, "Latest.*did not converge"):
                publish_release.wait_for_latest(REPOSITORY, TAG)
        self.assertGreater(api.call_count, 1)
        self.assertLessEqual(api.call_count, 10)

    def test_invalid_latest_metadata_is_not_treated_as_superseded_release(self):
        for data in ({}, {"tag_name": "not-a-stable-version"}):
            with self.subTest(data=data), patch.object(publish_release, "github_api", return_value=data):
                with self.assertRaises(ValueError):
                    publish_release.wait_for_latest(REPOSITORY, TAG)

    def test_cli_rejects_invalid_commit_or_tag_before_external_commands(self):
        for tag, commit, diagnostic in (
            (TAG, "short", "full checked source commit"),
            ("v0.1.0-rc.1", COMMIT, "stable vX.Y.Z"),
        ):
            with self.subTest(tag=tag, commit=commit):
                result = subprocess.run(
                    [sys.executable, publish_release.__file__, "--module", "root", "--tag", tag, "--commit", commit],
                    env=dict(os.environ, PATH=""), capture_output=True, text=True, check=False,
                )
                self.assertEqual(2, result.returncode)
                self.assertIn(diagnostic, result.stderr)
                self.assertEqual("", result.stdout)

    def test_github_api_uses_http_status_and_serializes_post_payload(self):
        response = 'HTTP/2.0 201 Created\nContent-Type: application/json\n\n{"id": 42}\n'
        completed = subprocess.CompletedProcess([], 0, response, "")
        with patch.object(publish_release.subprocess, "run", return_value=completed) as run:
            self.assertEqual({"id": 42}, publish_release.github_api(f"{PREFIX}/releases", payload=self.payload))
        command = run.call_args.args[0]
        self.assertEqual(["gh", "api", "--include", "--method", "POST", f"{PREFIX}/releases", "--input", "-"], command)
        self.assertEqual(self.payload, json.loads(run.call_args.kwargs["input"]))

    def test_github_api_distinguishes_404_from_auth_and_transport_failures(self):
        for status in (404, 403, 502):
            response = f'HTTP/2.0 {status} Error\nContent-Type: application/json\n\n{{"message":"Denied"}}\n'
            completed = subprocess.CompletedProcess([], 1, response, "error")
            with self.subTest(status=status), patch.object(
                publish_release.subprocess, "run", return_value=completed
            ):
                with self.assertRaises(publish_release.GitHubError) as raised:
                    publish_release.github_api(f"{PREFIX}/releases/latest")
                self.assertEqual(status, raised.exception.status)
        with patch.object(publish_release.subprocess, "run", side_effect=subprocess.TimeoutExpired(["gh"], 20)):
            with self.assertRaisesRegex(ValueError, "GitHub request"):
                publish_release.github_api(f"{PREFIX}/releases/latest")

    def test_github_api_rejects_malformed_success_response(self):
        for response in ('HTTP/2.0 200 OK\n\nnot JSON', 'HTTP/2.0 200 OK\n\n[]', 'not HTTP'):
            completed = subprocess.CompletedProcess([], 0, response, "")
            with self.subTest(response=response), patch.object(
                publish_release.subprocess, "run", return_value=completed
            ):
                with self.assertRaisesRegex(ValueError, "GitHub response"):
                    publish_release.github_api(f"{PREFIX}/releases/latest")


if __name__ == "__main__":
    unittest.main()
