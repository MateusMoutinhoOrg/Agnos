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
  command/action split, the `entries.yaml` + `handler.go` command contract). If no pattern fits, first define the
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
agnos start --project-name <name> --module <go-module-path> [--path <dir>]   # scaffold; --path defaults to "."
agnos deps-init [--path <dir>]        # dependency subsystem: add sandbox/deps + adapters, then rebuild
agnos deps-purge [--path <dir>]       # dependency subsystem: remove them, then rebuild
agnos dep-install <dep> [--path <dir>]  # render assets/deplist/<dep>/** into the project, then rebuild
agnos dep-remove <dep> [--path <dir>]   # remove what that dep installed (and now-empty dirs), then rebuild
agnos dep-list [--path <dir>]         # list the dep names available under assets/deplist
agnos cli-init [--path <dir>]         # cli subsystem: install the std + verb deps, render the cli asset group, then rebuild
agnos cli-purge [--path <dir>]        # cli subsystem: remove every file the cli asset group installs, then rebuild
agnos verify [--path <dir>]           # check the project keeps the sandbox/adapter schema (no writes)
agnos build [--path <dir>] [--unsafe] # run verify (unless --unsafe), then re-render generated files from templates
agnos help | version
```

Every command that operates on a target project takes the target directory as the
**`--path` flag** (never a positional arg), defaulting to `.`. In the generated dispatch a
value flag with a `default:` is assigned that literal when absent, so the corresponding
`entries.Path` field is always populated. (Agnos-Cli's own commands still use the previous
`entries.GetFlagById("path")` accessor — see the migration note below.)

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
  `Deps.Requestdeps` / `Deps.Dbdeps` / `Deps.Argvdeps`, exactly like `Deps.Iodeps` and
  `Deps.Serializebles`. **Every `Deps` field name is mechanical: it is the title-cased
  sub-contract directory name** (`sandbox/deps/iodeps/` → `Deps.Iodeps`,
  `sandbox/deps/serializebles/` → `Deps.Serializebles`, …) — never a hand-picked alias —
  because `render_deps.go` regenerates `sandbox/deps/deps.go` by iterating those dirs. There
  is no separate `Factory` struct: `deps.Deps` references each `<x>.Lib` directly.
- **`adapters/`** — the only place `os`, `embed`, `net/http` etc. are touched. Split in two:
  - **`adapters/libs/<x>deps/`** — one isolated real implementation per sub-contract
    (`argvdeps`, `dbdeps`, `embeddeps`, `iodeps`, `requestdeps`, `serializebles`, `std`),
    each a package exporting a uniform `Bind(deps *deps.Deps)` that fills that one `deps`
    sub-contract field.
  - **`adapters/availables/<name>/new.go`** — a ready-made `Deps` assembly wiring the libs
    together; `adapters/availables/standard` (`standard.New()`) is the default. A user who
    wants a different mix composes their own `Deps` from `adapters/libs` directly.
- **`assets/`** — files embedded into the binary via `//go:embed all:*` in `assets/asset.go`.
  The embed directive must live in this package. Reached only through `deps.Embeddeps`.
  `assets/sandbox/**` are Go **text/templates** (`{{.Module}}`, `{{if .HasDeps}}`) that
  `agnos build` renders into a target project.
- **`cmd/main/main.go`** — wires `standard.New()` into `sandbox.New(&deps)` and calls
  `lib.Cli.CliMain(os.Args[1:])`. Holds no logic.

### Sandbox internals

`sandbox.New` runs two binders that populate an `api.Sandbox`:

- **`binds/cli.go`** only sets `Cli.CliMain` to the generated `sandbox/internal/cli.CliMain`.
  There is no command registry: `api.Cli` holds a single `CliMain func([]string) int`.
  Each command lives in `sandbox/internal/commands/<name>/` as three files — `entries.yaml`
  (hand-written declaration: `identifiers`, `category`, `help`, `long-description`,
  `examples`, `hidden`, `flags:`, `args:`), `entries.go` (the typed struct **generated** by
  `agnos build` from `entries.yaml`), and `handler.go` (hand-written:
  `func CommandHander(deps *deps.Deps, entries *Entries) int`). `agnos build` parses every
  `entries.yaml` (`parsables/commandconf`), regenerates each `entries.go`, then regenerates
  `sandbox/internal/cli/climain.go` with one `dispatch<Name>` arm per command that reads the
  argv via `argvdeps`, fills that command's `Entries`, and calls its `CommandHander`. A flag
  or arg with `array: true` becomes a `[]T` field collecting every occurrence. Exit codes
  (`ExitOk`/`ExitUsage`/`ExitFailure`) are consts in both `sandbox/api` and the generated
  `cli` package.
- **`binds/actions.go`** registers the reusable operations in `api.Sandbox.Actions`
  (`Build`, `Verify`, `Start`, `DepsInit`, `DepsPurge`, `DepInstall`, `DepRemove`, `DepList`,
  `CliInit`, `CliPurge`), each from
  `sandbox/internal/actions/<name>/`.

> **Migration note:** the generator (build action + `assets/cli/**`) emits this new
> three-file command shape, but Agnos-Cli's own `sandbox/internal/commands/**` and
> `sandbox/api/cli.go` still use the previous `api.CliCommand` / `NewCommand` registry —
> the self-hosting bootstrap has not been re-run against it yet.

**Two layers per operation:** `internal/commands/<name>` is the CLI surface
(`entries.yaml` declaration + `handler.go`); `internal/actions/<name>` is the logic. Command handlers currently
call the action package directly (e.g. `start` calls `startAction.Start` then
`buildAction.Build`); `binds/actions.go` also exposes the same actions as library API.
Note `build` runs as a follow-up step after `start`/`deps-init`/`deps-purge`.

**The verify gate.** `internal/actions/verify` reads the target tree through `smartio` and
performs **no writes** (it never calls `io.Persist`). It enforces the harness schema and
returns one error listing every violation:

- `sandbox/` holds only the `api`, `binds`, `deps`, `internal` directories plus a loose
  `new.go`.
- No file under `sandbox/` imports a module-internal package outside `sandbox/`.
- `sandbox/api/*` imports only other `sandbox/api` packages — nothing else at all (no
  stdlib, no external modules, not even `sandbox/deps`: api is pure contract). In the new
  command shape `api/cli.go` is just `type Cli struct { CliMain func([]string) int }` plus
  the exit-code consts; all parsing lives in the generated `sandbox/internal/cli`.
  `sandbox/deps/*` imports only the standard library and other `sandbox/deps` packages.
- Every file in `sandbox/binds/` mirrors a file of the same name in `sandbox/api/` and
  declares only functions (no top-level types/consts/vars).
- `adapters/` holds only the `availables` and `libs` directories.

The `build` **command** handler runs `verifyAction.Verify` before `buildAction.Build` unless
`--unsafe` is passed; the `build` **action** (and the follow-up `build` after
`start`/`deps-*`/`cli-*`) does not, so mid-refactor states can still regenerate. Checks live
in `check_sandbox.go` / `check_adapters.go`, each returning `[]string` violations; add new
rules as sibling `check_*.go` files.

**Asset groups.** Templates live under `assets/<group>/**`, each file at the path it will be
written to inside a target project (`assets/all/sandbox/new.go` → `sandbox/new.go`).
`utils.RenderGroup(deps, io, "<group>", vars)` renders every file in one group as a Go
`text/template` with the same `vars` and writes each to its stripped path (via
`io.WriteFileOverwrite`). The groups: `start` (config skeleton written by `agnos start`),
`all` (always rendered by `agnos build`), `deps` (rendered by `agnos build` only when
`sandbox/deps` exists), `cli` (the CLI layer — `cmd/main/main.go`, `sandbox/api/cli.go`, `sandbox/binds/cli.go`,
`sandbox/internal/cli/climain.go` (generated from every command's `entries.yaml`), the
generated `help` command, and the `version` command as `entries.yaml` + `handler.go` —
rendered by `agnos cli-init` and every later `agnos build` when `sandbox/internal/cli`
exists, deleted by `agnos cli-purge`). `utils.RenderTemplateToDest` renders one file to one
dest — it also emits each command's `sandbox/internal/commands/<name>/entries.go` from
`assets/templates/entries.go` (see `build/generate_command_entries.go`).

**Deps (`assets/deplist/<dep>/**`).** Each sub-directory of `assets/deplist/` is one
installable dep; the tree under it mirrors the target-project layout (an asset at
`assets/deplist/embed/sandbox/deps/embeddeps/embeddeps.go` installs to
`sandbox/deps/embeddeps/embeddeps.go`). `dep-install` renders that subtree as one asset
group via `utils.RenderGroup(deps, io, "deplist/"+dep, {"Module": …})`, then runs `build`
as a follow-up step so the collectors pick up the new sub-contract dirs; `dep-remove`
deletes those files plus any directory the removal emptied, then runs `build`. Both also
sync the target project's `go.mod`: `assets/depsversion.yaml` (parsed by
`parsables/depsversionconf`) maps a dep name to the `<module>@<version>` it needs, and if
the dep is listed there `dep-install` adds/updates that `require` entry (via
`moduleconf.AddRequire` + `moduleconf.Render`, written back to `<path>/go.mod`) while
`dep-remove` strips it (`moduleconf.RemoveRequire`). Deps that bundle only sandbox-copy
code are absent from `depsversion.yaml` and leave `go.mod` untouched. Both persist
their own filesystem changes *before* invoking `buildAction.Build` (a fresh transaction),
because `build`'s collectors list dirs from disk and would not see pending writes.
`dep-list` reads `deplist/` and returns one entry per first path segment.

Every lib this repo ships is mirrored there as a dep, each bundling its
`sandbox/deps/<contract>/` copy with its `adapters/libs/<lib>/` implementation and
`{{.Module}}` substituted for the import path: `embed` (embeddeps + `assets/asset.go`),
`iodeps`, `keep` (dbdeps), `requestdeps`, `serializebles`, `std`, `verb` (argvdeps). After
`dep-install`, `build`'s collectors regenerate `sandbox/deps/deps.go` and
`adapters/availables/standard/new.go` to include the new sub-contract; `keep` and `verb`
pull external modules, so `dep-install` writes their pinned `require` into the target
`go.mod` (from `assets/depsversion.yaml`) and the user still runs `go mod tidy` there to
download them.

**Collectors (`internal/actions/build/collect_*.go`)** all follow one shape: list a directory
with `io.List*`, derive names from the last path segment, title-case them, return the slice.
`BuildInternal` gathers them into one `vars` map handed to both `RenderGroup` calls:
`CollectBinds` iterates `sandbox/binds` into `{{range .Binds}}` (`sandbox/new.go`);
`CollectConstructors` iterates `sandbox/api` into `{{range .Constructors}}`
(`sandbox/api/sandbox.go`); `CollectDepsLibs` iterates the `sandbox/deps/<x>/` sub-contract
dirs, emitting one `{{.Title}} {{.Name}}.Lib` field (and its import) per dir in
`{{range .DepsLibs}}` (`sandbox/deps/deps.go`); `CollectAdapterLibs` iterates the
`adapters/libs/<x>/` dirs, emitting one `{{.Name}}.Bind(&deps)` call (and its import) per lib
in `{{range .AdapterLibs}}` (`adapters/availables/standard/new.go`). Every `adapters/libs/<x>`
package therefore exposes its binder under the single uniform name `Bind(deps *deps.Deps)`.
`CollectCommands` is the exception to the one-line shape: it reads each
`sandbox/internal/commands/<x>/entries.yaml` through `parsables/commandconf` and returns a
rich `map[string]any` per command (identifiers, category, help text, and a precomputed
`Flags`/`Args` list carrying Go field names, types, getter names and default literals) for
the `{{range .Commands}}` loops in the generated `sandbox/internal/cli/climain.go` and
`help` package. The `help` directory is skipped — it is generated, not declared. Every
command package therefore exposes its handler under the single uniform name
`CommandHander(deps *deps.Deps, entries *Entries) int`.

`BuildInternal` also sets `HasCli` (`sandbox/internal/cli` exists) alongside `HasDeps`;
when true it calls `GenerateCommandEntries` (one `entries.go` per command) and renders the
`cli` group — so once `cli-init` has run, every later `agnos build` regenerates
`climain.go` and the `help` package from the current `entries.yaml` set.

**This repo bootstraps itself.** `agnos build .` runs against Agnos-Cli's own tree
(`AgnosConfig/` holds its `project.yaml` / `themes.yaml` / `ignore.yaml` / `paths.yaml`) and
regenerates `sandbox/deps/deps.go`, `adapters/availables/standard/new.go` and
`sandbox/binds/cli.go` in place, and the result still compiles (`go build ./cmd/... ./sandbox/... ./adapters/...`). Keep it that
way: `agnos build .` must stay idempotent and compilable. In particular, all sandbox and
adapter code must reference `Deps` fields by their mechanical names (`deps.Iodeps`,
`deps.Serializebles`, `deps.Argvdeps`, `deps.Dbdeps`, `deps.Requestdeps`, `deps.Embeddeps`,
`deps.Std`) so the regenerated struct keeps matching its callers. After changing anything
under `sandbox/deps/<x>/`, `adapters/libs/<x>/`, `sandbox/binds/`, or `sandbox/api/`, run
`agnos build .` and rebuild.

### SmartIO — transactional filesystem

`sandbox/internal/smartio` wraps `deps.Iodeps` with an in-memory transaction layer. Actions
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
directly. `agnos start` renders the `start` asset group (`project.yaml`, `themes.yaml`, `ignore.yaml`,
`paths.yaml`) into `AgnosConfig/`; `agnos build` reads them back.

`sandbox/internal/config/config.go` holds the CLI's own constants: `ProjectName = "Agnos"`
(used to name the generated `AgnosConfig/` dir) and `Version`. It is itself part of the
`all` asset group (`assets/all/sandbox/internal/config/config.go`): every `agnos build`
reads `<ProjectName>Config/project.yaml` back through `projectconf` and regenerates it, with
`ProjectName` set to the title-cased `name` and `Version` to the config's `version`. A
missing/unparsable `project.yaml` is a hard error — `agnos start` is a prerequisite for
`agnos build`, so `loadProjectConf` never falls back to defaults. For this repo the
rendered result is unchanged (`Agnos` / `0.0.1`), keeping `agnos build .` idempotent.

## Adding a command (generated-project contract)

1. `sandbox/internal/commands/<name>/entries.yaml` — declare `identifiers`, `category`,
   `help`, optional `long-description`, `examples`, `hidden`, and `flags:` / `args:` maps
   (each entry: `identifiers`, `description`, `type` (`string`/`boolean`/`int`/`float`),
   optional `default`, `required`, `array`).
2. `sandbox/internal/commands/<name>/handler.go` — `package <name>`, one function
   `func CommandHander(deps *deps.Deps, entries *Entries) int`.
3. Run `agnos build` — it generates `entries.go` and re-wires `climain.go` + `help`.
   Nothing to append anywhere.
4. If it needs logic reusable as library API, add `sandbox/internal/actions/<name>/` and
   register it in `sandbox/binds/actions.go` + `sandbox/api/actions.go`.
