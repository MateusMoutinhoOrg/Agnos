# Project Structure

This document maps the project **schema** — the kinds of files a project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/Development/References/Specs.md) to get its description and sample. Use them to build a project from scratch, even before any code exists.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and quick-start guide | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition | |

---

## `/docs/`
Documentation of the project, one directory per theme.

### `/docs/<Theme>/`

| File | Description | Spec |
|------|-------------|------|
| `Index.md` | The theme's entry point: one table of its Protocols, one of its References | Index |
| `Protocols/<Goal>.md` | One page per workflow, its title phrased as the action it performs | TutorialDocs |
| `References/<Name>.md` | One page per lookup table or explained mechanic | ReferenceDocs |

#### `/docs/Development/References/`
The theme holding the schema, the rules, and the specifications.

| File | Description | Spec |
|------|-------------|------|
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `Specs/<Spec>/` | One directory per specification, holding its `Specs.md` and `sample` | |

---

## `/sandbox/`
The closed sandbox — the library, its public entry point, its contracts, and its internal logic. It imports nothing outside itself.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor taking a `Deps` and returning an `api.Lib` | |

### `/sandbox/contracts/`
The public contract structs the project is wired through — the only part of the sandbox the outside imports.

| File | Description | Spec |
|------|-------------|------|
| `deps/deps.go` | The `Deps` struct, one function field per injectable behavior | Deps |
| `api/api.go` | The structs the library hands back to callers | Outputs |

### `/sandbox/internal/`
The factories filling the `api` structs' function fields, unreachable from outside `sandbox/`. Declares no types.

| File | Description | Spec |
|------|-------------|------|
| `lib/lib.go` | `Lib`: reaches deps via `l.Deps.<Field>(...)` and creates the lib's objects | LibFunctions |
| `<object>/<object>.go` | The factories filling an `api` struct's fields, plus its `New` constructor propagating `Deps` | LibObjects |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontracts) contract, and the only place OS-bound and third-party code is allowed.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory` per `Deps` field, and the `New(...) deps.Deps` factory running them all | Adapters |

---

## `/cmd/`
Outside the sandbox. The executables the project ships.

### `/cmd/main/`

| File | Description | Spec |
|------|-------------|------|
| `main.go` | Self-contained `package main` wiring an adapter into the lib, calling `api.Lib.Sandboxmain`, and exiting with its return | CliMain |

---

## `/examples/cliExamples/`
Outside the sandbox. Shell scripts driving the built binary the way a user would.

| File | Description | Spec |
|------|-------------|------|
| `example<N>.sh` | Self-contained shell script demonstrating one goal against the built CLI | CliExamples |

---

## `/examples/libraryExamples/`
Outside the sandbox. Runnable Go examples demonstrating how to use the library from code.

| File | Description | Spec |
|------|-------------|------|
| `<example>/<example>.go` | Self-contained `package main` wiring an adapter into the lib | LibraryExamples |
