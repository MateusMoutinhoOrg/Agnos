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
`docs/Structure/doc.md` (schema), `docs/BuildPipeline/doc.md` (what `build` does), and
`docs/Contributing/doc.md` (recipes for every kind of change).

## Audience: LLMs, not humans

The primary reader and writer of this repo is an LLM. Every choice — code shape, docs,
naming, file layout — optimizes for machine reading and machine writing, and token cost is a
first-class constraint. Humans are a secondary audience.

- **Generate over hand-write.** If a file can be rendered from a template, a collector or a
  declaration, it must be — hand-written code is only contracts, adapters and `handler.go`.
  A new hand-written file needs a reason why generation cannot cover it.
- **Generate over document.** The same applies to docs: `README.md` (its documentation index
  included), every `Index.md`, `docs/LibUsage/`, `docs/PublicApi/` and `docs/Structure/` are
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
./release/bootstrap.bin build                                    # verify + regenerate + compile
./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
./release/bootstrap.bin local-install                            # install the result
```

Compile scope is always `./cmd/... ./sandbox/... ./adapters/...` — **never `go build ./...`**,
because `assets/` holds Go templates, not compilable Go.

There is no test suite. `verify` (schema check, no writes) plus a compiling, idempotent
`build` is the check that a change is correct. Release: bump `version` in
`AgnosConfig/project.yaml`, then `agnos publish`.

## Architecture

```
adapters/  -->  sandbox/  <--  cmd/main/        assets/ (templates, read via Deps.Embeddeps)
(reaches OS)    (closed)       (wires them)
```

- **`sandbox/`** — the closed core. Imports nothing outside `sandbox/` and no OS packages.
  `api/` holds contracts only, `deps/` holds dependency contracts (stdlib imports only),
  `binds/` holds one function file per `api/` file, `internal/` holds the logic.
- **`adapters/`** — the only place OS-bound and third-party code lives. `libs/<x>/` exports
  `Bind(deps *deps.Deps)`; `availables/standard/new.go` is generated from that dir listing.
- **`assets/`** — every generated file's template. Groups `start`, `all`, `deps`, `cli` render
  `assets/<group>/<path>` to `<path>`; `deplist/<dep>/**` is one installable dep;
  `templates/*` are single-file scaffolds.
- **`cmd/main/`** — generated; wires an adapter into the sandbox, holds no logic.
- **`AgnosConfig/`** — written once by `start`, read by every `build`.

Two layers per feature: an **action** (`sandbox/internal/actions/<name>/`) with `<name>.go`
(opens SmartIO, persists, runs the follow-up `build`) plus `<name>_internal.go` (pure logic on
an already-open SmartIO), and a **command**
(`sandbox/internal/commands/<name>/`) with `entries.yaml` (declaration), `entries.go`
(generated) and `handler.go` (hand-written). Only `handler.go` and contract/adapter pairs are
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
- Never hand-edit a command's `entries.yaml`; use `add-flag` / `add-arg` / `set-command`.
- `assets/deplist/<dep>/**` must render byte-for-byte to the copy this repo runs on.
- Output: `deps.Std.Printf` -> stdout, `deps.Std.Log` -> stderr (silenced by `--quiet`),
  `deps.Std.Error` -> stderr. Never `fmt.Printf`.
- Docs: create/delete with `agnos add-doc` / `agnos remove-doc`. `README.md`, every `Index.md`,
  and `docs/{LibUsage,PublicApi,Structure,Rules}/` are generated. `docs/Commands/doc.md` stays
  hand-written and `verify` requires it to name every visible command and declared flag.
  `docs/Structure/doc.md`'s tree comes from `AgnosConfig/structure.yaml`.
- A pattern changed here is mirrored in `docs/Contributing/doc.md` in the same commit, and the
  reverse.
