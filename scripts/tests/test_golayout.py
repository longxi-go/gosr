import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

SCRIPTS = Path(__file__).resolve().parents[1]
FIXTURES = Path(__file__).parent / "testdata/golayout"


class GolayoutCLITests(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory(prefix="gosr formatter tests ")
        self.addCleanup(temporary.cleanup)
        self.root = Path(temporary.name)
        for directory in ("scripts", "tools", "temporary files", "formatter/scripts"):
            (self.root / directory).mkdir(parents=True)
        for name in ("golayout.py", "tools.py", "go_module.py"):
            shutil.copyfile(SCRIPTS / name, self.root / "scripts" / name)
        (self.root / "go.mod").write_text(
            "module github.com/longxi-go/gosr\n\ngo\t1.27.0\n", encoding="utf-8"
        )
        (self.root / "tools/go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\ngo 1.27.0\n"
            "require github.com/longxi-go/gosr v0.1.0\n", encoding="utf-8"
        )
        go = self.root / "go"
        shutil.copyfile(FIXTURES / "go", go)
        go.chmod(0o755)
        shutil.copyfile(FIXTURES / "format.sh", self.root / "formatter/scripts/format.sh")
        self.env = dict(
            os.environ, GO=str(go), GOFLAGS="", GOENV="off", GOWORK="off",
            GOPROXY="off", GOTOOLCHAIN="local", TMPDIR=str(self.root / "temporary files"),
            GOLAYOUT_TEST_GO=shutil.which(os.environ.get("GO", "go")),
            GOLAYOUT_TEST_MODULE=str(self.root / "formatter"),
            GOLAYOUT_TEST_TOOLS=str(self.root / "tools"),
            GOLAYOUT_TEST_LANGUAGE=str(self.root / "language"),
        )
        self.originals = {name: (self.root / name).read_bytes() for name in ("go.mod", "tools/go.mod")}

    def run_formatter(self):
        result = subprocess.run(
            [sys.executable, str(self.root / "scripts/golayout.py"), "setup"],
            env=self.env, capture_output=True, text=True, check=False, timeout=30,
        )
        for name, original in self.originals.items():
            self.assertEqual(original, (self.root / name).read_bytes())
        self.assertEqual([], list((self.root / "temporary files").iterdir()))
        return result

    def test_paths_with_spaces_and_tab_separated_go_directive(self):
        result = self.run_formatter()
        self.assertEqual(0, result.returncode, result.stdout + result.stderr)
        self.assertEqual("go1.27.0\n", (self.root / "language").read_text(encoding="utf-8"))

    def test_malformed_download_metadata_reports_a_diagnostic(self):
        for response in ("not JSON", "[]", "{}", json.dumps({"Dir": 3})):
            with self.subTest(response=response):
                self.env["GOLAYOUT_TEST_RESPONSE"] = response
                result = self.run_formatter()
                self.assertNotEqual(0, result.returncode)
                self.assertIn("Invalid golayout download metadata", result.stderr)
                self.assertNotIn("Traceback", result.stderr)


if __name__ == "__main__":
    unittest.main()
