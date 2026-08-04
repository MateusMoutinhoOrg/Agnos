# Project Structure

This document maps the project **schema** — the kinds of files the project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/References/Specs.md) to get its description and sample.

The project is split into three top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  examples/
(reaches the OS)  (closed)     (wires the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the pure library. Nothing inside it may import `adapters/`, `examples/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each adapter imports `sandbox/contracts/deps` and nothing else from the sandbox.
- **`/examples/`** sits outside the sandbox too, and is the only place where an adapter and the sandbox meet.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview, quick start, Doc Index, and Samples section | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition and dependencies | |
| `.gitignore` | Intentionally untracked files to ignore | |

---

## `/sandbox/`
The closed sandbox — the pure library. It holds its own entry point, the contracts everything is wired through, and the internal implementation. It reaches nothing outside itself: every OS-bound or third-party effect arrives through the injected `Deps`. Its package is named `lib`, so consumers import it as `lib "…/sandbox"` and call `lib.New`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor storing `Deps` on `api.Lib` and running the internal factories over it | |

### `/sandbox/contracts/`
The structs the rest of the project is wired through — the only part of the sandbox anything outside it may import. Contracts hold the project's **public types** and are structs of function fields, never interfaces; see [StructContracts.md](/docs/Explanations/StructContracts.md). Contracts import nothing from `adapters/` or `sandbox/internal/`.

#### `/sandbox/contracts/deps/`
The contract every adapter must fill.

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | The `Deps` struct, one function field per injectable behavior, plus one plain field per embedded library | Deps |

##### `/sandbox/contracts/deps/verbdeps/`
The sandbox's copy of the embedded [Verb](https://github.com/MateusMoutinhoOrg/Verb) argv-parser library's public api. The sandbox may not import Verb — that would be a third-party import — so it restates the shape it needs, field for field; the adapter, outside the sandbox, is what fills it. Same mechanic as `bootstrap/sandbox/contracts/deps/agnosdeps/`.

| File | Description | Spec |
|------|-------------|------|
| `verbdeps.go` | Copy of the embedded Verb library's `api.Lib` struct, injected whole as the `Deps.VerbLib` field | |

#### `/sandbox/contracts/api/`
The structs the library hands back to callers.

| File | Description | Spec |
|------|-------------|------|
| `api.go` | The `Lib` entry-point struct plus one struct per object the lib creates, each carrying a `Deps` field | Outputs |

### `/sandbox/internal/`
**Factories only** — no types. Each package here holds the functions that take a pointer to an [`api`](#sandboxcontractsapi) struct and return closures reading that struct's `Deps`, which the package's `New` constructor assigns into the matching function fields. Types never live here; they stay in `contracts/`. Go's `internal/` rule makes this tree unreachable from outside `sandbox/`, so neither consumers nor `adapters/` can reach in — the sandbox wall is enforced by the compiler, not by convention alone.

#### `/sandbox/internal/lib/`
The entry-point implementation. The `internal/` parent already marks it private, so the package carries no `internal_` prefix.

| File | Description | Spec |
|------|-------------|------|
| `lib.go` | One `<Field>Factory(l *api.Lib)` per lib function, each returning a closure, plus the `New(d deps.Deps) api.Lib` constructor that assigns every factory's return value and runs them all | LibFunctions |

#### `/sandbox/internal/<object>/`
One package per object the library creates, named after the object itself.

| File | Description | Spec |
|------|-------------|------|
| `<object>.go` | The object's `<Field>Factory` functions, each returning a closure, plus the `New(d deps.Deps, …) api.<Object>` constructor that propagates `Deps` and assigns every factory's return value | LibObjects |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract, each providing a distinct concrete behavior. This is where OS-bound and third-party code lives; an adapter imports `sandbox/contracts/deps` and nothing else from `sandbox/`. An adapter fills its contract with the same **factories** [`sandbox/internal/`](#sandboxinternal) uses — the carrier is the adapter struct, which declares the `Deps` field the factories' return values are assigned into.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory(a *<Name>Adapter)` per `Deps` field returning a closure, plus the `New(...) deps.Deps` constructor that assigns every factory's return value and runs them all | Adapters |

---

## `/examples/`
Outside the sandbox. Runnable examples demonstrating how to use the library — the only place an adapter and the sandbox are wired together.

### `/examples/<example>/`

| File | Description | Spec |
|------|-------------|------|
| `<example>.go` | Self-contained `package main` wiring an adapter into the lib | Examples |

**Run an example:**
```sh
go run ./examples/<example>/<example>.go
```

---

## `/bootstrap/`
A second, self-contained Agnos library — same three trees (`sandbox/`, `adapters/`, `examples/`) and the same rules — demonstrating how one Agnos-compliant library **embeds** another. Its sandbox reaches nothing outside itself, so it never imports the root library: the embedded library arrives as one plain `Deps` field.

| Path | Description |
|------|-------------|
| `sandbox/contracts/deps/deps.go` | The `Deps` struct, including `CacheLib` — the embedded library, held as a locally declared contract struct |
| `sandbox/contracts/deps/agnosdeps/agnosdeps.go` | Copy of the embedded library's `api` structs, declared inside the sandbox so the sandbox never imports the embedded library |
| `adapters/<name>/<name>.go` | Its `CacheLibFactory` initializes the embedded library with the embedded library's own adapter, and copies its `api` fields onto the local `agnosdeps` ones |
| `examples/<example>/<example>.go` | Self-contained `package main` wiring a bootstrap adapter into the bootstrap lib |

The copying lives in the adapter because only code outside the sandbox may import the embedded library. Because both sides are structs of function fields, the copy is field assignment: a wrapper is needed only where a named type differs between the two declarations. See [StructContracts.md](/docs/Explanations/StructContracts.md).

---

## `/docs/`
Documentation of the project.

### `/docs/References/`
Listable material — structures, rules, specifications, and the public API index.

| File | Description | Spec |
|------|-------------|------|
| `RULES.md` | Rules to follow when contributing to this project | Rules |
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `PublicApi.md` | Index of all public-facing components, linking to their detail pages | ReferenceDocs |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |
| `TemplateFileActions.md` | The action each template file takes when forking or adapting a library | ReferenceDocs |
| `<Name>.md` | Any other reference page the library needs | ReferenceDocs |

#### `/docs/References/Meta/`
The specifications describing how each kind of file in the project must be shaped. Never browse this directory — locate a specification by reading `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/Specs.md` | The required shape of the artifact the specification governs | |
| `<Spec>/sample.<ext>` | Concrete reference implementation of the specification | |

#### `/docs/References/PublicApi/`
Detailed documentation for each individual public API entry.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | One detail page per public struct, function, or field | ReferenceDocs |

### `/docs/Explanations/`
Explanations of the project's mechanics and features.

| File | Description | Spec |
|------|-------------|------|
| `<Topic>.md` | One page per mechanic the library needs explained | ExplanationDocs |

### `/docs/Tutorials/`
Workflow guides explaining how to use, extend, and maintain the project. Each file covers a single goal.

| File | Description | Spec |
|------|-------------|------|
| `<Goal>.md` | One page per workflow the library's maintainers repeat | TutorialDocs |
