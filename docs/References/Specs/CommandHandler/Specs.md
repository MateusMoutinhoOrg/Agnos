# CommandHandler Specification

## Description
Defines the required shape of `sandbox/internal/commands/<name>/handler.go` — the one hand-written Go file of a command, holding the function the generated dispatch calls once the command line has been parsed into `Entries`. Writing one is [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md); how it is reached is [CommandDispatch.md](/docs/References/CommandDispatch.md).

### Rules
- `package <name>`, where `<name>` is the command's package name — the identifier with hyphens replaced by underscores (`add-flag` → `add_flag`).
- Exactly one exported function, `CommandHandler(deps *deps.Deps, entries *Entries) int`. `Entries` is the generated struct in the sibling `entries.go`; nothing else in the package is exported.
- It imports `sandbox/api` (for the exit constants), `sandbox/deps`, and packages under `sandbox/internal/`; nothing outside `sandbox/`.
- Every effect goes through a `deps` field by its mechanical name. The three output channels are `deps.Std.Printf` (the result, stdout), `deps.Std.Log` (progress, stderr, silenced by `--quiet`) and `deps.Std.Error` (failures, stderr).
- It returns `api.ExitOk` or `api.ExitFailure`, never `api.ExitUsage`: a wrong command line never reaches it.
- It parses nothing: `entries` is typed, defaulted and range-checked. Logic worth reusing lives in an action or a `sandbox/internal/` package; the handler calls it and reports.
- In Agnos's own commands, the handler calls the matching action package (`<name>Action.<Name>(deps, …)`), writes the error through `Error` and returns `ExitFailure` on failure, prints any returned result through `Printf`, and returns `ExitOk`.

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: `api`, `deps`, and the action or internal packages called.
3. **`CommandHandler`**: the one function, with a doc comment naming the verb it backs and what each `Entries` field carries.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/CommandHandler/sample.go).
