# Generated Files

## Description
Every file `agnos` writes into a project, grouped by the command that first writes it, with the asset group it comes from and whether later builds overwrite it. A file marked **overwritten** is regenerated on every `agnos build` and must never be edited by hand; a file marked **once** is written on first render and then belongs to the project. The rendering order is [BuildPipeline](/docs/BuildPipeline/doc.md); the templates live under `assets/` — see [Structure](/docs/Structure/doc.md#assets).

---

## Written by `start`

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `AgnosConfig/project.yaml` | `start` | once | `name`, `version`, `description` of the project — read back by every build to regenerate `config.go`. |
| `AgnosConfig/themes.yaml` | `start` | once | The documentation themes of the project. Parsed by `themesconf`; every build renders one `docs/Index/<id>.md` per entry and one README row. |
| `AgnosConfig/ignore.yaml` | `start` | once | Paths SmartIO hides from listings. Empty by default. |
| `AgnosConfig/paths.yaml` | `start` | once | Path rewrites SmartIO applies to listings. Empty by default. |
| `AgnosConfig/docs/ReadmeHeader.md` | `start` | once | The body of the project's `README.md`, itself a Go `text/template` — every `build` renders it (with the full `vars` map) into `README.md`. Edit this, not `README.md`. |
| `go.mod` | — | once | Written from `--module` and the installed toolchain's version when the directory has none (`--force` overwrites). `dep-install` / `dep-remove` edit its `require` block. |

---

## Written by every `build`

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `sandbox/new.go` | `all` | overwritten | `New(deps *deps.Deps) *api.Sandbox` — or `New()` without deps — calling one `binds.<X>Bind` per file of `sandbox/binds/`. |
| `sandbox/api/sandbox.go` | `all` | overwritten | `type Sandbox struct` with one field per package of `sandbox/api/`. |
| `sandbox/internal/config/config.go` | `all` | overwritten | `ProjectName` (title-cased `name`) and `Version` from `project.yaml`. |
| `README.md` | `all` | overwritten | `{{ render "<ProjectName>Config/docs/ReadmeHeader.md" }}` then `{{ copy "LICENSE" }}` — the `all` template renders `ReadmeHeader.md` as a template and appends the project's `LICENSE` verbatim (`copy` native function). Change the text in `ReadmeHeader.md`. |
| `docs/Index/<theme-id>.md` | `templates/theme_index.md` | overwritten | One index per theme, listing the first-level docs whose `props.yaml` names it. `docs/Index/` is removed and rewritten whole, so a deleted theme's index never survives. Written only when `docs/` exists. |
| `docs/**/Index.md` | `templates/doc_index.md` | overwritten | One index per doc directory that has sub-docs, listing its direct sub-docs. |
| `go.sum` | — | — | Written by the `go mod tidy` the `go` runtime runs. |

---

## Written when `sandbox/deps/` exists

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `sandbox/deps/deps.go` | `deps` | overwritten | `type Deps struct` with one `<Title> <dir>.Lib` field per directory of `sandbox/deps/`. |
| `adapters/availables/standard/new.go` | `deps` | overwritten | `standard.New() deps.Deps` calling `<lib>.Bind(&deps)` for every directory of `adapters/libs/`. |

---

## Written by `dep-install <dep>`

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `sandbox/deps/<dep>/<dep>.go` | `deplist/<dep>` | once | The contract — a `Lib` struct of function fields. Owned by the project after install; `dep-remove` deletes it. |
| `adapters/libs/<lib>/<File>.go` | `deplist/<dep>` | once | The implementation exporting `Bind(deps *deps.Deps)`. |
| `assets/asset.go` | `deplist/embeddeps` | once | The `//go:embed all:*` directive and the `Files` filesystem; only the `embeddeps` dep carries it. |

---

## Written when `sandbox/internal/cli/` exists

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `cmd/main/main.go` | `cli` | overwritten | Wires `standard.New()` into `sandbox.New` and exits with `lib.Cli.CliMain(os.Args[1:])`. |
| `sandbox/api/cli.go` | `cli` | overwritten | `type Cli struct { CliMain func([]string) int }` and the exit-code constants. |
| `sandbox/binds/cli.go` | `cli` | overwritten | `CliBind`, assigning the generated `cli.CliMain`. |
| `sandbox/internal/cli/climain.go` | `cli` | overwritten | `CliMain` and one `dispatch<Name>` per command, generated from every `entries.yaml`. |
| `sandbox/internal/commands/help/entries.yaml` | `templates/help_entries.yaml` | once | The `help` command's declaration — created if missing, then editable like any other. |
| `sandbox/internal/commands/help/entries.go` | `templates/entries.go` | overwritten | Generated like every command's. |
| `sandbox/internal/commands/help/handler.go` | `cli` | overwritten | The help screens, with every command's metadata baked into its `helpCommands` table. |
| `sandbox/internal/commands/version/entries.yaml` | `cli` | overwritten | The `version` command's declaration. |
| `sandbox/internal/commands/version/handler.go` | `cli` | overwritten | Prints `config.Version`. |
| `sandbox/internal/commands/<name>/entries.go` | `templates/entries.go` | overwritten | One typed `Entries` struct per command, from its `entries.yaml`. |

`cli-init` writes this group first, and every later `build` rewrites it. `cli-purge` deletes the group's files and drops `sandbox/internal/cli/` and `sandbox/internal/commands/` whole.

---

## Written by `add-command <name>`

| File | Group | Rewritten by `build` | Description |
|------|-------|----------------------|-------------|
| `sandbox/internal/commands/<name>/entries.yaml` | `templates/command_entries.yaml` | once | `identifiers`, `category`, `help` — grown afterwards by `add-flag`, `add-arg`, `set-command`, which rewrite it in canonical (alphabetical) key order. |
| `sandbox/internal/commands/<name>/handler.go` | `templates/command_handler.go` | once | A stub `CommandHandler` printing `<name> called`. The project's to fill. |
