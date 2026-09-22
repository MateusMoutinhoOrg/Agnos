# Contributing

Recipes for changing **agnos itself**. Everything an agnos project does through agnos commands — declare a command, its flags and args, add a doc, install a dep, describe a path — is in [Workflow](../Workflow/doc.md), and the rules every recipe holds to are in [Rules](../Rules/doc.md). Only what is specific to this repo is below.

## Bootstrap

Never run an installed `agnos build` on this repo after changing templates, collectors, deps or adapters: it rewrites the tree to its older shape.

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build                                   # verify + regenerate + compile
./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
./release/bootstrap.bin local-install                           # install the result
```

Release: bump `version` in `AgnosConfig/project.yaml`, then `./release/bootstrap.bin build` and `./release/bootstrap.bin exec-test --update` — the bumped version is rendered into `docs/Requirements/doc.md`, so every golden carrying that page moves — then `agnos publish` (or `agnos compile --target all` for the binaries alone).

## Add an action

1. `sandbox/internal/actions/<name>/<name>_internal.go`: `func <Name>Internal(sandbox, io *smartio.SmartIO, ...) error`. Project-relative paths only. Log via `sandbox.Deps.Std.Log`, fail via `sandbox.Deps.Std.Errorf`, never `Printf`.
2. `<name>.go`: `func <Name>(sandbox, ...) error` = `smartio.New(sandbox, path, sandbox.Config.ProjectName)` -> internal -> `io.Persist()` -> `buildAction.Build(sandbox, api.BuildProps{Path, Runtime})` (`RuntimeGo` if it adds, `RuntimeNone` if it removes). Props with more than three values go in a struct in `sandbox/api/actions.go`.
3. Add the field to `api.Actions` and the assignment to `NewActions` in `sandbox/internal/actions/new.go`.
4. Comment the new field: its row in [PublicApi](../PublicApi/doc.md) is generated from that comment.

## Add a command to agnos

Declare it with the bootstrap binary, as in [Workflow](../Workflow/doc.md#change-the-command-surface), with a `--category` this repo already uses (Core Commands, Cli System, Server System, Front System, Database System, Dependencies, Dependency System, Info) and the two flags every agnos command carries:

```bash
./release/bootstrap.bin add-command <name> --help "..." --category "Core Commands"
./release/bootstrap.bin add-flag path --command <name> --default . --description "the dir holding the project (defaults to the current directory)"
./release/bootstrap.bin add-flag quiet --command <name> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
```

`handler.go` calls the action, returns `api.ExitFailure` on error and `Printf`s any result.

`interview` needs nothing for the new command: it generates its questions from the declaration. It is the one feature of agnos written for a beginner rather than for an llm ([Interview](../Interview/doc.md)), and seven of its tables take an entry — only for a command that is one of these:

| Teach | When |
|---|---|
| `interview/suggest.go` | one of the command's fields names something that already exists: the entry says where to read that list, off the project at `--path`, or off the thing a required `--route` / `--command` names — `scopeFields` is what gets that one asked first |
| `interview/state.go` | the command opens or fills a layer: a row in `areas` for a new category, one in `extensionInit` for a new `<x>-init`, a rung in `nextSteps` for the step that offers it, and a name in `scaffoldedUnits` for a unit its init writes |
| `interview/ruled_out.go` | two of the command's fields exclude each other, or one only applies to some `--type`: the entry is the same rule the action refuses with, so the question is never asked |
| `interview/required.go` | a field the declaration calls optional that some projects make mandatory: the entry is the condition the handler rejects on, so the question is asked as the required one it is |
| `interview/danger.go` | the command removes or overwrites something: the entry keeps it off the row a blind enter lands on, and — for a purge — names the units it takes with it |
| `interview/normalize.go` | the command rewrites a name before writing it down: the entry is the noun, and the confirm screen says what the typed name becomes |
| `interview/followup.go` | the command leaves a thing half-declared, or finishes one another command started: the entry is the command offered next and the answer it inherits, so the route just declared is not typed again |

`add-flag`/`add-arg` read the command line they are typed on, so never pass a value that is *exactly* one of their own flag spellings (`--identifier --example`, `--example --required`): the parser counts it as an occurrence and the declaration comes out polluted. `--help` joins that set on any command that does not declare a flag of that name, where the dispatch reads it as a request for that command's help screen. Declare such a flag without the short alias instead, and word its `--example` as a whole command line.

## Add a layer (cli, server, …)

A layer is an extension plus an `<x>-init`/`<x>-purge` pair, and the server layer is the pattern to copy — file for file, it mirrors the cli one. The front layer is the third column: it declares no unit of its own, because a page **is** a route. The database layer is the fourth, and the one that breaks the shape: its methods are typed by table, so there is no generic dispatch and no field of the `Sandbox` — the declaration is spelled out into three generated files instead of being read back at runtime.

| Concept | CLI | Server | Front | Database |
|---|---|---|---|---|
| External input contract | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` | — (`embeddeps` + `templatedeps`) | `sandbox/deps/database/` (a remote dep, not a catalog one) |
| Surface + constructor | `sandbox/api/cli.go`, `sandbox/api/command.go`, `sandbox/internal/cli/new.go` -> `Cli.Commands` | `sandbox/api/server.go`, `sandbox/api/route.go`, `sandbox/internal/server/new.go` -> `Server.Routes` | — (served through the server's) | — (typed by table: `<db>.New(sandbox)` on the spot) |
| Constructor package | `sandbox/constructors/cli/` | `sandbox/constructors/server/` | — | — |
| Dispatch (generic) | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/servermain.go` | — | — (the methods are generated, not dispatched) |
| Shared package | — | `sandbox/internal/routeio/` | `sandbox/internal/pageio/` | `sandbox/internal/databaseio/` |
| Answer to bad input | the dispatch, exit 2 | `server/handle_*.go`, written once by `build` | — (the server's) | — |
| Declared unit | `commands/<name>/entries.yaml` -> generated `new.go` | `routes/<name>/route.yaml` -> generated `new.go` | `routes/<page>/route.yaml` + `assets/frontend/pages/<page>.html` | `databases/<db>/specs.yaml` -> generated `api.go`, `new.go`, `methods.go` (+ hand-written `methods_custom.go`) |
| Parsable | `parsables/commandconf/` | `parsables/routeconf/` | — (`routeconf`) | `parsables/databaseconf/` |
| Collectors | `collect_commands.go`, `collect_command_docs.go` | `collect_routes.go`, `collect_route_docs.go` | `utils/front_mount.go` (`RenderExtensionCode` reads it too, so it sits in `utils`) | `collect_databases.go`, `collect_database_docs.go` |
| Per-unit generator | `generate_command_new.go` | `generate_route_new.go` (+ `generate_error_handlers.go`, once) | — (`generate_route_new.go`) | `generate_database_new.go` (three files per unit) |
| Asset groups | `assets/sandbox-cli/`, `assets/doc-cli/`, `assets/doc-example-cli/` | `assets/sandbox-server/`, `assets/doc-server/` | `assets/sandbox-front/`, `assets/doc-front/` | `assets/sandbox-database/`, `assets/doc-database/` |
| Extension key | `sandbox-cli` | `sandbox-server` | `sandbox-front` | `sandbox-database` |
| Init / purge | `cli-init` / `cli-purge` | `server-init` / `server-purge` | `front-init` / `front-purge` | `database-init` / `database-purge` |
| Interview gate (`interview/state.go`) | `Cli System` -> `sandbox-cli`, step `cli-init` | `Server System` -> `sandbox-server`, step `server-init` | `Front System` -> `sandbox-front`, step `front-init` | `Database System` -> `sandbox-database`, step `database-init` |
| Verify | `check_sandbox.go` et al | `check_routes.go` | — (`check_routes.go`) | `check_databases.go` |

A layer is an extension, so adding one is [Add an extension](#add-an-extension) plus the rows above. A layer whose init needs another layer calls the other one's `<X>InitInternal` on the *same* open SmartIO — `server_init` does that with `cli_init`, `front_init` with `server_init` — so there is no intermediate `Persist` and no intermediate `build`. The dep installs are the exception: `<X>InitInternal` writes nothing to `go.mod`, so a composing init calls the other's exported `InstallDeps` first.

`<X>InitInternal` renders no group of its own: it flips the key with `utils.SetExtension`, which renders the mechanic's code group into the same transaction, and the follow-up `build` renders the rest. `<X>PurgeInternal` removes `utils.ExtensionFiles(sandbox, <key>)` — every group the mechanic owns, its pages included — plus the directories the layer owns whole — its `sandbox/constructors/<x>/` included, since that package names what is being removed — then writes the key back as `false`.

`sandbox/internal/databaseio/` is the same third package one layer over: a database may not import another one, so the readers every generated `methods.go` shares live beside them. Its one rule is the layer's: nothing converts a stored value without a comma-ok, so a malformed record is an error and never a panic.

`sandbox/internal/routeio/` exists because a route may not import `internal/server`: shared route code goes in a third package both may import. It is also where the two readers that put the dep names back on a bound route live — `RequestOf` and `ResponseOf`, because `api.Route` carries the request and the response as `any`: `sandbox/api/` may name no type of `sandbox/deps`. `sandbox/internal/pageio/` is the same shape one layer up.

That same import rule is why `routeio.Fail` reaches the project's `handle_*.go` through the `Fail` field of `api.Server` rather than by calling them: `internal/server` imports every route package, so nothing under `routes/` may import it back, and a function field on the api is how this repo already crosses that line everywhere else. The six files are written by `build`, not by `server-init`, for the same reason `sandbox/constructors/<x>/constructor.go` is — a project that gained the layer before they existed picks them up on its next build, and the generated `server/new.go` always has something to call.

## Add an extension

An extension is one generation mechanic and it is declared in three places:

1. `sandbox/internal/utils/extensions_conf.go`: an `Extension<X>` const and a row in `ExtensionCatalog()` (key, default for a fresh project, the line `list-extensions` prints).
2. `sandbox/internal/utils/asset_groups.go`: one row per group in `AssetGroups()` — the mechanic's own group (`Code: true` when it carries Go the collectors read back) and any `doc-<x>` group it gates.
3. `assets/<group>/**`: the templates, at the path each renders to.

Then `assets/start/AgnosConfig/extensions.yaml` gains the key at its default, and the page listing it — `assets/doc/docs/Extensions/doc.md` — gains its row. A mechanic with files to scaffold or to delete also gets an `<x>-init` / `<x>-purge` pair; one that only turns generation on and off is reached through `enable-extension` / `disable-extension` and needs no command of its own.

An existing project does not have to be touched: `build` fills a key the catalog gained with its default and writes `extensions.yaml` back.

The one thing still read off the tree is `hasAssets` (`assets/sandbox/` exists), because "this project is itself a generator" is a property of the tree, not a mechanic anyone turns on.

## Add a contract + adapter lib

The two halves are in [Workflow](../Workflow/doc.md#add-a-dependency), the three units they belong to in [Adapters](../Adapters/doc.md). Per-call objects use a constructor field (`New func(...) Thing`) like `argvdeps`/`requestdeps`. Doc-comment every exported declaration and every field you want described in [PublicApi](../PublicApi/doc.md). To make it installable elsewhere, mirror it into `assets/deplist/` below.

A second implementation of a contract that already exists is an adapter alone: only step 2 below, plus `add-adapter` and `set-adapter` to bind it.

Whatever you put in `sandbox/api/` has to stay convertible — `verify` runs `check_api_shape.go` on every run, because another repo may install this one as a dep and copy that package. The rule and the converter planner are in `sandbox/internal/apishape/`; the shim is rendered from `assets/templates/remote_shim.go`.

## Add an installable dep

The catalog is two catalogs: `assets/deplist/<dep>/` holds the contract, `assets/adapterlist/<adapter>/` holds one implementation of it, so a dep can gain a second adapter without moving.

1. Mirror the contract under `assets/deplist/<dep>/sandbox/deps/<dep>/`, replacing this module path with `{{.Module}}`, and write `assets/deplist/<dep>/dep.yaml` (`name`, `field`, `help`, `default-adapter`).
2. Mirror the adapter under `assets/adapterlist/<adapter>/adapters/libs/<adapter>/`, same substitution, and write `assets/adapterlist/<adapter>/adapter.yaml` (`name`, `dep`, `help`, `module` — `""` when it needs nothing beyond the stdlib — and `origin: catalog`).
3. Bootstrap, test with `add-dep`/`remove-dep` on a scratch project. Add a row to [DepList](../DepList/doc.md).

Neither `dep.yaml` nor `adapter.yaml` is part of the mirror: the first is installed nowhere, and the second is installed to `adapters/libs/<adapter>/adapter.yaml`, which is what tells the tree later which dep that adapter fills.

## Add a template or collector

- Template: `assets/<group>/<target path>`, a `text/template` over the vars in [BuildPipeline](../BuildPipeline/doc.md#buildinternal). The groups are in [Asset groups](../BuildPipeline/doc.md#asset-groups); a new one is a new extension. A scaffold that renders to a file which is itself a template — `page_html.html`, or any `{{` inside `front_main.js` — escapes its own braces (`{{ "{{ .Title }}" }}`), or the outer render eats them. Single-destination scaffolds go in `assets/templates/` and are rendered with `utils.RenderTemplateToDest`. Add a row for the new destination to `assets/doc/docs/GeneratedFiles/doc.md`.
- Collector: `sandbox/internal/actions/build/collect_<x>.go`, `func Collect<X>(sandbox, io) []string` listing one dir and title-casing the last segment; add `"<X>": Collect<X>(sandbox, io)` to the vars map in `build_internal.go`. A collector that has to look inside Go sources reads them through `sandbox.Deps.Goimportsdeps.Parse`, returning `([]map[string]any, error)` like `CollectPublicApi`.
- Page generator: a doc that grows with the project is split into one page per unit instead of one page that holds them all, so a lookup costs the unit asked about. `sandbox/internal/actions/build/generate_<x>_pages.go` renders `assets/templates/<x>_page.md` once per unit into `docs/<Doc>/<unit>.md`, and the doc's own `doc.md` becomes the index that links them. A page is an asset of the doc directory, not a sub-doc: `CollectDocTree` walks directories, so a plain `.md` beside `doc.md` is ignored by the index and by `verify`, exactly as the generated `Index.md` is. The collector names the page (`CommandDoc`'s identifier, `PublicApiPageOf`, `DepsApiPageOf`) so the index link and the written file cannot disagree, and the generator ends by calling `removeStaleDocPages`, which drops every page this build did not write. `docs/Commands/`, `docs/Routes/` and `docs/PublicApi/` work this way. A layer whose purge owns the doc directory must list it among the directories it owns whole (`cliDirs`, `serverDirs`): the asset group installs only `doc.md` and `props.yaml`, so removing those two alone leaves a directory of pages with no `props.yaml`, which every later build reads as a doc that fails to load.
- Bootstrap twice; the second run must change nothing.

## Add an example to agnos

Examples are declared with the bootstrap binary and run against this tree, never against an installed `agnos`:

```bash
./release/bootstrap.bin add-cli-example <name>
./release/bootstrap.bin add-lib-example <name>
./release/bootstrap.bin update-test <name>
```

`exec-test` writes `release/exec-test/agnos` — `exec go run <repo>/cmd/main "$@"` — and puts it in front of the PATH, so an `example.sh` typing `agnos` runs this source tree. An example that reaches the go runtime is slow (`go mod tidy` + `go build` per run); `--only` narrows the suite. `examples/lib/*/example.go` is `package main` inside the module but outside the compile scope, so it is checked by `exec-test` alone, never by `build`.

What each example asserts is the set it copies into `AssertDir`, and only `start` copies the whole tree: a change to a `start` template must move that one golden and no other. `exec-test --update` rewrites all of them at once and is for a shape change alone.

## Add a parsable

`sandbox/internal/parsables/<name>conf/`: `api.go` (struct: data fields, func fields, `Render` last), `new.go` (`New(sandbox, content) (*T, error)` via `sandbox.Deps.Serializables.ParseYaml`), `new_empty.go` (`NewEmpty(sandbox) *T`), `bind_methods.go` (`bindMethods(sandbox, self)`), `render.go` (`SerializeToYaml`). `Render` must round-trip through `New`. Shared loaders go in `utils/` (`LoadXConf`/`SaveXConf`).

## Docs

`add-doc` / `remove-doc` and `AgnosConfig/structure.yaml` work as in [Workflow](../Workflow/doc.md#add-a-doc), driven by `./release/bootstrap.bin`. What is specific to this repo: [Workflow](../Workflow/doc.md), [Rules](../Rules/doc.md), [Structure](../Structure/doc.md), [EntriesYaml](../EntriesYaml/doc.md), [DepList](../DepList/doc.md), [GeneratedFiles](../GeneratedFiles/doc.md), [LibUsage](../LibUsage/doc.md), [LibExamples](../LibExamples/doc.md) and [PublicApi](../PublicApi/doc.md) are rendered into *every* agnos project from `assets/doc/docs/`, and [CliInstall](../CliInstall/doc.md) and [Commands](../Commands/doc.md) into every one with a cli from `assets/doc-cli/docs/` ([CliExamples](../CliExamples/doc.md) from `assets/doc-example-cli/docs/`), so a change to one of them is a change to that template — and must read correctly in a scaffolded project, not only here. Guard a line that only holds for this repo with `{{ if .HasAssets }}`.

Two different names live in these templates and must never be swapped: `{{.GeneratorName}}` is the cli running the build (`agnos`) and prefixes every command agnos owns — `agnos build`, `agnos add-command`, `agnos add-route`, `agnos exec-test`; `{{.Name}}` is the project being generated and prefixes only the commands that project answers itself — `<name> help`, `<name> version`, `<name> start-server`, plus any command its own `add-command` declared. Never hardcode `agnos` in a template and never reach for `{{.Name}}` to spell an agnos command: in this repo both render `agnos`, so the mistake is invisible here and only shows up in a scaffolded project. The same pair exists for the version: `{{.GeneratorVersion}}` is the release of the binary running the build — the floor [Requirements](../Requirements/doc.md) names, since the tree in front of the reader was rendered by it — and `{{.Version}}` is what the generated project releases under.

A doc an LLM re-reads on every task is split, not grown: [Commands](../Commands/doc.md), [PublicApi](../PublicApi/doc.md) and the server layer's `docs/Routes/` are one page per unit behind an index, written by the page generators above. A `doc.md` that would list every unit of a growing set is the shape to avoid.

A new rule goes in `assets/doc/docs/Rules/doc.md`, never in the page that happens to touch it. A pattern changed in `CLAUDE.md` is mirrored here in the same commit, and the reverse.

## Write-once files

Some files are written once and never rewritten, each by a different mechanic — copy the nearest one rather than inventing a fourth:

| File | Written by | Guard |
|---|---|---|
| `sandbox/constructors/<x>/constructor.go` | `build`, `generate_constructors.go` | `io.IsFile(dest)` |
| `sandbox/internal/server/handle_*.go` | `build`, `generate_error_handlers.go` | `io.IsFile(dest)` |
| `sandbox/internal/commands/start_server/*` | `server-init` | `io.IsDir(dir)` |
| `routes/<name>/{route.yaml,handler.go}` | `add-route` | `io.WriteFile`, which refuses an existing path |

A write-once file may **not** live in an asset group: `utils.RenderGroupExcept` writes every file of a group with `WriteFileOverwrite` on every build, so a group is the one place it cannot go. Put its template in `assets/templates/` and render it by name. Writing it from `build` rather than from an `<x>-init` is what carries a project that gained the layer before the file existed.
