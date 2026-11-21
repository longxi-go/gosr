# Development

## Quality gates {#quality-gates}

Use the Go versions and linter selected by the Quality workflow; the language
minimum is declared in go.mod. Install Git, Bash, Make, Python 3.12 and the linter:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

From the gosr repository root:

```sh
make setup
make license-headers format
make check
make docs-build
```

`setup` downloads the golayout version declared in tools/go.mod into the module
cache and verifies its checksums. Formatting and source-header operations use
that module's scripts. The full source check covers module tidiness, headers,
layout, formatter regressions, vet, coverage, shuffled race tests, lint and script
tests. A single checkout is sufficient.

`make check` runs root checks with GOWORK disabled and checks tools against this
root checkout through a temporary module replacement. Repository module files
are not changed. `make verify-tools-published` checks tools module tidiness and
`make check-published` runs the full gate against its declared dependency versions
without local overrides; those versions must be resolvable.

Run module checks with GOFLAGS unset. They reject inherited module/build flags,
including settings persisted in GOENV. Published checks reject local directory
replacements and all replacements of the root gosr dependency in tools/go.mod.

## Formatting and source notices

`make format` and `make check-format` use golayout's structural layout and gofumpt
composition. Nonempty executable blocks and compound declarations expand. Empty
`struct{}`, `interface{}`, `select {}` and functions may remain compact. File
responsibilities and logical line breaks require review. Generated Go follows
the same layout policy as handwritten source.

`codequality/HEADER` supplies Go source notices. `make license-headers` updates
owned files; `make check-license-headers` checks them without writes. Complete
`testdata/**/*.go.tmpl` fixtures participate in formatting; source-header updates
exclude testdata. Backups, caches, vendor and symlinks are excluded.

## Local workspace {#local-workspace}

To develop a separate module against a gosr checkout, create a workspace in their
common parent directory. This example assumes your module is in `my-app`:

```sh
go work init ./gosr ./my-app
```

Use `go work use` when a workspace already exists. Select `./gosr/tools` explicitly
if your tool imports its adapters; selecting the parent module does not include
nested modules. Use local source selection for development, then verify declared
module dependencies independently before releasing a consumer.

Keep go.work and go.work.sum local, and keep local directory replacements out of
production module files. Workspace selection tests source integration; the
published-dependency gates test standalone module resolution.

## Documentation

Python 3.12 and uv build the manual:

```sh
make docs-setup
make docs-serve
make docs-build
```

Documentation sources are under `docs/design/docs`; navigation and language
settings are in `docs/design/mkdocs.yml`. Pages use paired `page.md` and
`page.zh.md` files. Generated HTML and the Python environment are excluded from
version control.
