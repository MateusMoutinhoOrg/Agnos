# GeneratedFiles

`once` = written the first time, then yours to edit. `always` = rewritten by every
`{{.GeneratorName}} build`, so an edit to it is lost — change the declaration it is rendered from instead.

| File | Written by | Rewrite |
|---|---|---|
| `{{.ConfigDir}}/{project,themes,structure,paths}.yaml` | `start` | once |
| `{{.ConfigDir}}/extensions.yaml` | `start` | once, then rewritten by `enable-extension` / `disable-extension` and every `<x>-init` / `<x>-purge` — never by hand |
| `{{.ConfigDir}}/docs/ReadmeHeader.md` | `start` | once. The whole of `README.md` above the doc index, itself a template |
| `go.mod` | `start` | once. `add-dep` / `remove-dep` edit `require` |
| `go.sum` | `go mod tidy` | - |
| `LICENSE` | `start` | once. A placeholder; its text is pasted into `README.md`'s License section |
{{- if .HasReadme }}
| `README.md` | `build` | always. `ReadmeHeader.md` + one index section per theme of `themes.yaml` |
{{- end }}
{{- if .HasSandbox }}
| `sandbox/new.go` | `build` | always. One `<x>.Constructor(&self)` per directory of `sandbox/constructors/` |
| `sandbox/constructors/<x>/constructor.go` | `build` | once, per contract of `sandbox/api/` that has a `sandbox/internal/<x>/new.go`. Then yours — write your own package there and `new.go` calls it too |
| `sandbox/api/sandbox.go` | `build` | always. One field per other file of `sandbox/api/`, plus `Deps` while the project carries the deps layer |
| `sandbox/internal/config/config.go` | `build` | always. `ProjectName`, `Version` from `project.yaml` |
{{- end }}
{{- if .HasDoc }}
| `docs/{Requirements,Workflow,Rules,Extensions,Structure,EntriesYaml,DepList,GeneratedFiles,LibUsage,PublicApi}/` | `build` | always. Both `doc.md` and `props.yaml` |
| `docs/**/Index.md` | `build` | always, for every doc that has sub-docs |
{{- end }}
{{- if and .HasDoc .HasExample }}
| `docs/LibExamples/` | `build` | always. Both `doc.md` and `props.yaml` |
{{- end }}
{{- if .HasDeps }}
| `sandbox/deps/deps.go` | `build` | always. One `<Title> <dir>.Sandbox` per dir of `sandbox/deps/` |
| `adapters/availables/<name>/new.go` | `build` | always. One `<adapter>.Bind(&deps)` per entry of that available's `available.yaml`; an available with no `available.yaml` is hand-written and left alone |
{{- end }}
| `adapters/availables/<name>/available.yaml` | `deps-init` | once, then rewritten by `add-dep` / `remove-dep` — never by hand |
| `sandbox/deps/<dep>/*.go`, `adapters/libs/<adapter>/*.go` | `add-dep` | once |
| `adapters/libs/<adapter>/adapter.yaml` | `add-dep` | once |
| `sandbox/deps/<dep>/*.go` of a remote dep | `add-dep <module>` | rewritten by `set-dep`; a copy of that module's `sandbox/api/` |
| `adapters/libs/<dep>/<dep>.go` of a remote dep | `add-dep <module>` | rewritten by `set-dep`; the generated shim |
| `assets/asset.go` | `add-dep embeddeps` | once |
{{- if .HasCli }}
| `cmd/main/main.go` | `build` | always |
{{- if .HasDoc }}
| `docs/{CliInstall,Commands}/` | `build` | always. Both `doc.md` and `props.yaml` |
{{- end }}
{{- if and .HasDoc .HasExample }}
| `docs/CliExamples/` | `build` | always. Both `doc.md` and `props.yaml` |
{{- end }}
| `sandbox/api/cli.go`, `sandbox/api/command.go` | `build` | always |
| `sandbox/internal/cli/new.go` | `build` | always. `NewCli` builds `Cli.Commands` from every command's `NewCommand` |
| `sandbox/internal/cli/climain.go` | `build` | always. `CliMain`, the one dispatch every command goes through |
| `sandbox/internal/commands/help/{entries.yaml,handler.go}` | `build` | always |
| `sandbox/internal/commands/version/{entries.yaml,handler.go}` | `build` | always |
| `sandbox/internal/commands/<name>/new.go` | `build` | always. `NewCommand`, that command's `api.Command` |
| `sandbox/internal/commands/<name>/entries.yaml` | `add-command` | once, then rewritten by `add-flag` / `add-arg` / `set-command` — never by hand |
| `sandbox/internal/commands/<name>/handler.go` | `add-command` | once. A stub; the command's whole hand-written half |
{{- end }}
{{- if .HasServer }}
| `sandbox/api/{server.go,route.go}` | `build` | always |
| `sandbox/internal/server/new.go` | `build` | always. `NewServer` builds `Server.Routes` from every route's `NewRoute` |
| `sandbox/internal/server/servermain.go` | `build` | always. `ServerMain` + the one dispatch that binds a request against `Server.Routes` |
| `sandbox/internal/routeio/*.go` | `build` | always |
| `sandbox/internal/routes/health/{route.yaml,handler.go}` | `build` | always |
| `sandbox/internal/routes/<name>/new.go` | `build` | always. `NewRoute`, that route's `api.Route`, and its `ReadBody` |
| `docs/{RouteYaml,Routes,ServerUsage}/` | `build` | always. Both `doc.md` and `props.yaml` |
| `sandbox/internal/routes/<name>/route.yaml` | `add-route` | once, then rewritten by `set-route` / `add-segment` / `add-header` / `add-param` / `set-body` / `add-body-field` and their inverses — never by hand |
| `sandbox/internal/routes/<name>/handler.go` | `add-route` | once. A stub; the route's whole hand-written half |
| `sandbox/internal/commands/start_server/{entries.yaml,handler.go}` | `server-init` | once |
{{- end }}
{{- if .HasFront }}
| `sandbox/internal/pageio/templates.go` | `build` | always. `Render` + the asset helpers; `StaticMount` from the `static` route's first segment |
| `docs/FrontUsage/` | `build` | always. Both `doc.md` and `props.yaml` |
| `sandbox/internal/routes/static/{route.yaml,handler.go}` | `front-init` | once. Keep `safeSegments` if you edit it |
| `assets/frontend/static/{styles/main.css,scripts/main.js}` | `front-init` | once |
| `sandbox/internal/routes/<page>/{route.yaml,handler.go}` | `add-page` | once |
| `assets/frontend/pages/<page>.html` | `add-page` | once. Kept as is by a second `add-page` |
{{- end }}
| `docs/<Name>/{props.yaml,doc.md}` | `add-doc` | once |
{{- if .HasExample }}
{{- if .HasCli }}
| `examples/cli/<name>/example.sh` | `add-cli-example` | once. A stub that already runs |
{{- end }}
| `examples/lib/<name>/example.go` | `add-lib-example` | once. A stub that already runs |
| `examples/<side>/<name>/result.yaml` | `exec-test` | on `update-test <name>`, on `--update` or when absent — never by hand |
{{- end }}

Everything not listed is yours: `sandbox/internal/<pkg>/`, the contracts under `sandbox/api/`
and `sandbox/deps/` that you write, their `sandbox/internal/<x>/new.go` and `adapters/libs/`
halves, and any
directory of `adapters/availables/` other than `standard`.
