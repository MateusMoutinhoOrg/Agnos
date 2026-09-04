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

- **Every file is an instance of a pattern.** New code copies an existing sibling exactly:
  same filenames, same function names, same ordering. If no pattern fits, define and document
  the pattern first. Never add a one-off — `verify` and the collectors read shape by convention.
- Naming is load-bearing: a `Deps` field is the title-cased `sandbox/deps/<dir>` (`iodeps` ->
  `deps.Iodeps`); an adapter binder is always `Bind(deps *deps.Deps)` in
  `adapters/libs/<lib>/<lib>.go`; a command handler is always
  `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- Generated `.go` is written through `deps.Goimportsdeps.Format` (gofmt), so the tree is always
  formatted and a template's own indentation only has to be parsable.
- `assets/deplist/<dep>/**` must render byte-for-byte to the copy this repo runs on; `verify`
  compares them, so re-mirror whenever either side changes.
- Generated files (`(gen)` in `docs/Structure/doc.md`, the `always` rows of
  `docs/GeneratedFiles/doc.md`) are never edited — change the template under `assets/` and
  bootstrap.
- Never hand-edit a command's `entries.yaml`; use `add-flag` / `add-arg` / `set-command`
  (they re-render it alphabetically and drop comments).
- Output channels: `deps.Std.Printf` -> stdout (the result), `deps.Std.Log` -> stderr
  (progress, silenced by `--quiet`), `deps.Std.Error` -> stderr (failures). Never `fmt.Printf`.
- Exit codes: `api.ExitOk` 0, `api.ExitFailure` 1 (well-formed command failed),
  `api.ExitUsage` 2 (bad command line).
- Docs: create/delete with `agnos add-doc` / `agnos remove-doc`, never by hand. `README.md`
  (the documentation index included), `Index.md`, `docs/LibUsage/`, `docs/PublicApi/` and
  `docs/Structure/` are generated. A theme only groups a doc into a README section; a theme no doc names renders
  nothing and is not an error. PublicApi is rendered
  from the doc comments of `sandbox/api/` and `sandbox/deps/` (parsed through
  `deps.Goimportsdeps`), and `verify` requires one on every exported declaration there — so
  the public api is documented by commenting the contract, never by editing the page.
  `docs/Structure/doc.md`'s tree is rendered from `AgnosConfig/structure.yaml`, so an element
  is described by an entry there — and `verify` rejects an entry whose path is gone (a ghost
  spec).
  `docs/Commands/doc.md` stays hand-written (its columns carry judgement the declarations do
  not hold), but `verify` requires it to name every visible command and every declared flag —
  so `add-command` / `add-flag` are followed by a row there. Links inside `docs/` are relative
  to the doc; keep pages short — tables and commands over prose.
- A pattern changed here is mirrored in `docs/Contributing/doc.md` in the same commit, and the
  reverse.
