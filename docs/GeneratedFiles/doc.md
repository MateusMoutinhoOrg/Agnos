# GeneratedFiles

`once` = written the first time, then yours to edit. `always` = rewritten by every
`agnos build`, so an edit to it is lost — change the declaration it is rendered from instead.

| File | Written by | Rewrite |
|---|---|---|
| `AgnosConfig/{project,themes,structure,ignore,paths}.yaml` | `start` | once |
| `AgnosConfig/docs/ReadmeHeader.md` | `start` | once. The whole of `README.md` above the doc index, itself a template |
| `go.mod` | `start` | once. `dep-install` / `dep-remove` edit `require` |
| `go.sum` | `go mod tidy` | - |
| `LICENSE` | `start` | once. A placeholder; its text is pasted into `README.md`'s License section |
| `README.md` | `build` | always. `ReadmeHeader.md` + one index section per theme of `themes.yaml` |
| `sandbox/new.go` | `build` | always. One `binds.<X>Bind` per file of `sandbox/binds/` |
| `sandbox/api/sandbox.go` | `build` | always. One field per other file of `sandbox/api/` |
| `sandbox/internal/config/config.go` | `build` | always. `ProjectName`, `Version` from `project.yaml` |
| `docs/{Requirements,Workflow,Rules,Structure,EntriesYaml,DepList,GeneratedFiles,LibUsage,PublicApi,Commands}/` | `build` | always. Both `doc.md` and `props.yaml` |
| `docs/**/Index.md` | `build` | always, for every doc that has sub-docs |
| `sandbox/deps/deps.go` | `build` | always. One `<Title> <dir>.Lib` per dir of `sandbox/deps/` |
| `adapters/availables/standard/new.go` | `build` | always. One `<lib>.Bind(&deps)` per dir of `adapters/libs/` |
| `sandbox/deps/<dep>/*.go`, `adapters/libs/<lib>/*.go` | `dep-install` | once |
| `assets/asset.go` | `dep-install embeddeps` | once |
| `cmd/main/main.go` | `build` | always |
| `docs/CliInstall/` | `build` | always. Both `doc.md` and `props.yaml` |
| `sandbox/api/cli.go`, `sandbox/binds/cli.go` | `build` | always |
| `sandbox/internal/cli/climain.go` | `build` | always. `CliMain` + one `dispatch<Name>` per command |
| `sandbox/internal/commands/help/{entries.yaml,handler.go}` | `build` | always |
| `sandbox/internal/commands/version/{entries.yaml,handler.go}` | `build` | always |
| `sandbox/internal/commands/<name>/entries.go` | `build` | always. The `Entries` struct of that command |
| `sandbox/internal/commands/<name>/entries.yaml` | `add-command` | once, then rewritten by `add-flag` / `add-arg` / `set-command` — never by hand |
| `sandbox/internal/commands/<name>/handler.go` | `add-command` | once. A stub; the command's whole hand-written half |
| `docs/<Name>/{props.yaml,doc.md}` | `add-doc` | once |

Everything not listed is yours: `sandbox/internal/<pkg>/`, the contracts under `sandbox/api/`
and `sandbox/deps/` that you write, their `sandbox/binds/` and `adapters/libs/` halves, and any
directory of `adapters/availables/` other than `standard`.
