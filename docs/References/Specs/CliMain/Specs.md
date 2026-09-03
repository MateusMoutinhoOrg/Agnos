# CliMain Specification

## Description
Defines the required shape of `cmd/main/main.go` — the installed binary — as rendered from `assets/cli/cmd/main/main.go` by every `agnos build` of a project with a CLI layer. It is the one place an adapter and the sandbox are named in the same file, and it holds nothing else.

### Rules
- `package main`, importing `os`, the standard assembly as `agnosadapter "{{.Module}}/adapters/availables/standard"`, and the sandbox as `agnoslib "{{.Module}}/sandbox"`. Nothing else.
- `main` does four things and no more: build the deps (`agnosadapter.New()`), inject them (`agnoslib.New(&deps)`), run the interface over `os.Args[1:]`, and `os.Exit` with its return.
- No command, no flag, no output, no logic of its own: every word the binary prints comes from inside the sandbox, so the same interface runs against a fixed vector and a buffer in a test.
- The file is generated; a second executable with a different assembly is a hand-written `cmd/<name>/main.go` of the same four steps, which `build` leaves alone.

## Structure
1. **Package clause**: `package main`.
2. **Imports**: `os`, `agnosadapter`, `agnoslib`.
3. **`main`**: deps, lib, `CliMain(os.Args[1:])`, `os.Exit`.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/CliMain/sample.go).
