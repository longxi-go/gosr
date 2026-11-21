import os
from pathlib import Path
import sys
import tempfile
import unittest
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from go_module import read_module


class ModuleMetadataTests(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory(prefix="gosr module tests ")
        self.addCleanup(temporary.cleanup)
        self.root = Path(temporary.name)
        self.module = self.root / "go.mod"

    def test_reads_only_the_selected_file_without_resolving_dependencies(self):
        source = (
            '// module example.com/comment\nmodule "example.com/selected"\n\n'
            'go\t1.27.0\nrequire "example.invalid/unpublished" v1.4.7\n'
        )
        self.module.write_text(source, encoding="utf-8")
        # These overrides would break a normal Go command. Metadata parsing must
        # ignore them and must not contact the unavailable dependency proxy.
        with patch.dict(os.environ, {
            "GOWORK": str(self.root / "missing.work"), "GOFLAGS": "-invalid-flag",
            "GOENV": str(self.root / "missing.env"), "GOPROXY": "off",
            "GOTOOLCHAIN": "invalid-toolchain",
        }):
            data = read_module(self.module)
        self.assertEqual("example.com/selected", data["Module"]["Path"])
        self.assertEqual("1.27.0", data["Go"])
        self.assertEqual(1, len(data["Require"]))
        self.assertEqual("example.invalid/unpublished", data["Require"][0]["Path"])
        self.assertEqual("v1.4.7", data["Require"][0]["Version"])
        self.assertEqual(source, self.module.read_text(encoding="utf-8"))

    def test_missing_file_and_invalid_syntax_include_the_selected_path(self):
        with self.assertRaisesRegex(ValueError, "Cannot read Go module") as failure:
            read_module(self.module)
        self.assertIn(str(self.module), str(failure.exception))
        self.module.write_text("not a module directive\n", encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "Cannot read Go module") as failure:
            read_module(self.module)
        self.assertIn(str(self.module), str(failure.exception))

    def test_missing_go_executable_reports_a_diagnostic(self):
        with patch.dict(os.environ, {"GO": str(self.root / "missing-go")}):
            with self.assertRaisesRegex(ValueError, "Cannot read Go module"):
                read_module(self.module)


if __name__ == "__main__":
    unittest.main()
