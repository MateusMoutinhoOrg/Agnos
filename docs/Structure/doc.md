# Project Structure

This document maps the project **schema** — the kinds of files the project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs](/docs/Specs/doc.md) to get its description and sample. A slot marked *generated* is written by `agnos build` and never edited by hand — see [GeneratedFiles](/docs/GeneratedFiles/doc.md).

Agnos is a **factory**: `agnos start` scaffolds a project of this same schema, and `agnos build` regenerates its generated slots. The repository is therefore also the reference instance of what it produces — a generated project has every slot below except `assets/deplist/`, `assets/templates/`, `sandbox/internal/actions/`, `parsables/`, `smartio/` and `utils/`, which are Agnos's own logic. The dependency flow between the code trees is one-way:

```
adapters/  ──▶  sandbox/  ◀──  cmd/
(reaches the OS)  (closed)     (wires the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the library and, with it, the command-line interface. Nothing inside it may import `adapters/`, `cmd/`, `assets/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation](/docs/SandboxIsolation/doc.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each lib imports `sandbox/deps` and nothing else from the sandbox.
- **`/cmd/`** sits outside the sandbox too, and is the only place where an adapter and the sandbox meet.
- **`/assets/`** sits outside the sandbox as well: the Go templates compiled into the binary and reached only through `Deps.Embeddeps`.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and the Doc Index pointing at each theme index under `docs/Index/` | Readme |
| `CLAUDE.md` | The harness specification for the coding agent: every pattern in this document, stated as a rule | |
| `LICENSE` | License terms for the project | |
| `go.mod`, `go.sum` | Go module definition and dependencies; `agnos start` writes the first, `go mod tidy` the second | |
| `.gitignore` | Intentionally untracked files: `release/`, the local build scripts, scratch dirs | |
| `Relatorio.md` | The bug and developer-experience report against Agnos 0.0.2 that motivated the exit codes, output channels and runtime gate | |

---

## `/AgnosConfig/`
The project's configuration, written once by `agnos start` from the `start` asset group and read back by every `agnos build`. Each file is parsed by one `parsables/<name>conf` package.

| File | Description | Spec |
|------|-------------|------|
| `project.yaml` | `name`, `version`, `description` — regenerated into `sandbox/internal/config/config.go` (`projectconf`) | |
| `themes.yaml` | The documentation themes, one `name` / `id` / `description` each — each `id` names a generated `docs/Index/<id>.md` (`themesconf`) | |
| `ignore.yaml` | Paths SmartIO hides from listings (`ignorableconf`) | |
| `paths.yaml` | Path rewrites SmartIO applies to listings (`pathreplacerconf`) | |

---

## `/sandbox/`
The closed sandbox. `verify` allows exactly the four directories below plus `new.go`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | *generated* — `New(deps *deps.Deps) *api.Sandbox`, calling one `binds.<X>Bind` per file of `binds/` | |

### `/sandbox/api/`
The contract structs the sandbox hands back — the only part of the sandbox a consumer imports. Imports **nothing** but other `sandbox/api` packages: no stdlib, no `sandbox/deps`.

| File | Description | Spec |
|------|-------------|------|
| `sandbox.go` | *generated* — `type Sandbox struct` with one field per other file of this directory | |
| `actions.go` | `type Actions struct` of function fields, one per action, plus the props structs (`BuildProps`, `StartProps`, `FieldProps`, `CommandProps`) and the `RuntimeGo` / `RuntimeNone` constants | Contract |
| `cli.go` | *generated* — `type Cli struct { CliMain func([]string) int }` and the `ExitOk` / `ExitFailure` / `ExitUsage` constants | Contract |

### `/sandbox/binds/`
The seam between contract and implementation. Every file mirrors a file of the same name in `api/` and declares only functions.

| File | Description | Spec |
|------|-------------|------|
| `actions.go` | `ActionsBind(deps, sandbox)`: one assignment per `Actions` field, each calling an action package | Binder |
| `cli.go` | *generated* — `CliBind(deps, sandbox)`, assigning the generated `cli.CliMain` | Binder |

### `/sandbox/deps/`
The dependency contract. Each sub-directory is the sandbox's **copy** of an outside library's api; the field name on `Deps` is the title-cased directory name. Imports only the standard library and other `sandbox/deps` packages.

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | *generated* — `type Deps struct` with one `<Title> <dir>.Lib` field per sub-directory | |
| `<x>/<x>.go` | One `Lib` struct of function fields per sub-contract — `argvdeps`, `embeddeps`, `goimportsdeps`, `iodeps`, `requestdeps`, `rundeps`, `serializables`, `std` | DepsContract |

### `/sandbox/internal/`
Every implementation, unreachable from outside `sandbox/` by Go's `internal/` rule.

| File | Description | Spec |
|------|-------------|------|
| `config/config.go` | *generated* — `ProjectName` and `Version` from `AgnosConfig/project.yaml` | |

#### `/sandbox/internal/commands/`
The CLI surface: one package per command, three files each. `help` and `version` are written by `agnos build` itself.

| File | Description | Spec |
|------|-------------|------|
| `<name>/entries.yaml` | The command's declaration: identifiers, category, help, flags, args | CommandEntries |
| `<name>/entries.go` | *generated* — the typed `Entries` struct the dispatch fills | |
| `<name>/handler.go` | `CommandHandler(deps *deps.Deps, entries *Entries) int` — for Agnos's commands, a call into the matching action | CommandHandler |

#### `/sandbox/internal/cli/`
The generated dispatch layer.

| File | Description | Spec |
|------|-------------|------|
| `climain.go` | *generated* — `CliMain` plus one `dispatch<Name>` per command; see [CommandDispatch](/docs/CommandDispatch/doc.md) | |

#### `/sandbox/internal/actions/`
The logic layer: one package per reusable operation, exposed on `api.Actions` and called by the matching command's handler.

| File | Description | Spec |
|------|-------------|------|
| `<name>/<name>.go` | The public entry: opens a SmartIO, calls the internal function, persists, runs any follow-up build | Action |
| `<name>/<name>_internal.go` | The logic against an open `*smartio.SmartIO`, so other actions can compose it | Action |
| `build/collect_*.go` | The collectors: one directory listed, names title-cased, a slice returned | Collector |
| `build/generate_*.go`, `build/run_runtime.go` | The per-command `entries.go`, the `help` declaration and the documentation indexes, and the runtime hand-off | |
| `verify/check_*.go` | One schema rule set per file, each returning `[]string` violations | |

#### `/sandbox/internal/parsables/`
One small parser per configuration file, five files each.

| File | Description | Spec |
|------|-------------|------|
| `<name>conf/api.go` | The struct: data fields, then method fields, `Render` last | Parsable |
| `<name>conf/new.go`, `new_empty.go` | Parse from a string; build the defaults | Parsable |
| `<name>conf/bind_methods.go`, `render.go` | Assign every method field; render back to text | Parsable |

`commandconf` (entries.yaml), `depsversionconf` (depsversion.yaml), `docpropsconf` (docs/**/props.yaml), `ignorableconf`, `moduleconf` (go.mod, parsed by hand), `pathreplacerconf`, `projectconf`, `themesconf`.

#### `/sandbox/internal/smartio/`
The transactional filesystem rooted at `--path` — see [SmartIO](/docs/SmartIO/doc.md).

| File | Description | Spec |
|------|-------------|------|
| `api.go` | The `SmartIO` struct: `Root`, the pending sets, one function field per operation | |
| `new.go`, `bind_methods.go`, `persist.go`, `<operation>.go` | The constructor, the binder, the flush, and one file per operation | |

#### `/sandbox/internal/utils/`
Helpers shared by the actions: template rendering (`RenderGroup`, `RenderTemplateToDest`), `go.mod` and `themes.yaml` loading, the documentation tree walk (`CollectDocTree`, `LoadDocProps`, `SortDocs`), and the `entries.yaml` editing helpers (name normalization, `NewField`, `FindField` / `InsertField` / `RemoveField` / `AppendUnique`).

---

## `/adapters/`
Outside the sandbox. `verify` allows exactly the two directories below.

### `/adapters/libs/`

| File | Description | Spec |
|------|-------------|------|
| `<lib>/<File>.go` | One package per sub-contract, exporting `Bind(deps *deps.Deps)` and filling that one field — `verb` (fills `Argvdeps`), `embeddeps`, `goimportsdeps`, `iodeps`, `requestdeps`, `rundeps`, `serializables`, `std` | AdapterLib |

### `/adapters/availables/`

| File | Description | Spec |
|------|-------------|------|
| `standard/new.go` | *generated* — `New() deps.Deps` calling every lib's `Bind` | |
| `<name>/new.go` | A hand-written assembly with a different mix; left alone by `build` | Available |

---

## `/assets/`
Outside the sandbox. Every file is a Go `text/template` compiled into the binary by `asset.go` and reached through `Deps.Embeddeps`. `assets/<group>/<path>` renders to `<path>` inside the target project. Never part of `go build ./...`.

| File | Description | Spec |
|------|-------------|------|
| `asset.go` | Package `assets`: the `//go:embed all:*` directive and the `Files` filesystem | |
| `depsversion.yaml` | `<dep>: <module>@<version>` for every dep that pulls an external module | |
| `start/AgnosConfig/*.yaml` | The configuration `agnos start` writes once | AssetTemplate |
| `all/**` | Rendered by every build: `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/internal/config/config.go` | AssetTemplate |
| `deps/**` | Rendered when `sandbox/deps/` exists: `sandbox/deps/deps.go`, `adapters/availables/standard/new.go` | AssetTemplate |
| `cli/**` | Rendered when `sandbox/internal/cli/` exists: `cmd/main/main.go`, `sandbox/api/cli.go`, `sandbox/binds/cli.go`, `sandbox/internal/cli/climain.go`, the `help` handler, the `version` command | AssetTemplate |
| `deplist/<dep>/**` | One installable dep, mirroring the target layout: its `sandbox/deps/<dep>/` contract and `adapters/libs/<lib>/` implementation | Dep |
| `templates/*` | Single-file scaffolds rendered outside any group: `command_entries.yaml`, `command_handler.go` (by `add-command`), `entries.go` (per command, by `build`), `help_entries.yaml` (once, by `build`), `theme_index.md` and `doc_index.md` (per index, by `build`) | AssetTemplate |

---

## `/cmd/`
Outside the sandbox. The executables the project ships.

### `/cmd/main/`

| File | Description | Spec |
|------|-------------|------|
| `main.go` | *generated* — self-contained `package main` wiring `standard.New()` into `sandbox.New`, calling `lib.Cli.CliMain(os.Args[1:])`, and exiting with its return | CliMain |

**Run the CLI from source:**
```sh
go run ./cmd/main <command> [flags] [args]
```

---

## `/release/`
Git-ignored. The cross-compiled binaries `Build.md` produces (`mac86.bin`, `linux86.out`, `windows86.exe`, …) and the transient `bootstrap.bin` of [BootstrapAgnos](/docs/BootstrapAgnos/doc.md).

---

## `/old/`
The previous incarnation of this repository — a Go CLI **template** with a financial-tracker example — kept as reference material with its full documentation under `old/docs/`. Not part of the build; nothing links into it.

---

## `/docs/`
Documentation of the project. Every doc is a **directory**: `docs/<DocName>/` holding the `doc.md` a reader opens and the `props.yaml` that declares it, plus — recursively — the sub-doc directories nested under it. A **theme** — what the reader wants to accomplish — is not a directory: it is a `themes.yaml` entry a doc names in its `props.yaml`, and the index that groups them is generated. `Index` is reserved at the first level. The [README](/README.md) links to the theme indexes and to nothing else inside `docs/`.

| Directory | Description |
|-----------|-------------|
| `Index/` | *generated* — one theme index per entry of `themes.yaml` |
| `<DocName>/` | One doc: `doc.md`, `props.yaml`, its assets, and its sub-docs |

### `/docs/Index/`
Written by every build and owned by it: the directory is rewritten whole, so an index of a deleted theme never survives.

| File | Description | Spec |
|------|-------------|------|
| `<theme-id>.md` | *generated* — the theme's entry point: one row per first-level doc naming that theme | Index |

### `/docs/<DocName>/`
One doc, at any depth. A doc is a workflow, a lookup, or an explanation — the kind picks its specification, not its location.

| File | Description | Spec |
|------|-------------|------|
| `doc.md` | The documentation itself, titled with the doc's name | TutorialDocs / ReferenceDocs / ExplanationDocs |
| `props.yaml` | The declaration: `name`, `description`, `themes` (first level only), optional `order` | DocProps |
| `Index.md` | *generated*, only when the doc has sub-docs — one row per direct sub-doc | Index |
| `<asset>` | Any other file (a `sample.go`, an image): carried by the doc, ignored by the indexes | |
| `<SubDocName>/` | A sub-doc, of this same shape, listed only by this doc's `Index.md` | |

The docs of this repository: `Commands`, `EntriesYaml`, `DepList`, `GeneratedFiles` and `PublicApi` are lookups; `Adapters` is a lookup with its own AdaptersDoc spec; `Structure` (this page), `BuildPipeline`, `SmartIO`, `SandboxIsolation`, `CommandDispatch` and `StructContracts` are explanations; every other first-level doc is a workflow. Two carry sub-docs:

| Doc | Sub-docs |
|-----|----------|
| `PublicApi/` | `<pkg>.<Symbol>/` — one doc per public struct, function, or field, named after the package the symbol is declared in |
| `Specs/` | `<Spec>/` — one doc per specification, its `sample.<ext>` sitting next to `doc.md` as an asset |
