import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

SCRIPT = Path(__file__).resolve().parents[1] / "release.py"


class ReleaseCLITests(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.root = Path(temporary.name)
        (self.root / "tools").mkdir()
        (self.root / "go.mod").write_text(
            "module github.com/longxi-go/gosr\n\ngo 1.27.0\n", encoding="utf-8"
        )
        (self.root / "tools/go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\ngo 1.27.0\n\n"
            "require github.com/longxi-go/gosr v0.1.0\n",
            encoding="utf-8",
        )
        for tag, title in (("v0.1.0", "gosr"), ("tools/v0.1.0", "gosr tools")):
            notes = self.root / "docs/releases" / f"{tag}.md"
            notes.parent.mkdir(parents=True, exist_ok=True)
            notes.write_text(f"# {title} v0.1.0\n\nRelease notes.\n", encoding="utf-8")
        self.output = self.root / "github-output"

    def run_release(self, module, *flags, tag=None, env=None):
        if tag is None:
            tag = "v0.1.0" if module == "root" else "tools/v0.1.0"
        return subprocess.run(
            [
                sys.executable,
                str(SCRIPT),
                "--root", str(self.root),
                "--module", module,
                "--tag", tag,
                "--github-output", str(self.output),
                *flags,
            ],
            env=env,
            capture_output=True,
            text=True,
            check=False,
            timeout=30,
        )

    def test_root_metadata_preserves_existing_step_outputs(self):
        self.output.write_text("existing=value\n", encoding="utf-8")
        result = self.run_release("root")
        self.assertEqual(0, result.returncode, result.stderr)
        self.assertEqual("0.1.0\n", result.stdout)
        self.assertEqual(
            "existing=value\nmodule=root\ntitle=gosr v0.1.0\n",
            self.output.read_text(encoding="utf-8"),
        )

    def test_tools_metadata_names_the_exact_root_dependency(self):
        result = self.run_release("tools")
        self.assertEqual(0, result.returncode, result.stderr)
        self.assertEqual(
            "module=tools\ntitle=gosr tools v0.1.0\nroot_version=v0.1.0\n",
            self.output.read_text(encoding="utf-8"),
        )

    def test_invalid_notes_do_not_emit_release_metadata(self):
        (self.root / "docs/releases/tools/v0.1.0.md").unlink()
        result = self.run_release("tools")
        self.assertEqual(2, result.returncode)
        self.assertFalse(self.output.exists())
        self.assertIn("Missing release notes", result.stderr)

    def github_environment(self, version, repository):
        cli = self.root / "bin"
        cli.mkdir(exist_ok=True)
        fake = cli / "gh"
        shutil.copyfile(Path(__file__).parent / "testdata/gh", fake)
        fake.chmod(0o755)
        env = dict(os.environ, PATH=f"{cli}{os.pathsep}{os.environ['PATH']}")
        env["RELEASE_TEST_RESPONSE"] = str(self.root / "response.json")
        env["RELEASE_TEST_EXPECTED_REQUEST"] = json.dumps([
            "release", "view", version, "--repo", repository,
            "--json", "tagName,isDraft,isPrerelease",
        ])
        return env

    def test_tools_release_requires_a_matching_published_stable_root_release(self):
        env = self.github_environment("v0.1.0", "longxi-go/gosr")
        response = Path(env["RELEASE_TEST_RESPONSE"])
        cases = [
            ("published", {"tagName": "v0.1.0", "isDraft": False, "isPrerelease": False}, 0),
            ("draft", {"tagName": "v0.1.0", "isDraft": True, "isPrerelease": False}, 2),
            ("prerelease", {"tagName": "v0.1.0", "isDraft": False, "isPrerelease": True}, 2),
            ("wrong-tag", {"tagName": "v0.2.0", "isDraft": False, "isPrerelease": False}, 2),
            ("missing-fields", {"tagName": "v0.1.0"}, 2),
            ("missing-release", None, 2),
            ("malformed-response", "invalid JSON", 2),
        ]
        for name, data, expected in cases:
            with self.subTest(scenario=name):
                self.output.unlink(missing_ok=True)
                response.write_text(json.dumps(data), encoding="utf-8")
                result = self.run_release(
                    "tools", "--check-root-release", "--repository", "longxi-go/gosr", env=env
                )
                self.assertEqual(expected, result.returncode, result.stderr)
                self.assertEqual(expected == 0, self.output.exists())
                if expected:
                    self.assertIn("published stable root release v0.1.0", result.stderr)

    def test_tools_version_and_repository_are_independent_of_root_release(self):
        (self.root / "tools/go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\ngo 1.27.0\n\n"
            "require github.com/longxi-go/gosr v1.4.7\n", encoding="utf-8"
        )
        (self.root / "docs/releases/tools/v0.9.2.md").write_text(
            "# gosr tools v0.9.2\n\nAdapter release.\n", encoding="utf-8"
        )
        env = self.github_environment("v1.4.7", "example-org/gosr-fork")
        Path(env["RELEASE_TEST_RESPONSE"]).write_text(
            json.dumps({"tagName": "v1.4.7", "isDraft": False, "isPrerelease": False}),
            encoding="utf-8",
        )
        result = self.run_release(
            "tools", "--check-root-release", "--repository", "example-org/gosr-fork",
            tag="tools/v0.9.2", env=env,
        )
        self.assertEqual(0, result.returncode, result.stderr)
        self.assertEqual("0.9.2\n", result.stdout)
        self.assertEqual(
            "module=tools\ntitle=gosr tools v0.9.2\nroot_version=v1.4.7\n",
            self.output.read_text(encoding="utf-8"),
        )


if __name__ == "__main__":
    unittest.main()
