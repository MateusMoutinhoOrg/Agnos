# GeneratedFiles

`once` = written the first time, then yours to edit. `always` = rewritten by every
`agnos build`, so an edit to it is lost — change the declaration it is rendered from instead.

| File | Written by | Rewrite |
|---|---|---|
| `AgnosConfig/{project,themes,structure,paths}.yaml` | `start` | once |
| `AgnosConfig/extensions.yaml` | `start` | once, then rewritten by `enable-extension` / `disable-extension` and every `<x>-init` / `<x>-purge` — never by hand |
| `AgnosConfig/docs/ReadmeHeader.md` | `start` | once. The whole of `README.md` above the doc index, itself a template |
| `go.mod` | `start` | once. `add-dep` / `remove-dep` edit `require` |
| `go.sum` | `go mod tidy` | - |
| `LICENSE` | `start` | once. A placeholder; its text is pasted into `README.md`'s License section |
| `README.md` | `build` | always. `ReadmeHeader.md` + one index section per theme of `themes.yaml` |
| `sandbox/new.go` | `build` | always. One `<x>.Constructor(&self)` per directory of `sandbox/constructors/` |
| `sandbox/constructors/<x>/constructor.go` | `build` | once, per contract of `sandbox/api/` that has a `sandbox/internal/<x>/new.go`. Then yours — write your own package there and `new.go` calls it too |
| `sandbox/api/sandbox.go` | `build` | always. Every struct of `sandbox/api/<x>sandbox.go` embedded, plus `Config` and `Deps` while the project carries the deps layer |
| `sandbox/api/config.go` | `build` | always. The `Config` contract: every struct of `sandbox/api/<x>config.go` embedded, `ProjectName`, `Version` |
| `sandbox/api/projectsandbox.go` | `start` | once. `api.ProjectSandbox`, embedded in `api.Sandbox` — declare the project's own fields of the sandbox there |
| `sandbox/api/projectconfig.go` | `start` | once. `api.ProjectConfig`, embedded in `api.Config` — declare the project's own config fields there |
| `sandbox/internal/generated/config/new.go` | `build` | always. `NewConfig`, filled with `ProjectName` and `Version` from `project.yaml` |
| `docs/{Requirements,Workflow,Rules,Extensions,Structure,DepList,GeneratedFiles,LibUsage,PublicApi}/` | `build` | always. Both `doc.md` and `doc.yaml` |
| `docs/**/Index.md` | `build` | always, for every doc that has sub-docs |
| `docs/PublicApi/<contract>.md` | `build` | always. One page per file of `sandbox/api/` and per contract of `sandbox/deps/`; `docs/PublicApi/doc.md` indexes them by the symbols each declares |
| `docs/LibExamples/` | `build` | always. Both `doc.md` and `doc.yaml` |
| `sandbox/deps/deps.go` | `build` | always. One `<Field> <dir>.Contract` per dir of `sandbox/deps/` (`<dir>.Sandbox` for a remote dep, which keeps the name of the api it copies) |
| `adapters/bindings/<name>/new.go` | `build` | always. One `<adapter>.Bind(&deps)` per entry of that binding's `binding.yaml`; a binding with no `binding.yaml` is hand-written and left alone |
| `adapters/bindings/<name>/binding.yaml` | `deps-init` | once, then rewritten by `add-dep` / `remove-dep` — never by hand |
| `sandbox/deps/<dep>/*.go`, `adapters/impls/<adapter>/*.go` | `add-dep` | once |
| `adapters/impls/<adapter>/adapter.yaml` | `add-dep` | once |
| `sandbox/deps/<dep>/*.go` of a remote dep | `add-dep <module>` | rewritten by `set-dep`; a copy of that module's `sandbox/api/` |
| `adapters/impls/<dep>/<dep>.go` of a remote dep | `add-dep <module>` | rewritten by `set-dep`; the generated shim |
| `assets/asset.go` | `add-dep embeddeps` | once |
| `cmd/main/main.go` | `build` | always |
| `docs/{CliInstall,Commands}/` | `build` | always. Both `doc.md` and `doc.yaml` |
| `docs/Commands/<command>.md` | `build` | always. One page per visible command; `docs/Commands/doc.md` indexes them |
| `docs/CliExamples/` | `build` | always. Both `doc.md` and `doc.yaml` |
| `sandbox/api/{cli,command,trigger}.go` | `build` | always. Aliases of the `OpinionatedAgnosCli` contract's types |
| `sandbox/api/clisandbox.go` | `build` | always. `api.CliSandbox`, the part of `api.Sandbox` holding `Cli` |
| `sandbox/internal/generated/cli/new.go` | `build` | always. `NewCli` builds `Cli.Commands` from every command's `NewCommand`, and `Cli.Main`, which hands the line to `OpinionatedAgnosCli.Main` |
| `sandbox/deps/OpinionatedAgnosCli/`, `adapters/impls/OpinionatedAgnosCli/` | `cli-init` | once, like any dep. The dispatch, the binder, the matcher and the failures every command runs through |
| `sandbox/internal/commands/{info/help,info/version,middleware/help_flag}/{command.yaml,handler.go}` | `build` | always |
| `sandbox/internal/commands/<name>/new.go` | `build` | always. `NewCommand`, that command's `api.Command`, a 1:1 image of `command.yaml` |
| `sandbox/internal/commands/<name>/input.go` | `build` | always. `Input`, one field per arg and flag |
| `sandbox/internal/commands/<name>/command.yaml` | `add-command` | once, then rewritten by `add-flag` / `add-arg` / `set-command`, their `set-` editors and their inverses — never by hand |
| `sandbox/internal/commands/<name>/handler.go` | `add-command` | once. A stub; the command's whole hand-written half |
| `sandbox/internal/cli/errors/handle_*.go` | `build` | once. Five files, one per failure — what this project answers when no command does |
| `sandbox/internal/commandprops/commandprops.go` | `build` | always. `commandprops.CommandProps`, what one command line's chain of commands shares: every struct of the package embedded |
| `sandbox/internal/commandprops/project.go` | `build` | once, while the package has no other part. `Project`, the project's own fields of `CommandProps` |
| `docs/<Name>/{doc.yaml,doc.md}` | `add-doc` | once |
| `examples/cli/<name>/example.sh` | `add-cli-example` | once. A stub that already runs |
| `examples/lib/<name>/example.go` | `add-lib-example` | once. A stub that already runs |
| `examples/<side>/<name>/result.yaml` | `run-examples` | on `update-example <name>`, on `--update` or when absent — never by hand |

Everything under `sandbox/internal/generated/` is `always`. Everything not listed is yours:
`sandbox/internal/<pkg>/`, the contracts under `sandbox/api/`
and `sandbox/deps/` that you write, their `sandbox/internal/<x>/new.go` and `adapters/impls/`
halves, and any
directory of `adapters/bindings/` other than `standard`.
