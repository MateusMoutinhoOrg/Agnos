# LibraryExamples Specification

## Description
Defines the required shape of a runnable library example in `examples/libraryExamples/<example>/<example>.go`. A library example is a self-contained `package main` program that wires an adapter into the lib to demonstrate real usage from Go code. The command-line counterpart — a shell script driving the installed binary — is governed by [CliExamples](/docs/Meta/CliExamples/Specs.md) instead.

### Rules
- Each example lives in its own directory under `examples/libraryExamples/` named after the feature it demonstrates (e.g. `examples/libraryExamples/ExampleSample/`).
- The file is named after its directory (`<example>/<example>.go`) and declares `package main` with a `main` function.
- An example wires the two layers together: it builds a `deps.Deps` through an adapter's `New(...)` factory, then passes it to `lib.New(...)`, which returns an `api.Lib`.
- An example may import `adapters/<name>` (aliased `agnosadapter`), `sandbox` (aliased `agnoslib`), and `sandbox/contracts/api` (aliased `agnostypes`); it must never import `sandbox/internal/` — Go's `internal/` rule rejects it — nor reconstruct dependencies by hand, which is the adapter's job.
- Examples live outside the sandbox and are the only place an adapter and the library are named in the same file.
- Keep examples minimal and runnable via `go run ./examples/libraryExamples/<example>/<example>.go`; add explanatory comments on the key wiring steps.
- Adding, renaming, or deleting an example requires updating the Library Examples section of [README.md](/README.md) — see [HandleSamples.md](/docs/HandleSamples.md).

## Structure
1. **Package clause**: `package main`.
2. **Imports**: every import of this module is aliased with the `agnos` prefix, so a reader sees at a glance which layer a call belongs to — the adapter as `agnosadapter` (e.g. `agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"`), the sandbox entry point as `agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"`, and, when the example names an output type, `agnostypes "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"`.
3. **`main` function**: build deps via `agnosadapter.New(...)`, inject them with `agnoslib.New`, then exercise the returned `agnostypes.Lib`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
