# Rules

Every rule of this project, in one page. `verify` enforces the ones marked **(verify)**;
the rest are read by the generators or by whoever writes the hand-written files.
Nothing here is repeated elsewhere in `docs/` — other pages link here. The command that
makes each kind of change is in [Workflow](../Workflow/doc.md).

## Authoring

- **Generate over hand-write.** A file that can be rendered from a template, a collector or a
  declaration must be. Hand-written code is contracts, adapters, `sandbox/internal/` and
  `handler.go` only; a new hand-written file needs a reason why generation cannot cover it.
- **Every file is an instance of a pattern.** New code copies an existing sibling exactly:
  same filenames, same function names, same ordering. If no pattern fits, define and document
  the pattern first — `verify` and the collectors read shape by convention, so a one-off
  breaks them.
- **Deterministic and idempotent.** Same input, same bytes out: `agnos build` run twice must
  leave the tree unchanged.
- A generated file is never edited — the `always` rows of
  [GeneratedFiles](../GeneratedFiles/doc.md), `(gen)` in [Structure](../Structure/doc.md).
  Change the declaration it is rendered from{{ if .HasAssets }}, or the template under `assets/`
  when this project carries one for it{{ end }}, then run `build`.
- Generated `.go` is gofmt'ed as it is written, so a regenerated tree diffs to zero against one
  a formatting editor has saved.
- `build` compiles `./cmd/... ./sandbox/... ./adapters/...`, never `./...`{{ if .HasAssets }}:
  `assets/` holds Go templates, not compilable Go{{ end }}.

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
{{- if .HasAssets }}
- Every `assets/deplist/<dep>/<path>`, rendered with this module, equals `<path>` whenever that
  file exists here — re-mirror whenever either side changes. **(verify)**
{{- end }}

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
- A command's `entries.yaml` is written by `add-flag` / `add-arg` / `set-command`, never by
  hand: they re-render it with keys in alphabetical order and drop comments.
{{ end }}
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

## Docs

- A doc is `docs/<Name>/{doc.md,props.yaml}`; sub-docs nest as `docs/<Name>/<Sub>/`. Other
  files in a doc dir are assets. Create and delete them with `add-doc` / `remove-doc`.
- Every `docs/**` dir has a parsable `props.yaml`; a first-level doc names at least one theme of
  `{{.ConfigDir}}/themes.yaml`, a sub-doc names none. A theme no doc names renders no README
  section and is not an error. **(verify)**
- A theme only groups a doc into a section of `README.md`.
- Every entry of `{{.ConfigDir}}/{{.StructureConfFile}}` names a path that exists — a directory
  when it declares `dir: true`, a file otherwise. A path holding `<`, `*` or `?` stands for a
  family, and only its literal head has to exist. Drop the entry when the path goes. **(verify)**
- A generated page is changed at its source, never on the page:
  [PublicApi](../PublicApi/doc.md) from the doc comments of `sandbox/api/` and `sandbox/deps/`,
  [Commands](../Commands/doc.md) from each `entries.yaml`, [Structure](../Structure/doc.md) from
  `{{.ConfigDir}}/{{.StructureConfFile}}`, `README.md` from
  `{{.ConfigDir}}/docs/ReadmeHeader.md` and every `props.yaml`.
- Docs are short, objective and dense: tables, commands, file paths and rules — no prose, no
  narrative, no tutorials, no motivation sections. One page per topic; no sub-doc unless the
  content is a real list of independent items.
- Say a rule once, in this page, and link to it. Links are relative to the file that carries
  them: `../X/doc.md` inside `docs/`, `docs/X/doc.md` in `README.md` and `ReadmeHeader.md`.

## Examples

- An example is `examples/<side>/<name>/`, holding exactly one `example.go` under `lib/`{{ if .HasCli }}
  or one `example.sh` under `cli/`{{ end }}. Create and delete them with `add-lib-example` /
  `remove-lib-example`{{ if .HasCli }} and `add-cli-example` / `remove-cli-example`{{ end }},
  never by hand — the same rule as `add-doc` / `remove-doc`.
- An example runs with its own directory as the working directory and writes only inside its own
  `TestDir`, which `exec-test` removes before every run.
- `result.yaml` is generated by `exec-test`. Refresh a golden with `--update` or by deleting it;
  never edit one.
- An example's output carries no absolute path other than its own directory, no timestamp and no
  resolved version: those are normalized away or make the golden machine-specific.
{{- if .HasCli }}
- An `example.sh` types the project's `name` exactly as `{{.ConfigDir}}/project.yaml` spells it —
  that name is the alias `exec-test` puts on the PATH, and a case mismatch passes on macOS and
  fails on Linux.
- A `<name>` declared on both sides leaves the same `tree` and exits the same way; `cli-output`
  is compared per side only.
{{- end }}

Details: [LibExamples](../LibExamples/doc.md){{ if .HasCli }} and
[CliExamples](../CliExamples/doc.md){{ end }}.
