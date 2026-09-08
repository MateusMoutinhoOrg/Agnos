# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Agnos (`agnos`) is a Go CLI that **scaffolds and regenerates other Go CLIs**. `agnos start`
writes a project skeleton; `agnos build` re-renders every generated file from `text/template`
assets embedded in the binary; commands like `add-command`, `add-flag`, `dep-install` declare
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
are rendered from `assets/all/docs/` into **every** agnos project, this one included, and
`docs/{CliInstall,CliExamples}/` from `assets/cli/docs/`: editing one means editing that
template, and it has to read correctly in a scaffolded project, not only here. Guard a line
that holds for this repo alone with `{{ if .HasAssets }}`.

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
  so text, sorting, hashing and templating come from `deps.<Contract>` too.
  `api/` holds contracts only, `deps/` holds dependency contracts (each `deps/<x>/` imports
  nothing at all; only the loose `deps/deps.go` names them),
  `binds/` holds one function file per `api/` file, `internal/` holds the logic.
- **`adapters/`** — the only place OS-bound and third-party code lives. `libs/<x>/` exports
  `Bind(deps *deps.Deps)`; `availables/standard/new.go` is generated from that dir listing.
- **`assets/`** — every generated file's template. Groups `start`, `all`, `deps`, `cli` render
  `assets/<group>/<path>` to `<path>`; `deplist/<dep>/**` is one installable dep;
  `templates/*` are single-file scaffolds rendered with `utils.RenderTemplateToDest`.
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
- Naming is load-bearing: `deps.Iodeps` from `sandbox/deps/iodeps`, `Bind(deps *deps.Deps)`,
  `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- Generated files are never edited — change the template under `assets/` and bootstrap.
  `docs/GeneratedFiles/doc.md` lists which files are rewritten by every build.
- Never hand-edit a command's `entries.yaml`; use `add-flag` / `add-arg` / `set-command`.
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
