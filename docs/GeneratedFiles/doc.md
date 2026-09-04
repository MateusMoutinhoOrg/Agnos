# GeneratedFiles

`once` = written on first render, then yours. `always` = rewritten by every `build`, never edit.

| File | Written by | Template | Rewrite |
|---|---|---|---|
| `AgnosConfig/{project,themes,ignore,paths}.yaml` | `start` | `start/` | once |
| `AgnosConfig/docs/ReadmeHeader.md` | `start` | `start/` | once. The README body (itself a template) |
| `go.mod` | `start` | - | once. `dep-install/remove` edit `require` |
| `go.sum` | `go mod tidy` (runtime) | - | - |
| `README.md` | `build` | `all/` | always. `render ReadmeHeader.md` + theme index + `copy LICENSE` |
| `sandbox/new.go` | `build` | `all/` | always. One `binds.<X>Bind` per file of `sandbox/binds/` |
| `sandbox/api/sandbox.go` | `build` | `all/` | always. One field per file of `sandbox/api/` |
| `sandbox/internal/config/config.go` | `build` | `all/` | always. `ProjectName`, `Version` from `project.yaml` |
| `docs/LibUsage/{doc.md,props.yaml}` | `build` | `all/` | always. Generic "use this project as a Go module" doc |
| `docs/PublicApi/{doc.md,props.yaml}` | `build` | `all/` | always. Every exported symbol of `sandbox/api/` and `sandbox/deps/`, described by its own doc comment |
| `docs/Index/<theme-id>.md` | `build` (if `docs/` exists) | `templates/theme_index.md` | always. `docs/Index/` is removed and rewritten whole |
| `docs/**/Index.md` | `build` (docs with sub-docs) | `templates/doc_index.md` | always |
| `sandbox/deps/deps.go` | `build` (if `sandbox/deps/`) | `deps/` | always. One `<Title> <dir>.Lib` per dir |
| `adapters/availables/standard/new.go` | `build` (if `sandbox/deps/`) | `deps/` | always. One `<lib>.Bind(&deps)` per dir of `adapters/libs/` |
| `sandbox/deps/<dep>/*.go`, `adapters/libs/<lib>/*.go` | `dep-install` | `deplist/<dep>/` | once |
| `assets/asset.go` | `dep-install embeddeps` | `deplist/embeddeps/` | once |
| `cmd/main/main.go` | `build` (if `sandbox/internal/cli/`) | `cli/` | always |
| `sandbox/api/cli.go`, `sandbox/binds/cli.go` | `build` (cli) | `cli/` | always |
| `sandbox/internal/cli/climain.go` | `build` (cli) | `cli/` | always. `CliMain` + one `dispatch<Name>` per command |
| `sandbox/internal/commands/help/entries.yaml` | `build` (cli) | `templates/help_entries.yaml` | once |
| `sandbox/internal/commands/help/handler.go` | `build` (cli) | `cli/` | always |
| `sandbox/internal/commands/version/{entries.yaml,handler.go}` | `build` (cli) | `cli/` | always |
| `sandbox/internal/commands/<name>/entries.go` | `build` (cli) | `templates/entries.go` | always |
| `sandbox/internal/commands/<name>/{entries.yaml,handler.go}` | `add-command` | `templates/command_*` | once. `entries.yaml` is rewritten by the field editors |
| `docs/<name>/{props.yaml,doc.md}` | `add-doc` | `templates/doc_doc.md` | once |
