# Repository scripts

[English](README.md) | [简体中文](README.zh_CN.md)

These Python scripts support repository checks, shared formatting and release
validation. Run the Make targets from the repository root. Use Python 3.12 and a
Go toolchain compatible with both [go.mod](../go.mod) and
[tools/go.mod](../tools/go.mod). Formatting also requires Git and Bash.
Checking a published root Release requires an authenticated GitHub CLI (`gh`).

## Entrypoints

| Script | Responsibility | Make targets |
| --- | --- | --- |
| `golayout.py` | Resolve the golayout version declared by tools and delegate to its shared scripts | `setup`, `format`, `check-format`, `test-format`, `license-headers`, `check-license-headers` |
| `tools.py` | Run tools-module commands with explicit source or published dependency policy | `test`, `race`, `vet`, `lint`, `verify-tools-source`, `verify-tools-published` |
| `release.py` | Validate module tags, identities and release notes; optionally check the required root Release | `validate-release`, `prepare-tools-release` |
| `go_module.py` | Read module metadata through `go mod edit -json`; imported by the other scripts | No CLI entrypoint |

Common checks:

```sh
make setup
make check-format
make test-scripts
```

`make format` rewrites owned Go source; `make license-headers` updates its file
headers. Repository setup and quality gates are described in the
[contributor guide](../docs/design/docs/development.md). Tag conventions,
dependency preparation and publication order are covered by the
[release guide](../docs/maintainers/releases.md).

## Module and environment policy

`GO` selects the Go executable, defaulting to `go`. The formatter's language
version comes from the root go.mod; the formatter dependency comes from tools.

`tools.py` runs commands in the tools directory:

- `source` uses temporary module files with the root dependency replaced by this
  checkout. Go checks receive `-modfile` as an argument. External tools and shared
  formatting scripts use a private workspace so the flag does not leak into
  their Go subprocesses. Temporary files are removed when the invocation exits.
- `published` disables workspaces and uses the declared dependencies. It rejects
  local directory replacements and every replacement of the root module,
  including versioned replacements. Commands such as `go mod tidy` can update
  tools/go.mod and tools/go.sum in this mode.

Both modes reject nonempty `GOFLAGS`, including flags from persistent `GOENV`
configuration. This policy controls module resolution; it does not sandbox the
executed command or disable dependency downloads.

Metadata reads disable workspaces, user flags, persistent Go configuration and
automatic toolchain selection. They parse the selected module file without
resolving its dependency graph. A successful local release validation therefore
does not establish that a dependency has been published.

## Tests and external commands

`make test-scripts` runs the Python unittest suite. It requires real Go, but does
not require GitHub authentication or a published gosr root version. Selected
external operations use fixtures under `tests/testdata`:

| Fixture | Behavior |
| --- | --- |
| `gh` | Check the complete expected request and return scenario-specific Release data or an error. No GitHub request is sent. |
| `golayout/go` | Intercept only the golayout download request. Run real Go to check its module environment; delegate all other commands to real Go. |
| `golayout/format.sh` | Check wrapper arguments, module paths and environment, then record the language version. It does not format source. |

Tests select these fixtures through `GO` or `PATH` in child-process environments.
The fixtures cover wrapper behavior and failure paths. Shared formatting behavior
is checked with real golayout through `make test-format` and `make check-format`;
published dependency resolution requires `make check-published`.
