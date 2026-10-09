# CLAUDE.md

## What this is

Agnos (`agnos`) is a Go CLI that **scaffolds and regenerates other Go CLIs**. `agnos start`
writes a project skeleton; `agnos build` re-renders every generated file from `text/template`
assets embedded in the binary; commands like `add-command`, `add-flag`, `add-dep` declare the
project's command surface without a file being hand-edited.

Agnos is built with itself: `agnos build` regenerates this repo in place, and the result must
compile and be idempotent. That self-hosting constraint drives every rule below.

## Which doc to read

`docs/` is the source of truth, indexed in `README.md`. Read the one the change needs — not all
of them. In doubt, `docs/Workflow/doc.md`.

| Changing | Read |
|---|---|
| a template under `assets/`, a collector, the build order | `docs/BuildPipeline/doc.md` |
| adding a command, action, route, layer, extension, dep, example or doc **to agnos itself** | `docs/Contributing/doc.md` |
| the recipe every agnos project follows for a change | `docs/Workflow/doc.md` |
| a declaration's schema | `docs/Structure/doc.md`, `docs/CommandYaml/doc.md`, `assets/doc-server/docs/RouteYaml/doc.md`, `assets/doc-database/docs/Databases/doc.md` |
| the interactive session | `docs/Interview/doc.md` |
| the example suite | `docs/CliExamples/doc.md` |
| which files a build rewrites | `docs/GeneratedFiles/doc.md` |
| unsure whether something is a rule | `docs/Rules/doc.md` — every rule, one page |
| what a command declares | `docs/Commands/<command>.md` — the path is the command's own name |
| what a contract declares | `docs/PublicApi/doc.md` — its index names the page per symbol |

A rule is added or changed in `assets/doc/docs/Rules/doc.md`, never in the rendered copy.
`docs/{Adapters,Requirements,Workflow,Rules,Extensions,Structure,DepList,GeneratedFiles,LibUsage,PublicApi}/`
render from `assets/doc/docs/` into **every** agnos project, this one included; the cli, server,
front, database, backoffice and example docs render from `assets/doc-cli/`, `assets/doc-server/`,
`assets/doc-front/`, `assets/doc-database/`, `assets/doc-backoffice/`, `assets/doc-example{,-cli}/`. Editing one means editing that template, and it has to read
correctly in a scaffolded project, not only here.

## Traps

None of these is caught by a `verify` check: each one is invisible in this repo and surfaces
only later, or in a generated project.

### Never run an installed `agnos build` on this repo

It rewrites the tree to the older binary's shape. Always bootstrap:

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build                 # verify + regenerate + go mod tidy + compile
./release/bootstrap.bin verify                # the schema check alone, writes nothing
./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
./release/bootstrap.bin local-install          # install the result
```

Compile scope is always `./cmd/... ./sandbox/... ./adapters/...` — **never `go build ./...`**,
because `assets/` holds Go templates, not compilable Go. Every command takes `--path <dir>`
(default `.`) and `-q`, and runs `build` for you.

### Two names, never swapped

`{{.GeneratorName}}` is the cli running the build — agnos — and prefixes every command agnos
owns (`agnos build`, `agnos add-command`, `agnos add-route`, `agnos add-dep`, `agnos run-examples`).
`{{.ProjectName}}` is the project being generated and prefixes only what that project answers itself
(`<name> help`, `<name> version`, `<name> start-server`, and whatever its own `add-command`
declared). Never hardcode `agnos` in a template, and never use `{{.ProjectName}}` to spell an agnos
command: in this repo both render `agnos`, so the mistake is invisible here and surfaces only in
a scaffolded project. The version is the same pair: `{{.GeneratorVersion}}` is the release of the
binary that rendered the tree, `{{.Version}}` what the generated project releases under. No
template reads a bare `.Name` at its top level: a unit's own name is `.CommandName`,
`.RouteName`, `.DatabaseName`, `.ExampleName`, `.DocName`.

Guard a doc line that holds for this repo alone with `{{ if .HasAssets }}`.

### Never hand-edit a generated file or a declaration

Change the template under `assets/` and bootstrap; `docs/GeneratedFiles/doc.md` is the list of
what every build rewrites. One editor per place a declaration holds something:

| File | Its editors |
|---|---|
| `AgnosConfig/extensions.yaml` | `enable-extension` / `disable-extension`, or the `<x>-init` / `<x>-purge` pair that owns the key |
| a command's `command.yaml` | `add-flag` / `add-arg` / `set-command`, the `set-` editor and the `remove-` inverse of each |
| a route's `route.yaml` | `add-route`, `set-route`, `add-path`, `add-parameter`, `set-body`, `add-body-field`, `import-body`, each with its `set-` and `remove-` pair, plus `rename-route` and `rebalance-routes`; `show-route`, `list-routes` and `explain-route` read it back and write nothing |
| `sandbox/internal/server/errors/handle_*.go` | nobody — eight files `build` writes **once**, then the project's, like a route's `handler.go` |
| a database's `database.yaml` | `add-database`, `add-table`, `add-table-field`, each with its `set-` and `remove-` pair; `show-database` reads it back and writes nothing |
| `examples/<side>/<name>/result.yaml` | `update-example <name>`, or `run-examples --update` |
| an example | `add-cli-example` / `add-lib-example` and their `remove-` pair |
| a doc | `agnos add-doc` / `agnos remove-doc` |

A missing `extensions.yaml` is a hard error, not a default; `verify` also rejects an unknown key
and any mechanic that renders into the sandbox on with `sandbox` off.

When declaring one of agnos's own flags, never pass a value that is exactly one of `add-flag`'s
own spellings (`--identifier --example`) — the argv parser counts it as an occurrence and
pollutes the declaration. `--help` is in that set on every command that does not declare a flag
of that name: the dispatch reads it as `help <command>` and prints a screen instead of running.

### Output channels

`deps.StdDeps.Printf` -> stdout, `deps.StdDeps.Logf` -> stderr (silenced by `--quiet`), `deps.StdDeps.Eprintf`
-> stderr. Never `fmt.Printf`. A handler returns `api.ExitOk` or `api.ExitFailure`, never
`api.ExitUsage` — the dispatch rejects bad input before it runs.

### Naming is load-bearing

`sandbox.Deps.IoDeps` from `sandbox/deps/iodeps` (a dep's field is its directory title-cased, a
trailing `deps` spelled `Deps`), `Bind(deps *deps.Deps)` (an adapter fills `deps.Deps` directly),
`Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error`
for a command and
`Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error`
for a route (each in `handler.go`, beside its `generated.new.go` and `generated.input.go`), and
`(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error` for each of the
eight `handle_*.go`. **Every file is an instance of a pattern**: new code copies an
existing sibling exactly — same filenames, same function names, same ordering. `verify` and the
collectors read shape by convention, so a one-off breaks the machine reader.

### Searching this repo

`git grep` hits the `assets/` mirror of nearly every file, so a symbol search fans out to
hundreds of hits. Narrow with `':!assets' ':!examples'` to **locate** a file — never to judge
impact. Once you have it, the twin under `assets/` is what the next build renders from, and
editing only the rendered copy is undone in silence.

### The rest

- Every exported declaration of `sandbox/api/` and `sandbox/deps/` carries a doc comment —
  `docs/PublicApi/` is generated from those comments and `verify` fails without them.
- `assets/dep-catalog/<dep>/**` must render byte-for-byte to the copy this repo runs on.
- Any new path worth naming gets an entry in `AgnosConfig/structure.yaml`; `verify` fails on an
  entry whose path does not exist.
- A database's `link` field names a `target` that is a table of the same database, an `object`
  field carries `fields` and nests no further, and no table declares a field named `id`, nor
  an `object` named `position` or `values` (Keep reserves both).
  `remove-database` refuses a package carrying a `methods_custom.go`.
- A database's `key-prefix` is a folder inside `sandbox.Config.DatabaseDir`, which the generated
  `database_dir` middleware reads from `--database` (default `data`) on every command line. The
  store (Keep) resolves every key under the directory the program runs from and escapes any byte
  but `a-z0-9-_`, so `--database` is a relative path of such segments, enforced by its `pattern`;
  `verify` refuses a `key-prefix` starting with `data/`, a declaration from before the flag.
- Every `route.yaml` declares `methods` (or `ANY` alone), `priority` and `response-type`;
  `add-route` writes `100`, `10` for a `--middleware`. A path reads the segments `start`..`end`
  (inclusive, `-1` the last) into `Input.<id>` in its `type` (`string`, `integer`, `number`,
  `uuid`); a parameter reads `key` from the first of its `sources` (`query`, `header`, `cookie`).
  Every `id` is an exported Go name, unique across both, never `FullRoute` or `Body`. The
  `OpinionatedAgnosServer` lib fills `Input` by its `id` tags, by reflection. `--pattern`
  compiles to these paths plus `segments`; the yaml never holds the pattern.
- On a path, `prefix` is segment-wise (`/admin` never matches `/administrator`), `text-prefix` the
  plain one. `utils/route_match.go` is the `OpinionatedAgnosServer` lib's `matches.go` read
  against a `route.yaml`, for `explain-route`: a change to one is a change to the other.
- **The server is a chain.** Every route matching a request runs, lowest `priority` first, and
  the first one to answer — `SetStatus`, or a `Write`, which sends a `200` — ends it; a handler
  that does neither has declined and the next runs, which is the whole of what a middleware is.
  One `props *routeprops.RouteProps` per request is shared by the whole chain. `routeprops.go` is
  generated: it embeds every struct of the other files of `sandbox/internal/routeprops/` —
  `project.go` (the project's, written once) plus one file per mechanic (`backoffice.go`). It is
  not in `sandbox/api` so a field may name a project type (a database record), and
  `api.Route.Props` holds it as `any`. `api.Config`/`api.Sandbox` embed every struct of
  `sandbox/api/<x>config.go`/`<x>sandbox.go` the same way: a mechanic adds a part
  (`clisandbox.go`, `serversandbox.go`, `backofficeconfig.go`), never edits the project's.
  `api.Sandbox` declares only `Deps` and `Config` itself; a contract the project writes is a
  field of its `ProjectSandbox` (agnos's `Actions`), which `verify` demands.
  `commandprops.CommandProps` is the cli's mirror. A handler is handed no request, so what
  it reads is declared (`Input`, the body on `Input.Body`). A `Handle` returns
  `error`, never a status; it refuses a request by returning
  `sandbox.Deps.OpinionatedAgnosServer.Fail`. A path type or a
  `trigger` is part of what the route matches on, so failing one is a non-match, not a `400`.
  There is one chain and no `after` phase: a `phase` key is an old declaration `verify` names.
- Nothing in the dispatch writes a response: only the lib's dispatch raises, and every failure
  reaches one of the eight `sandbox/internal/server/errors/handle_*.go`, which `build` writes once
  and never rewrites, through the `Fail` field of `api.Server` the generated registry fills. A
  failure the dispatch raises with nothing to add carries no message, so the wording is the one
  that file spells; that is what makes editing it change what the server says.
- A page is a file of `assets/front/` and nothing else — no route, no declaration. The
  `front` route `front-init` writes (priority `1000`, after every api route) serves the whole
  tree through the `OpinionatedAgnosFront` lib, whose `SafePath` keeps a caller's path
  inside it; a path naming no file is answered `404` with the `assets/front/404.html`
  `front-init` writes, and declined only when that file is gone.
- A pattern changed here is mirrored in `docs/Contributing/doc.md` in the same commit, and the
  reverse.

## Architecture

```
adapters/  -->  sandbox/  <--  cmd/main/        assets/ (templates, read via Deps.EmbedDeps)
(reaches OS)    (closed)       (wires them)
```

- **`sandbox/`** — the closed core. It imports only `sandbox/` packages, the stdlib included, so
  text, sorting, hashing and templating come from `sandbox.Deps.<Contract>` too. `api/` holds
  contracts, `deps/` dependency contracts, `internal/` the logic, and `constructors/` is the one
  open list. A Go file `build` rewrites whole is named `generated.<name>.go` and opens with
  `// Code generated by agnos. DO NOT EDIT.`, wherever it sits — the registries
  `internal/{cli,server}/generated.new.go`, `internal/config/generated.new.go`, a command's
  `generated.new.go` beside its `handler.go`; a file without the prefix is the project's. **Every function of `internal/` takes `sandbox *api.Sandbox` first**, and nothing
  else standing for the outside world: holding the api is holding everything.
  `api.Sandbox` embeds `api.ProjectSandbox` and `api.Config` embeds `api.ProjectConfig`
  (`api/projectsandbox.go`, `api/projectconfig.go`): written once by `start`, never by `build` —
  the project's own public fields. `apishape` admits an embedded field only for such a struct.
- **`adapters/`** — the only place OS-bound and third-party code lives.
- **`assets/`** — every generated file's template, one group per extension; `assets/<group>/<path>`
  renders to `<path>`. `utils.AssetGroups()` is the list.
- **`cmd/main/`** — generated; wires an adapter into the sandbox, holds no logic.
- **`AgnosConfig/`** — written once by `start`, read by every `build`.
- **`examples/`** — one directory per example, on the `cli` and the `lib` side.

**Two layers per feature**: an **action** (`sandbox/internal/actions/<name>/`) with `<name>.go`
(opens StagedFS, persists, runs the follow-up `build`) plus `<name>_internal.go` (pure logic on an
already-open StagedFS), and a **command** (`sandbox/internal/commands/<category>/<name>/`, the
folder its category snake_cased: `core/`, `cli/`, `server/`, `middleware/`, …) with
`command.yaml`, a `generated.new.go` and `generated.input.go`, and a hand-written `handler.go`. Both directories are
snake_case for a kebab-case command (`add-command` -> `add_command/`). Only `handler.go` and
contract/adapter pairs are hand-written; everything else is generated.

**`sandbox.Cli.Commands` is the command surface.** `OpinionatedAgnosCli.Main` is one dispatch that
binds a command line onto a copy of the matched declaration, and a handler reads its values off
the `Input` its `command.yaml` declares — `input.Path`, `input.Quiet`.
`sandbox.Server.Routes` is the http surface the same way, and the server layer mirrors the cli
layer file for file. The front layer declares no unit of its own: a page **is** a file of
`assets/front/`, served by one route.

**A database is the one unit with no surface.** `sandbox/internal/databases/<db>/database.yaml`
declares tables and fields; `build` renders `generated.api.go` (the `<T>Record`/`<T>Input`/`<T>Filter`
records and the `<Db>` struct of function fields), `generated.new.go` (the `databasedeps.Props`
and the wiring) and `generated.methods.go` (every body) from it, with `methods_custom.go` the one hand-written
escape no build reads. There is no field in `api.Sandbox` and no package in
`sandbox/constructors/`: the methods are typed by table, so whoever needs one calls
`<db>.New(sandbox)` on the spot — which touches no key. A `Find<T>By<Field>` is generated for a
`key` field and for no other, because that is the only one the store indexes; every other plain
field is reached through `List<T>s` and its filter.

**Extensions** are declared in `AgnosConfig/extensions.yaml` and nowhere else — `build` never
infers a mechanic from a directory being present. Ten keys: `sandbox`, `deps`,
`cli`, `server`, `front`, `database`, `backoffice`,
`example`, `doc`, `readme`. `backoffice` generates nothing: `backoffice-init`
writes `assets/templates/backoffice/**` once — the backups included: the `backup` database,
`sandbox/internal/snapshots/` and the `*-backup*` routes — and the key gates its doc. `false`
means **stop generating**, never **delete**: what the mechanic wrote stays and becomes the
project's, and removing it is what `<x>-purge` does.

**The deps layer is three units**: a **dep** is the contract (`sandbox/deps/<dep>/`), an
**adapter** one implementation of it (`adapters/impls/<adapter>/`), and a **binding** a
selection (`adapters/bindings/<name>/`). `verify` demands every binding fill every field
exactly once — zero panics on first use, two overwrite in silence.

**The opinionated libs** — `OpinionatedAgnosCli`, `OpinionatedAgnosServer`, `OpinionatedAgnosFront`,
`OpinionatedAgnosDatabase` — are the one kind of dep that carries an agnos mechanic rather than a
library's raw capability: the dispatch, the binders, the matchers, json-schema, the file layer,
the database readers. Code that is the same in every project lives there, never in a
`generated.*.go`. Each mechanic's `-init` installs its lib; the mechanic's `sandbox/api/`
files are type aliases of the lib's contract, so `api.Command` still reads the same. A lib holds
no dep: an entry point takes a `MainProps` the registry builds, any other function the one dep it
needs first. This repo runs on its own `OpinionatedAgnosCli`: change
`assets/{dep-catalog,adapter-catalog}/OpinionatedAgnosCli/` and re-mirror the installed copy.

**StagedFS** (`sandbox/internal/stagedfs/`) is a transactional filesystem rooted at `--path`.
Actions pass project-relative paths only. Writes buffer until `Persist`, but `List*` reads disk —
so an action that runs `build` as a follow-up must `Persist` first.

`docs/Contributing/doc.md` holds the recipe for adding any of these.

## Audience: LLMs, not humans

The primary reader and writer of this repo is an LLM. Every choice optimizes for machine reading
and machine writing, and token cost is a first-class constraint.

- **Generate over hand-write.** If a file can be rendered from a template, a collector or a
  declaration, it must be. A new hand-written file needs a reason why generation cannot cover it.
- **Generate over document.** `README.md`, every `Index.md`, `docs/Commands/`, `docs/PublicApi/`
  and `docs/Structure/` are rendered, never typed. Document by commenting the contract.
- **Docs are short, objective and dense.** Tables, commands, file paths and rules — no prose, no
  tutorials, no repetition across pages. Say the rule once, in the one place it belongs, and link
  with a relative path. A page an LLM re-reads on every task costs tokens each time.
- **One page per unit.** A doc that grows with the project is split so a lookup costs the unit
  asked about: one page per command, one per contract, indexed by a `doc.md` that routes.
- **Convention over configuration.** Uniformity is what makes generation possible.
- **Deterministic and idempotent.** Same input, same bytes out.

**Two exceptions, and only two, both written for a person.** `interview`: an LLM drives agnos
through the plain cli; the interactive session is what a *person* uses, and it is written for a
beginner who has never read a page of this repo — plain words instead of agnos vocabulary, the
next step suggested first, and no row on a menu that the project in front of them cannot run.
`docs/Interview/doc.md` is its page. And a server project's `docs/Routes/`: it is read by whoever
calls the server, so every page speaks plain words and carries `curl` requests that run as they
are, generated from the `route.yaml` by `CollectRouteDocs` and `assets/templates/route_page.md`.
Density, token cost and "read the declaration" do not apply to either. Every rule above still
binds their *code*; the exceptions cover only who their screens and pages are written for.

## Testing

There are no Go tests. The checks are `verify` (one `check_*.go` per rule set in
`sandbox/internal/actions/verify/`), a compiling and idempotent `build`, and the example suite:

```bash
./release/bootstrap.bin run-examples               # every example, checked against its golden
./release/bootstrap.bin run-examples --only start  # one example, both sides
./release/bootstrap.bin update-example start       # rewrite one golden, printing what it changes
./release/bootstrap.bin run-examples --update      # rewrite every golden; for a shape change alone
```

Each run that reaches the go runtime pays a `go mod tidy` + `go build`, so prefer `--only`. A
`<name>` on both sides must leave the same tree and exit the same way. `docs/CliExamples/doc.md`
has the rest.

**Status: pre-beta.** A project scaffolded by an earlier release must keep building on the next
one: no rename or removal of a command, flag, declaration key or exported contract without
keeping the old one working, or `verify` naming the migration. `docs/Contributing/doc.md#breaking-changes`
is the table; an unavoidable break is named in the release.

Release: bump `version` in `AgnosConfig/project.yaml`, then `build` + `run-examples --update` (the
bumped version renders into `docs/Requirements/doc.md`, so every golden holding that page moves),
then `agnos publish`.
