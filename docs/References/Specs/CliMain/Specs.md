# CliMain Specification

## Description
Defines the required shape of the executable entry point in `cmd/main/main.go` — the binary a user installs. It is the thinnest file in the project: it wires an adapter into the library, hands the command line to `api.Lib.Sandboxmain`, and exits with the code that returns. The interface itself is **not** written here; it lives inside the sandbox, so it can be run over injected dependencies and never depends on the process it happens to be hosted in.

### Rules
- `cmd/main/main.go` declares `package main` with a `func main()`, and is the only executable of the library outside `examples/libraryExamples/`.
- `main` does exactly three things: build a `deps.Deps` through an adapter's `New(...)` constructor, inject it with `lib.New(...)`, and call `os.Exit(l.Sandboxmain(os.Args[1:]))`.
- It must never branch on a command, parse a flag, format a result, or print anything. Every one of those belongs to `api.Lib.Sandboxmain` inside the sandbox — a flag handled here would be a flag no other front end can offer. Deciding **where** state lives is the exception: that is an OS-bound choice, so resolving the data directory (and any environment override of it) belongs here.
- The argument vector passed to `Sandboxmain` must be the same one the adapter wired `deps.Deps.VerbLib` over — `os.Args[1:]` on both sides — or the interface and the parser disagree about the command line.
- It imports the library under the `agnos`-prefixed aliases every outside consumer uses; see the Import Aliases rule in [RULES.md](/docs/References/RULES.md).
- The exit code is whatever `Sandboxmain` returns, unmapped: the constants in `sandbox/contracts/api` are the process's contract with its caller.
- Renaming, moving, or adding an entry point requires updating [Structure.md](/docs/References/Structure.md) and the install command in [InstallCli.md](/docs/Tutorials/InstallCli.md).

## Structure
1. **Package clause**: `package main`.
2. **Imports**: `os` and whatever resolving the data directory needs, plus `agnosadapter "…/adapters/standard"` and `agnoslib "…/sandbox"`.
3. **Configuration constants**: the data directory's name and the environment variable overriding it.
4. **`main` function**: build deps through the adapter, inject with `agnoslib.New`, and `os.Exit(l.Sandboxmain(os.Args[1:]))`.
5. **Unexported helpers**: only for the OS-bound choices `main` makes, such as resolving the data path.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
