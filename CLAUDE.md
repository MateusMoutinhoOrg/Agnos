# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Agnos (`agnos`) is a Go CLI that **scaffolds and regenerates other Go CLIs**. `agnos start`
writes a project skeleton; `agnos build` re-renders every generated file from `text/template`
assets embedded in the binary; commands like `add-command`, `add-flag`, `add-dep` declare
the project's command surface without a file being hand-edited.

Agnos is built with itself: `agnos build` regenerates this repo in place, and the result must
compile and be idempotent. That self-hosting constraint drives every rule below.

Full documentation lives under `docs/` and is the source of truth, indexed in `README.md`
(one section per theme of `AgnosConfig/themes.yaml`; there is no index file between the
README and a doc). Start with
`docs/Rules/doc.md` (every rule, in one page — generated from
`assets/all/docs/Rules/doc.md`, which is where a rule is added or changed),
`docs/Structure/doc.md` (schema), `docs/BuildPipeline/doc.md` (what `build` does),
`docs/Workflow/doc.md` (the recipe for every change any agnos project takes) and
`docs/Contributing/doc.md` (what is specific to changing agnos itself).

`docs/{Requirements,Workflow,Rules,Structure,EntriesYaml,DepList,GeneratedFiles,Commands,LibUsage,LibExamples,PublicApi}/`
are rendered from `assets/all/docs/` into **every** agnos project, this one included,
`docs/{CliInstall,CliExamples}/` from `assets/cli/docs/`,
`docs/{RouteYaml,Routes,ServerUsage}/` from `assets/server/docs/` and `docs/FrontUsage/` from
`assets/front/docs/`: editing one means editing that
template, and it has to read correctly in a scaffolded project, not only here. Guard a line
that holds for this repo alone with `{{ if .HasAssets }}`.

**Two names, never swapped.** `{{.GeneratorName}}` is the cli running the build — agnos — and
prefixes every command agnos owns (`agnos build`, `agnos add-command`, `agnos add-route`,
`agnos add-dep`, `agnos exec-test`). `{{.Name}}` is the project being generated and prefixes
only what that project answers itself (`<name> help`, `<name> version`, `<name> start-server`,
and whatever its own `add-command` declared). Never hardcode `agnos` in a template, and never
use `{{.Name}}` to spell an agnos command: in this repo both render `agnos`, so the mistake is
invisible here and surfaces only in a scaffolded project.

## Audience: LLMs, not humans

The primary reader and writer of this repo is an LLM. Every choice — code shape, docs,
naming, file layout — optimizes for machine reading and machine writing, and token cost is a
first-class constraint. Humans are a secondary audience.

- **Generate over hand-write.** If a file can be rendered from a template, a collector or a
  declaration, it must be — hand-written code is only contracts, adapters, `sandbox/internal/`
  and `handler.go`. A new hand-written file needs a reason why generation cannot cover it.
- **Generate over document.** The same applies to docs: `README.md` (its documentation index
  included), every `Index.md`, `docs/Commands/`, `docs/PublicApi/` and `docs/Structure/` are
  rendered, never typed. Document by commenting the contract, not by writing a page.
- **Docs are short, objective and dense.** Tables, commands, file paths and rules — no prose,
  no narrative, no tutorials, no motivation sections, no repetition across pages. Shorter is
  strictly better: a page an LLM re-reads on every task costs tokens each time. Say the rule
  once, in the one place it belongs, and link with a relative path from elsewhere.
- **Convention over configuration.** Shape is read by convention (`verify` and the collectors
  parse it), so uniformity is what makes generation possible. Every file is an instance of a
  pattern; one-offs break the machine reader.
- **Deterministic and idempotent.** Same input, same bytes out. An LLM must be able to
  regenerate the tree and diff it to zero.

## Bootstrap workflow

**Never run an installed `agnos build` on this repo after touching templates, collectors,
deps or adapters** — it rewrites the tree to the older binary's shape. Always bootstrap:

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build                                    # verify + regenerate + go mod tidy + compile
./release/bootstrap.bin verify                                   # the schema check alone, writes nothing
./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
./release/bootstrap.bin local-install                            # install the result
```

Compile scope is always `./cmd/... ./sandbox/... ./adapters/...` — **never `go build ./...`**,
because `assets/` holds Go templates, not compilable Go. Every command takes `--path <dir>`
(default `.`) and `-q`, and runs `build` for you.

There are no Go tests. The checks are `verify` (layers, contracts, adapters, deplist mirror,
docs, structure — one `check_*.go` per rule set in `sandbox/internal/actions/verify/`), a
compiling and idempotent `build`, and the example suite:

```bash
./release/bootstrap.bin exec-test               # every example, checked against its golden
./release/bootstrap.bin exec-test --only start  # one example, both sides
./release/bootstrap.bin update-test start       # rewrite one golden, printing what it changes
./release/bootstrap.bin exec-test --update      # rewrite every golden; for a shape change alone
```

`examples/{cli,lib}/<name>/` holds one `example.sh` / `example.go` that runs with its own
directory as the cwd and writes only into `TestDir`, then copies out of it into `AssertDir` the
paths it asserts; `result.yaml` (output, exit code, sha256 of every `AssertDir` file) is the
golden, written by `exec-test`, never by hand. An example that copies nothing out fails, and
only `start` copies the whole tree — every other one asserts what its own command touched, so a
`start` template change moves one golden and not all of them. `exec-test` writes an `agnos`
alias (`go run ./cmd/main`) into `release/exec-test/` and puts it in front of the PATH, so an
example always runs against this tree; each run that reaches the go runtime pays a
`go mod tidy` + `go build`, so prefer `--only`. A `<name>` on both sides must leave the same
tree and exit the same way, so both sides copy the same set. Create and delete examples with
`add-cli-example` / `add-lib-example` / `remove-cli-example` / `remove-lib-example` only.
Release: bump `version` in `AgnosConfig/project.yaml`, then `agnos publish`.

## Architecture

```
adapters/  -->  sandbox/  <--  cmd/main/        assets/ (templates, read via Deps.Embeddeps)
(reaches OS)    (closed)       (wires them)
```

- **`sandbox/`** — the closed core. It imports only `sandbox/` packages — the stdlib included,
  so text, sorting, hashing and templating come from `sandbox.Deps.<Contract>` too.
  `api/` holds contracts only, `deps/` holds dependency contracts (each `deps/<x>/` imports
  nothing at all; only the loose `deps/deps.go` names them),
  `binds/` holds one function file per `api/` file, `internal/` holds the logic.
  **Every function of `binds/` and `internal/` takes `sandbox *api.Sandbox` first**, and
  nothing else standing for the outside world: `api.Sandbox` carries `Deps`, so holding the
  api is holding everything — one part of the api can call another, and a field a caller
  replaced takes effect everywhere. `sandbox/api/` is the one place allowed to import
  `sandbox/deps`, and only for that field.
- **`adapters/`** — the only place OS-bound and third-party code lives. `libs/<adapter>/`
  exports `Bind(deps *deps.Deps)` beside an `adapter.yaml` naming the dep it fills;
  `availables/<name>/new.go` is generated from `availables/<name>/available.yaml`, never from a
  dir listing.
- **`assets/`** — every generated file's template. Groups `start`, `all`, `deps`, `cli`,
  `server`, `front` render `assets/<group>/<path>` to `<path>`; `deplist/<dep>/**` is one
  installable contract and `adapterlist/<adapter>/**` one installable adapter, each with its
  own declaration at the root of the group and installed nowhere (`dep.yaml`) or into the
  package (`adapter.yaml`); `templates/*` are single-file scaffolds rendered with
  `utils.RenderTemplateToDest`. A scaffold rendering to a file that is *itself* a template
  (`page_html.html`) escapes its own braces: `{{ "{{ .Title }}" }}`.
- **`cmd/main/`** — generated; wires an adapter into the sandbox, holds no logic.
- **`AgnosConfig/`** — written once by `start`, read by every `build`.
- **`examples/`** — one directory per example, on the `cli` and the `lib` side; `exec-test` runs
  each and diffs it against the `result.yaml` beside it.

Two layers per feature: an **action** (`sandbox/internal/actions/<name>/`) with `<name>.go`
(opens SmartIO, persists, runs the follow-up `build`) plus `<name>_internal.go` (pure logic on
an already-open SmartIO), and a **command**
(`sandbox/internal/commands/<name>/`) with `entries.yaml` (declaration), `entries.go`
(generated) and `handler.go` (hand-written). Both directories are snake_case for a kebab-case
command (`add-command` -> `add_command/`). Only `handler.go` and contract/adapter pairs are
written by hand; everything else is generated.

The **server layer** mirrors the cli layer file for file, and is the pattern to copy when a
layer is added: `serverdeps` mirrors `argvdeps`, `api/server.go` mirrors `api/cli.go`,
`internal/server/servermain.go` mirrors `internal/cli/climain.go`, `internal/routes/<name>/`
(`route.yaml` + generated `entries.go` + hand-written `handler.go` -> `RouteHandler`) mirrors
`internal/commands/<name>/`, `parsables/routeconf/` mirrors `commandconf/`, `assets/server/`
mirrors `assets/cli/`, and `server-init`/`server-purge` mirror `cli-init`/`cli-purge`. It is
rendered when `sandbox/internal/server/` exists, exactly as the cli group is rendered when
`sandbox/internal/cli/` does. The dispatch settles everything but the body (404/405/415/413/400)
before a handler runs; the body is read on demand by the generated `Entries.ReadBody`.
`sandbox/internal/routeio/` holds what both `servermain.go` and the routes need, because
`servermain.go` imports every route and a route may not import it back.

The **front layer** is the third column, and declares no unit of its own: a page **is** a
route. `sandbox/internal/pageio/` (generated; `Render` plus the `staticref`/`cssref`/`jsref`/
`dirref`/`inline`/`include` helpers) mirrors `routeio/` and is what `build` reads `hasFront`
from — never `assets/frontend/`, which is the project's own content and may be empty.
`assets/front/` is the group, `front-init`/`front-purge` the pair, and a page is
`routes/<page>/route.yaml` plus `assets/frontend/pages/<page>.html`, written by
`add-page`/`remove-page`. Both halves are written **once** and are then the project's; only
`pageio` and `docs/FrontUsage` are rewritten by every build. `pageio.StaticMount` is rendered
from the first segment of the `static` route's declaration (`collect_front_mount.go`), so
renaming the mount moves every generated link instead of breaking it in silence.

The **deps layer** is three units, not one: a **dep** is the contract (`sandbox/deps/<dep>/`,
one field of `deps.Deps`), an **adapter** is one implementation of it
(`adapters/libs/<adapter>/`, `adapter.yaml` says which dep), and an **available** is a selection
(`adapters/availables/<name>/available.yaml`, exactly one adapter per field). One dep may have
many adapters; `verify` demands every available fill every field exactly once — zero panics on
first use, two overwrite in silence. `add-dep`/`remove-dep` own the contract, `add-adapter`/
`remove-adapter` one implementation, `set-adapter` the choice, `add-available`/`remove-available`
the selection itself. `cmd/main/main.go` imports one available, and that import is how a program
picks its implementations. Another agnos repo installs as a dep too: `add-dep
<module>@<version> --as <name>` copies its `sandbox/api/` into `sandbox/deps/<name>/`
with the package clause changed and `Sandbox.Deps` dropped — a consumer installs
the api of a repo, never the wiring behind it, and that is what keeps the copy
self-contained (`apishape.DepsField`) — and **generates** `adapters/libs/<name>/`
— the shim that builds the remote sandbox from the remote repo's own adapters
and converts it, because
a top-level cast cannot work (Go's type identity is not recursive through named
types) and the conversion names a type whose import path lives in the consumer.
`sandbox/internal/apishape/` holds the convertibility rule and the converter
plan; `verify` applies that rule to every repo's `sandbox/api/` unconditionally,
which is what makes every agnos repo installable by construction.

**SmartIO** (`sandbox/internal/smartio/`) is a transactional filesystem rooted at `--path`.
Actions pass project-relative paths only; `Root` is joined at the `deps.Iodeps` boundary.
Writes buffer until `Persist`, but `List*` reads disk — so an action that runs `build` as a
follow-up must `Persist` first. Actions compose by sharing one open `*SmartIO` through their
`*Internal` function.

## Rules that generated code depends on

The full list is `docs/Rules/doc.md`, rendered from `assets/all/docs/Rules/doc.md`; add or
change a rule there and nowhere else. The ones most easily broken:

- **Every file is an instance of a pattern.** New code copies an existing sibling exactly:
  same filenames, same function names, same ordering. Never add a one-off — `verify` and the
  collectors read shape by convention.
- Naming is load-bearing: `sandbox.Deps.Iodeps` from `sandbox/deps/iodeps`,
  `Bind(deps *deps.Deps)` (an adapter still fills `deps.Deps` directly),
  `CommandHandler(sandbox *api.Sandbox, entries *Entries) int`.
- Generated files are never edited — change the template under `assets/` and bootstrap.
  `docs/GeneratedFiles/doc.md` lists which files are rewritten by every build.
- Never hand-edit a command's `entries.yaml`; use `add-flag` / `add-arg` / `set-command`. The
  same holds for a route's `route.yaml`: `add-route`/`remove-route`, `set-route`,
  `add-segment`/`remove-segment`, `add-header`/`remove-header`, `add-param`/`remove-param`,
  `set-body`, `add-body-field`/`remove-body-field` — one editor per place the file holds
  something, so no key of the declaration needs a hand edit.
  When declaring one of agnos's own flags, never pass a value that is exactly one of
  `add-flag`'s own spellings (`--identifier --example`) — the argv parser counts it as an
  occurrence and pollutes the declaration.
- A page is a route with `assets/frontend/pages/<page>.html` beside it — that file is the whole
  of what tells one from any other route. `add-page`/`remove-page` are its editors, and
  `remove-route` refuses a route that has one. `front-purge` drops the routes (their handlers
  import `pageio`) and never `assets/frontend/`.
- Every `identifier` of a route's `paths` starts with `/` and spells one segment; a captured
  segment is always `required: true`. `array: true` on the last capture makes it take every
  segment left in the path (`[]T`, one or more), and only there. Match order is by specificity
  (most identifiers, then longest, then fixed length before catch-all), settled by the
  collector, not the template.
- Every exported declaration of `sandbox/api/` and `sandbox/deps/` carries a doc comment —
  `docs/PublicApi/doc.md` is generated from those comments, and `verify` fails without them.
- `assets/deplist/<dep>/**` must render byte-for-byte to the copy this repo runs on.
- Output: `deps.Std.Printf` -> stdout, `deps.Std.Log` -> stderr (silenced by `--quiet`),
  `deps.Std.Error` -> stderr. Never `fmt.Printf`. A handler returns `api.ExitOk` or
  `api.ExitFailure`, never `api.ExitUsage` — the dispatch rejects bad input before it runs.
- An example is never created or deleted by hand, and `result.yaml` is never edited — refresh one
  golden with `update-test <name>`, the whole suite with `exec-test --update`, or delete it.
- Docs: create/delete with `agnos add-doc` / `agnos remove-doc`. `README.md`, every `Index.md`,
  and the docs of `assets/{all,cli}/docs/` are generated — `docs/Commands` from every command's
  `entries.yaml`, `docs/Structure/doc.md`'s tree from `AgnosConfig/structure.yaml`. Any new
  path worth naming gets an entry in `structure.yaml`, and `verify` fails on an entry whose
  path does not exist.
- A pattern changed here is mirrored in `docs/Contributing/doc.md` in the same commit, and the
  reverse.
