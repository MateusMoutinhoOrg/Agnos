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

1. `sandbox/internal/actions/<name>/<name>_internal.go`: `func <Name>Internal(deps, io *smartio.SmartIO, ...) error`. Project-relative paths only. Log via `deps.Std.Log`, fail via `deps.Std.Errorf`, never `Printf`.
2. `<name>.go`: `func <Name>(deps, ...) error` = `smartio.New(deps, path, config.ProjectName)` -> internal -> `io.Persist()` -> `buildAction.Build(deps, api.BuildProps{Path, Runtime})` (`RuntimeGo` if it adds, `RuntimeNone` if it removes). Props with more than three values go in a struct in `sandbox/api/actions.go`.
3. Add the field to `api.Actions` and the assignment to `sandbox/binds/actions.go`.
4. Comment the new field: its row in [PublicApi](../PublicApi/doc.md) is generated from that comment.

## Add a command to agnos

Declare it with the bootstrap binary, as in [Workflow](../Workflow/doc.md#change-the-command-surface), with a `--category` this repo already uses (Core Commands, Cli System, Dependencies, Dependency System, Info) and the two flags every agnos command carries:

```bash
./release/bootstrap.bin add-command <name> --help "..." --category "Core Commands"
./release/bootstrap.bin add-flag path --command <name> --default . --description "the dir holding the project (defaults to the current directory)"
./release/bootstrap.bin add-flag quiet --command <name> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
```

`handler.go` calls the action, returns `api.ExitFailure` on error and `Printf`s any result.

## Add a contract + adapter lib

The two halves are in [Workflow](../Workflow/doc.md#add-a-dependency). Per-call objects use a constructor field (`New func(...) Thing`) like `argvdeps`/`requestdeps`. Doc-comment every exported declaration and every field you want described in [PublicApi](../PublicApi/doc.md). To make it installable elsewhere, mirror it into `assets/deplist/` below.

## Add an installable dep

1. Mirror the two files under `assets/deplist/<contract>/sandbox/deps/<contract>/` and `assets/deplist/<contract>/adapters/libs/<lib>/`, replacing this module path with `{{.Module}}`.
2. If it pulls a module: `assets/depsversion.yaml` `<dep>: <module>@<version>`.
3. Bootstrap, test with `dep-install`/`dep-remove` on a scratch project. Add a row to [DepList](../DepList/doc.md).

## Add a template or collector

- Template: `assets/<group>/<target path>`, a `text/template` over the vars in [BuildPipeline](../BuildPipeline/doc.md#buildinternal). Groups: `all`, `deps`, `cli`, `start`. Single-destination scaffolds go in `assets/templates/` and are rendered with `utils.RenderTemplateToDest`. Add a row for the new destination to `assets/all/docs/GeneratedFiles/doc.md`.
- Collector: `sandbox/internal/actions/build/collect_<x>.go`, `func Collect<X>(io) []string` listing one dir and title-casing the last segment; add `"<X>": Collect<X>(io)` to the vars map in `build_internal.go`. A collector that has to look inside Go sources takes `deps` too and reads them through `deps.Goimportsdeps.Parse`, returning `([]map[string]any, error)` like `CollectPublicApi`.
- Bootstrap twice; the second run must change nothing.

## Add an example to agnos

Examples are declared with the bootstrap binary and run against this tree, never against an installed `agnos`:

```bash
./release/bootstrap.bin add-cli-example <name>
./release/bootstrap.bin add-lib-example <name>
./release/bootstrap.bin exec-test --only <name> --update
```

`exec-test` writes `release/exec-test/agnos` — `exec go run <repo>/cmd/main "$@"` — and puts it in front of the PATH, so an `example.sh` typing `agnos` runs this source tree. An example that reaches the go runtime is slow (`go mod tidy` + `go build` per run); `--only` narrows the suite. `examples/lib/*/example.go` is `package main` inside the module but outside the compile scope, so it is checked by `exec-test` alone, never by `build`.

## Add a parsable

`sandbox/internal/parsables/<name>conf/`: `api.go` (struct: data fields, func fields, `Render` last), `new.go` (`New(deps, content) (*T, error)` via `deps.Serializables.ParseYaml`), `new_empty.go` (`NewEmpty(deps) *T`), `bind_methods.go` (`bindMethods(deps, self)`), `render.go` (`SerializeToYaml`). `Render` must round-trip through `New`. Shared loaders go in `utils/` (`LoadXConf`/`SaveXConf`).

## Docs

`add-doc` / `remove-doc` and `AgnosConfig/structure.yaml` work as in [Workflow](../Workflow/doc.md#add-a-doc), driven by `./release/bootstrap.bin`. What is specific to this repo: [Workflow](../Workflow/doc.md), [Rules](../Rules/doc.md), [Structure](../Structure/doc.md), [EntriesYaml](../EntriesYaml/doc.md), [DepList](../DepList/doc.md), [GeneratedFiles](../GeneratedFiles/doc.md), [Commands](../Commands/doc.md), [LibUsage](../LibUsage/doc.md), [LibExamples](../LibExamples/doc.md) and [PublicApi](../PublicApi/doc.md) are rendered into *every* agnos project from `assets/all/docs/`, and [CliInstall](../CliInstall/doc.md) and [CliExamples](../CliExamples/doc.md) into every one with a cli from `assets/cli/docs/`, so a change to one of them is a change to that template — and must read correctly in a scaffolded project, not only here. Guard a line that only holds for this repo with `{{ if .HasAssets }}`.

A new rule goes in `assets/all/docs/Rules/doc.md`, never in the page that happens to touch it. A pattern changed in `CLAUDE.md` is mirrored here in the same commit, and the reverse.
