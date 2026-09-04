# Contributing

Everything here is a pattern; new code copies an existing sibling. Shape is enforced by `verify` and by the generators, so match names exactly. Compile scope is always `./cmd/... ./sandbox/... ./adapters/...`.

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
4. Add a row to [PublicApi](../PublicApi/doc.md#actions).

## Add a command to agnos

```bash
./release/bootstrap.bin add-command <name> --help "..." --category "Core Commands"   # categories: Core Commands, Cli System, Dependencies, Dependency System, Info
./release/bootstrap.bin add-flag path --command <name> --default . --description "..."
./release/bootstrap.bin add-flag quiet --command <name> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
```

Write `handler.go`: call the action, return `api.ExitFailure` on error, `Printf` any result. Add a row to [Commands](../Commands/doc.md).

## Add a contract + adapter lib

1. `sandbox/deps/<x>/<x>.go`: `type Lib struct { ... func fields }`, stdlib imports only. Per-call objects use a constructor field (`New func(...) Thing`) like `argvdeps`/`requestdeps`.
2. `adapters/libs/<x>/<File>.go`: `func Bind(deps *deps.Deps) { deps.<X> = <x>.Lib{...} }`. Any import allowed.
3. Bootstrap: `deps.go` and `standard/new.go` pick it up from the dir listing. Reach it as `deps.<X>`.

## Add an installable dep

1. Mirror the two files under `assets/deplist/<contract>/sandbox/deps/<contract>/` and `assets/deplist/<contract>/adapters/libs/<lib>/`, replacing this module path with `{{.Module}}`. Keep the copies identical to the repo's.
2. If it pulls a module: `assets/depsversion.yaml` `<dep>: <module>@<version>`.
3. Bootstrap, test with `dep-install`/`dep-remove` on a scratch project. Add a row to [DepList](../DepList/doc.md).

## Add a template or collector

- Template: `assets/<group>/<target path>`, a `text/template` over the vars in [BuildPipeline](../BuildPipeline/doc.md#buildinternal). Groups: `all`, `deps`, `cli`, `start`. Single-destination scaffolds go in `assets/templates/` and are rendered with `utils.RenderTemplateToDest`. Register in [GeneratedFiles](../GeneratedFiles/doc.md).
- Collector: `sandbox/internal/actions/build/collect_<x>.go`, `func Collect<X>(io) []string` listing one dir and title-casing the last segment; add `"<X>": Collect<X>(io)` to the vars map in `build_internal.go`.
- Bootstrap twice; the second run must change nothing.

## Add a parsable

`sandbox/internal/parsables/<name>conf/`: `api.go` (struct: data fields, func fields, `Render` last), `new.go` (`New(deps, content) (*T, error)` via `deps.Serializables.ParseYaml`), `new_empty.go` (`NewEmpty(deps) *T`), `bind_methods.go` (`bindMethods(deps, self)`), `render.go` (`SerializeToYaml`). `Render` must round-trip through `New`. Shared loaders go in `utils/` (`LoadXConf`/`SaveXConf`).

## Docs

- A doc is `docs/<Name>/{doc.md,props.yaml}`; sub-docs nest as `docs/<Name>/<Sub>/`. Other files in a doc dir are assets.
- Create/delete with `agnos add-doc <Name> --description "..." --theme <id>` / `agnos remove-doc <Name>`; a theme is an entry in `AgnosConfig/themes.yaml` (`cli-usage`, `lib-usage`, `development`). Indexes and `README.md` are generated on build; never edit them.
- Keep docs short: tables and commands over prose. One page per topic; no sub-doc unless the content is a real list of independent items.
- Links inside `docs/` are repo-rooted (`/docs/X/doc.md`); in `README.md` / `ReadmeHeader.md` they are relative (`docs/X/doc.md`). Never link into `old/`.
- A pattern changed in `CLAUDE.md` is mirrored here in the same commit, and the reverse.
