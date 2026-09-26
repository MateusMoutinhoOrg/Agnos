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

## Extensions

- What this project generates is declared in `{{.ConfigDir}}/extensions.yaml`, one key per
  mechanic, and nowhere else: `build` never infers a mechanic from a directory being present.
  A missing declaration is a hard error, not a default. **(verify)**
- Only the keys of the catalog may appear, and no `sandbox-<x>` mechanic is on while `sandbox`
  is off. **(verify)**
- `false` means *stop generating*, never *delete*: agnos leaves what the mechanic already
  wrote exactly as it is, for the project to keep or edit by hand. Removing those files is
  what an `<x>-purge` does — and it is the same command that writes the `false`.
- The declaration is written by `{{.GeneratorName}} enable-extension` / `disable-extension` and
  by every `<x>-init` / `<x>-purge` pair, never by hand.
- Every key is in [Extensions](../Extensions/doc.md).

## Layers

- `sandbox/` is closed: a file there imports only `sandbox/` packages — the stdlib included. A
  capability from outside (io, text, sorting, hashing, templating) is restated as a contract
  under `sandbox/deps/` and reached as `sandbox.Deps.<Contract>`. **(verify)**
- `sandbox/` holds only `api`, `constructors`, `deps`, `internal` and `new.go`. **(verify)**
- `sandbox/api/*` imports nothing but the loose `sandbox/deps` package, and imports that
  only for `Sandbox.Deps`. **(verify)**
- Every function of `sandbox/internal/` takes `sandbox *api.Sandbox` as
  its first parameter and nothing else standing for the outside world: deps is reached as
  `sandbox.Deps.<Contract>`, and the rest of the api as `sandbox.<Field>`. Holding the api is
  what lets one part of it call another, and what makes a field a caller replaced take effect
  everywhere.
- `sandbox/deps/<x>/` imports nothing at all: a contract is written in Go's builtin types only,
  and the adapter converts. The loose `sandbox/deps/*.go` is the one exception — it may name
  `sandbox/deps` packages, to compose `deps.Deps`. **(verify)**
- Every `sandbox/api/<x>.go` other than `sandbox.go`, `command.go` and `route.go` is a field of
  the `Sandbox`, built by the `New<X>(sandbox) api.<X>` its `sandbox/internal/<x>/new.go` —
  or `sandbox/internal/<x>/<x>/new.go`, for a layer split into packages — declares — the one name `sandbox/constructors/<x>/constructor.go` calls. A contract with no
  such file is a field nothing fills, and no constructor is written for it. **(verify)**
- `sandbox/new.go` is one `<x>.Constructor(&self)` per directory of `sandbox/constructors/`,
  in name order, and nothing else. The directories are the list, so a constructor written by
  hand is called exactly like a generated one.
- Every directory under `sandbox/constructors/` holds a `constructor.go` declaring
  `Constructor(sandbox *api.Sandbox)`, and is named after the package it declares.
  **(verify)**
- `sandbox/constructors/<x>/constructor.go` is written **once**, by the first `build` that
  finds the contract, and no build rewrites it: how a field of the `Sandbox` is built — wrapped,
  decorated, swapped for another implementation — is the project's, not the generator's.
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
- Every type of `sandbox/api/` is convertible: its underlying type is identical
  in a copy of the package made elsewhere, or it is a struct the generator can
  write a converter for. No generics, no `chan`, no anonymous struct or
  interface, no embedded field, and no identifier that is neither predeclared
  nor declared in the package. `Sandbox.Deps` is the one field exempt, because
  it is the one field that does not cross: a consumer installs the api of a
  repo, never its wiring, so the copy drops it. This is what makes every agnos
  repo installable as a dep. **(verify)**
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
- A command handler is always `CommandHandler(sandbox *api.Sandbox, command *api.Command) int`.
- A package's first file is named after the package (`sandbox/deps/iodeps/iodeps.go`,
  `adapters/libs/iodeps/iodeps.go`); a second file is named after what it holds.
- A dep is named after the contract it installs; an adapter after what backs it (`sortdeps`,
  adapter `reflectsort`). The two are separate names because one dep may have several adapters.
- Reusable logic goes in `sandbox/internal/<pkg>/`, one directory per concern.
{{ if .HasCli }}
## Handlers

- A command is `sandbox/internal/commands/<name>/`, holding `entries.yaml` (the declaration),
  `new.go` (generated) and `handler.go` (hand-written), snake_case for a kebab-case name.
- Only `CommandHandler(sandbox *api.Sandbox, command *api.Command) int` is exported. Every flag
  and arg of `entries.yaml` is read off `command` by id (`command.GetString("path")`), already
  typed, defaulted and range-checked.
- Import nothing outside `sandbox/`, the stdlib included. Every effect and every helper goes
  through `sandbox.Deps.<Contract>` — see [PublicApi](../PublicApi/doc.md).
- Return `api.ExitOk` or `api.ExitFailure`, never `api.ExitUsage`: the dispatch rejects bad
  input before the handler runs.
- Reusable logic goes in `sandbox/internal/<pkg>/`, not in the handler.
- A command's `entries.yaml` is written by `add-flag` / `add-arg` / `set-command`, never by
  hand: they re-render it with keys in alphabetical order and drop comments.
- `Cli.Commands` is the whole command surface, one `api.Command` per declared command, built by
  `sandbox/internal/cli/new.go` from each package's generated `NewCommand`. The dispatch and
  both help screens read it; nothing about the command set is generated per command anywhere
  else. Each run binds to its own copy of the declaration, made by `api.BindCommand`, so what
  the slice holds is never written to.
{{ end }}{{ if .HasServer }}
## Routes

- A route is `sandbox/internal/routeslist/<name>/`, holding `route.yaml` (the declaration),
  `new.go` and `entries.go` (generated) and `InternalPureHandler.go` (hand-written) — the
  server layer's mirror of a command package, snake_case for a kebab-case name. **(verify)**
- Only `InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error`
  is exported from a route's hand-written half. **(verify)**
- `new.go` is a 1:1 image of `route.yaml`, built on the generic
  `sandbox/internal/server/route.NewRoute`; `entries.go` is the `Entries` struct — `FullRoute`,
  one field per path, one per parameter, each tagged `id:"<id>"` — plus the `ReadBody` a body
  calls for. The generic `RequestHandler` fills `Entries` by those tags through
  `Deps.Reflectdeps`.
- Setting a status or writing a byte on the response is what answers a request and ends the
  chain — a write sends a `200` ahead of it. A handler that does neither has declined, and the
  next route matching that request runs; a handler that returns a non-nil error without
  answering has failed, and `handle_server_error.go` answers for it. What a handler returns is
  never the status. `SetHeader` alone answers nothing.
- Every route of one request shares `route.Locals`; a middleware hands what it learned to the
  routes after it there, through `routeio.SetLocal` / `routeio.GetLocal`.
- A route of the `after` phase runs once the request has been answered, on a frozen response:
  it never answers, and it declares no body. **(verify)**
- A route's `route.yaml` is written by `add-route` and rewritten by `set-route`,
  `add-path` / `set-path` / `remove-path`, `add-parameter` / `set-parameter` /
  `remove-parameter`, `set-body` and `add-body-field` / `set-body-field` / `remove-body-field` /
  `import-body` — one editor per place the file holds something and one `set-` per `add-`, and
  never by hand: they re-render it with keys in alphabetical order and drop comments.
  `rename-route` and `rebalance-routes` rewrite whole routes; `show-route`, `list-routes` and
  `explain-route` read them and write nothing.
- `methods`, `priority` and `response-type` are required on every route; `priority` is never
  negative, and `methods` holds known methods only — or `ANY`, alone. `segments`, when
  declared, is at least `1`; `phase` is `before` or `after`. **(verify)**
- `add-route` lands a route on rung `100` and a `--middleware` on `10`, so a guard goes in front
  of the routes it guards without renumbering them; `--before` / `--after` place one next to
  another, and `rebalance-routes` makes room again.
- A route declares at least one path. A path's `start` is never negative and its `end` is `-1`
  or not before `start`; its `type` is `string`, `integer`, `number` or `uuid`, and anything but
  `string` reads one segment (`start == end`); a `trigger` has a known type (`equal`, `prefix`,
  `text-prefix`, `suffix`, `regex`), a value, and — for a regex — one that compiles. **(verify)**
- On a path a `prefix` holds on a segment boundary — `/admin` is `/admin` or `/admin/…`, never
  `/administrator`; `text-prefix` is the plain one. On a parameter value the two are the same.
- Every `id` of `paths` and `parameters` is an exported Go name, unique across both and never
  `FullRoute` or `Body`: each one names one field of `Entries`. **(verify)**
- A parameter declares a known type and at least one known font; it is never both `required`
  and defaulted, and a `boolean` is never `required`. **(verify)**
- A trigger — and a path's type — decides whether the route runs at all. A request that fails
  one is not a bad request: that route is simply not the one for it.
- A `405` is answered when a route with explicit `methods` matched the path under another method
  and no route with explicit `methods` ran; an `ANY` route running does not hide it. A `HEAD`
  nothing declares runs the chain again as a `GET`.
- No two routes declare the same method and path pattern *on the same rung of the same phase*.
  Sharing a pattern across rungs is what a middleware in front of a route is; sharing a rung as
  well would leave the order between them undeclared. **(verify)**
- A `json-schema` is declared on a `type: json` body alone, and only with the keywords of the
  subset — `$ref`, `oneOf`, `allOf`, `anyOf` and `patternProperties` fail the build. **(verify)**
- `Server.Routes` is the whole http surface, one `*api.Route` per declared route, built by
  `sandbox/internal/server/server/new.go` from each package's generated `NewRoute`. The dispatch
  reads it and nothing about the route set is generated per route anywhere else; each request
  runs on its copy of the declaration, made by `api.BindRoute`, so nothing bound is ever shared.
- Run order is the collector's, not the directory's: the `before` phase, then the `after` one,
  each lowest `priority` first, then by name.
- Nothing in the dispatch writes a response. Every way a request ends without a route answering
  it is handed to one of the eight `sandbox/internal/server/errors/handle_*.go` — one per status.
  They are written **once**, by the first `build` that finds the server layer, and no build
  rewrites them: what a project answers when nothing matches is the project's. **(verify)**
- A failure is raised with `routeio.Fail` — from the dispatch, from a generated `ReadBody` or
  from a handler — which reaches the right file through the `Fail` field of `api.Server`,
  because a route package may not import `sandbox/internal/server/server`. A `Handle*` file
  answers a failure and never raises one.
- A failure the dispatch raises with nothing to add — nothing matched, method not allowed —
  carries no message, so the wording is the one its `Handle*` file spells. One that knows
  something that file could not — which parameter would not bind, and why — carries its own.
  `routeio.FailureOf` is the one reading of that rule.
- A response body for a failure is written by `routeio.WriteError` alone, so every route answers
  one JSON shape.

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
- `pageio.StaticMount` is generated from the trigger value of the first path of the `static`
  route: the mount is declared in one place, and renaming it moves every link.
- Whoever edits `sandbox/internal/routeslist/static/InternalPureHandler.go` keeps `safeSegments`: it is the
  only thing between a caller's path and the rest of the embedded asset tree.

Every helper and every var is in [FrontUsage](../FrontUsage/doc.md).
{{ end }}{{ if .HasDatabase }}
## Databases

- A database is `sandbox/internal/databases/<db>/`, declared by `specs.yaml` and generated
  whole from it. `add-database`, `add-table`, `add-table-field`, their `set-` editors and
  their inverses are its only editors — never by hand. **(verify)**
- `api.go`, `new.go` and `methods.go` are rewritten by every build. `methods_custom.go` is the
  one escape: hand-written, in the same package, read and rewritten by nothing. A name it
  shares with a generated one is a violation. **(verify)**
- A database is **not** a surface of `sandbox/api/`: its methods are typed by table, so there
  is no `[]Database` standing where `Cli.Commands` stands. Whoever needs one builds it with
  `<db>.New(sandbox)`, which touches no key — building one is free and creates nothing until
  the first record is written.
- A `Find<T>By<Field>` is generated for a `key` field and for no other: it is the only field
  the storage indexes. Every other plain field is reached through `List<T>` and its
  `<T>Filtrage`, so a scan is never sold with the face of an indexed lookup.
- A search answers `(<T>Item, bool)` and a write answers `error`: what failed is an error,
  what is absent is a `false`, and no single `nil` ever stands for both.
- No stored value is asserted into a type without `ok`. A value of the wrong type is an error;
  a field a record never carried reads as the zero value of its type.
- Every `link` names a `target` that is a table of the same database, and every `database`
  field carries `fields` of its own and nests no further — one level is what is generated.
  **(verify)**
- No table declares a field named `id`: every record already carries its permanent one.
  **(verify)**

Every key of a declaration is in [Databases](../Databases/doc.md).
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
{{- if .HasCli }}
  [Commands](../Commands/doc.md) from each `entries.yaml`,
{{- end }}
  [Structure](../Structure/doc.md) from
  `{{.ConfigDir}}/{{.StructureConfFile}}`, `README.md` from
  `{{.ConfigDir}}/docs/ReadmeHeader.md` and every `props.yaml`.
- Docs are short, objective and dense: tables, commands, file paths and rules — no prose, no
  narrative, no tutorials, no motivation sections. One page per topic; no sub-doc unless the
  content is a real list of independent items.
- Say a rule once, in this page, and link to it. Links are relative to the file that carries
  them: `../X/doc.md` inside `docs/`, `docs/X/doc.md` in `README.md` and `ReadmeHeader.md`.

{{ if .HasExample }}
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
{{ end }}
