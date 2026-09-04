# Rules

Every rule of this project, in one page. `verify` enforces the ones marked **(verify)**;
the rest are read by the generators or by whoever writes the hand-written files.
Nothing here is repeated elsewhere in `docs/` — other pages link here.

## Authoring

- **Generate over hand-write.** A file that can be rendered from a template, a collector or a
  declaration must be. Hand-written code is contracts, adapters and `handler.go` only; a new
  hand-written file needs a reason why generation cannot cover it.
- **Every file is an instance of a pattern.** New code copies an existing sibling exactly:
  same filenames, same function names, same ordering. If no pattern fits, define and document
  the pattern first — `verify` and the collectors read shape by convention, so a one-off
  breaks them.
- **Deterministic and idempotent.** Same input, same bytes out: `build` run twice must leave
  the tree unchanged.
- Generated `.go` is written through `deps.Goimportsdeps.Format` (`gofmt`), so a template's own
  indentation only has to be parsable.
- Generated files — `(gen)` in [Structure](../Structure/doc.md), the `always` rows of the
  generated-file listing — are never edited. Change the template under `assets/` and rebuild.
- Compile scope is always `./cmd/... ./sandbox/... ./adapters/...`, never `./...`: `assets/`
  holds Go templates, not compilable Go.

## Layers

- `sandbox/` is closed: it imports nothing outside `sandbox/` and no OS package. **(verify)**
- `sandbox/` holds only `api`, `binds`, `deps`, `internal` and `new.go`. **(verify)**
- `sandbox/api/*` imports only `sandbox/api`. **(verify)**
- `sandbox/deps/*` imports only the stdlib and `sandbox/deps`. **(verify)**
- Every `sandbox/binds/` file mirrors one `api/` file and declares only functions. **(verify)**
- Every file of `sandbox/api/` and `sandbox/deps/` parses, and every exported type, func, const
  and var in them carries a doc comment — [PublicApi](../PublicApi/doc.md) is generated from
  those comments. **(verify)**
- `adapters/` is the only place OS-bound and third-party code lives, and holds only
  `availables` and `libs`. **(verify)**
- Every `adapters/libs/<lib>/` exports `Bind(deps *deps.Deps)`, and every `sandbox/deps/<x>`
  contract has a lib filling its `deps.<X>` field. **(verify)**
- `cmd/main/` wires an adapter into the sandbox and holds no logic.
- Every `assets/deplist/<dep>/<path>`, rendered with this module, equals `<path>` whenever that
  file exists here — re-mirror whenever either side changes. **(verify)**

## Naming

- A `Deps` field is the title-cased `sandbox/deps/<dir>` (`iodeps` -> `deps.Iodeps`). Always
  use that spelling; an added contract never renames an existing one.
- An adapter lib's binder is always `Bind(deps *deps.Deps)` in `adapters/libs/<lib>/<lib>.go`.
- A command handler is always `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- A package's first file is named after the package (`sandbox/deps/iodeps/iodeps.go`,
  `adapters/libs/iodeps/iodeps.go`); a second file is named after what it holds.
- A dep is named after the contract it installs, not after its lib (`argvdeps`, lib `verb`).
- Reusable logic goes in `sandbox/internal/<pkg>/`, one directory per concern.
{{ if .HasCli }}
## Handlers

- Only `CommandHandler(deps *deps.Deps, entries *Entries) int` is exported. `Entries` is
  generated from `entries.yaml` (flags first, then args, in declaration order), already typed,
  defaulted and range-checked.
- Import nothing outside `sandbox/`. Every effect goes through `deps.<Contract>` — see
  [PublicApi](../PublicApi/doc.md).
- Return `api.ExitOk` or `api.ExitFailure`, never `api.ExitUsage`: the dispatch rejects bad
  input before the handler runs.
- Reusable logic goes in `sandbox/internal/<pkg>/`, not in the handler.
- Never hand-edit a command's `entries.yaml`: use `add-flag` / `add-arg` / `set-command`, which
  re-render it with keys in alphabetical order and drop comments.

## Output channels

| Channel | Stream | Carries | `--quiet` |
|---|---|---|---|
| `deps.Std.Printf` | stdout | The result (listings, version, help) | kept |
| `deps.Std.Log` | stderr | Progress | silenced |
| `deps.Std.Error` | stderr | Usage errors and failures | kept |

Never `fmt.Printf`.

## Exit codes

| Code | Const | Meaning |
|---|---|---|
| 0 | `api.ExitOk` | Done |
| 1 | `api.ExitFailure` | A well-formed command failed |
| 2 | `api.ExitUsage` | Bad command line: unknown command or flag, leftover positional, missing required, bad or out-of-range number |
{{ end }}
## Callers

| Rule | Why it matters to a caller |
| --- | --- |
| `sandbox/api` is pure contract | Import it freely; it pulls in nothing else. |
| `sandbox/` never touches the OS | Every effect goes through a `deps.Deps` field you can replace. |
| `deps.Deps` field names are mechanical | A field is the title-cased contract directory name. |
| Generated files are overwritten | Build your own entry point in a package of your own. |
| Deps are patched before `sandbox.New` | The binders capture the pointer. |

## Docs

- A doc is `docs/<Name>/{doc.md,props.yaml}`; sub-docs nest as `docs/<Name>/<Sub>/`. Other
  files in a doc dir are assets.
- Create and delete docs with `add-doc` / `remove-doc`, never by hand.
- `README.md` (its documentation index included), every `Index.md`, and the docs rendered from
  `assets/all/docs/` are generated — never edit them.
- Every `docs/**` dir has a parsable `props.yaml`; a first-level doc names at least one theme of
  `{{.ConfigDir}}/themes.yaml`, a sub-doc names none. A theme no doc names renders no README
  section and is not an error. **(verify)**
- A theme only groups a doc into a section of `README.md`.
- [PublicApi](../PublicApi/doc.md) is rendered from the doc comments of `sandbox/api/` and
  `sandbox/deps/`: change a comment, not the page.
- [Structure](../Structure/doc.md)'s tree is rendered from `{{.ConfigDir}}/{{.StructureConfFile}}`:
  describe an element by adding `<path>: {description: "..."}` there — nested under `children:`
  of its parent, with `dir: true` on a directory, `gen: true` on a file `build` rewrites, and
  `order:` to place it among its siblings (unordered siblings follow, alphabetically).
- Every entry of `{{.ConfigDir}}/{{.StructureConfFile}}` names a path that exists — a directory
  when it declares `dir: true`, a file otherwise. A path holding `<`, `*` or `?` stands for a
  family, and only its literal head has to exist. Drop the entry when the path goes. **(verify)**
- [Commands](../Commands/doc.md) is rendered from every
  `sandbox/internal/commands/<name>/entries.yaml`: document a command by declaring it —
  `set-command` for its help, long description and examples, `add-flag` / `add-arg` for its
  fields — never by editing the page.
- Docs are short, objective and dense: tables, commands, file paths and rules — no prose, no
  narrative, no tutorials, no motivation sections. One page per topic; no sub-doc unless the
  content is a real list of independent items.
- Say a rule once, in this page, and link to it. Links are relative to the file that carries
  them: `../X/doc.md` inside `docs/`, `docs/X/doc.md` in `README.md` and `ReadmeHeader.md`.
