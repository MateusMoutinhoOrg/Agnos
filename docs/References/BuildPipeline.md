# The Build Pipeline

## Description
Explains what happens between `agnos build` (or any command that ends by running it) and a compiling project: the schema gate, the collectors, the asset groups rendered in order, the single persist, and the runtime handed the result. The command's flags are in [Commands.md](/docs/References/Commands.md#core-commands); every file it writes is in [GeneratedFiles.md](/docs/References/GeneratedFiles.md); the transactional filesystem it writes through is [SmartIO.md](/docs/References/SmartIO.md).

---

## The Verify Gate

The `build` **command** runs `verify` first, unless `--unsafe` is passed. `verify` reads the tree through SmartIO, writes nothing, and returns one error listing every violation of the harness schema:

| Rule | Where it is checked |
|------|---------------------|
| `sandbox/` holds only the `api`, `binds`, `deps`, `internal` directories plus a loose `new.go`. | `check_sandbox.go` |
| No file under `sandbox/` imports a module-internal package outside `sandbox/`. | `check_sandbox.go` |
| `sandbox/api/*` imports only other `sandbox/api` packages — no stdlib, no external module, not even `sandbox/deps`: api is pure contract. | `check_sandbox.go` |
| `sandbox/deps/*` imports only the standard library and other `sandbox/deps` packages. | `check_sandbox.go` |
| Every file in `sandbox/binds/` mirrors a file of the same name in `sandbox/api/` and declares only functions — no top-level types, consts or vars. | `check_sandbox.go` |
| `adapters/` holds only the `availables` and `libs` directories. | `check_adapters.go` |

The `build` **action** — and the follow-up build run by `start`, `deps-*`, `dep-*`, `cli-*` and the command editors — does not run the gate, so a mid-refactor state can still regenerate. Checks live in `check_*.go` files each returning `[]string`; a new rule is a sibling file.

---

## Collectors

`BuildInternal` first reads `AgnosConfig/project.yaml` (a missing or unparsable one is a hard error: `start` is a prerequisite) and `go.mod`, then decides two switches: `HasDeps` (`sandbox/deps/` exists) and `HasCli` (`sandbox/internal/cli/` exists). When `HasCli`, it renders `help`'s `entries.yaml` if missing, so the collector below sees it.

Then the collectors run. Each one lists one directory, takes the last path segment, title-cases it, and returns the slice — so a generated file's contents are predictable from a directory listing:

| Collector | Lists | Feeds | Into |
|-----------|-------|-------|------|
| `CollectBinds` | `sandbox/binds/*.go` | `{{range .Binds}}` — `CliBind`, `ActionsBind` | `sandbox/new.go` |
| `CollectConstructors` | `sandbox/api/*` | `{{range .Constructors}}` — `Actions`, `Cli` | `sandbox/api/sandbox.go` |
| `CollectDepsLibs` | `sandbox/deps/<x>/` | `{{range .DepsLibs}}` — `{Title, Name}` per contract | `sandbox/deps/deps.go` |
| `CollectAdapterLibs` | `adapters/libs/<x>/` | `{{range .AdapterLibs}}` — `{Name}` per lib | `adapters/availables/standard/new.go` |
| `CollectCommands` | `sandbox/internal/commands/<x>/entries.yaml` | `{{range .Commands}}` — a rich map per command | `climain.go`, `help/handler.go`, each `entries.go` |

`CollectCommands` is the one exception to the one-line shape: it parses each `entries.yaml` through `parsables/commandconf` and precomputes, per command, the Go name, a `MatchExpr` over its identifiers, and per field the Go field name, Go type, getter names, default literal and a `RangeCheck` snippet.

---

## Template Variables

Every group is rendered with one `vars` map:

| Key | Value |
|-----|-------|
| `Module` | The module path from `go.mod` — every generated import is `{{.Module}}/…`. |
| `Name`, `Description`, `Version` | From `AgnosConfig/project.yaml`. |
| `ProjectName` | `Name` title-cased; becomes `config.ProjectName`, and lowercased the binary name `help` prints. |
| `ConfigDir` | `config.ProjectName + "Config"` — the config directory this binary reads (`AgnosConfig`). Used by `README.md` to locate `docs/ReadmeHeader.md`. |
| `HasDeps`, `HasCli` | The two switches above; `sandbox/new.go` takes a `*deps.Deps` only when `HasDeps`. |
| `Binds`, `Constructors`, `DepsLibs`, `AdapterLibs`, `Commands` | The collectors' output. |

---

## Asset Groups, in Order

Templates live at `assets/<group>/<path>` and render to `<path>` inside the project through `utils.RenderGroup`, which writes each with `WriteFileOverwrite`. Groups render in this order, each only when its condition holds:

| Group | Condition | Renders |
|-------|-----------|---------|
| `all` | always | `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/internal/config/config.go`, `README.md` |
| `deps` | `HasDeps` | `sandbox/deps/deps.go`, `adapters/availables/standard/new.go` |
| *(per command)* | `HasCli` | `sandbox/internal/commands/<x>/entries.go` from `assets/templates/entries.go` |
| `cli` | `HasCli` | `cmd/main/main.go`, `sandbox/api/cli.go`, `sandbox/binds/cli.go`, `sandbox/internal/cli/climain.go`, `help/handler.go`, `version/{entries.yaml,handler.go}` |

Two groups are rendered by other commands: `start` (the `AgnosConfig/` files, once) and `deplist/<dep>` (by `dep-install`). `assets/templates/` holds single-file scaffolds rendered outside any group.

Every asset template is rendered by `utils.RenderGroup` / `utils.RenderTemplateToDest` with two native functions in scope: **`render "<project-relative path>"`** reads that file from the target project (through the transaction-aware `io`, so a file written earlier in the same build is visible), renders it as a Go `text/template` with the *same* `vars` and the same native functions (so nesting works to any depth), and returns the result; **`copy "<project-relative path>"`** reads a file through the same `io` and returns its contents verbatim, with no rendering, for embedding a non-template file. `README.md`'s `all` template is a `render` call against `<ConfigDir>/docs/ReadmeHeader.md` followed by `copy "LICENSE"`. A missing target is a hard build error (unparsable too, for `render`).

---

## Persist, then Runtime

Nothing above touches the disk. Every write is buffered in the SmartIO transaction, and `io.Persist()` flushes it once: pending removals, then pending directory creations, then file writes. Only then does `RunRuntime` hand the project to a toolchain, because the toolchain reads the disk, not the transaction:

| Runtime | Steps |
|---------|-------|
| `go` (default) | `go mod tidy` — the step that writes `go.sum`, so a fresh scaffold is left compilable — then `go build` over whichever of `./cmd/...`, `./sandbox/...`, `./adapters/...` exist. Never `./...`: `assets/` holds Go **templates**. |
| `none` | Nothing. Renders only. |

Any other name is a usage error. A non-zero toolchain exit is reported with the toolchain's output and exits `1`: `build` and `verify` cannot report success over a project that does not compile.

Every follow-up build names its runtime through `api.BuildProps`. Commands that add something (`start`, `deps-init`, `dep-install`, `cli-init`, `add-command`, `add-flag`, `add-arg`, `set-command`) pass `go`; commands that remove something (`deps-purge`, `dep-remove`, `cli-purge`, `remove-command`, `remove-flag`, `remove-arg`) pass `none`, because dropping a field or a command can legitimately leave hand-written code referring to what is gone.

---

## Deps

`dep-install <dep>` and `dep-remove <dep>` sit beside the pipeline rather than inside it. Install lists `deplist/<dep>` in the embedded assets (an empty listing is `unknown dep`), adds the pinned `require` from `assets/depsversion.yaml` to `go.mod` when the dep is listed there, renders the dep's subtree as one group with `{"Module": …}`, **persists**, and only then runs `build` — because `build`'s collectors list directories from disk and would not see the pending writes. Remove deletes those same files and any directory left empty, strips the `require`, persists, and runs `build` with `none`.

---

## Self-Hosting

Agnos is one of the projects this pipeline renders: its own `sandbox/deps/deps.go`, `adapters/availables/standard/new.go`, `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/binds/cli.go`, `climain.go`, every `entries.go` and the `help` command are generated in place, and the result compiles. Two consequences follow. `agnos build` must stay **idempotent** over this tree, and every sandbox and adapter file must reach `Deps` fields by their mechanical, title-cased directory names, because the struct is regenerated from the listing. Regenerating the checkout with itself is [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md).
