# Structure

`(gen)` = written by `build`, never edited — the full list is in
[GeneratedFiles](../GeneratedFiles/doc.md).

```
adapters/  -->  sandbox/  <--  cmd/        assets/ (templates, reached via Deps.Embeddeps)
(reaches OS)    (closed)       (wires)
```

Every line below is one entry of `AgnosConfig/structure.yaml` — add `<path>:
{description: "..."}` there, nested under `children:` of its parent, with `dir: true` on a
directory, `gen: true` on a file `build` rewrites, and `order:` to place it among its siblings
(unordered siblings follow, alphabetically).

```
AgnosConfig/                     written once by `start`, read by every `build`
  project.yaml                   name, version, description  (projectconf)
  themes.yaml                    doc themes: name, id, description  (themesconf)
  structure.yaml                 this tree  (structureconf)
  ignore.yaml                    SmartIO listing filters  (ignorableconf)
  paths.yaml                     SmartIO listing rewrites  (pathreplacerconf)
  docs/ReadmeHeader.md           README body, a template
sandbox/                         closed: imports nothing outside sandbox/, no OS packages
  new.go                         (gen) New(deps) *api.Sandbox, one binds.<X>Bind per binds/ file
  api/                           contracts only; imports nothing but sandbox/api
    sandbox.go                   (gen) Sandbox struct, one field per api/ file
    actions.go                   Actions struct + props structs + Runtime consts
    cli.go                       (gen) Cli struct + exit consts
  binds/                         one file per api/ file, functions only
    actions.go                   ActionsBind(deps, sandbox): one assignment per action
    cli.go                       (gen) CliBind
  deps/                          contracts; imports only stdlib + sandbox/deps
    deps.go                      (gen) Deps struct, one <Title> <dir>.Lib per dir
    <x>/<x>.go                   type Lib struct of func fields
  internal/                      the logic; unreachable from outside the sandbox
    config/config.go             (gen) ProjectName, Version
    cli/climain.go               (gen) CliMain + dispatch<Name> per command
    commands/<name>/             entries.yaml (decl), entries.go (gen), handler.go (hand)
    actions/<name>/              <name>.go (opens SmartIO, persists, follow-up build) + <name>_internal.go (logic on an open SmartIO)
    actions/build/collect_*.go   collectors: list one dir, title-case names
    actions/build/generate_*.go  entries.go per command, help entries.yaml, doc indexes
    actions/verify/check_*.go    one rule set per file, each returns []string
    parsables/<name>conf/        api.go, new.go, new_empty.go, bind_methods.go, render.go
    smartio/                     transactional fs rooted at --path
    utils/                       RenderGroup, RenderTemplateToDest, Load*Conf, CollectDocTree, FlattenStructure, entries.yaml field helpers
adapters/                        the only place OS-bound and third-party code lives
  libs/<lib>/<lib>.go            one package per contract, exports Bind(deps *deps.Deps)
  availables/standard/new.go     (gen) New() deps.Deps calling every lib's Bind
  availables/<name>/new.go       hand-written mix, left alone
assets/                          Go text/templates embedded by asset.go; never `go build ./...`
  start/                         written once, on `start`
  all/                           rendered by every `build`
  deps/                          rendered when sandbox/deps/ exists
  cli/                           rendered when sandbox/internal/cli/ exists
  deplist/<dep>/                 one installable dep, mirroring the target layout
  templates/                     single-file scaffolds (entries.go, command_*, help_entries.yaml, doc_doc.md, *_index.md)
  depsversion.yaml               <dep>: <module>@<version>
cmd/main/main.go                 (gen) standard.New() -> sandbox.New -> CliMain(os.Args[1:])
docs/                            one dir per doc, holding doc.md + props.yaml (+ assets, + sub-docs). README.md indexes them all
  **/Index.md                    (gen) written for every doc that has sub-docs
examples/                        one dir per example; `exec-test` runs each and diffs it against its golden
  cli/<name>/example.sh          the example, run with `sh` and its own dir as cwd
  lib/<name>/example.go          the example, run with `go run` and its own dir as cwd
  <side>/<name>/result.yaml      (gen) golden: cli-output, exit-code, sha256 of every TestDir file
  <side>/<name>/TestDir/         the only place an example writes; removed before every run
release/                         git-ignored binaries, and the exec-test cli alias
```

Every rule this shape has to hold to — layers, naming, generated files, docs — is in
[Rules](../Rules/doc.md); the command that makes each change is in
[Workflow](../Workflow/doc.md).
