import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import tools


class ToolEnvironmentTests(unittest.TestCase):
    def setUp(self):
        isolated = patch.dict(os.environ, {
            "GOENV": "off", "GOFLAGS": "", "GOWORK": "off",
            "GOPROXY": "off", "GOSUMDB": "off", "GOTOOLCHAIN": "local",
        })
        isolated.start()
        self.addCleanup(isolated.stop)
        self.temporary = tempfile.TemporaryDirectory(prefix="gosr tool tests ")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        (self.root / "go.mod").write_text(
            "module github.com/longxi-go/gosr\n\ngo 1.27.0\n", encoding="utf-8"
        )
        (self.root / "tools").mkdir()
        (self.root / "tools/go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\ngo 1.27.0\n\n"
            "require github.com/longxi-go/gosr v0.1.0\n",
            encoding="utf-8",
        )

    def test_source_scope_uses_a_temporary_modfile_and_preserves_the_original(self):
        original = (self.root / "tools/go.mod").read_bytes()
        with tools.environment(self.root, published=False) as env:
            self.assertEqual("off", env["GOWORK"])
            flags = env["GOFLAGS"]
            self.assertIn("-modfile=", flags)
            modfile = Path(flags.removeprefix("-modfile="))
            self.assertIn(str(self.root), modfile.read_text(encoding="utf-8"))
            self.assertEqual(original, (self.root / "tools/go.mod").read_bytes())
        self.assertFalse(modfile.exists())

    def test_published_scope_never_injects_local_module_resolution(self):
        with tools.environment(self.root, published=True) as env:
            self.assertEqual("off", env["GOWORK"])
            self.assertNotIn("-modfile=", env.get("GOFLAGS", ""))

    def test_source_tidy_preserves_published_root_sums_in_repository(self):
        (self.root / "tools/adapter.go").write_text(
            'package adapter\n\nimport _ "github.com/longxi-go/gosr"\n',
            encoding="utf-8",
        )
        (self.root / "root.go").write_text("package gosr\n", encoding="utf-8")
        sumfile = self.root / "tools/go.sum"
        # Sentinels must be removed from the scratch copy before Go reads sums.
        # They are not fabricated publication checksums.
        original = (
            "github.com/longxi-go/gosr v0.1.0 h1:fixture\n"
            "github.com/longxi-go/gosr v0.1.0/go.mod h1:fixture\n"
        )
        sumfile.write_text(original, encoding="utf-8")

        go = os.environ.get("GO", "go")
        with tools.environment(self.root, published=False) as env:
            command, child_env = tools.command_environment([go, "mod", "tidy", "-diff"], env)
            result = subprocess.run(
                command, cwd=self.root / "tools", env=child_env,
                capture_output=True, text=True, check=False,
            )
        self.assertEqual(0, result.returncode, result.stdout + result.stderr)
        self.assertEqual(original, sumfile.read_text(encoding="utf-8"))

    def test_external_module_override_cannot_contaminate_published_checks(self):
        with patch.dict(os.environ, {"GOFLAGS": f"-modfile={self.root / 'unrelated.mod'}"}):
            with self.assertRaisesRegex(ValueError, "GOFLAGS"):
                with tools.environment(self.root, published=True):
                    self.fail("external module override accepted")

    def test_go_command_keeps_module_flags_out_of_child_process_environment(self):
        flag = f"-modfile={self.root / 'source.mod'}"
        env = {"GOWORK": "off", "GOFLAGS": flag}
        go = os.environ.get("GO", "go")
        command, child_env = tools.command_environment([go, "test", "./..."], env)
        self.assertEqual([go, "test", flag, "./..."], command)
        self.assertEqual("", child_env["GOFLAGS"])
        self.assertEqual(flag, env["GOFLAGS"])

    def test_linter_scope_uses_a_workspace_without_modfile_flags(self):
        with tools.environment(self.root, published=False) as env:
            lint_env = tools.workspace_environment(self.root, env)
            self.assertEqual("", lint_env["GOFLAGS"])
            work = Path(lint_env["GOWORK"])
            text = work.read_text(encoding="utf-8")
            self.assertIn(str(self.root / "tools"), text)
            self.assertIn(str(self.root), text)
        self.assertFalse(work.exists())

    def test_persistent_go_flags_do_not_override_published_checks(self):
        config = self.root / "goenv"
        config.write_text(f"GOFLAGS=-modfile={self.root / 'unrelated.mod'}\n", encoding="utf-8")
        with patch.dict(os.environ, {"GOENV": str(config)}):
            with self.assertRaisesRegex(ValueError, "go env GOFLAGS"):
                with tools.environment(self.root, published=True):
                    self.fail("persistent modfile override accepted")

    def test_published_checks_reject_a_local_replacement_in_the_module_file(self):
        module = self.root / "tools/go.mod"
        with module.open("a", encoding="utf-8") as target:
            target.write(f"\nreplace github.com/longxi-go/gosr => {json.dumps(str(self.root))}\n")
        with self.assertRaisesRegex(ValueError, "local module replacements"):
            with tools.environment(self.root, published=True):
                self.fail("local module dependency accepted for publication")

    def test_published_checks_reject_versioned_root_replacements(self):
        module = self.root / "tools/go.mod"
        original = module.read_text(encoding="utf-8")
        for replacement in (
            "github.com/longxi-go/gosr v0.2.0-rc.1",
            "example.com/alternate-gosr v0.1.0",
        ):
            with self.subTest(replacement=replacement):
                module.write_text(
                    original + f"\nreplace github.com/longxi-go/gosr => {replacement}\n",
                    encoding="utf-8",
                )
                with self.assertRaisesRegex(ValueError, "root module replacements"):
                    with tools.environment(self.root, published=True):
                        self.fail("root release validation bypassed by a versioned replacement")

if __name__ == "__main__":
    unittest.main()
