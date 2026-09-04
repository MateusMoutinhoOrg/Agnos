# Structure

Agnos scaffolds projects of this same schema and regenerates itself with it. `(gen)` = written by `build`, never edited. A generated project has everything below except `actions/`, `parsables/`, `smartio/`, `utils/`, `assets/deplist/` and `assets/templates/`, which are Agnos's own logic.

```
adapters/  -->  sandbox/  <--  cmd/        assets/ (templates, reached via Deps.Embeddeps)
(reaches OS)    (closed)       (wires)
```

```
AgnosConfig/                      written once by `start`, read by every `build`
  project.yaml                    name, version, description  (projectconf)
  themes.yaml                     doc themes: name, id, description  (themesconf)
  ignore.yaml, paths.yaml         SmartIO listing filters / rewrites
  docs/ReadmeHeader.md            README body, a template
sandbox/                          closed: imports nothing outside sandbox/, no OS packages
  new.go                          (gen) New(deps) *api.Sandbox, one binds.<X>Bind per binds/ file
  api/                            contracts only; imports nothing but sandbox/api
    sandbox.go                    (gen) Sandbox struct, one field per api/ file
    actions.go                    Actions struct + props structs + Runtime consts
    cli.go                        (gen) Cli struct + exit consts
  binds/                          one file per api/ file, functions only
    actions.go                    ActionsBind(deps, sandbox): one assignment per action
    cli.go                        (gen) CliBind
  deps/                           contracts; imports only stdlib + sandbox/deps
    deps.go                       (gen) Deps struct, one <Title> <dir>.Lib per dir
    <x>/<x>.go                    type Lib struct of func fields
  internal/
    config/config.go              (gen) ProjectName, Version
    cli/climain.go                (gen) CliMain + dispatch<Name> per command
    commands/<name>/              entries.yaml (decl), entries.go (gen), handler.go (hand)
    actions/<name>/               <name>.go (opens SmartIO, persists, follow-up build) + <name>_internal.go (logic on an open SmartIO)
    actions/build/collect_*.go    collectors: list one dir, title-case names
    actions/build/generate_*.go   entries.go per command, help entries.yaml, doc indexes
    actions/verify/check_*.go     one rule set per file, each returns []string
    parsables/<name>conf/         api.go, new.go, new_empty.go, bind_methods.go, render.go
    smartio/                      transactional fs rooted at --path
    utils/                        RenderGroup, RenderTemplateToDest, Load/Save*Conf, CollectDocTree, entries.yaml field helpers
adapters/
  libs/<lib>/<lib>.go             one package per contract, exports Bind(deps *deps.Deps)
  availables/standard/new.go      (gen) New() deps.Deps calling every lib's Bind
  availables/<name>/new.go        hand-written mix, left alone
assets/                           Go text/templates embedded by asset.go; never `go build ./...`
  start/, all/, deps/, cli/       asset groups; assets/<group>/<path> renders to <path>
  deplist/<dep>/**                one installable dep, mirroring the target layout
  templates/*                     single-file scaffolds (entries.go, command_*, help_entries.yaml, doc_doc.md, *_index.md)
  depsversion.yaml                <dep>: <module>@<version>
cmd/main/main.go                  (gen) standard.New() -> sandbox.New -> CliMain(os.Args[1:])
docs/<Doc>/                       doc.md + props.yaml (+ assets, + sub-doc dirs); **/Index.md is (gen). README.md indexes them all
release/                          git-ignored binaries
```

## Verify rules

- `sandbox/` holds only `api`, `binds`, `deps`, `internal` and `new.go`.
- Nothing under `sandbox/` imports a module package outside `sandbox/`.
- `sandbox/api/*` imports only `sandbox/api`. `sandbox/deps/*` imports only stdlib and `sandbox/deps`.
- Every file of `sandbox/api/` and `sandbox/deps/` parses, and every exported type, func, const and var in them carries a doc comment — `docs/PublicApi/doc.md` is generated from those comments.
- Every `sandbox/binds/` file mirrors an `api/` file and declares only functions.
- `adapters/` holds only `availables` and `libs`. Every `libs/<lib>/` exports `Bind(deps *deps.Deps)`, and every `sandbox/deps/<x>` contract has a lib mentioning its `deps.<X>` field.
- Every `assets/deplist/<dep>/<path>`, rendered with this module, equals `<path>` whenever that file exists here.
- `docs/Commands/doc.md`, when present, names every visible command and every flag it declares.
- Every `docs/**` dir has a parsable `props.yaml`; first-level docs name at least one theme from `themes.yaml`; sub-docs name none. A theme no doc names renders no README section and is not an error.

## Naming

- `Deps` field = title-cased `sandbox/deps/<dir>` (`iodeps` -> `deps.Iodeps`). Always use that spelling.
- Adapter lib binder is always `Bind(deps *deps.Deps)`. Command handler is always `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- An adapter lib's file is named after its package (`adapters/libs/iodeps/iodeps.go`), like a contract's (`sandbox/deps/iodeps/iodeps.go`). A second file in the package is named after what it holds.
- A dep is named after the contract it installs, not its lib (`argvdeps`, lib `verb`).
- Parsables: struct with data fields, then func fields, `Render` last; parsing through `deps.Serializables` (except `moduleconf`, hand-parsed `go.mod`).
