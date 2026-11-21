"""Check native repository links and diagram markup in the bilingual build."""

import re
import unittest
from html.parser import HTMLParser
from pathlib import Path

from mkdocs.config import load_config


ROOT = Path(__file__).resolve().parents[1]


class Markup(HTMLParser):
    def __init__(self):
        super().__init__()
        self.repositories = []
        self.diagram_count = 0
        self.copyright = []
        self.in_copyright = False

    def handle_starttag(self, tag, attrs):
        attributes = dict(attrs)
        if "mermaid" in attributes.get("class", "").split():
            self.diagram_count += 1
        if tag == "a" and attributes.get("data-md-component") == "source":
            self.repositories.append(attributes.get("href"))
        if "md-copyright__highlight" in attributes.get("class", "").split():
            self.in_copyright = True

    def handle_endtag(self, tag):
        if tag == "div":
            self.in_copyright = False

    def handle_data(self, data):
        if self.in_copyright:
            self.copyright.append(data)


class DocumentationBuild(unittest.TestCase):
    def test_project_authorship_in_both_languages(self):
        config = load_config(config_file=str(ROOT / "mkdocs.yml"))
        owner = f"The {config.repo_name} Authors"
        self.assertEqual(config.site_author, owner)
        for locale in ("", "zh"):
            with self.subTest(locale=locale):
                html = (ROOT / "site" / locale / "index.html").read_text()
                markup = Markup()
                markup.feed(html)
                self.assertEqual(
                    "".join(markup.copyright).strip(),
                    f"Copyright © 2026 - present {owner}",
                )

    def test_native_repository_links_in_both_languages(self):
        repository = load_config(config_file=str(ROOT / "mkdocs.yml")).repo_url
        for locale in ("", "zh"):
            with self.subTest(locale=locale):
                html = (ROOT / "site" / locale / "index.html").read_text()
                markup = Markup()
                markup.feed(html)
                self.assertGreater(len(markup.repositories), 0)
                self.assertEqual(set(markup.repositories), {repository})

    def test_diagram_types_match_between_languages(self):
        for english in (ROOT / "docs").rglob("*.md"):
            if english.name.endswith(".zh.md"):
                continue
            diagrams = re.findall(r"```mermaid\n(.*?)```", english.read_text(), re.S)
            if not diagrams:
                continue
            chinese = english.with_name(english.stem + ".zh.md")
            translated = re.findall(r"```mermaid\n(.*?)```", chinese.read_text(), re.S)
            with self.subTest(page=english.name):
                self.assertEqual(
                    [block.splitlines()[0] for block in diagrams],
                    [block.splitlines()[0] for block in translated],
                )

    def test_all_diagram_fences_reach_the_site(self):
        diagrams = 0
        for path in (ROOT / "docs").rglob("*.md"):
            count = path.read_text().count("```mermaid\n")
            if not count:
                continue
            diagrams += count
            relative = path.relative_to(ROOT / "docs")
            chinese = path.name.endswith(".zh.md")
            stem = path.name.removesuffix(".zh.md" if chinese else ".md")
            target = ROOT / "site" / ("zh" if chinese else "") / relative.parent
            if stem != "index":
                target /= stem
            html = (target / "index.html").read_text()
            with self.subTest(page=str(relative)):
                markup = Markup()
                markup.feed(html)
                self.assertEqual(markup.diagram_count, count)
        self.assertGreater(diagrams, 0)


if __name__ == "__main__":
    unittest.main()
