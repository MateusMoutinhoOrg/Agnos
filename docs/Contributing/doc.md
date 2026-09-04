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
4. Comment the new field: `verify` requires a doc comment on every exported declaration of `sandbox/api/`, and its row in [PublicApi](../PublicApi/doc.md) is generated from that comment.

## Add a command to agnos

```bash
./release/bootstrap.bin add-command <name> --help "..." --category "Core Commands"   # categories: Core Commands, Cli System, Dependencies, Dependency System, Info
./release/bootstrap.bin add-flag path --command <name> --default . --description "..."
./release/bootstrap.bin add-flag quiet --command <name> --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
```

Write `handler.go`: call the action, return `api.ExitFailure` on error, `Printf` any result. Add a row to [Commands](../Commands/doc.md) — `verify` requires that page to name every visible command and every flag it declares.

## Add a contract + adapter lib

1. `sandbox/deps/<x>/<x>.go`: `type Lib struct { ... func fields }`, stdlib imports only. Per-call objects use a constructor field (`New func(...) Thing`) like `argvdeps`/`requestdeps`. Doc-comment every exported declaration (`verify` requires it) and every field you want described in [PublicApi](../PublicApi/doc.md) — that page is rendered from these comments.
2. `adapters/libs/<x>/<x>.go` (file named after the package): `func Bind(deps *deps.Deps) { deps.<X> = <x>.Lib{...} }`. Any import allowed. `verify` requires the binder and requires some lib to fill every `deps.<X>`.
3. Bootstrap: `deps.go` and `standard/new.go` pick it up from the dir listing. Reach it as `deps.<X>`.

## Add an installable dep

1. Mirror the two files under `assets/deplist/<contract>/sandbox/deps/<contract>/` and `assets/deplist/<contract>/adapters/libs/<lib>/`, replacing this module path with `{{.Module}}`. `verify` requires the copies to render byte-for-byte to the repo's, so re-mirror on every change to either side.
2. If it pulls a module: `assets/depsversion.yaml` `<dep>: <module>@<version>`.
3. Bootstrap, test with `dep-install`/`dep-remove` on a scratch project. Add a row to [DepList](../DepList/doc.md).

## Add a template or collector

- Template: `assets/<group>/<target path>`, a `text/template` over the vars in [BuildPipeline](../BuildPipeline/doc.md#buildinternal). Groups: `all`, `deps`, `cli`, `start`. Single-destination scaffolds go in `assets/templates/` and are rendered with `utils.RenderTemplateToDest`. Register in [GeneratedFiles](../GeneratedFiles/doc.md).
- Generated `.go` is formatted by `deps.Goimportsdeps.Format` at write time, so a template's own indentation only has to be parsable.
- Collector: `sandbox/internal/actions/build/collect_<x>.go`, `func Collect<X>(io) []string` listing one dir and title-casing the last segment; add `"<X>": Collect<X>(io)` to the vars map in `build_internal.go`. A collector that has to look inside Go sources takes `deps` too and reads them through `deps.Goimportsdeps.Parse`, returning `([]map[string]any, error)` like `CollectPublicApi`.
- Bootstrap twice; the second run must change nothing.

## Add a parsable

`sandbox/internal/parsables/<name>conf/`: `api.go` (struct: data fields, func fields, `Render` last), `new.go` (`New(deps, content) (*T, error)` via `deps.Serializables.ParseYaml`), `new_empty.go` (`NewEmpty(deps) *T`), `bind_methods.go` (`bindMethods(deps, self)`), `render.go` (`SerializeToYaml`). `Render` must round-trip through `New`. Shared loaders go in `utils/` (`LoadXConf`/`SaveXConf`).

## Docs

- A doc is `docs/<Name>/{doc.md,props.yaml}`; sub-docs nest as `docs/<Name>/<Sub>/`. Other files in a doc dir are assets.
- Create/delete with `agnos add-doc <Name> --description "..." --theme <id>` / `agnos remove-doc <Name>`; a theme is an entry in `AgnosConfig/themes.yaml` (`cli-usage`, `lib-usage`, `development`) and only groups the doc into a section of `README.md`. `README.md` and every `Index.md` are generated on build; never edit them.
- `docs/LibUsage/` and `docs/PublicApi/` are generated on every build (`assets/all/docs/`); PublicApi's content is the doc comments of `sandbox/api/` and `sandbox/deps/`, so it is changed by editing those comments, never the page.
- Keep docs short: tables and commands over prose. One page per topic; no sub-doc unless the content is a real list of independent items.
- Links are relative to the file that carries them: `../X/doc.md` inside `docs/`, `docs/X/doc.md` in `README.md` / `ReadmeHeader.md`.
- A pattern changed in `CLAUDE.md` is mirrored here in the same commit, and the reverse.
