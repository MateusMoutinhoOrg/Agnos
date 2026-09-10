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
- **Deterministic and idempotent.** Same input, same bytes out: `{{.GeneratorName}} build` run twice must
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

- `sandbox/` is closed: a file there imports only `sandbox/` packages — the stdlib included. A
  capability from outside (io, text, sorting, hashing, templating) is restated as a contract
  under `sandbox/deps/` and reached as `deps.<Contract>`. **(verify)**
- `sandbox/` holds only `api`, `binds`, `deps`, `internal` and `new.go`. **(verify)**
- `sandbox/api/*` imports nothing at all. **(verify)**
- `sandbox/deps/<x>/` imports nothing at all: a contract is written in Go's builtin types only,
  and the adapter converts. The loose `sandbox/deps/*.go` is the one exception — it may name
  `sandbox/deps` packages, to compose `deps.Deps`. **(verify)**
- Every `sandbox/binds/` file mirrors one `api/` file and declares only functions. **(verify)**
- Every file of `sandbox/api/` and `sandbox/deps/` parses, and every exported type, func, const
  and var in them carries a doc comment — [PublicApi](../PublicApi/doc.md) is generated from
  those comments. **(verify)**
- `adapters/` is the only place OS-bound and third-party code lives, and holds only
  `availables` and `libs`. **(verify)**
- Every `adapters/libs/<adapter>/` exports `Bind(deps *deps.Deps)` and carries the
  `adapter.yaml` naming the dep it fills. **(verify)**
- Every available fills every field of `Deps` **exactly once**: zero is a nil func that panics
  on first use, two is a silent overwrite in which the last binder wins. Which adapter fills
  which field is read from `adapter.yaml`, never from the body of a `Bind`. **(verify)**
- `adapters/availables/<name>/available.yaml` is the only place the choice of adapter is
  recorded; `set-adapter` is its only editor. An available with no `available.yaml` is
  hand-written and no build touches it.
- `cmd/main/` wires an adapter into the sandbox and holds no logic.
{{- if .HasAssets }}
- Every `assets/deplist/<dep>/<path>` and `assets/adapterlist/<adapter>/<path>`, rendered with
  this module, equals `<path>` whenever that file exists here — re-mirror whenever either side
  changes. **(verify)**
{{- end }}

## Naming

- A `Deps` field is the title-cased `sandbox/deps/<dir>` (`iodeps` -> `deps.Iodeps`). Always
  use that spelling; an added contract never renames an existing one.
- An adapter's binder is always `Bind(deps *deps.Deps)` in `adapters/libs/<adapter>/<adapter>.go`.
- A command handler is always `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- A package's first file is named after the package (`sandbox/deps/iodeps/iodeps.go`,
  `adapters/libs/iodeps/iodeps.go`); a second file is named after what it holds.
- A dep is named after the contract it installs; an adapter after what backs it (`argvdeps`,
  adapter `verb`). The two are separate names because one dep may have several adapters.
- Reusable logic goes in `sandbox/internal/<pkg>/`, one directory per concern.
{{ if .HasCli }}
## Handlers

- Only `CommandHandler(deps *deps.Deps, entries *Entries) int` is exported. `Entries` is
  generated from `entries.yaml` (flags first, then args, in declaration order), already typed,
  defaulted and range-checked.
- Import nothing outside `sandbox/`, the stdlib included. Every effect and every helper goes
  through `deps.<Contract>` — see [PublicApi](../PublicApi/doc.md).
- Return `api.ExitOk` or `api.ExitFailure`, never `api.ExitUsage`: the dispatch rejects bad
  input before the handler runs.
- Reusable logic goes in `sandbox/internal/<pkg>/`, not in the handler.
- A command's `entries.yaml` is written by `add-flag` / `add-arg` / `set-command`, never by
  hand: they re-render it with keys in alphabetical order and drop comments.
{{ end }}{{ if .HasServer }}
## Routes

- A route is `sandbox/internal/routes/<name>/`, holding `route.yaml` (the declaration),
  `entries.go` (generated) and `handler.go` (hand-written) — the server layer's mirror of a
  command package, snake_case for a kebab-case name. **(verify)**
- Only `RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int` is
  exported from a route. It returns the status it answered with, and reaches `400`/`413`/`415`
  only by propagating one from `ReadBody`: the dispatch settles everything but the body before
  the handler runs. **(verify)**
- A route's `route.yaml` is written by `add-route` and rewritten by `set-route`,
  `add-segment` / `remove-segment`, `add-header` / `remove-header`, `add-param` /
  `remove-param`, `set-body` and `add-body-field` / `remove-body-field` — one editor per place
  the file holds something, and never by hand: they re-render it with keys in alphabetical
  order and drop comments.
- Every `identifier` of `paths` starts with `/` and spells exactly one segment; `/` alone is
  the root. A route declares at least one of them, and every entry of `paths` carries an
  `identifier` or a `name`, never both. **(verify)**
- A captured segment is always `required: true` and never defaulted: it is present whenever
  the route matched. `array: true` on it takes every segment left in the path into a `[]T`
  field, which only the last entry of `paths` may do, and which needs at least one segment to
  match — so its name is declared nowhere else. **(verify)**
- A name is declared once per origin, and the origins declaring the same name agree on its
  type — the `Entries` field is written once. **(verify)**
- No two routes declare the same method and path pattern. **(verify)**
- A `json-schema` is declared on a `type: json` body alone, and only with the keywords of the
  subset — `$ref`, `oneOf`, `allOf`, `anyOf` and `patternProperties` fail the build. **(verify)**
- Match order is the collector's, not the directory's: most `identifier`s first, then the
  longest ones, then the routes of fixed length before the ones taking the rest of the path,
  then the pattern alphabetically. Without it a route on `/` would swallow one on `/home`.
- A failure is written by `routeio.WriteError` alone, so every route answers one JSON shape.

Every key of a declaration is in [RouteYaml](../RouteYaml/doc.md).
{{ end }}{{ if .HasFront }}
## Pages

- A page is a route with an html template beside it: `assets/frontend/pages/<page>.html` is
  the whole of what tells one from any other route, and `add-page` / `remove-page` are its
  editors — `remove-route` refuses a route that has one, so no html is ever orphaned.
- A page's route and its html are written once and then the project's, `add-page` keeping an
  html that is already there. Only `sandbox/internal/pageio/` is rewritten by every build, so
  a fix to the scaffolded route or page reaches a project by
  `remove-page <p> && add-page <p>`, never on its own.
- Everything under `assets/frontend/` is hand-written content and survives `front-purge`; the
  routes reading it do not, because their handlers import `pageio`.
- A page renders through `pageio.Render` alone, which is what registers `staticref`, `cssref`,
  `jsref`, `dirref`, `inline` and `include`. A helper pointed at an asset that is not there
  fails the render, so a dead link is a `500` and never a silent `404`.
- `pageio.StaticMount` is generated from the `identifier` of the first segment of the `static`
  route: the mount is declared in one place, and renaming it moves every link.
- Whoever edits `sandbox/internal/routes/static/handler.go` keeps `safeSegments`: it is the
  only thing between a caller's path and the rest of the embedded asset tree.

Every helper and every var is in [FrontUsage](../FrontUsage/doc.md).
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
- An example ends by copying out of `TestDir` into `AssertDir` the paths it asserts, each keeping
  the place it holds in the tree — `AssertDir` is what `result.yaml` records, and `exec-test`
  removes it before every run too. Copying is not moving: `TestDir` stays whole, for reading.
- An example that copies nothing out fails. Assert the paths the example is about and no more:
  a golden holding the whole project breaks on every unrelated template change.
- `result.yaml` is generated by `exec-test`. Refresh one golden with `update-test <name>`, the
  whole suite with `exec-test --update`, or delete it; never edit one.
- An example's output carries no absolute path other than its own directory, no timestamp and no
  resolved version: those are normalized away or make the golden machine-specific.
{{- if .HasCli }}
- An `example.sh` types the project's `name` exactly as `{{.ConfigDir}}/project.yaml` spells it —
  that name is the alias `exec-test` puts on the PATH, and a case mismatch passes on macOS and
  fails on Linux.
- A `<name>` declared on both sides leaves the same `tree` and exits the same way — so the two
  sides copy the same set into `AssertDir`; `cli-output` is compared per side only.
{{- end }}

Details: [LibExamples](../LibExamples/doc.md){{ if .HasCli }} and
[CliExamples](../CliExamples/doc.md){{ end }}.
