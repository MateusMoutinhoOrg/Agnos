# BuildPipeline

`build` command = `verify` (unless `--unsafe`) -> `BuildInternal` -> `Persist` -> runtime. The `build` **action** (the follow-up every other command runs) skips verify so mid-refactor states still regenerate.

## BuildInternal

1. Read `go.mod`, `AgnosConfig/project.yaml` and `AgnosConfig/extensions.yaml` (hard error if either is missing). Fill any catalog key the declaration lacks with its default and write it back. `HasSandbox`, `HasDeps`, `HasCli`, `HasServer`, `HasFront`, `HasDatabase`, `HasExample`, `HasDoc` and `HasReadme` are read straight off those keys — nothing is inferred from a directory. `HasAssets` is the one exception and stays a probe (`assets/sandbox/` exists — the project is itself an agnos-style generator, so its docs name its templates and its own bootstrap).
2. Load `themes.yaml`; `CollectDocs`, merge in `CollectGeneratedDocs` (the docs the enabled groups themselves write — listings read disk, so on a first build they are not there yet), then `GenerateSubdocIndexes` (one `Index.md` per doc with sub-docs; deletes `docs/Index/` left by older versions), which `HasDoc` gates.
3. If `HasCli`: write `info/help/command.yaml` if missing, then `CollectCommands`, then one `generated.new.go` and `generated.input.go` per command. If `HasServer`: `GenerateBuiltinRouteYamls` renders the `route.yaml` of every `utils.GeneratedRoutes` (`health`, `openapi`) into the transaction, and `utils.RouteDirs` lists one still pending, so a server group that just gained a route collects it — and renders its `generated.new.go` and `generated.input.go` — in the same build.
4. Collectors, then the per-unit generators (`GenerateBindingNewFiles`, `GenerateCommandNew`, `GenerateRouteNew`, `GenerateDatabaseNew`, each behind its own key), then `utils.RenderGroup` over `utils.RenderableGroups(extensions)` — see [Asset groups](#asset-groups).

| Collector | Lists | Var | Feeds |
|---|---|---|---|
| `CollectConstructors` | `sandbox/api/*` minus `utils.ConstructorExempt` and every `<x>sandbox.go` / `<x>config.go` | `Constructors` (`Name`, `Package`, `HasNew`; `constructor.go` reads them as `ContractName`, `Package`) | `sandbox/constructors/`, `sandbox/generated.new.go` |
| `CollectEmbeddedStructs` | `sandbox/api/<x>sandbox.go`, `<x>config.go` — after `GenerateApiParts` renders the enabled groups' parts, so a first build sees them | `SandboxStructs`, `ConfigStructs` | `sandbox/api/generated.{sandbox,config}.go` |
| `CollectDepLibs` | `sandbox/deps/<x>/` | `DepLibs` (`Name`, `Title` — the `Deps` field — and `Type`, `Contract` or a remote dep's `Sandbox`) | `sandbox/deps/generated.deps.go` |
| `CollectAdapterImpls` | `adapters/impls/<x>/` | `AdapterImpls` (`Name`) | `docs/LibUsage/doc.md` |
| `CollectBindings` | `adapters/bindings/<x>/binding.yaml` | `Bindings` (`BindingName`, `Adapters`) | `GenerateBindingNewFiles` -> `adapters/bindings/<x>/generated.new.go` |
| `CollectCommands` | `commands/<category>/<x>/command.yaml` | `Commands` (the declaration itself: `CommandName`, identifiers, category, summary, `Flags`/`Args` with ids, types, defaults, bounds) | `generated.new.go`, `internal/cli/generated.new.go` |
| `CollectDocs` | `docs/**/doc.yaml` | doc tree sorted by `order` then name | `**/Index.md`, `DocIndex` |
| `CollectGeneratedDocs` | `assets/doc*/docs/*/doc.yaml` over `utils.DocGroups(extensions)`, rendered | merged into the doc tree | same |
| `CollectDocIndex` | the merged tree grouped by theme | `DocIndex` (per theme: `Name`, `Description`, `Docs`) | `README.md`. A theme no doc names renders no section |
| `CollectPublicApi` | `sandbox/api/*.go` parsed by `deps.GoimportsDeps` | `PublicApi` (per file: `Path`, `Doc`, `Types`, `Constants`, `Variables`, `Functions`; exported only, doc comments flattened to one table line) | `docs/PublicApi/doc.md` |
| `CollectDepContracts` | `sandbox/deps/<x>/*.go`, same parse | `DepContracts` (`Name`, `Title`, `Files`) | `docs/PublicApi/doc.md` |
| `CollectRoutes` | `routes/<x>/route.yaml` | `Routes` (the declaration itself: `RouteName`, `Methods`, `ResponseType`, `Paths`, `Parameters`, `Body`, `SchemaJson`, `BodyStructs`), **ordered to run**: `Priority`, then name | `route_new.go`, `route_input.go`, `internal/server/generated.new.go` |
| `CollectRouteDocs` | `routes/<x>/route.yaml` (visible ones), grouped by category in first-seen order, each crossed with the routes in front of it | `RouteDocs` (per category: `Routes` with `Method`, `Pattern`, `Summary`, `Description`, `Address`/`Parameters` as table rows, `Body` with its schema rows and a json sample, `Requests` as generated `curl` calls, `Examples`, `Statuses`, `Middlewares`) — plain words, for whoever calls the server | `docs/Routes/doc.md`, `GenerateRoutePages` |
| `CollectOpenApi` | the same declarations, in run order, each crossed with the routes in front of it | `OpenApi` (one OpenAPI 3.0.3 document as text, a `routeDocJson` tree printed by `routeDocJsonText`; mapping in [RouteYaml](../../assets/doc-server/docs/RouteYaml/doc.md#openapi)) | `docs/Routes/openapi.json`, the `openapi` route's `handler.go` |
| `CollectCommandDocs` | `commands/<category>/<x>/command.yaml` (visible ones), grouped by category in first-seen order | `CommandDocs` (per category: `Commands` with `Identifier`, `Aliases`, `Summary`, `Description`, `Usage`, `Flags`/`Args` as table rows, `Examples`) | `docs/Commands/doc.md` |
| `CollectDatabases` | `databases/<x>/database.yaml` | `Databases` (per database: `DatabaseName`, `Package`, `Type`, `KeyPrefix`, `Tables` as `databasedeps.Item` literals, `Records` as the Go structs, `Methods` with `Kind`, `Params`, `Results`, `Args`) | `GenerateDatabaseNew` -> `databases/<x>/generated.{api,new,methods}.go` |
| `CollectDatabaseDocs` | `databases/<x>/database.yaml` | `DatabaseDocs` (per database: `Tables` with their fields as table rows, `Methods` with the whole signature) | `docs/Databases/doc.md`, `GenerateDatabasePages` |
| `CollectStructure` | `AgnosConfig/structure.yaml` (structureconf) | `Structure` (one `Line` per item, depth-indented and padded to a common description column) | `docs/Structure/doc.md` |

A Go file the build rewrites whole is written as `generated.<name>.go` (`utils.GeneratedPrefix`) through `utils.WriteGenerated`, which puts `// Code generated by <generator>. DO NOT EDIT.` on top and drops the `<name>.go` an older build left beside it; until that write, `utils.DropSuperseded` keeps a collector from counting both. `MigrateGeneratedNames`, first in `BuildInternal`, moves an older tree: the import of `internal/generated/{cli,server,config}` in the project's files, the old packages, and a remote dep's copy and shim.

Every render whose destination ends in `.go` is passed through `deps.GoimportsDeps.Format` (`go/format`, i.e. `gofmt`) before it is written, so generated Go is byte-identical to what a formatting editor saves and a regenerated tree diffs to zero. An unparsable render is written unformatted and reported by the runtime compile, not by the renderer.

Template vars: `Module`, `ProjectName` (the project being generated, as `project-name` of `project.yaml` spells it), `GeneratorName` (the cli running the build — `agnos`; every command agnos owns is spelled with it, never with `ProjectName`), `Version`, `GeneratorVersion`, `ConfigDir`, `StructureConfFile`, `HasSandbox`, `HasDeps`, `HasCli`, `HasServer`, `HasFront`, `HasDatabase`, `HasExample`, `HasDoc`, `HasReadme`, `HasAssets`, `Themes`, plus the collector outputs. The two parsing collectors read the sources as they are on disk at collect time, so a doc comment added to a *generated* contract file shows up on the next build. Native template funcs: `render "<path>"` (read a project file through the transaction, render it with the same vars, nestable) and `copy "<path>"` (verbatim). Missing target = hard error. `README.md` = `render ConfigDir/docs/ReadmeHeader.md` + the `DocIndex` sections + a link to `LICENSE`. It is the single entry point to the docs: there is no index file between it and a doc.

## Asset groups

`assets/<group>/` is one group; the group's name is the condition under which `build` renders
it, and `utils.AssetGroups()` is the whole list. A group named after an extension renders when
that extension is on; a group named `doc-<a>-<b>` renders when `doc` and every `<a>`,
`<b>` it names are on, and one named `<a>-<b>` (`database-cli`) when `<a>` and `<b>` are.
`assets/start/` is outside this: it is written once, by `start`.

| Group | Renders when | Holds |
|---|---|---|
| `sandbox` | `sandbox` | `sandbox/generated.new.go`, `api/generated.{sandbox,config}.go` (embedding every struct of `api/<x>sandbox.go` / `api/<x>config.go`; native: `Deps`, `Config`), `internal/config/generated.new.go` |
| `deps` | `deps` | `sandbox/deps/generated.deps.go` |
| `cli` | `cli` | `cmd/main`, `api/generated.{cli,command,trigger,clisandbox}.go` (aliases of `OpinionatedAgnosCli`), `internal/cli/generated.new.go` (the registry), `info/help`, `info/version`, `middleware/help_flag` |
| `server` | `server` | `api/generated.{server,route,serversandbox}.go` (aliases of `OpinionatedAgnosServer`), `internal/server/generated.new.go` (the registry), the health and openapi routes |
| `database` | `database` | `api/generated.databaseconfig.go`: `Config.DatabaseDir`, the folder every key-prefix is a path inside, and `DefaultDatabaseDir` (`data`), which `config/generated.new.go` starts it at |
| `database-cli` | `database` + `cli` | `middleware/database_dir`: reads `--database` into `Config.DatabaseDir` in front of every command line |
| `doc` | `doc` | `docs/{Adapters,DepList,Extensions,GeneratedFiles,LibUsage,PublicApi,Requirements,Rules,Structure,Workflow}` |
| `doc-cli` | `doc` + `cli` | `docs/{CliInstall,Commands}` |
| `doc-server` | `doc` + `server` | `docs/{RouteYaml,Routes,ServerUsage}`, `docs/Routes/openapi.json` |
| `doc-front` | `doc` + `front` | `docs/FrontUsage` |
| `doc-database` | `doc` + `database` | `docs/Databases` |
| `doc-backoffice` | `doc` + `backoffice` | `docs/{Backoffice,Backups}`. The backoffice has no code group: `backoffice-init` writes `assets/templates/backoffice/**` once |
| `doc-example` | `doc` + `example` | `docs/LibExamples` |
| `doc-example-cli` | `doc` + `example` + `cli` | `docs/CliExamples` |
| `readme` | `readme` | `README.md` |

`front` has no code group: its code is the `OpinionatedAgnosFront` lib `front-init` installs, so
it renders its pages alone. `database`'s code is the `OpinionatedAgnosDatabase` lib the same way;
its groups hold only where the data lives (`database`, `database-cli`).

A code group is rendered ahead of the build by `utils.RenderExtensionCode`, which `SetExtension`
calls when a key turns on: every code group requiring that key whose requirements are now all on —
so `database-cli` lands on `database-init` with the cli on and on `cli-init` with the database on,
before the build's collectors read its `command.yaml`. The build renders groups after its
collectors, so a change to a group's `command.yaml` reaches its `generated.new.go` one build later.

The code every project shares — the dispatch, the binders, the matchers, json-schema, the file
layer, the database readers — is not rendered at all: it is the four `OpinionatedAgnos<X>` catalog
deps (`assets/{dep-catalog,adapter-catalog}/OpinionatedAgnos<X>/`), installed by each mechanic's `-init`.
What the build still renders is what changes per project: the registries
`internal/{cli,server}/generated.new.go` and `internal/config/generated.new.go`.
`utils.RemoveRetiredGenerated` drops what an older build generated under `internal/generated/`,
for every mechanic that is on.

A group marked `Code: true` — the four `sandbox*` ones — carries Go the collectors read back
off disk, so `utils.SetExtension` renders it once into the transaction that turns the mechanic
on. Without that pass the build that follows would collect `sandbox/api/` as it was before the
mechanic existed and write a `sandbox.go` missing the field `cmd/main/generated.main.go` already uses.

## StagedFS

`stagedfs.New(sandbox, path, projectName)`: `Root` = `--path` (normalized; `""`/`.`/`./` = no prefix). Every path an action passes is project-relative; `Root` is joined only at the `deps.IoDeps` boundary, so nothing escapes `--path`. Loads `paths.yaml` to rewrite the paths it is handed.

| Call | Effect |
|---|---|
| `CreateFile` | Buffers; refuses to overwrite (disk or pending) |
| `WriteFile` | Buffers, replaces. Every generated file uses it |
| `CreateDir`, `RemoveDir` | Pending sets (`RemoveDir` takes files too) |
| `ReadFile`, `Exists`, `IsFile`, `IsDir` | Transaction-aware |
| `List*` | **Disk only** (with ignore/paths applied) |
| `Persist` | Removals, then dir creations, then file writes |

Because listings read disk, an action that runs `build` as a follow-up must `Persist` first. Actions compose by sharing one open `*StagedFS` through their `*Internal` function.

## Runtime

After `Persist`, `RunRuntime(deps, path, runtime)`: `go` = `go mod tidy` (writes `go.sum`) then `go build` over whichever of `./cmd/... ./sandbox/... ./adapters/...` exist (never `./...`, `assets/` is templates); `none` = nothing. Commands that add pass `go`; commands that remove pass `none`.

## Deps install

`add-dep <name>`: read `dep-catalog/<dep>/dep.yaml` in embedded assets (missing = unknown dep), pick `--adapter` or its `default-adapter`, read `adapter-catalog/<adapter>/adapter.yaml` (its `dep:` must match), `RenderGroupExcept` both catalogs minus their own declaration, write the adapter's declaration to `adapters/impls/<adapter>/adapter.yaml`, add the `require` its `module:` pins, enroll the adapter in every binding, persist, then `build`. `remove-dep` is the inverse: every adapter whose declaration names the dep, then the contract.

`add-dep <module>` (an argument holding a `/`): resolve the module through `go list -m -json` and, failing that, `go mod download -json` (`RemoteModule`); parse `<dir>/sandbox/api/*.go` into an `apishape.Api`; reject it with `apishape.Violations`; copy each file to `sandbox/deps/<name>/` with only the package clause rewritten, then `Format`; plan the converters with `apishape.Converters` and render `templates/remote_shim.go` to `adapters/impls/<name>/<name>.go`; write its `adapter.yaml` with `origin: generated`; `AddRequire`; enroll; `build`. `set-dep <name> --version` reads the module back out of that `adapter.yaml` and runs the same path again.

## Dispatch (`OpinionatedAgnosCli.Main`)

`Cli.Main(args)` — filled by the generated registry `internal/cli/generated.new.go` — hands `MainProps{Cli: &sandbox.Cli, Args, NewProps, StdDeps: &sandbox.Deps.StdDeps, ArgvDeps}` to the lib's `Main`. It runs every `api.Command` of `Cli.Commands` in run order, each on a copy made by `BindCommand`: `Matches` (segments, arg and flag triggers), then the binder — args off the segments, flags through `ArgvDeps`, defaults, conversion and bounds, `Input` filled by its `id` tags by reflection — then the command's `Handle`. The first command that answers (a status, or a `Printf`) ends the chain; a strict one that answers nothing exits `0`; a middleware that answers nothing hands the line on. Every unanswered ending — no command, a value that will not bind, an unconsumed token, a returned failure, a panic — is raised through `Cli.Fail` to one of the project's `handle_*.go`. Nothing in the lib is generated per command — the registry is where the set is spelled, out of each package's `NewCommand`. The server's `OpinionatedAgnosServer.Main` is the same chain over `Server.Routes`.

## Self-hosting

Agnos regenerates its own `generated.deps.go`, `standard/generated.new.go`, `sandbox/generated.new.go`, `generated.sandbox.go`, `generated.command.go`, `internal/cli/generated.new.go`, every command's `generated.new.go` and `generated.input.go`, `info/help`, and runs on its own `OpinionatedAgnosCli` (installed from the catalog, so `check_dep_catalog`/`check_adapter_catalog` hold the catalog to the copy it runs on). It turns `database` on nowhere: agnos declares no database of its own, so the mechanic is exercised by `examples/{cli,lib}/database` rather than by this tree. `build` must stay idempotent and compilable over this tree. See [Contributing](../Contributing/doc.md#bootstrap).
