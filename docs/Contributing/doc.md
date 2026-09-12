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

Release: bump `version` in `AgnosConfig/project.yaml`, then `agnos publish` (or `agnos compile --target all` for the binaries alone).

## Add an action

1. `sandbox/internal/actions/<name>/<name>_internal.go`: `func <Name>Internal(sandbox, io *smartio.SmartIO, ...) error`. Project-relative paths only. Log via `sandbox.Deps.Std.Log`, fail via `sandbox.Deps.Std.Errorf`, never `Printf`.
2. `<name>.go`: `func <Name>(sandbox, ...) error` = `smartio.New(sandbox, path, config.ProjectName)` -> internal -> `io.Persist()` -> `buildAction.Build(sandbox, api.BuildProps{Path, Runtime})` (`RuntimeGo` if it adds, `RuntimeNone` if it removes). Props with more than three values go in a struct in `sandbox/api/actions.go`.
3. Add the field to `api.Actions` and the assignment to `sandbox/binds/actions.go`.
4. Comment the new field: its row in [PublicApi](../PublicApi/doc.md) is generated from that comment.

## Add a command to agnos

Declare it with the bootstrap binary, as in [Workflow](../Workflow/doc.md#change-the-command-surface), with a `--category` this repo already uses (Core Commands, Cli System, Server System, Front System, Dependencies, Dependency System, Info) and the two flags every agnos command carries:

```bash
./release/bootstrap.bin add-command <name> --help "..." --category "Core Commands"
./release/bootstrap.bin add-flag path --command <name> --default . --description "the dir holding the project (defaults to the current directory)"
./release/bootstrap.bin add-flag quiet --command <name> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
```

`handler.go` calls the action, returns `api.ExitFailure` on error and `Printf`s any result.

`add-flag`/`add-arg` read the command line they are typed on, so never pass a value that is *exactly* one of their own flag spellings (`--identifier --example`, `--example --required`): the parser counts it as an occurrence and the declaration comes out polluted. Declare such a flag without the short alias instead, and word its `--example` as a whole command line.

## Add a layer (cli, server, …)

A layer is an asset group plus an `<x>-init`/`<x>-purge` pair, and the server layer is the pattern to copy — file for file, it mirrors the cli one. The front layer is the third column: it declares no unit of its own, because a page **is** a route.

| Concept | CLI | Server | Front |
|---|---|---|---|
| External input contract | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` | — (`embeddeps` + `templatedeps`) |
| Surface + bind | `sandbox/api/cli.go`, `sandbox/binds/cli.go` | `sandbox/api/server.go`, `sandbox/binds/server.go` | — (served through the server's) |
| Generated dispatch | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/servermain.go` | — |
| Shared package | — | `sandbox/internal/routeio/` | `sandbox/internal/pageio/` |
| Declared unit | `commands/<name>/entries.yaml` | `routes/<name>/route.yaml` | `routes/<page>/route.yaml` + `assets/frontend/pages/<page>.html` |
| Parsable | `parsables/commandconf/` | `parsables/routeconf/` | — (`routeconf`) |
| Collectors | `collect_commands.go`, `collect_command_docs.go` | `collect_routes.go`, `collect_route_docs.go` | `collect_front_mount.go` |
| Per-unit generator | `generate_command_entries.go` | `generate_route_entries.go` | — (`generate_route_entries.go`) |
| Asset group | `assets/cli/` | `assets/server/` | `assets/front/` |
| Build trigger | `hasCli := io.IsDir("sandbox/internal/cli")` | `hasServer := io.IsDir("sandbox/internal/server")` | `hasFront := io.IsDir("sandbox/internal/pageio")` |
| Init / purge | `cli-init` / `cli-purge` | `server-init` / `server-purge` | `front-init` / `front-purge` |
| Verify | `check_sandbox.go` et al | `check_routes.go` | — (`check_routes.go`) |

Both halves of the pair go in `GeneratedDocsGroups`, in the vars map of `build_internal.go`, and in the `{{ if .Has<X> }}` blocks of `assets/all/docs/{Rules,GeneratedFiles,Workflow}/doc.md`. A layer whose init needs another layer calls the other one's `<X>InitInternal` on the *same* open SmartIO — `server_init` does that with `cli_init`, `front_init` with `server_init` — so there is no intermediate `Persist` and no intermediate `build`. The dep installs are the exception: `<X>InitInternal` renders assets and writes nothing to `go.mod`, so a composing init calls the other's exported `InstallDeps` first.

The build trigger is never the content tree: `assets/frontend/` is the project's own and may legitimately be empty, so `hasFront` reads `pageio/`. `hasAssets` stays `assets/all`, so a front project does not read as a generator.

`sandbox/internal/routeio/` exists because `servermain.go` imports every route, so a route cannot import `internal/server` back: shared route code goes in a third package both may import. `sandbox/internal/pageio/` is the same shape one layer up.

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

- Template: `assets/<group>/<target path>`, a `text/template` over the vars in [BuildPipeline](../BuildPipeline/doc.md#buildinternal). Groups: `all`, `deps`, `cli`, `server`, `front`, `start`. A scaffold that renders to a file which is itself a template — `page_html.html`, or any `{{` inside `front_main.js` — escapes its own braces (`{{ "{{ .Title }}" }}`), or the outer render eats them. Single-destination scaffolds go in `assets/templates/` and are rendered with `utils.RenderTemplateToDest`. Add a row for the new destination to `assets/all/docs/GeneratedFiles/doc.md`.
- Collector: `sandbox/internal/actions/build/collect_<x>.go`, `func Collect<X>(sandbox, io) []string` listing one dir and title-casing the last segment; add `"<X>": Collect<X>(sandbox, io)` to the vars map in `build_internal.go`. A collector that has to look inside Go sources reads them through `sandbox.Deps.Goimportsdeps.Parse`, returning `([]map[string]any, error)` like `CollectPublicApi`.
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

`add-doc` / `remove-doc` and `AgnosConfig/structure.yaml` work as in [Workflow](../Workflow/doc.md#add-a-doc), driven by `./release/bootstrap.bin`. What is specific to this repo: [Workflow](../Workflow/doc.md), [Rules](../Rules/doc.md), [Structure](../Structure/doc.md), [EntriesYaml](../EntriesYaml/doc.md), [DepList](../DepList/doc.md), [GeneratedFiles](../GeneratedFiles/doc.md), [Commands](../Commands/doc.md), [LibUsage](../LibUsage/doc.md), [LibExamples](../LibExamples/doc.md) and [PublicApi](../PublicApi/doc.md) are rendered into *every* agnos project from `assets/all/docs/`, and [CliInstall](../CliInstall/doc.md) and [CliExamples](../CliExamples/doc.md) into every one with a cli from `assets/cli/docs/`, so a change to one of them is a change to that template — and must read correctly in a scaffolded project, not only here. Guard a line that only holds for this repo with `{{ if .HasAssets }}`.

Two different names live in these templates and must never be swapped: `{{.GeneratorName}}` is the cli running the build (`agnos`) and prefixes every command agnos owns — `agnos build`, `agnos add-command`, `agnos add-route`, `agnos exec-test`; `{{.Name}}` is the project being generated and prefixes only the commands that project answers itself — `<name> help`, `<name> version`, `<name> start-server`, plus any command its own `add-command` declared. Never hardcode `agnos` in a template and never reach for `{{.Name}}` to spell an agnos command: in this repo both render `agnos`, so the mistake is invisible here and only shows up in a scaffolded project.

A new rule goes in `assets/all/docs/Rules/doc.md`, never in the page that happens to touch it. A pattern changed in `CLAUDE.md` is mirrored here in the same commit, and the reverse.
