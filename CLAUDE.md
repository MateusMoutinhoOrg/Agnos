# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Agnos is a **CLI template**: a command-line financial tracker — categories holding spend and received transactions, persisted through an injected schema database — demonstrating OS-independent Dependency Injection. The CLI is the main face; consuming the same behavior as a Go library still works and is fully documented, but is the background feature. The value is as much in the *structure and documentation conventions* as in the domain code.

The whole command-line interface lives **inside the sandbox**, as the `Sandboxmain` field of `api.Lib`: it reads the command line through the injected Verb parser (`Deps.VerbLib`) and prints through the injected `Deps.Printf`. The installed binary, `cmd/main/main.go`, holds no command, no flag, and no output of its own — it wires an adapter into the library, calls that one field, and exits with its return.

## Commands

```bash
go build ./...                                   # build everything
go vet ./...                                     # the only linter available; currently clean
go test ./...                                    # run tests (none exist yet)

AGNOS_DATA=./scratch go run ./cmd/main <command> [args]           # run the CLI from source
go run ./examples/libraryExamples/TrackSpendSample/TrackSpendSample.go # run a library example
bash ./examples/cliExamples/ManageCategories.sh                   # run a CLI example
go run ./bootstrap/libraryExamples/Test/test.go                   # run the embedding lib's example
```

There is no lint config, CI, or test suite in the repo yet. When adding adapters or lib functions, `go build ./...` is the primary verification step, and `go vet ./...` must stay silent.

**Where each entry point writes its records** — none of this is decided inside the sandbox; `cmd/main/main.go:43` and every example pick the path themselves:

- `go run ./cmd/main` writes to `~/.agnos` unless `AGNOS_DATA` is set. Always export `AGNOS_DATA` when driving the CLI from source, or you edit the user's real budget.
- Library examples write `trackerdata/` in the working directory, `bootstrap/` writes `bootstrap-trackerdata/`; both are gitignored.
- CLI examples build the binary into a `mktemp -d` dir and point `AGNOS_DATA` at it, so they touch nothing else — required by the `CliExamples` spec.

## Architecture

Three top-level trees, wired through **structs of function fields** (never interfaces), with a strict one-way dependency flow:

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)

standard.New()  ──▶  deps.Deps  ──▶  lib.New(deps)  ──▶  api.Lib  ──▶  api.Category  ──▶  api.Transaction
(opinionated impl)   (contract)      (entry point)       (output structs, filled by sandbox/internal/ factories)
```

Contracts are structs whose fields hold functions, and **every** one of them is filled by **factories** — `func <Field>Factory(carrier *T) <FieldType>` bodies that return one closure reading the carrier at call time, with the assignment made explicitly by the caller. Inside the sandbox the carrier is the `api` struct, which carries its own `Deps` field, and `New` assigns the result (`l.GetCategory = GetCategoryFactory(&l)`, reading `l.Deps` inside the closure); inside `adapters/` the carrier is the adapter struct, which declares a `Deps deps.Deps` field its `New` assigns into from each factory's return value (`s.Deps.Now = NowFactory(s)`). No methods bound into fields, no internal mirror type, no `Api()` projection. This is a binding rule — see `docs/Development/References/RULES.md#factory-pattern`, the `Factories` spec (`docs/Development/References/Specs/Factories/`), and `docs/Development/References/StructContracts.md`.

Two trade-offs, neither caught by the compiler: **completeness is unchecked** — a field no factory fills is nil and panics on first call, so every factory must be called from its package's `New` constructor; and **`Deps` is read-only after construction** — the closures captured the struct the factories ran over, so patch `deps.Deps` before calling `lib.New`, never on the returned struct.

`sandbox/` is a **closed sandbox**: nothing in it may import `adapters/`, `cmd/`, `examples/libraryExamples/`, a third-party module, or an OS-bound stdlib package (`os`, `net`, `syscall`, …). Every such effect is a `Deps` field reached through `l.Deps`. This is a binding rule — see `docs/Development/References/RULES.md` and `docs/Development/References/SandboxIsolation.md`.

- **`sandbox/new.go`** — package `lib`, the only wiring point consumers touch: `New(deps.Deps) api.Lib`. Never imports `adapters/`. Importers alias it: `agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"`.
- **`sandbox/contracts/deps/deps.go`** — the `Deps` **struct**. Adding a requirement = adding a function field here. This is the contract every adapter must fill. Two fields are not functions: `VerbLib` (the embedded Verb argv-parser library) and `KeepLib` (the embedded Keep schema-database library), each injected whole.
- **`sandbox/contracts/deps/verbdeps/verbdeps.go`** — the sandbox's *copy* of the embedded Verb library's `api.Lib`, declared here because the sandbox may not import Verb. Each adapter's `VerbLibFactory` initializes the real library and assigns its fields onto this copy — the same mechanic `bootstrap/` uses for `agnosdeps`, with the roles reversed.
- **`sandbox/contracts/deps/keepdeps/keepdeps.go`** — the same thing for the embedded Keep schema-database library: a copy of its whole api (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`, `Props`, `Schema`, `Item`, `Error` and the constants), filled by each adapter's `KeepLibFactory`. Because Keep's fields hand back further api structs, the factory wraps them in closures that convert each returned struct too.
- **`sandbox/contracts/api/api.go`** — the output structs the lib hands back (`Lib`, `Category`, `Transaction`) plus the `Spend`/`Received` kind constants and the `ExitOk`/`ExitUsage`/`ExitFailure` codes `Sandboxmain` returns, each struct leading with a `Deps deps.Deps` field. Every field must be exported, or `sandbox/internal/` cannot fill it. **Every type in the project is declared here**, never in `internal/`. Types only — no function bodies.
- **`sandbox/internal/lib/`** — the lib's factories: one `<Field>Factory(l *api.Lib)` per function field of `api.Lib`, each returning the field's closure, plus `New(d deps.Deps) api.Lib` assigning every factory's return value into the matching field (`sandbox/new.go` just delegates to it). Go's `internal/` rule keeps it unreachable from `adapters/`, `examples/libraryExamples/`, and consumers.
- **`sandbox/internal/<object>/`** — one package per object the lib creates (`category/`, `transaction/`), holding that object's `<Field>Factory` functions plus a `New(d deps.Deps, …) api.<Object>` constructor that runs them all. There is no separate `Factory` aggregate — `New` is the aggregate. **Factories only, no type declarations.** Packages here take no `internal_` prefix — the `internal/` parent already says it.
- **`sandbox/internal/store/`** — shared helpers over the injected Keep database: the tracker's `Props` (a `category` collection owning a nested `transactions` collection), the field-name constants, and the record readers. Neither an object nor the entry point, so no spec governs it; it declares no types. Keep offers unique string keys and integers only, so a transaction's non-unique description travels inside its unique `reference`, composed as `<sequence>|<description>`.
- **`adapters/<name>/`** — outside the sandbox; the only place OS-bound and third-party code is allowed. Each declares a struct carrying a `Deps deps.Deps` field, one `<Field>Factory(a *<Name>Adapter)` per `Deps` field returning that field's value, and a `New(...) deps.Deps` constructor that assigns each factory's return value into `a.Deps` and returns it — the populated **contract struct**, never the adapter type. `standard` is the default adapter (Go stdlib only).
- **`sandbox/internal/cli/`** — the command-line interface itself: the `Run(l *api.Lib, args []string) int` dispatch `SandboxmainFactory` delegates to, the usage screen, and the amount parser. Like `store/` it is neither an object nor the entry point, so no spec governs it and it declares no types and no factories. Adding a command means editing here **and** `docs/CliUsage/References/Commands.md` — see `docs/Development/Tutorials/HandleCliCommands.md`.
- **`cmd/main/main.go`** — outside the sandbox; the installed binary. Wires the standard adapter into the lib, calls `l.Sandboxmain(os.Args[1:])`, and `os.Exit`s with its return. The argument vector must be the same one the adapter wired `Deps.VerbLib` over. Governed by the `CliMain` spec.
- **`examples/cliExamples/<Name>.sh`** — outside the sandbox; shell scripts that build the CLI into a scratch dir, point it at `AGNOS_DATA` of their own, and drive it as a user would. Governed by the `CliExamples` spec.
- **`examples/libraryExamples/<name>/<name>.go`** — outside the sandbox; self-contained `package main` programs wiring an adapter into the lib.
- **`bootstrap/`** — a second Agnos library embedding the root one, demonstrating the pattern when a lib's dependency is another lib built the same way: its sandbox declares a *copy* of the embedded api structs (`sandbox/contracts/deps/agnosdeps/`) and its adapter's `TrackerLibFactory` fills them by field assignment — the case where a factory assigns a value rather than a closure, because the field is a struct.

Every object propagates `Deps` to the objects it creates: a lib factory's closure calls `<object>.New(l.Deps, …)`, which stores the deps on the new api struct before running that object's factories — see `GetCategoryFactory` in `sandbox/internal/lib/lib.go`. Objects hold no database handle: `Category` and `Transaction` carry only their identifying data and re-find their record through `store.FindCategory` on every call, so a value never goes stale.

## Critical: this repo is documentation-driven

Changes are governed by required-reading docs, and several actions **must** update companion files in the same commit. `docs/` is split into four themes — `CliUsage/`, `LibUsage/`, `Development/`, `Templating/` — each with an `Index.md` listing its `Tutorials/` (one workflow per page) and its `References/` (lookups and explanations). The `README.md` links to the four indexes and nothing else. Read the tutorial matching your change:

| If you... | Read | And keep in sync |
|-----------|------|------------------|
| write or edit any `<Field>Factory` (sandbox **or** adapter) | `docs/Development/References/Specs/Factories/Specs.md` | the `New` constructor that must call it |
| add/rename/delete any file or dir | `docs/Development/References/Structure.md` | `docs/Development/References/Structure.md` |
| add/rename/delete a `.md` file | `docs/Development/Tutorials/HandleDocuments.md` | the theme's `Index.md` + `docs/Development/References/Structure.md` |
| add a lib function/object | `docs/Development/Tutorials/HandleLibElements.md` | `docs/LibUsage/References/PublicApi.md` (+ detail page `docs/LibUsage/References/PublicApi/<pkg>.<Symbol>.md`) |
| add a `Deps` field | `docs/Development/Tutorials/HandleDependencies.md` | **every** adapter in `adapters/` (and `bootstrap/adapters/`) |
| add an adapter | `docs/Development/Tutorials/HandleAdapters.md` | `docs/Development/References/Structure.md`, `docs/LibUsage/References/Adapters.md` |
| need an OS/third-party call inside `sandbox/` | `docs/Development/References/SandboxIsolation.md`, `docs/Development/Tutorials/HandleDependencies.md` | `sandbox/contracts/deps/deps.go` + **every** adapter |
| add/rename/delete a library sample | `docs/Development/Tutorials/HandleSamples.md` | `docs/LibUsage/References/ApiSamplesList.md` |
| add/rename/delete a CLI sample | `docs/Development/Tutorials/HandleCliExamples.md` | `docs/CliUsage/References/SamplesList.md` |
| add or change a CLI command or flag | `docs/Development/Tutorials/HandleCliCommands.md` | the `Usage` screen in `sandbox/internal/cli/cli.go` + `docs/CliUsage/References/Commands.md` |
| fork or adapt the template into a real library | `docs/Templating/Tutorials/ForkTemplate.md`, `AdaptExistingLib.md` | `docs/Templating/References/TemplateFileActions.md` (the per-file copy/create/rewrite/delete list both tutorials follow) |

`docs/Development/References/RULES.md` is the binding rule set and `docs/Development/References/Specs.md` is the index of every file specification; `AGENTS.md` points here. Adding a `Deps` field without filling it in all adapters breaks every consumer at **runtime**, not at build time — that's the most common footgun, and `go build` will not catch it. A new public lib function or object must be declared in `sandbox/contracts/api/api.go`, given a factory in `sandbox/internal/`, **and** called from that package's `New` constructor, or callers get a nil field.

## Conventions

- Code that consumes the library from outside it (`cmd/`, `examples/libraryExamples/`, the `bootstrap/` adapter, third-party callers) aliases every import with the `agnos` prefix: `agnosadapter` (`adapters/<name>`), `agnoslib` (`sandbox`), `agnostypes` (`sandbox/contracts/api`), `agnosdeps` (`sandbox/contracts/deps`). Files belonging to the library itself — `sandbox/` and `adapters/` — keep the plain package names. See the Import Aliases rule in `docs/Development/References/RULES.md`.
- **Money is always an `int64` in the smallest currency unit (cents), never a float.** `84.50` is `8450` across the whole library. Only two places convert: `cli.ParseAmount` turns the typed `84.50` into cents, and `store.Money` renders cents back for display. `Transaction.Amount` is always positive — direction lives in `Kind` (`Spend`/`Received`), and `Transaction.SignedAmount()` applies the sign, so balances are a plain sum.
- The CLI version reported by `agnos version` and `--version` is the `Version` const at the top of `sandbox/internal/cli/cli.go`; it is independent of the `@v0.0.3` install tag pinned in `docs/CliUsage/Tutorials/QuickStart.md`, `docs/CliUsage/Tutorials/InstallCli.md`, `docs/LibUsage/Tutorials/QuickStart.md`, `docs/LibUsage/Tutorials/LibInitialization.md`, and `docs/Development/References/Structure.md`. Bumping a release means touching both.
- Module path is `github.com/MateusMoutinhoOrg/Agnos-Cli`; renaming it is a documented procedure — see `docs/Templating/Tutorials/RenameModule.md`.
- Public-facing lib API entries each get a detail page in `docs/LibUsage/References/PublicApi/` named `<pkg>.<Symbol>.md`.
- `docs/Development/References/Specs/` holds the specifications: one directory per kind of file, each pairing a `Specs.md` (how the file must be shaped) with a `sample`. Never browse it — always locate a spec through `docs/Development/References/Specs.md`.
