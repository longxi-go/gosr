# Documentation site

[English](README.md) | [简体中文](README.zh.md)

The gosr manual uses MkDocs Material and contains English and Simplified Chinese
pages. Building it requires Python 3.12 and uv; dependencies are locked in `uv.lock`.

## Build and preview

Run from the repository root:

```sh
make docs-setup
make docs-serve
make docs-build
```

`docs-serve` previews English at `http://127.0.0.1:8000/gosr/` and Chinese at
`http://127.0.0.1:8000/gosr/zh/`. `docs-build` runs the strict build and site tests.
Mermaid diagrams render in the browser during preview.

## Source files

| Path within `docs/design` | Purpose |
| --- | --- |
| `mkdocs.yml` | Theme, navigation and language configuration |
| `docs/` | English pages and matching `*.zh.md` translations |
| `uv.lock` | Locked documentation dependencies |
| `tests/` | Built-site regression checks |

Page files are paired by basename. Add the page to both navigation trees in
`mkdocs.yml`. Generated `site/` and `.venv/` directories are excluded from version
control. Set `DOCS_SITE_URL` to a URL ending in `/` to override the base URL.
