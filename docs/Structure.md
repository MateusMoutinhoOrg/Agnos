# Project Structure

This document maps the project **schema** — the kinds of files the project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/Specs.md) to get its description and sample.

The project is a **CLI** whose interface lives inside the library. It is split into three top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the pure library, and the command-line interface with it. Nothing inside it may import `adapters/`, `cmd/`, `examples/libraryExamples/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation.md](/docs/SandboxIsolation.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each adapter imports `sandbox/contracts/deps` and nothing else from the sandbox.
- **`/cmd/`** and **`/examples/libraryExamples/`** sit outside the sandbox too, and are the only places where an adapter and the sandbox meet — `cmd/` for the installable binary, `examples/libraryExamples/` for the runnable Go samples.

Because the interface is `api.Lib.Sandboxmain` — one field of the library like any other — the binary in `cmd/main/` holds no command, no flag, and no output of its own: it wires, runs, and exits.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview, quick starts, Doc Index, and the two Examples sections | Readme |
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
The structs the rest of the project is wired through — the only part of the sandbox anything outside it may import. Contracts hold the project's **public types** and are structs of function fields, never interfaces; see [StructContracts.md](/docs/StructContracts.md). Contracts import nothing from `adapters/` or `sandbox/internal/`.

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

##### `/sandbox/contracts/deps/keepdeps/`
The sandbox's copy of the embedded [Keep](https://github.com/MateusMoutinhoOrg/Keep) schema-database library's public api, for the same reason `verbdeps/` exists. Keep's api is a tree of structs — `Lib` hands back a `KeepDatabase`, which hands back a `SchemaInstance`, which hands back `SchemaItem`s — so the copy restates the whole tree, and the adapter converts each returned struct as it crosses.

| File | Description | Spec |
|------|-------------|------|
| `keepdeps.go` | Copy of the embedded Keep library's api structs and constants, injected whole as the `Deps.KeepLib` field | |

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
| `lib.go` | One `<Field>Factory(l *api.Lib)` per lib function — `Sandboxmain` among them — each returning a closure, plus the `New(d deps.Deps) api.Lib` constructor that assigns every factory's return value and runs them all | LibFunctions |

#### `/sandbox/internal/cli/`
The command-line interface itself: the command dispatch `Sandboxmain` delegates to, the usage screen, and the reading of amounts off the command line. It reads the command line through `deps.Deps.VerbLib` and writes every line through `deps.Deps.Printf`, so the whole interface stays inside the closed sandbox. Like `store/`, it is neither an object nor the entry point, so no specification governs it, and it declares **no types and no factories**.

| File | Description | Spec |
|------|-------------|------|
| `cli.go` | The `Run(l *api.Lib, args []string) int` dispatch, one helper per command group, the usage screen, and the amount parser | |

#### `/sandbox/internal/<object>/`
One package per object the library creates, named after the object itself — `category/` and `transaction/` for this library.

| File | Description | Spec |
|------|-------------|------|
| `<object>.go` | The object's `<Field>Factory` functions, each returning a closure, plus the `New(d deps.Deps, …) api.<Object>` constructor that propagates `Deps` and assigns every factory's return value | LibObjects |

#### `/sandbox/internal/store/`
Shared helpers over the injected database, used by `internal/lib/` and by every object package: the schema the tracker's records are persisted under, the lookups that reach it, and the encoding of a transaction's reference. It declares **no types and no factories** — it is the one internal package that is neither an object nor the entry point, so no specification governs it.

| File | Description | Spec |
|------|-------------|------|
| `store.go` | The database `Props`, the field-name constants, and the helpers that read a stored record through `deps.Deps.KeepLib` | |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract, each providing a distinct concrete behavior. This is where OS-bound and third-party code lives; an adapter imports `sandbox/contracts/deps` and nothing else from `sandbox/`. An adapter fills its contract with the same **factories** [`sandbox/internal/`](#sandboxinternal) uses — the carrier is the adapter struct, which declares the `Deps` field the factories' return values are assigned into.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory(a *<Name>Adapter)` per `Deps` field returning a closure, plus the `New(...) deps.Deps` constructor that assigns every factory's return value and runs them all | Adapters |

---

## `/cmd/`
Outside the sandbox. The executables the project ships — for this project, the single binary a user installs.

### `/cmd/main/`

| File | Description | Spec |
|------|-------------|------|
| `main.go` | Self-contained `package main` that wires an adapter into the lib, calls `api.Lib.Sandboxmain(os.Args[1:])`, and exits with its return | CliMain |

**Run the CLI from source:**
```sh
go run ./cmd/main <command> [arguments]
```

**Install it globally:**
```sh
go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3
```

---

## `/examples/cliExamples/`
Outside the sandbox. Shell scripts driving the built binary the way a user would from a terminal. Each one builds the CLI into a scratch directory and points it at a budget of its own, so nothing a script does touches the records in the user's home.

| File | Description | Spec |
|------|-------------|------|
| `<Name>.sh` | Self-contained shell script demonstrating one goal against the built CLI | CliExamples |

**Run a CLI example:**
```sh
bash ./examples/cliExamples/ManageCategories.sh
```

---

## `/examples/libraryExamples/`
Outside the sandbox. Runnable Go examples demonstrating how to use the library from code, when the CLI is not what the caller wants.

### `/examples/libraryExamples/<example>/`

| File | Description | Spec |
|------|-------------|------|
| `<example>.go` | Self-contained `package main` wiring an adapter into the lib | LibraryExamples |

**Run an example:**
```sh
go run ./examples/libraryExamples/<example>/<example>.go
```

---

## `/bootstrap/`
A second, self-contained Agnos library — same three trees (`sandbox/`, `adapters/`, `examples/libraryExamples/`) and the same rules — demonstrating how one Agnos-compliant library **embeds** another. Its sandbox reaches nothing outside itself, so it never imports the root library: the embedded library arrives as one plain `Deps` field.

| Path | Description |
|------|-------------|
| `sandbox/contracts/deps/deps.go` | The `Deps` struct, including `TrackerLib` — the embedded library, held as a locally declared contract struct |
| `sandbox/contracts/deps/agnosdeps/agnosdeps.go` | Copy of the embedded library's `api` structs, declared inside the sandbox so the sandbox never imports the embedded library |
| `adapters/<name>/<name>.go` | Its `TrackerLibFactory` initializes the embedded library with the embedded library's own adapter, and copies its `api` fields onto the local `agnosdeps` ones |
| `examples/libraryExamples/<example>/<example>.go` | Self-contained `package main` wiring a bootstrap adapter into the bootstrap lib |

The copying lives in the adapter because only code outside the sandbox may import the embedded library. Because both sides are structs of function fields, the copy is field assignment: a wrapper is needed only where a named type differs between the two declarations. See [StructContracts.md](/docs/StructContracts.md).

---

## `/docs/`
Documentation of the project, kept **flat**: reference pages, tutorials, and explanations sit side by side, indexed by the [README](/README.md) Doc Index. The kind of a page decides its specification, not its path. Only listable pages covered by an index page get a directory of their own — the specifications under `Meta/` and the API detail pages under `PublicApi/`.

| File | Description | Spec |
|------|-------------|------|
| `RULES.md` | Rules to follow when contributing to this project | Rules |
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `PublicApi.md` | Index of all public-facing components, linking to their detail pages | ReferenceDocs |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |
| `TemplateFileActions.md` | The action each template file takes when forking or adapting a library | ReferenceDocs |
| `Cli.md` | Every command, flag, and exit code of the command-line interface | ReferenceDocs |
| `<Name>.md` | Any other reference page the library needs | ReferenceDocs |
| `<Topic>.md` | One page per mechanic the library needs explained | ExplanationDocs |
| `<Goal>.md` | One page per workflow the library's maintainers repeat | TutorialDocs |

### `/docs/PublicApi/`
One detail page per public-facing component, indexed by `PublicApi.md`. Reach a page through that index rather than by browsing the directory.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | One detail page per public struct, function, or field, named after the package the symbol is declared in | ReferenceDocs |

### `/docs/Meta/`
The specifications describing how each kind of file in the project must be shaped. Never browse this directory — locate a specification by reading `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/Specs.md` | The required shape of the artifact the specification governs | |
| `<Spec>/sample.<ext>` | Concrete reference implementation of the specification | |
