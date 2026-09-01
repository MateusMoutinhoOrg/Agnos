# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Agnos-Cli (`agnos`) is a Go CLI that **scaffolds and regenerates other Go applications**
built around a command/action architecture. It began as a Go CLI template but is being
refactored (branch `Cli-Refactor`) into a "factory builder": `agnos start` creates a new
project skeleton, and `agnos build` (re)renders parts of that project from embedded
templates according to the project's `*Config/` files. See `AGENTS.md` for the intent.

`old/` (previous template + full docs under `old/docs/`) and `local_test/` (a generated
sample project) are reference material, not part of the build. `main` is a committed
prebuilt binary.

## North star: self-hosting / bootstrap

The long-term goal is for **`agnos` to build `agnos` itself** — this repo becomes one of
the projects that `agnos start` scaffolds and `agnos build` regenerates. Treat that as the
design constraint behind every change, even while it is still aspirational:

- **Everything is a pattern.** Any new file, package, or operation must follow an existing
  repeatable shape (see the "fixed shape" of `parsables/`, the two-layer
  command/action split, the `NewCommand` contract). If no pattern fits, first define the
  pattern — document it here — *then* implement against it. Never add a one-off.
- **Uniformity over cleverness.** When two packages do a similar job they should be
  structurally identical (same filenames, same function names, same ordering). A reader or
  a template renderer should be able to predict a package's contents from its role.
- **Anything hand-written today is a template tomorrow.** Prefer mechanical, regular code
  that a Go `text/template` in `assets/sandbox/**` could emit. Avoid constructs that would
  be awkward to generate (irregular formatting, ad-hoc special cases, hidden state).
- **The sandbox rule is absolute.** No OS effects outside `adapters/`. Self-hosting depends
  on the core staying pure so it can be rendered and re-rendered safely.
- **Keep this file and `AGENTS.md` current.** They are the harness's spec. When you
  establish or change a pattern, update the relevant section in the same change.

This is a harness: consistency, predictability, and documented conventions matter more than
any individual feature.

## Commands

```bash
# Build the CLI for the current machine (macOS Intel example) into release/
bash ./scripts/mac86.sh          # other targets: macarm64, linux86, linuxarm64, linuxi32, windows86, windowsi32
bash ./scripts/all.sh            # every target into release/

go build -C . -o release/agnos ./cmd/main   # quick local build
go vet ./...

# End-to-end smoke test: builds, installs to /usr/local/bin/agnos, scaffolds a fresh project
sh ./local_test.sh
```

There is currently **no test suite** (`*_test.go` files) and no lint config beyond `go vet`.
`go` 1.25.0 is required (see `go.mod`).

### Running the CLI

```bash
agnos start --project-name <name> --module <go-module-path> [path]   # scaffold; path defaults to "."
agnos enable-deps [path]      # add sandbox/deps + adapters, then rebuild
agnos remove-deps [path]      # remove them, then rebuild
agnos build [path]            # re-render generated files from templates
agnos help | version
```

## Architecture

### Sandbox / adapter / cmd split (dependency injection)

- **`sandbox/`** — the closed core. It performs **no OS effects directly**. Everything
  external (filesystem, time, stdout/stderr, HTTP, embedded assets, YAML) arrives through
  `sandbox/deps.Deps`, a struct of **only** sub-contract structs — no bare function
  fields (`std`, `iodeps`, `argvdeps`, `dbdeps`, `embeddeps`, `requestdeps`,
  `serializebles`). Each sub-contract is *restated* inside `sandbox/deps/<x>/` precisely
  because its real implementation must live outside the sandbox. The loose runtime
  functions (`Now`, `Printf`, `Error`, `Errorf`) are gathered into `sandbox/deps/std.Lib`,
  injected as `Deps.Std`. Sub-contracts whose real object is created per call expose that
  constructor as a field on their own `Lib` struct — `requestdeps.Lib.NewRequest(url)`,
  `dbdeps.Lib.NewDatabase(props)` (rooted at `props.Path`), `argvdeps.Lib.New(args)`
  (which hands back an `argvdeps.Parser` bound to that argv) — injected whole as
  `Deps.RequestLib` / `Deps.DatabaseLib` / `Deps.ArgvLib`, exactly like `IoLib` and
  `SerializeLib`. There is no separate `Factory` struct: `deps.Deps` references each
  `<x>.Lib` directly.
- **`adapters/standard/`** — fills `Deps` with real implementations (`standard.New()`).
  The only place `os`, `embed`, `net/http` etc. are touched.
- **`assets/`** — files embedded into the binary via `//go:embed all:*` in `assets/asset.go`.
  The embed directive must live in this package. Reached only through `deps.EmbedDeps`.
  `assets/sandbox/**` are Go **text/templates** (`{{.Module}}`, `{{if .HasDeps}}`) that
  `agnos build` renders into a target project.
- **`cmd/main/main.go`** — wires `standard.New()` into `sandbox.New(&deps)` and calls
  `lib.Cli.CliMain(os.Args[1:])`. Holds no logic.

### Sandbox internals

`sandbox.New` runs two binders that populate an `api.Sandbox`:

- **`binds/cli.go`** registers `api.CliCommand` values from `sandbox/internal/commands/<name>/`
  and sets `Cli.CliMain` to `sandbox/internal/cli.CliMain` — a declarative parser: each
  command declares `Args`/`Flags` (`api.CliArg`/`api.Cliflag`), the parser collects and
  validates them from the argv via `argvdeps`, then calls `Handler(deps, entries)` which
  returns an exit code (`api.ExitOk` / `ExitUsage` / `ExitFailure`).
- **`binds/actions.go`** registers the reusable operations in `api.Sandbox.Actions`
  (`Build`, `Start`, `EnableDeps`, `RemoveDeps`), each from `sandbox/internal/actions/<name>/`.

**Two layers per operation:** `internal/commands/<name>` is the CLI surface (flag/arg
declaration + handler); `internal/actions/<name>` is the logic. Command handlers currently
call the action package directly (e.g. `start` calls `startAction.Start` then
`buildAction.Build`); `binds/actions.go` also exposes the same actions as library API.
Note `build` runs as a follow-up step after `start`/`enable-deps`/`remove-deps`.

### SmartIO — transactional filesystem

`sandbox/internal/smartio` wraps `deps.IoLib` with an in-memory transaction layer. Actions
call `io.WriteFile` / `io.CreateDir` / `io.RemoveDir` etc.; nothing hits disk until
`io.Persist()` (removes pending dirs, then creates pending dirs, then flushes buffered file
writes). Reads and existence checks are transaction-aware (pending creations count as
existing, pending removals as gone). `WriteFile` refuses to overwrite; `WriteFileOverwrite`
replaces. `smartio.New` also loads `<ProjectName>Config/ignore.yaml` and `paths.yaml` to
filter/rewrite paths during listing.

### Config files (`parsables`)

Each `sandbox/internal/parsables/<name>conf/` package is a small YAML-or-gomod parser with a
fixed shape: `api.go` (struct + `Render func() string`), `new.go` (parse from string),
`new_empty.go` (defaults), `bind_methods.go`, `render.go`. `moduleconf` parses `go.mod`
directly. `agnos start` writes `project.yaml`, `themes.yaml`, `ignore.yaml`, `paths.yaml`
into `<ProjectName>Config/`; `agnos build` reads them back.

`sandbox/internal/config/config.go` holds the CLI's own constants: `ProjectName = "Agnos"`
(used to name the generated `AgnosConfig/` dir) and `Version`.

## Adding a command

1. `sandbox/internal/commands/<name>/<name>.go` — `NewCommand(deps, sandbox) api.CliCommand`
   with `ValidStartIdentifiers`, `Category`, `Args`, `Flags`, `Description`, `Handler`.
2. If it needs logic reusable as library API, add `sandbox/internal/actions/<name>/` and
   register it in `sandbox/binds/actions.go` + `sandbox/api/actions.go`.
3. Append it in `sandbox/binds/cli.go`.
