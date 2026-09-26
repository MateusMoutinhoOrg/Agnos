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
| a declaration's schema | `docs/Structure/doc.md`, `docs/EntriesYaml/doc.md`, `assets/doc-server/docs/RouteYaml/doc.md`, `assets/doc-database/docs/Databases/doc.md` |
| the interactive session | `docs/Interview/doc.md` |
| the example suite | `docs/CliExamples/doc.md` |
| which files a build rewrites | `docs/GeneratedFiles/doc.md` |
| unsure whether something is a rule | `docs/Rules/doc.md` — every rule, one page |
| what a command declares | `docs/Commands/<command>.md` — the path is the command's own name |
| what a contract declares | `docs/PublicApi/doc.md` — its index names the page per symbol |

A rule is added or changed in `assets/doc/docs/Rules/doc.md`, never in the rendered copy.
`docs/{Requirements,Workflow,Rules,Extensions,Structure,EntriesYaml,DepList,GeneratedFiles,LibUsage,PublicApi}/`
render from `assets/doc/docs/` into **every** agnos project, this one included; the cli, server,
front, database and example docs render from `assets/doc-cli/`, `assets/doc-server/`,
`assets/doc-front/`, `assets/doc-database/`, `assets/doc-example{,-cli}/`. Editing one means editing that template, and it has to read
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
owns (`agnos build`, `agnos add-command`, `agnos add-route`, `agnos add-dep`, `agnos exec-test`).
`{{.Name}}` is the project being generated and prefixes only what that project answers itself
(`<name> help`, `<name> version`, `<name> start-server`, and whatever its own `add-command`
declared). Never hardcode `agnos` in a template, and never use `{{.Name}}` to spell an agnos
command: in this repo both render `agnos`, so the mistake is invisible here and surfaces only in
a scaffolded project. The version is the same pair: `{{.GeneratorVersion}}` is the release of the
binary that rendered the tree, `{{.Version}}` what the generated project releases under.

Guard a doc line that holds for this repo alone with `{{ if .HasAssets }}`.

### Never hand-edit a generated file or a declaration

Change the template under `assets/` and bootstrap; `docs/GeneratedFiles/doc.md` is the list of
what every build rewrites. One editor per place a declaration holds something:

| File | Its editors |
|---|---|
| `AgnosConfig/extensions.yaml` | `enable-extension` / `disable-extension`, or the `<x>-init` / `<x>-purge` pair that owns the key |
| a command's `entries.yaml` | `add-flag` / `add-arg` / `set-command` and their inverses |
| a route's `route.yaml` | `add-route`, `set-route`, `add-path`, `add-parameter`, `set-body`, `add-body-field`, `import-body`, each with its `set-` and `remove-` pair, plus `rename-route` and `rebalance-routes`; `show-route`, `list-routes` and `explain-route` read it back and write nothing |
| `sandbox/internal/server/errors/handle_*.go` | nobody — eight files `build` writes **once**, then the project's, like a route's `InternalPureHandler.go` |
| a database's `specs.yaml` | `add-database`, `add-table`, `add-table-field`, each with its `set-` and `remove-` pair; `show-database` reads it back and writes nothing |
| `examples/<side>/<name>/result.yaml` | `update-test <name>`, or `exec-test --update` |
| an example | `add-cli-example` / `add-lib-example` and their `remove-` pair |
| a doc | `agnos add-doc` / `agnos remove-doc` |

A missing `extensions.yaml` is a hard error, not a default; `verify` also rejects an unknown key
and any `sandbox-*` on with `sandbox` off.

When declaring one of agnos's own flags, never pass a value that is exactly one of `add-flag`'s
own spellings (`--identifier --example`) — the argv parser counts it as an occurrence and
pollutes the declaration. `--help` is in that set on every command that does not declare a flag
of that name: the dispatch reads it as `help <command>` and prints a screen instead of running.

### Output channels

`deps.Std.Printf` -> stdout, `deps.Std.Log` -> stderr (silenced by `--quiet`), `deps.Std.Error`
-> stderr. Never `fmt.Printf`. A handler returns `api.ExitOk` or `api.ExitFailure`, never
`api.ExitUsage` — the dispatch rejects bad input before it runs.

### Naming is load-bearing

`sandbox.Deps.Iodeps` from `sandbox/deps/iodeps`, `Bind(deps *deps.Deps)` (an adapter fills
`deps.Deps` directly), `CommandHandler(sandbox *api.Sandbox, command *api.Command) int`,
`InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error`
for a route (in `InternalPureHandler.go`, beside its generated `new.go` and `entries.go`), and
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
- `assets/deplist/<dep>/**` must render byte-for-byte to the copy this repo runs on.
- Any new path worth naming gets an entry in `AgnosConfig/structure.yaml`; `verify` fails on an
  entry whose path does not exist.
- A database's `link` field names a `target` that is a table of the same database, a `database`
  field carries `fields` and nests no further, and no table declares a field named `id`.
  `remove-database` refuses a package carrying a `methods_custom.go`.
- Every `route.yaml` declares `methods` (or `ANY` alone), `priority` and `response-type`;
  `add-route` writes `100`, `10` for a `--middleware`. A path reads the segments `start`..`end`
  (inclusive, `-1` the last) into `Entries.<id>` in its `type` (`string`, `integer`, `number`,
  `uuid`); a parameter reads `key` from the first of its `fonts` (`query`, `header`, `cookie`).
  Every `id` is an exported Go name, unique across both, never `FullRoute` or `Body`. The generic
  `RequestHandler` fills `Entries` by its `id` tags through `Deps.Reflectdeps`. `--pattern`
  compiles to these paths plus `segments`; the yaml never holds the pattern.
- On a path, `prefix` is segment-wise (`/admin` never matches `/administrator`), `text-prefix` the
  plain one. `utils/route_match.go` is the generated `IsActionable.go` read against a
  `route.yaml`, for `explain-route`: a change to one is a change to the other.
- **The server is a chain.** Every route matching a request runs, lowest `priority` first, and
  the first one to answer — `SetStatus`, or a `Write`, which sends a `200` — ends it; a handler
  that does neither has declined and the next runs, which is the whole of what a middleware is.
  `route.Locals` is shared by the whole chain of one request. An `InternalPureHandler` returns
  `error`, never a status. A path type or a `trigger` is part of what the route matches on, so
  failing one is a non-match, not a `400`. `phase: after` routes run once it is answered, on a
  frozen response.
- Nothing in the dispatch writes a response: every failure goes through `routeio.Fail` to one of
  the eight `sandbox/internal/server/errors/handle_*.go`, which `build` writes once and never
  rewrites. `Fail` reaches them through the `Fail` field of `api.Server` because a route package
  may not import `sandbox/internal/generated/server/server` — that package imports every route. A failure the dispatch
  raises with nothing to add carries no message, so the wording is the one that file spells;
  that is what makes editing it change what the server says.
- A page is a file of `assets/frontend/` and nothing else — no route, no declaration. The
  `frontend` route `front-init` writes (priority `1000`, after every api route) serves the whole
  tree through the generated `sandbox/internal/generated/frontio/`, whose `SafePath` keeps a caller's path
  inside it; a path naming no file is answered `404` with the `assets/frontend/404.html`
  `front-init` writes, and declined only when that file is gone.
- A pattern changed here is mirrored in `docs/Contributing/doc.md` in the same commit, and the
  reverse.

## Architecture

```
adapters/  -->  sandbox/  <--  cmd/main/        assets/ (templates, read via Deps.Embeddeps)
(reaches OS)    (closed)       (wires them)
```

- **`sandbox/`** — the closed core. It imports only `sandbox/` packages, the stdlib included, so
  text, sorting, hashing and templating come from `sandbox.Deps.<Contract>` too. `api/` holds
  contracts, `deps/` dependency contracts, `internal/` the logic, and `constructors/` is the one
  open list. `internal/generated/` holds every package `build` rewrites whole (cli, config,
  routeio, frontio, databaseio, server/route, server/server) and nothing hand-written; a package
  mixing both — a command, a route, a database — stays under `internal/`. **Every function of `internal/` takes `sandbox *api.Sandbox` first**, and nothing
  else standing for the outside world: holding the api is holding everything.
- **`adapters/`** — the only place OS-bound and third-party code lives.
- **`assets/`** — every generated file's template, one group per extension; `assets/<group>/<path>`
  renders to `<path>`. `utils.AssetGroups()` is the list.
- **`cmd/main/`** — generated; wires an adapter into the sandbox, holds no logic.
- **`AgnosConfig/`** — written once by `start`, read by every `build`.
- **`examples/`** — one directory per example, on the `cli` and the `lib` side.

**Two layers per feature**: an **action** (`sandbox/internal/actions/<name>/`) with `<name>.go`
(opens SmartIO, persists, runs the follow-up `build`) plus `<name>_internal.go` (pure logic on an
already-open SmartIO), and a **command** (`sandbox/internal/commands/<name>/`) with
`entries.yaml`, a generated `new.go` and a hand-written `handler.go`. Both directories are
snake_case for a kebab-case command (`add-command` -> `add_command/`). Only `handler.go` and
contract/adapter pairs are hand-written; everything else is generated.

**`sandbox.Cli.Commands` is the command surface.** `climain.go` is one generic dispatch that
binds a command line onto a copy of the matched declaration, and a handler reads its values by
the id its `entries.yaml` declares — `command.GetString("path")`, `command.GetBool("quiet")`.
`sandbox.Server.Routes` is the http surface the same way, and the server layer mirrors the cli
layer file for file. The front layer declares no unit of its own: a page **is** a file of
`assets/frontend/`, served by one route.

**A database is the one unit with no surface.** `sandbox/internal/databases/<db>/specs.yaml`
declares tables and fields; `build` renders `api.go` (the `<T>Item`/`<T>New`/`<T>Filtrage`
records and the `<Db>` struct of function fields), `new.go` (the `database.Props` and the
wiring) and `methods.go` (every body) from it, with `methods_custom.go` the one hand-written
escape no build reads. There is no field in `api.Sandbox` and no package in
`sandbox/constructors/`: the methods are typed by table, so whoever needs one calls
`<db>.New(sandbox)` on the spot — which touches no key. A `Find<T>By<Field>` is generated for a
`key` field and for no other, because that is the only one the store indexes; every other plain
field is reached through `List<T>` and its filtrage.

**Extensions** are declared in `AgnosConfig/extensions.yaml` and nowhere else — `build` never
infers a mechanic from a directory being present. Nine keys: `sandbox`, `sandbox-deps`,
`sandbox-cli`, `sandbox-server`, `sandbox-front`, `sandbox-database`, `sandbox-example`, `doc`,
`readme`. `false`
means **stop generating**, never **delete**: what the mechanic wrote stays and becomes the
project's, and removing it is what `<x>-purge` does.

**The deps layer is three units**: a **dep** is the contract (`sandbox/deps/<dep>/`), an
**adapter** one implementation of it (`adapters/libs/<adapter>/`), and an **available** a
selection (`adapters/availables/<name>/`). `verify` demands every available fill every field
exactly once — zero panics on first use, two overwrite in silence.

**SmartIO** (`sandbox/internal/smartio/`) is a transactional filesystem rooted at `--path`.
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

**`interview` is the exception, and the only one.** An LLM drives agnos through the plain cli;
the interactive session is what a *person* uses, and it is written for a beginner who has never
read a page of this repo — plain words instead of agnos vocabulary, the next step suggested
first, and no row on a menu that the project in front of them cannot run. Density, token cost and
"read the declaration" do not apply there. Every rule above still binds its *code*; the exception
covers only who its screens are written for. `docs/Interview/doc.md` is its page.

## Testing

There are no Go tests. The checks are `verify` (one `check_*.go` per rule set in
`sandbox/internal/actions/verify/`), a compiling and idempotent `build`, and the example suite:

```bash
./release/bootstrap.bin exec-test               # every example, checked against its golden
./release/bootstrap.bin exec-test --only start  # one example, both sides
./release/bootstrap.bin update-test start       # rewrite one golden, printing what it changes
./release/bootstrap.bin exec-test --update      # rewrite every golden; for a shape change alone
```

Each run that reaches the go runtime pays a `go mod tidy` + `go build`, so prefer `--only`. A
`<name>` on both sides must leave the same tree and exit the same way. `docs/CliExamples/doc.md`
has the rest.

Release: bump `version` in `AgnosConfig/project.yaml`, then `build` + `exec-test --update` (the
bumped version renders into `docs/Requirements/doc.md`, so every golden holding that page moves),
then `agnos publish`.
