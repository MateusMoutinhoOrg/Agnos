# BuildPipeline

`build` command = `verify` (unless `--unsafe`) -> `BuildInternal` -> `Persist` -> runtime. The `build` **action** (the follow-up every other command runs) skips verify so mid-refactor states still regenerate.

## BuildInternal

1. Read `AgnosConfig/project.yaml` (hard error if missing) and `go.mod`. Set `HasDeps` (`sandbox/deps/` exists) and `HasCli` (`sandbox/internal/cli/` exists).
2. Load `themes.yaml`; `CollectDocs`, merge in `CollectGeneratedDocs` (the docs the `all` group itself writes — listings read disk, so on a first build they are not there yet), then `GenerateSubdocIndexes` (one `Index.md` per doc with sub-docs; deletes `docs/Index/` left by older versions). Skipped when `docs/` is absent.
3. If `HasCli`: write `help/entries.yaml` if missing, then `CollectCommands`, then one `entries.go` per command.
4. Collectors, then render groups in order: `all` (always), `deps` (`HasDeps`), `cli` (`HasCli`).

| Collector | Lists | Var | Feeds |
|---|---|---|---|
| `CollectBinds` | `sandbox/binds/*.go` | `Binds` | `sandbox/new.go` |
| `CollectConstructors` | `sandbox/api/*` | `Constructors` | `sandbox/api/sandbox.go` |
| `CollectDepsLibs` | `sandbox/deps/<x>/` | `DepsLibs` (`Title`, `Name`) | `sandbox/deps/deps.go` |
| `CollectAdapterLibs` | `adapters/libs/<x>/` | `AdapterLibs` (`Name`) | `adapters/availables/standard/new.go` |
| `CollectCommands` | `commands/<x>/entries.yaml` | `Commands` (rich map: identifiers, category, help, `Flags`/`Args` with Go names, types, getters, defaults, `RangeCheck`) | `climain.go`, `help/handler.go`, `entries.go` |
| `CollectDocs` | `docs/**/props.yaml` | doc tree sorted by `order` then name | `**/Index.md`, `DocIndex` |
| `CollectGeneratedDocs` | `assets/all/docs/*/props.yaml`, rendered | merged into the doc tree | same |
| `CollectDocIndex` | the merged tree grouped by theme | `DocIndex` (per theme: `Name`, `Description`, `Docs`) | `README.md`. A theme no doc names renders no section |
| `CollectPublicApi` | `sandbox/api/*.go` parsed by `deps.Goimportsdeps` | `PublicApi` (per file: `Path`, `Doc`, `Types`, `Constants`, `Variables`, `Functions`; exported only, doc comments flattened to one table line) | `docs/PublicApi/doc.md` |
| `CollectDepsApi` | `sandbox/deps/<x>/*.go`, same parse | `DepsApi` (`Name`, `Title`, `Files`) | `docs/PublicApi/doc.md` |
| `CollectCommandDocs` | `commands/<x>/entries.yaml` (visible ones), grouped by category in first-seen order | `CommandDocs` (per category: `Commands` with `Identifier`, `Aliases`, `Help`, `LongDescription`, `Usage`, `Flags`/`Args` as table rows, `Examples`) | `docs/Commands/doc.md` |
| `CollectStructure` | `AgnosConfig/structure.yaml` (structureconf) | `Structure` (one `Line` per item, depth-indented and padded to a common description column) | `docs/Structure/doc.md` |

Every render whose destination ends in `.go` is passed through `deps.Goimportsdeps.Format` (`go/format`, i.e. `gofmt`) before it is written, so generated Go is byte-identical to what a formatting editor saves and a regenerated tree diffs to zero. An unparsable render is written unformatted and reported by the runtime compile, not by the renderer.

Template vars: `Module`, `Name`, `Description`, `Version`, `ProjectName`, `ConfigDir`, `HasDeps`, `HasCli`, `Themes`, plus the collector outputs. The two parsing collectors read the sources as they are on disk at collect time, so a doc comment added to a *generated* contract file shows up on the next build. Native template funcs: `render "<path>"` (read a project file through the transaction, render it with the same vars, nestable) and `copy "<path>"` (verbatim). Missing target = hard error. `README.md` = `render ConfigDir/docs/ReadmeHeader.md` + the `DocIndex` sections + a link to `LICENSE`. It is the single entry point to the docs: there is no index file between it and a doc.

## SmartIO

`smartio.New(deps, path, projectName)`: `Root` = `--path` (normalized; `""`/`.`/`./` = no prefix). Every path an action passes is project-relative; `Root` is joined only at the `deps.Iodeps` boundary, so nothing escapes `--path`. Loads `ignore.yaml` / `paths.yaml` to filter and rewrite listings.

| Call | Effect |
|---|---|
| `WriteFile` | Buffers; refuses to overwrite (disk or pending) |
| `WriteFileOverwrite` | Buffers, replaces. Every generated file uses it |
| `CreateDir`, `RemoveDir` | Pending sets (`RemoveDir` takes files too) |
| `ReadFile`, `Exist`, `IsFile`, `IsDir` | Transaction-aware |
| `List*` | **Disk only** (with ignore/paths applied) |
| `Persist` | Removals, then dir creations, then file writes |

Because listings read disk, an action that runs `build` as a follow-up must `Persist` first. Actions compose by sharing one open `*SmartIO` through their `*Internal` function.

## Runtime

After `Persist`, `RunRuntime(deps, path, runtime)`: `go` = `go mod tidy` (writes `go.sum`) then `go build` over whichever of `./cmd/... ./sandbox/... ./adapters/...` exist (never `./...`, `assets/` is templates); `none` = nothing. Commands that add pass `go`; commands that remove pass `none`.

## Deps install

`dep-install`: list `deplist/<dep>` in embedded assets (empty = unknown dep), add pinned `require` from `depsversion.yaml`, `RenderGroup("deplist/"+dep, {Module})`, persist, then `build`. `dep-remove` is the inverse.

## Dispatch (`climain.go`)

`CliMain(args)`: empty -> general help, exit 2. Match `args[0]` against every command's identifiers; unknown -> exit 2. `dispatch<Name>`: `argvdeps.New(args[1:])`, read each declared flag (a boolean `quiet` replaces `deps.Std.Log` with a no-op immediately), assign defaults, convert and range-check ints/floats, then drain positionals in order. Any unread `-`-prefixed arg = unknown flag; any leftover arg = unexpected argument; missing required = usage error. All exit 2 before the handler. Then `CommandHandler(deps, &entries)`.

## Self-hosting

Agnos regenerates its own `deps.go`, `standard/new.go`, `new.go`, `sandbox.go`, `binds/cli.go`, `climain.go`, every `entries.go` and `help`. `build` must stay idempotent and compilable over this tree. See [Contributing](../Contributing/doc.md#bootstrap).
