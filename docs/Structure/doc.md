# Structure

`(gen)` = written by `build`, never edited — the full list is in
[GeneratedFiles](../GeneratedFiles/doc.md).

```
adapters/  -->  sandbox/  <--  cmd/        assets/ (templates, reached via Deps.EmbedDeps)
(reaches OS)    (closed)       (wires)
```

Every line below is one entry of `AgnosConfig/structure.yaml` — add `<path>:
{description: "..."}` there, nested under `children:` of its parent, with `dir: true` on a
directory, `gen: true` on a file `build` rewrites, and `order:` to place it among its siblings
(unordered siblings follow, alphabetically).

```
AgnosConfig/                                    written once by `start`, read by every `build`
  project.yaml                                  project-name, version, description  (projectconf)
  themes.yaml                                   doc themes: name, id, description  (themesconf)
  structure.yaml                                this tree  (structureconf)
  extensions.yaml                               which mechanics agnos generates  (extensionsconf)
  paths.yaml                                    StagedFS listing rewrites  (pathsconf)
  docs/ReadmeHeader.md                          README body, a template
sandbox/                                        closed: imports nothing outside sandbox/, no OS packages
  new.go                                        (gen) New(deps) *api.Sandbox, one <x>.Constructor(&self) per constructors/ dir
  api/                                          contracts only; imports nothing but sandbox/deps, for Sandbox.Deps, and the OpinionatedAgnos<X> contracts it aliases
    sandbox.go                                  (gen) Sandbox struct, one field per api/ file
    actions.go                                  Actions struct + props structs + Runtime consts
    cli.go                                      (gen) Cli + exit consts, aliases of the OpinionatedAgnosCli contract
    command.go                                  (gen) Command/CommandFlag/CommandArg/CommandResponse/CommandFailure, aliases of the OpinionatedAgnosCli contract
    config.go                                   (gen) Config struct, the project's own name and version
    trigger.go                                  (gen) Trigger/TriggerType, the condition a route and a command match on; aliases of the OpinionatedAgnosCli contract
  constructors/                                 one <x>/constructor.go per field of the Sandbox; written once, then yours
    <x>/constructor.go                          Constructor(sandbox): sandbox.<X> = <x>.New<X>(sandbox)
  deps/                                         contracts; each <x>/ imports nothing at all, an OpinionatedAgnos<X>/ nothing but other contracts
    deps.go                                     (gen) Deps struct, one <Field> <dir>.Contract per dir
    <x>/<x>.go                                  type Sandbox struct of func fields
    OpinionatedAgnosCli/opinionatedagnoscli.go  the cli lib: the command declaration types, MainProps, and Main, NewCommand, BindCommand, Fail, FailureOf, MatchTrigger
  internal/                                     the logic; unreachable from outside the sandbox
    generated/                                  (gen) every package the build rewrites whole — the registries and config; never edited by hand
      config/new.go                             (gen) NewConfig(sandbox) api.Config: ProjectName, Version
      cli/new.go                                (gen) NewCli(sandbox) api.Cli: Cli.Commands in run order + Cli.Fail + Cli.Main, which hands the line to Deps.OpinionatedAgnosCli.Main
    commands/<name>/                            command.yaml (decl), new.go + input.go (gen), handler.go (hand); under the folder of its category (core/, cli/, server/, middleware/, …) — the command.yaml is what makes a dir a command
    actions/new.go                              NewActions(sandbox) api.Actions: one assignment per action
    actions/<name>/                             <name>.go (opens StagedFS, persists, follow-up build) + <name>_internal.go (logic on an open StagedFS)
    actions/build/collect_*.go                  collectors: list one dir, title-case names
    actions/build/generate_*.go                 new.go + input.go per command and per route, help's command.yaml, doc indexes
    actions/verify/check_*.go                   one rule set per file, each returns []string
    declarations/<name>conf/                    api.go, new.go, new_empty.go, bind_methods.go, render.go
    apishape/                                   the sandbox/api convertibility rule and the converter plan the remote-dep shim is generated from
    stagedfs/                                   transactional fs rooted at --path
    utils/                                      RenderGroup, RenderTemplateToDest, Load*Conf, CollectDocTree, FlattenStructure, GoIdentifier, command.yaml and route.yaml field helpers
adapters/                                       the only place OS-bound and third-party code lives
  impls/<adapter>/<adapter>.go                  one package per adapter, exports Bind(deps *deps.Deps)
  impls/<adapter>/adapter.yaml                  which dep the adapter fills, which module it pins  (adapterconf)
  impls/OpinionatedAgnosCli/                    the cli lib: the dispatch chain, binder, matcher and trigger every command runs through, over the stdlib
  bindings/<name>/binding.yaml                  which adapters this binding binds, one per Deps field  (bindingconf)
  bindings/<name>/new.go                        (gen) New() deps.Deps calling the Bind of every adapter declared; a binding with no binding.yaml is hand-written and left alone
assets/                                         Go text/templates embedded by asset.go; never `go build ./...`
  start/                                        written once, on `start`
  sandbox/                                      rendered when the `sandbox` extension is on; one group per extension from here down
  deps/                                         rendered when `deps` is on
  cli/                                          rendered when `cli` is on
  server/                                       rendered when `server` is on
  readme/                                       rendered when `readme` is on
  doc/                                          rendered when `doc` is on
  doc-<x>/                                      rendered when `doc` and every `<x>` it names are on (doc-cli, doc-server, doc-front, doc-database, doc-backoffice, doc-example, doc-example-cli)
  dep-catalog/<dep>/                            one installable contract, dep.yaml beside the target layout it mirrors; dep-catalog/OpinionatedAgnos<X> are the opinionated libs, one per mechanic with code (Cli, Server, Front, Database)
  adapter-catalog/<adapter>/                    one installable adapter, adapter.yaml beside the target layout it mirrors
  templates/                                    single-file scaffolds (command_*, route_*, database_*, front_*, start_server_*, server_handle_*, cli_handle_*, help_command.yaml, doc_page.md, doc_index.md)
  templates/backoffice/                         the tree backoffice-init writes once, at the path each file holds in it; assets/ under it is copied verbatim
cmd/main/main.go                                (gen) standard.New() -> sandbox.New -> Main(os.Args[1:])
docs/                                           one dir per doc, holding doc.md + doc.yaml (+ assets, + sub-docs). README.md indexes them all
  **/Index.md                                   (gen) written for every doc that has sub-docs
  Commands/<command>.md                         (gen) one page per visible command, indexed by docs/Commands/doc.md
  PublicApi/<contract>.md                       (gen) one page per file of sandbox/api and per contract of sandbox/deps, indexed by docs/PublicApi/doc.md
examples/                                       one dir per example; `run-examples` runs each and diffs it against its golden
  cli/<name>/example.sh                         the example, run with `sh` and its own dir as cwd
  lib/<name>/example.go                         the example, run with `go run` and its own dir as cwd
  <side>/<name>/result.yaml                     (gen) golden: cli-output, exit-code, sha256 of every assert-dir file
  <side>/<name>/test-dir/                       the only place an example writes; removed before every run
  <side>/<name>/assert-dir/                     what the example copied out of test-dir to assert; the golden's tree
release/                                        git-ignored binaries, and the run-examples cli alias
```

Every rule this shape has to hold to — layers, naming, generated files, docs — is in
[Rules](../Rules/doc.md); the command that makes each change is in
[Workflow](../Workflow/doc.md).
