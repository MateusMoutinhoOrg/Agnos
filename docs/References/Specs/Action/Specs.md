# Action Specification

## Description
Defines the required shape of an **action** — a package under `sandbox/internal/actions/<name>/` holding one reusable operation in two files: the public entry that owns a SmartIO transaction, and the internal function that does the work against an open one, so another action can compose it. Writing one is [AddAction.md](/docs/Tutorials/AddAction.md); the transaction it writes through is [SmartIO.md](/docs/References/SmartIO.md).

### Rules
- `package <name>`, `<name>` in snake case (`dep_install`), matching the command's package name when a command exposes it.
- `<name>.go` holds exactly one exported function, `<Name>(deps *deps.Deps, …) error` (or `(T, error)` for an action with a result), that: opens `smartio.New(deps, path, config.ProjectName)`; calls `<Name>Internal`; calls `io.Persist()`; and runs any follow-up as a **fresh** step — `buildAction.Build(deps, api.BuildProps{Path, Runtime})` — after persisting, because `build`'s collectors list from disk. Nothing else.
- `<name>_internal.go` holds `<Name>Internal(deps *deps.Deps, io *smartio.SmartIO, path string, …) error` and its unexported helpers. It begins with `deps.Std.Log("<verb> started with path %s \n", path)`, reads and writes through `io` only, and never calls `io.Persist`.
- Every path handed to `io` is project-relative. `path` is passed along only for logging and for `RunRuntime`, never joined onto a file path.
- Failures come back as `error` built with `deps.Std.Errorf`; progress goes through `deps.Std.Log`; nothing goes through `Printf`.
- The follow-up build's runtime is `api.RuntimeGo` when the action adds something and `api.RuntimeNone` when it removes something.
- `build` and `verify` are the two exceptions to the two-file shape: `build/` holds its collectors, generators and `run_runtime.go` beside `build.go` / `build_internal.go`; `verify/` holds its `check_*.go` files and performs no write at all.
- Every action is exposed as a field of `api.Actions` and assigned in `sandbox/binds/actions.go` in the same commit.

## Structure
1. **`<name>.go`**: package clause, imports, the public entry — open, internal, persist, follow-up.
2. **`<name>_internal.go`**: package clause, imports, `<Name>Internal`, unexported helpers.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/Action/sample.go), which shows both files in one listing.
