# Generator tooling adapters

[English](README.md) | [简体中文](README.zh_CN.md)

The module `github.com/longxi-go/gosr/tools` provides optional package-loading and
annotation adapters for Go tools. It has its own module path and version, keeping
x/tools and controller-tools dependencies separate from the root module's
standard-library-only contracts and helpers.

| Package | Provides | Caller controls |
| --- | --- | --- |
| `pkg/packageload` | go/packages loading with copied configuration and overlays; deterministic diagnostics | Roots, loading mode, flags, environment, module policy and isolation |
| `pkg/annotation` | controller-tools marker parsing with source positions | Schemas, placement, repetition and domain validation |

`packageload.Load` accepts explicit roots. Nil config uses go/packages defaults.
Discovery may accept incomplete types; call `Errors` when correctness requires
validated packages. A context or overlay does not isolate filesystem access or
module downloads.

Annotation parsing uses supplied line-comment doc groups. The prefix must start
with `+` and contain no whitespace; `+` alone selects all markers. Each definition's
`Sample` must carry a type; use `struct{}{}` for a no-argument marker. A type's own
doc group overrides its declaration's doc group. Grouped-declaration fallback is
per selected type. These adapters do not define a domain annotation language or
a compiler pipeline.

Integer conversion follows controller-tools. For bounded numeric fields, use a
wide integer or string sample and validate the range before narrowing the value.

See [code-generation mechanisms](../docs/design/docs/contracts/codegen.md) for
API guarantees and integration boundaries. Executable examples are included in
each package's `example_test.go`.

See [contributing](../docs/design/docs/development.md) for repository setup and
source checks.

Module release procedures are documented in the [maintainer guide](../docs/maintainers/releases.md).
