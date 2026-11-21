# Gosr roadmap

[English](roadmap.md) | [简体中文](roadmap.zh.md)

This roadmap describes available contract areas, supporting libraries and planned
capabilities.

| Stage | Capability and usable outcome | State |
| --- | --- | --- |
| G1 | Portable IoC metadata and independent lifecycle contracts | Implemented; [IoC specification](design/docs/contracts/ioc-contracts.md) defines metadata, lookup, ownership and lifecycle behavior |
| G2 | Independent ordering, typed synchronous events and optional container notifications | Implemented; [Ordering and event specifications](design/docs/contracts/events-ordering.md) define ordering and event-delivery behavior |
| G3 | Environment values, source/resolver/profile/snapshot and typed binding contracts | Implemented; [Environment specification](design/docs/contracts/environment.md) defines data and runtime semantics |
| G4 | Optional ordered-step runner and cohesive code-generation helpers | Implemented; [pipeline](design/docs/contracts/pipeline.md) and [codegen](design/docs/contracts/codegen.md) define executable behavior and ownership |
| G5 | Reproducible source/dependency checks, bilingual specifications and independent root/tools release workflows | Workflows provided; source checks are described in [contributing](design/docs/development.md#quality-gates) |
| G6 | Versioned behavioral catalog with adapters for independent implementations | Planned |

## Scope and dependencies

G1–G3 define portable semantics that implementations can adopt independently.
G4 supplies optional mechanisms with separate responsibilities; it does not define
domain behavior or a universal compiler framework.

G5 provides source and dependency checks and separate root/tools release workflows.
The [contribution guide](design/docs/development.md#quality-gates) documents source
and dependency checks.

G6 is planned to associate specification requirements with executable adapters for
independent implementations. The specifications define the required behavior;
adapters do not impose a reference runtime. Coverage is domain-specific: a container
adapter does not cover all event and environment behavior.

Additional domains require demonstrated shared needs and a separately scoped
contract design. Runtime and tool implementation roadmaps remain with their
respective projects.
