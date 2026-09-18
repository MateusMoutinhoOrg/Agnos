# BuildPipeline

`build` command = `verify` (unless `--unsafe`) -> `BuildInternal` -> `Persist` -> runtime. The `build` **action** (the follow-up every other command runs) skips verify so mid-refactor states still regenerate.

## BuildInternal

1. Read `go.mod`, `AgnosConfig/project.yaml` and `AgnosConfig/extensions.yaml` (hard error if either is missing). Fill any catalog key the declaration lacks with its default and write it back. `HasSandbox`, `HasDeps`, `HasCli`, `HasServer`, `HasFront`, `HasDatabase`, `HasExample`, `HasDoc` and `HasReadme` are read straight off those keys — nothing is inferred from a directory. `HasAssets` is the one exception and stays a probe (`assets/sandbox/` exists — the project is itself an agnos-style generator, so its docs name its templates and its own bootstrap).
2. Load `themes.yaml`; `CollectDocs`, merge in `CollectGeneratedDocs` (the docs the enabled groups themselves write — listings read disk, so on a first build they are not there yet), then `GenerateSubdocIndexes` (one `Index.md` per doc with sub-docs; deletes `docs/Index/` left by older versions), which `HasDoc` gates.
3. If `HasCli`: write `help/entries.yaml` if missing, then `CollectCommands`, then one `new.go` per command.
4. Collectors, then the per-unit generators (`GenerateAvailableNews`, `GenerateCommandNew`, `GenerateRouteNew`, `GenerateDatabaseNew`, each behind its own key), then `utils.RenderGroup` over `utils.RenderableGroups(extensions)` — see [Asset groups](#asset-groups).

| Collector | Lists | Var | Feeds |
|---|---|---|---|
| `CollectConstructors` | `sandbox/api/*` minus `sandbox.go`, `command.go`, `route.go` | `Constructors` (`Name`, `Package`, `HasNew`) | `sandbox/api/sandbox.go`, `sandbox/new.go` |
| `CollectDepsLibs` | `sandbox/deps/<x>/` | `DepsLibs` (`Title`, `Name`) | `sandbox/deps/deps.go` |
| `CollectAdapterLibs` | `adapters/libs/<x>/` | `AdapterLibs` (`Name`) | `docs/LibUsage/doc.md` |
| `CollectAvailables` | `adapters/availables/<x>/available.yaml` | `Availables` (`Name`, `Adapters`) | `GenerateAvailableNews` -> `adapters/availables/<x>/new.go` |
| `CollectCommands` | `commands/<x>/entries.yaml` | `Commands` (the declaration itself: identifiers, category, help, `Flags`/`Args` with ids, types, defaults, bounds) | `new.go`, `internal/cli/new.go` |
| `CollectDocs` | `docs/**/props.yaml` | doc tree sorted by `order` then name | `**/Index.md`, `DocIndex` |
| `CollectGeneratedDocs` | `assets/doc*/docs/*/props.yaml` over `utils.DocGroups(extensions)`, rendered | merged into the doc tree | same |
| `CollectDocIndex` | the merged tree grouped by theme | `DocIndex` (per theme: `Name`, `Description`, `Docs`) | `README.md`. A theme no doc names renders no section |
| `CollectPublicApi` | `sandbox/api/*.go` parsed by `deps.Goimportsdeps` | `PublicApi` (per file: `Path`, `Doc`, `Types`, `Constants`, `Variables`, `Functions`; exported only, doc comments flattened to one table line) | `docs/PublicApi/doc.md` |
| `CollectDepsApi` | `sandbox/deps/<x>/*.go`, same parse | `DepsApi` (`Name`, `Title`, `Files`) | `docs/PublicApi/doc.md` |
| `CollectRoutes` | `routes/<x>/route.yaml` | `Routes` (the declaration itself: `Paths`, `Headers`, `Params`, `Body`, `SchemaJson`, `BodyStructs`), **ordered for matching**: most `identifier`s, then longest, then pattern | `route_new.go`, `internal/server/new.go` |
| `CollectRouteDocs` | `routes/<x>/route.yaml` (visible ones), grouped by category in first-seen order | `RouteDocs` (per category: `Routes` with `Method`, `Pattern`, `Help`, `LongDescription`, `Fields` as table rows, `Body`, `Examples`) | `docs/Routes/doc.md` |
| `CollectCommandDocs` | `commands/<x>/entries.yaml` (visible ones), grouped by category in first-seen order | `CommandDocs` (per category: `Commands` with `Identifier`, `Aliases`, `Help`, `LongDescription`, `Usage`, `Flags`/`Args` as table rows, `Examples`) | `docs/Commands/doc.md` |
| `CollectDatabases` | `databases/<x>/specs.yaml` | `Databases` (per database: `Package`, `Type`, `Prefix`, `Tables` as `database.Item` literals, `Records` as the Go structs, `Methods` with `Kind`, `Params`, `Results`, `Args`) | `GenerateDatabaseNew` -> `databases/<x>/{api.go,new.go,methods.go}` |
| `CollectDatabaseDocs` | `databases/<x>/specs.yaml` | `DatabaseDocs` (per database: `Tables` with their fields as table rows, `Methods` with the whole signature) | `docs/Databases/doc.md`, `GenerateDatabasePages` |
| `CollectStructure` | `AgnosConfig/structure.yaml` (structureconf) | `Structure` (one `Line` per item, depth-indented and padded to a common description column) | `docs/Structure/doc.md` |

Every render whose destination ends in `.go` is passed through `deps.Goimportsdeps.Format` (`go/format`, i.e. `gofmt`) before it is written, so generated Go is byte-identical to what a formatting editor saves and a regenerated tree diffs to zero. An unparsable render is written unformatted and reported by the runtime compile, not by the renderer.

Template vars: `Module`, `Name` (the project being generated), `GeneratorName` (the cli running the build — `agnos`; every command agnos owns is spelled with it, never with `Name`), `Version`, `ProjectName`, `ConfigDir`, `StructureConfFile`, `HasSandbox`, `HasDeps`, `HasCli`, `HasServer`, `HasFront`, `HasDatabase`, `HasExample`, `HasDoc`, `HasReadme`, `HasAssets`, `Themes`, plus the collector outputs. The two parsing collectors read the sources as they are on disk at collect time, so a doc comment added to a *generated* contract file shows up on the next build. Native template funcs: `render "<path>"` (read a project file through the transaction, render it with the same vars, nestable) and `copy "<path>"` (verbatim). Missing target = hard error. `README.md` = `render ConfigDir/docs/ReadmeHeader.md` + the `DocIndex` sections + a link to `LICENSE`. It is the single entry point to the docs: there is no index file between it and a doc.

## Asset groups

`assets/<group>/` is one group; the group's name is the condition under which `build` renders
it, and `utils.AssetGroups()` is the whole list. A group named after an extension renders when
that extension is on; a group named `doc-<a>-<b>` renders when `doc` and every `sandbox-<a>`,
`sandbox-<b>` it names are on. `assets/start/` is outside this: it is written once, by `start`.

| Group | Renders when | Holds |
|---|---|---|
| `sandbox` | `sandbox` | `sandbox/new.go`, `api/{sandbox,config}.go`, `internal/config/new.go` |
| `sandbox-deps` | `sandbox-deps` | `sandbox/deps/deps.go` |
| `sandbox-cli` | `sandbox-cli` | `cmd/main`, `api/{cli,command}.go`, `internal/cli/`, `help`, `version` |
| `sandbox-server` | `sandbox-server` | `api/{server,route}.go`, `internal/{server,routes/health,routeio}` |
| `sandbox-front` | `sandbox-front` | `internal/pageio/` |
| `sandbox-database` | `sandbox-database` | `internal/databaseio/` |
| `doc` | `doc` | `docs/{Adapters,DepList,EntriesYaml,Extensions,GeneratedFiles,LibUsage,PublicApi,Requirements,Rules,Structure,Workflow}` |
| `doc-cli` | `doc` + `sandbox-cli` | `docs/{CliInstall,Commands}` |
| `doc-server` | `doc` + `sandbox-server` | `docs/{RouteYaml,Routes,ServerUsage}` |
| `doc-front` | `doc` + `sandbox-front` | `docs/FrontUsage` |
| `doc-database` | `doc` + `sandbox-database` | `docs/Databases` |
| `doc-example` | `doc` + `sandbox-example` | `docs/LibExamples` |
| `doc-example-cli` | `doc` + `sandbox-example` + `sandbox-cli` | `docs/CliExamples` |
| `readme` | `readme` | `README.md` |

A group marked `Code: true` — the six `sandbox*` ones — carries Go the collectors read back
off disk, so `utils.SetExtension` renders it once into the transaction that turns the mechanic
on. Without that pass the build that follows would collect `sandbox/api/` as it was before the
mechanic existed and write a `sandbox.go` missing the field `cmd/main/main.go` already uses.

## SmartIO

`smartio.New(deps, path, projectName)`: `Root` = `--path` (normalized; `""`/`.`/`./` = no prefix). Every path an action passes is project-relative; `Root` is joined only at the `deps.Iodeps` boundary, so nothing escapes `--path`. Loads `paths.yaml` to rewrite the paths it is handed.

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

`add-dep <name>`: read `deplist/<dep>/dep.yaml` in embedded assets (missing = unknown dep), pick `--adapter` or its `default-adapter`, read `adapterlist/<adapter>/adapter.yaml` (its `dep:` must match), `RenderGroupExcept` both catalogs minus their own declaration, write the adapter's declaration to `adapters/libs/<adapter>/adapter.yaml`, add the `require` its `module:` pins, enroll the adapter in every available, persist, then `build`. `remove-dep` is the inverse: every adapter whose declaration names the dep, then the contract.

`add-dep <module>` (an argument holding a `/`): resolve the module through `go list -m -json` and, failing that, `go mod download -json` (`RemoteModule`); parse `<dir>/sandbox/api/*.go` into an `apishape.Api`; reject it with `apishape.Violations`; copy each file to `sandbox/deps/<name>/` with only the package clause rewritten, then `Format`; plan the converters with `apishape.Converters` and render `templates/remote_shim.go` to `adapters/libs/<name>/<name>.go`; write its `adapter.yaml` with `origin: generated`; `AddRequire`; enroll; `build`. `set-dep <name> --version` reads the module back out of that `adapter.yaml` and runs the same path again.

## Dispatch (`climain.go`)

`CliMain(args)`: empty -> general help, exit 2. Match `args[0]` against the identifiers of every `api.Command` of `sandbox.Cli.Commands`, then copy the match with `api.BindCommand`; unknown -> exit 2. Then one generic `runCommand`: `argvdeps.New(args[1:])`, read each declared flag (a boolean `quiet` replaces `sandbox.Deps.Std.Log` with a no-op immediately), assign defaults, convert and range-check ints/floats, then drain positionals in order, binding every value into `command.Items` under its declared id. Any unread `-`-prefixed arg = unknown flag; any leftover arg = unexpected argument; missing required = usage error. All exit 2 before the handler. Then `command.Handler(command)`, which is that package's `CommandHandler(sandbox, command)`. Nothing in this file is generated per command — `internal/cli/new.go` is where the set is spelled, out of each package's `NewCommand`.

## Self-hosting

Agnos regenerates its own `deps.go`, `standard/new.go`, `new.go`, `sandbox.go`, `command.go`, `internal/cli/new.go`, `climain.go`, every command's `new.go` and `help`. It turns `sandbox-database` on nowhere: agnos declares no database of its own, so the mechanic is exercised by `examples/{cli,lib}/database` rather than by this tree. `build` must stay idempotent and compilable over this tree. See [Contributing](../Contributing/doc.md#bootstrap).
