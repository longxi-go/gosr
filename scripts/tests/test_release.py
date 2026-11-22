import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import release


class ReleaseTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        (self.root / "go.mod").write_text(
            "module github.com/longxi-go/gosr\n\ngo 1.27.0\n", encoding="utf-8"
        )
        tools = self.root / "tools"
        tools.mkdir()
        (tools / "go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\n"
            "go 1.27.0\n\nrequire github.com/longxi-go/gosr v0.1.0\n",
            encoding="utf-8",
        )

    def write_notes(self, path, title):
        notes = self.root / "docs/releases" / path
        notes.parent.mkdir(parents=True, exist_ok=True)
        notes.write_text(f"# {title}\n\nInitial release.\n", encoding="utf-8")

    def test_accepts_root_module_tag_and_matching_notes(self):
        self.write_notes("v0.1.0.md", "gosr v0.1.0")
        self.assertEqual("0.1.0", release.validate_tag(self.root, "root", "v0.1.0"))

    def test_accepts_nested_tools_tag_and_matching_notes(self):
        self.write_notes("tools/v0.1.0.md", "gosr tools v0.1.0")
        self.assertEqual(
            "0.1.0", release.validate_tag(self.root, "tools", "tools/v0.1.0")
        )

    def test_payload_preserves_module_notes_and_root_tools_latest_policy(self):
        for module, tag, title, latest in (
            ("root", "v0.1.0", "gosr v0.1.0", "true"),
            ("tools", "tools/v0.1.0", "gosr tools v0.1.0", "false"),
        ):
            with self.subTest(module=module):
                self.write_notes(f"{tag}.md", title)
                self.assertEqual({
                    "tag_name": tag, "name": title,
                    "body": f"# {title}\n\nInitial release.\n",
                    "draft": False, "prerelease": False, "make_latest": latest,
                }, release.release_payload(self.root, module, tag))

    def test_rejects_wrong_module_tag_prefix_and_unstable_versions(self):
        for module, tag in [
            ("root", "tools/v0.1.0"),
            ("tools", "v0.1.0"),
            ("root", "v00.1.0"),
            ("root", "v0.1.0-rc.1"),
            ("tools", "tools/v2.0.0"),
        ]:
            with self.subTest(module=module, tag=tag), self.assertRaises(ValueError):
                release.validate_tag(self.root, module, tag)

    def test_tools_release_requires_a_stable_root_dependency_version(self):
        self.write_notes("tools/v0.1.0.md", "gosr tools v0.1.0")
        module = self.root / "tools/go.mod"
        original = module.read_text(encoding="utf-8")
        for version, message in (
            ("v0.0.0", "stable gosr dependency version"),
            ("v0.1.0-rc.1", "stable gosr dependency version"),
            ("v2.0.0", "Cannot read Go module"),
        ):
            with self.subTest(version=version):
                module.write_text(original.replace("v0.1.0", version), encoding="utf-8")
                with self.assertRaisesRegex(ValueError, message):
                    release.validate_tag(self.root, "tools", "tools/v0.1.0")

    def test_rejects_missing_or_mismatched_release_notes(self):
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "root", "v0.1.0")
        self.write_notes("v0.1.0.md", "gosr v0.2.0")
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "root", "v0.1.0")

    def test_invalid_module_syntax_cannot_supply_identity(self):
        self.write_notes("v0.1.0.md", "gosr v0.1.0")
        (self.root / "go.mod").write_text(
            "/*\nmodule github.com/longxi-go/gosr\n*/\n"
            "module example.com/different\n\ngo 1.27.0\n",
            encoding="utf-8",
        )
        # go.mod supports line comments, not block comments.
        with self.assertRaisesRegex(ValueError, "Cannot read Go module"):
            release.validate_tag(self.root, "root", "v0.1.0")

    def test_required_root_version_ignores_versioned_replacements(self):
        (self.root / "tools/go.mod").write_text(
            "module github.com/longxi-go/gosr/tools\n\ngo 1.27.0\n\n"
            "replace (\n"
            "github.com/longxi-go/gosr v0.1.0 => example.com/fork v0.1.0\n"
            ")\n"
            "require github.com/longxi-go/gosr v0.2.0-rc.1\n",
            encoding="utf-8",
        )
        with self.assertRaisesRegex(ValueError, "stable gosr dependency version"):
            release.required_root_version(self.root)

    def test_quoted_module_paths_and_independent_versions_are_supported(self):
        (self.root / "go.mod").write_text(
            'module "github.com/longxi-go/gosr"\n\ngo\t1.27.0\n', encoding="utf-8"
        )
        (self.root / "tools/go.mod").write_text(
            'module "github.com/longxi-go/gosr/tools"\n\ngo\t1.27.0\n\n'
            'require "github.com/longxi-go/gosr" v1.4.7\n',
            encoding="utf-8",
        )
        self.write_notes("v1.4.7.md", "gosr v1.4.7")
        self.write_notes("tools/v0.9.2.md", "gosr tools v0.9.2")
        self.assertEqual("1.4.7", release.validate_tag(self.root, "root", "v1.4.7"))
        self.assertEqual("0.9.2", release.validate_tag(self.root, "tools", "tools/v0.9.2"))
        self.assertEqual("v1.4.7", release.required_root_version(self.root))


if __name__ == "__main__":
    unittest.main()
