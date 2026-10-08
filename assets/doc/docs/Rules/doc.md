# Rules

Every rule of this project, in one page. `verify` enforces the ones marked **(verify)**;
the rest are read by the generators or by whoever writes the hand-written files.
Nothing here is repeated elsewhere in `docs/` — other pages link here. The command that
makes each kind of change is in [Workflow](../Workflow/doc.md).

## Authoring

- **Generate over hand-write.** A file that can be rendered from a template, a collector or a
  declaration must be. Hand-written code is contracts, adapters, `sandbox/internal/` and
  `handler.go` only; a new hand-written file needs a reason why generation cannot
  cover it.
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
- `sandbox/internal/generated/` holds every package `build` rewrites whole and nothing else: no
  file there is ever edited, and no hand-written package is ever put there. It holds the
  registries and config alone — code that is the same in every project is an `OpinionatedAgnos<X>`
  lib, not a generated package. Importing a package an older build generated there names its
  replacement. **(verify)** A package mixing a
  generated file with a hand-written one — a command, a route, a database — stays under
  `sandbox/internal/`.
- Generated `.go` is gofmt'ed as it is written, so a regenerated tree diffs to zero against one
  a formatting editor has saved.
- `build` compiles `./cmd/... ./sandbox/... ./adapters/...`, never `./...`{{ if .HasAssets }}:
  `assets/` holds Go templates, not compilable Go{{ end }}.

## Extensions

- What this project generates is declared in `{{.ConfigDir}}/extensions.yaml`, one key per
  mechanic, and nowhere else: `build` never infers a mechanic from a directory being present.
  A missing declaration is a hard error, not a default. **(verify)**
- Only the keys of the catalog may appear, and no mechanic that renders into the sandbox is on while `sandbox`
  is off. **(verify)**
- `false` means *stop generating*, never *delete*: `{{.GeneratorName}}` leaves what the mechanic already
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
- `sandbox/api/*` imports nothing but the loose `sandbox/deps` package, for `Sandbox.Deps`, and
  the `sandbox/deps/OpinionatedAgnos<X>` contracts, for the aliases a mechanic's api file is made
  of. **(verify)**
- Every function of `sandbox/internal/` takes `sandbox *api.Sandbox` as
  its first parameter and nothing else standing for the outside world: deps is reached as
  `sandbox.Deps.<Contract>`, and the rest of the api as `sandbox.<Field>`. Holding the api is
  what lets one part of it call another, and what makes a field a caller replaced take effect
  everywhere.
- `sandbox/deps/<x>/` imports nothing at all: a contract is written in Go's builtin types only,
  and the adapter converts. The loose `sandbox/deps/*.go` may name `sandbox/deps` packages, to
  compose `deps.Deps`, and an `OpinionatedAgnos<X>/` contract may import other contracts under
  `sandbox/deps/` and nothing else. **(verify)**
- A dep states a library's raw capability and never a decision of the project using it — except
  an **opinionated lib**, `OpinionatedAgnos<X>`, the one kind of dep that carries an agnos mechanic
  itself: `OpinionatedAgnosCli` (the command types, the dispatch chain, binding, failures, triggers),
  `OpinionatedAgnosServer` (the route types, the request chain, binding, json-schema, writers),
  `OpinionatedAgnosFront` (the file layer of `assets/front/`), `OpinionatedAgnosDatabase` (the readers
  every `methods.go` shares). Each mechanic's `-init` installs its lib, and a mechanic is never on
  without it. **(verify)** The lib holds no dep: what it reaches the outside world through is
  handed to it — a `MainProps` built by the generated registry, or the one dep a call needs as
  its first parameter. What stays in the sandbox is what the project declares or edits.
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
- `sandbox/api/projectsandbox.go` and `sandbox/api/projectconfig.go` are written **once**, by `start`,
  and no build rewrites them: `api.Sandbox` embeds `api.ProjectSandbox` and `api.Config` embeds
  `api.ProjectConfig`, so what the project declares there is read as `sandbox.<Field>` and
  `sandbox.Config.<Field>`. Neither is a field of its own, so neither gets a constructor: a
  `ProjectSandbox` field is filled by a package of the project's under `sandbox/constructors/`, a
  `ProjectConfig` one in `sandbox/constructors/config/constructor.go`. Each must be there while the
  file embedding it is. **(verify)**
- Every file of `sandbox/api/` and `sandbox/deps/` parses, and every exported type, func, const
  and var in them carries a doc comment — [PublicApi](../PublicApi/doc.md) is generated from
  those comments. **(verify)**
- `adapters/` is the only place OS-bound and third-party code lives, and holds only
  `bindings` and `impls`. **(verify)**
- Every `adapters/impls/<adapter>/` exports `Bind(deps *deps.Deps)` and carries the
  `adapter.yaml` naming the dep it fills. **(verify)**
- Every binding fills every field of `Deps` **exactly once**: zero is a nil func that panics
  on first use, two is a silent overwrite in which the last binder wins. Which adapter fills
  which field is read from `adapter.yaml`, never from the body of a `Bind`. **(verify)**
- `adapters/bindings/<name>/binding.yaml` is the only place the choice of adapter is
  recorded; `set-adapter` is its only editor. A binding with no `binding.yaml` is
  hand-written and no build touches it.
- Every type of `sandbox/api/` is convertible: its underlying type is identical
  in a copy of the package made elsewhere, or it is a struct the generator can
  write a converter for. No generics, no `chan`, no anonymous struct or
  interface, no embedded field but a struct the package declares, and no identifier that is neither predeclared
  nor declared in the package. `Sandbox.Deps` is the one field exempt, because
  it is the one field that does not cross: a consumer installs the api of a
  repo, never its wiring, so the copy drops it. A mechanic's surface — an alias of an
  `OpinionatedAgnos<X>` type, and a part holding only those — is exempt for the same reason: it is
  the lib's, and the copy drops it too. This is what makes every agnos
  repo installable as a dep. **(verify)**
- `cmd/main/` wires an adapter into the sandbox and holds no logic.
{{- if .HasAssets }}
- Every `assets/dep-catalog/<dep>/<path>` and `assets/adapter-catalog/<adapter>/<path>`, rendered with
  this module, equals `<path>` whenever that file exists here — re-mirror whenever either side
  changes. **(verify)**
{{- end }}

## Naming

- A `Deps` field is the title-cased `sandbox/deps/<dir>` (`iodeps` -> `deps.IoDeps`). Always
  use that spelling; an added contract never renames an existing one.
- An adapter's binder is always `Bind(deps *deps.Deps)` in `adapters/impls/<adapter>/<adapter>.go`.
- A command handler is always `Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error`.
- A package's first file is named after the package (`sandbox/deps/iodeps/iodeps.go`,
  `adapters/impls/osio/osio.go`); a second file is named after what it holds.
- A dep is named after the contract it installs, `<x>deps`; an adapter `<impl><x>`, after what
  backs it (`sortdeps`, adapters `stdsort` and `reflectsort`). The two are separate names because one dep may have several adapters.
- Reusable logic goes in `sandbox/internal/<pkg>/`, one directory per concern.
{{ if .HasCli }}
## Handlers

- A command is a directory under the folder of its category in `sandbox/internal/commands/`, holding
  `command.yaml` (the declaration), `new.go` and `input.go` (generated) and
  `handler.go` (hand-written), snake_case for a kebab-case name. The `command.yaml` is
  what makes it one: a directory without it is a folder grouping commands
  (`add-command <name> --dir <folder>`, `rename-command <name> <name> --dir <folder>`). A name is
  unique across every folder, and a directory holding the go files without a `command.yaml` is
  a violation. An `entries.yaml` is an old declaration. **(verify)**
- Only `Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error`
  is exported. Every flag and arg of `command.yaml` is a field of `Input` (`input.Name`),
  already typed, defaulted and range-checked.
- Import nothing outside `sandbox/`, the stdlib included. Every effect and every helper goes
  through `sandbox.Deps.<Contract>` — see [PublicApi](../PublicApi/doc.md).
- `response.Printf` (stdout) answers the command line with `api.ExitOk`; `response.Eprintf` and
  `response.Logf` (stderr) answer nothing. Refuse a command line by returning
  `sandbox.Deps.OpinionatedAgnosCli.Fail` — a
  returned error fails it even after a print. A strict command that returns `nil` without
  printing has run and exits `0`; only a middleware declines.
- Reusable logic goes in `sandbox/internal/<pkg>/`, not in the handler.
- A command's `command.yaml` is written by `add-flag` / `add-arg` / `set-command`, never by
  hand: they re-render it with keys in alphabetical order and drop comments.
- `Cli.Commands` is the whole command surface, one `api.Command` per declared command, built by
  `sandbox/internal/generated/cli/new.go` from each package's generated `NewCommand`. The dispatch and
  both help screens read it; nothing about the command set is generated per command anywhere
  else. The dispatch itself is `OpinionatedAgnosCli.Main`, handed `Cli.Commands` by the
  registry. Each run binds to its own copy of the declaration, made by `BindCommand`, so what
  the slice holds is never written to.
{{ end }}{{ if .HasServer }}
## Routes

- A route is a directory under `sandbox/internal/routes/`, at any depth, holding `route.yaml`
  (the declaration), `new.go` and `input.go` (generated) and `handler.go`
  (hand-written) — the server layer's mirror of a command package, snake_case for a kebab-case
  name. The `route.yaml` is what makes it one: a directory without it is a folder grouping routes
  (`add-route <name> --dir <folder>`, `rename-route <name> <name> --dir <folder>`). A name is
  unique across every folder. **(verify)**
- Only `Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error`
  is exported from a route's hand-written half. **(verify)**
- `new.go` is a 1:1 image of `route.yaml`, built on `OpinionatedAgnosServer.NewRoute`; `input.go`
  is the `Input` struct — `FullRoute`, one field per path, one per parameter and `Body` when a
  body is declared, each tagged `id:"<id>"` — plus the `ReadBody` a body calls for. The lib's
  dispatch reads the body before the handler runs and fills `Input` by those tags. A handler
  is handed no request: what it reads is declared.
- Setting a status or writing a byte on the response is what answers a request and ends the
  chain — a write sends a `200` ahead of it. A handler that does neither has declined, and the
  next route matching that request runs; a handler that returns a non-nil error without
  answering has failed, and `handle_internal_server_error.go` answers for it. What a handler returns is
  never the status. `SetHeader` alone answers nothing.
- Every route of one request is handed the same `props *routeprops.RouteProps`, built empty per
  request; a middleware hands what it learned to the routes after it by setting a field of it.
  `routeprops.go` is generated: it embeds every exported struct of every other file of
  `sandbox/internal/routeprops/` — `project.go`, written once and the project's, plus one file per
  mechanic that hands something on (`backoffice.go`). A field is declared in one of those parts,
  never in `routeprops.go`; two parts may not declare the same field (`verify`). It sits under
  `sandbox/internal`, never `sandbox/api`, so a field may name any type of the project — a
  database record — as long as that package imports no route; `api.Route.Props` holds it as
  `any`. `build` moves a hand-written `routeprops.go` (or an old `sandbox/api/routeprops.go`) to
  `project.go`, its struct renamed `Project`; the same holds for `commandprops.CommandProps` in
  `sandbox/internal/commandprops/`.
- `api.Config` and `api.Sandbox` are generated the same way: each embeds every struct of
  `sandbox/api/<x>config.go` / `sandbox/api/<x>sandbox.go`. `projectconfig.go` and
  `projectsandbox.go` are the project's; each mechanic writes its own part beside them
  (`clisandbox.go`, `serversandbox.go`, `backofficeconfig.go`), never edits another's. The
  `Sandbox` declares only `Deps` and `Config` itself: a contract's field is declared by the part
  of whoever owns it, so a contract the project writes is a field of `ProjectSandbox`. **(verify)**
  `build` renames an old `usersandbox_<x>.go` / `userconfig_<x>.go` to `<x>sandbox.go` /
  `<x>config.go`.
- A route's `route.yaml` is written by `add-route` and rewritten by `set-route`,
  `add-path` / `set-path` / `remove-path`, `add-parameter` / `set-parameter` /
  `remove-parameter`, `set-body` and `add-body-field` / `set-body-field` / `remove-body-field` /
  `import-body` — one editor per place the file holds something and one `set-` per `add-`, and
  never by hand: they re-render it with keys in alphabetical order and drop comments.
  `rename-route` and `rebalance-routes` rewrite whole routes; `show-route`, `list-routes` and
  `explain-route` read them and write nothing.
- `methods`, `priority` and `response-type` are required on every route; `priority` is never
  negative, and `methods` holds known methods only — or `ANY`, alone. `segments`, when
  declared, is at least `1`; a `phase` key is an old declaration. **(verify)**
- `add-route` lands a route on rung `100` and a `--middleware` on `10`, so a guard goes in front
  of the routes it guards without renumbering them; `--before` / `--after` place one next to
  another, and `rebalance-routes` makes room again.
- A route declares at least one path. A path's `start` is never negative and its `end` is `-1`
  or not before `start`; its `type` is `string`, `integer`, `number` or `uuid`, and anything but
  `string` reads one segment (`start == end`); a `trigger` has a known type (`equal`, `prefix`,
  `text-prefix`, `suffix`, `regex`, `one-of`), a value — `values` for a `one-of` — and — for a
  regex — one that compiles. **(verify)**
- On a path a `prefix` holds on a segment boundary — `/admin` is `/admin` or `/admin/…`, never
  `/administrator`; `text-prefix` is the plain one. On a parameter value the two are the same.
- Every `id` of `paths` and `parameters` is an exported Go name, unique across both and never
  `FullRoute` or `Body`: each one names one field of `Input`. **(verify)**
- A parameter declares a known type and at least one known source; it is never both `required`
  and defaulted, and a `boolean` is never `required`. **(verify)**
- A trigger — and a path's type — decides whether the route runs at all. A request that fails
  one is not a bad request: that route is simply not the one for it.
- A `405` is answered when a route with explicit `methods` matched the path under another method
  and no route with explicit `methods` ran; an `ANY` route running does not hide it. A `HEAD`
  nothing declares runs the chain again as a `GET`.
- No two routes declare the same method and path pattern *on the same rung*.
  Sharing a pattern across rungs is what a middleware in front of a route is; sharing a rung as
  well would leave the order between them undeclared. **(verify)**
- A `json-schema` is declared on a `type: json` body alone and a `form-schema` on a `type: form`
  one alone, each only with the keywords of the subset — `$ref`, `oneOf`, `allOf`, `anyOf` and
  `patternProperties` fail the build. A `form-schema` is flat: scalar properties or arrays of
  them, never an object or `nullable`. **(verify)**
- `Server.Routes` is the whole http surface, one `*api.Route` per declared route, built by
  `sandbox/internal/generated/server/new.go` from each package's generated `NewRoute`. The dispatch
  reads it and nothing about the route set is generated per route anywhere else. The dispatch
  itself is `OpinionatedAgnosServer.Main`, handed `Server.Routes` by the registry through
  `Server.Serve`; each request runs on its copy of the declaration, made by `BindRoute`, so
  nothing bound is ever shared.
- Run order is the collector's, not the directory's: lowest `priority` first, then by name.
- Nothing in the dispatch writes a response. Every way a request ends without a route answering
  it is handed to one of the eight `sandbox/internal/server/errors/handle_*.go` — one per status.
  They are written **once**, by the first `build` that finds the server layer, and no build
  rewrites them: what a project answers when nothing matches is the project's. **(verify)**
- A handler — or a generated `ReadBody` — refuses a request by returning
  `sandbox.Deps.OpinionatedAgnosServer.Fail`, a `*api.RouteFailure`; the dispatch raises it. Only
  the dispatch raises, and it reaches the right file through the `Fail` field of `api.Server`,
  which the registry fills. A `Handle*` file answers a failure and never raises one.
- A failure the dispatch raises with nothing to add — nothing matched, method not allowed —
  carries no message, so the wording is the one its `Handle*` file spells. One that knows
  something that file could not — which parameter would not bind, and why — carries its own.
  `OpinionatedAgnosServer.FailureOf` is the one reading of that rule.
- A response body for a failure is written by `OpinionatedAgnosServer.WriteError` alone, so every
  route answers one JSON shape.

Every key of a declaration is in [RouteYaml](../RouteYaml/doc.md).
{{ end }}{{ if .HasFront }}
## Front

- A page is a file of `assets/front/` and nothing else. The `front` route serves the
  whole tree, so a file dropped there by hand or by a bundler is as much a page as one
  `add-page` wrote; `add-page` / `remove-page` only write and delete the html.
- Everything under `assets/front/` is the project's content: no build writes there,
  `add-page` refuses an existing file, and `front-purge` leaves the tree whole.
- The `front` route is written once and then the project's. The file layer is the
  `OpinionatedAgnosFront` lib, and its `SafePath`, which `Resolve` runs first, is what keeps a
  caller's path inside `assets/front/`: the handler resolves every path through it.
- The `front` route runs at priority `1000`, after every api route, and answers a path that
  names no file with `assets/front/404.html` under a `404`; only with that file gone does it
  decline, so the `404` falls to `handle_not_found.go`.

How a path is resolved, and a bundler's build, is in [FrontUsage](../FrontUsage/doc.md).
{{ end }}{{ if .HasBackoffice }}
## Backoffice

- Every file `backoffice-init` wrote is the project's: no build rewrites it, and a second
  `backoffice-init` keeps each one already there. `backoffice-purge` removes them all.
- The backoffice edits nothing the project wrote. What it hands a route is
  `sandbox/internal/routeprops/backoffice.go`, what it reads at startup
  `sandbox/api/backofficeconfig.go`, both embedded by the generated aggregates, and the
  secret is read by the `backoffice-start-server` middleware in front of `start-server`.
- The session secret is the `{{.SecretEnv}}` environment variable, at least 32 characters, never
  a flag or a file; shorter, `start-server` refuses to start. Unset, `start-server` generates one
  for the run and warns: every session ends at a restart.
- Its records live in `backoffice-db`, a database of its own; the project's databases are never
  touched. `./data/backofficedb` is gitignored and survives a purge.

Routes, roles and the API are in [Backoffice](../Backoffice/doc.md).
{{ end }}{{ if .HasDatabase }}
## Databases

- A database is `sandbox/internal/databases/<db>/`, declared by `database.yaml` and generated
  whole from it. `add-database`, `add-table`, `add-table-field`, their `set-` editors and
  their inverses are its only editors — never by hand. **(verify)**
- `api.go`, `new.go` and `methods.go` are rewritten by every build; what every `methods.go`
  shares — resolving a table, reading a stored value, the filter filters — is the
  `OpinionatedAgnosDatabase` lib. `methods_custom.go` is the one escape: hand-written, in the same
  package, read and rewritten by nothing. A name it shares with a generated one is a violation.
  **(verify)**
- A database is **not** a surface of `sandbox/api/`: its methods are typed by table, so there
  is no `[]Database` standing where `Cli.Commands` stands. Whoever needs one builds it with
  `<db>.New(sandbox)`, which touches no key — building one is free and creates nothing until
  the first record is written.
- A `Find<T>By<Field>` is generated for a `key` field and for no other: it is the only field
  the storage indexes. Every other plain field is reached through `List<T>` and its
  `<T>Filter`, so a scan is never sold with the face of an indexed lookup.
- A search answers `(<T>Record, bool)` and a write answers `error`: what failed is an error,
  what is absent is a `false`, and no single `nil` ever stands for both.
- No stored value is asserted into a type without `ok`. A value of the wrong type is an error;
  a field a record never carried reads as the zero value of its type.
- Every `link` names a `target` that is a table of the same database, and every `object`
  field carries `fields` of its own and nests no further — one level is what is generated.
  **(verify)**
- No table declares a field named `id`: every record already carries its permanent one.
  **(verify)**

Every key of a declaration is in [Databases](../Databases/doc.md).
{{ end }}
## Output channels

| Channel | Stream | Carries | Silenced |
|---|---|---|---|
| `deps.StdDeps.Printf` | stdout | The result (listings, version, help) | never |
| `deps.StdDeps.Logf` | stderr | Progress | by a middleware that turns it off{{ if .HasAssets }} — `--quiet` here{{ end }} |
| `deps.StdDeps.Eprintf` | stderr | Usage errors and failures | never |

Never `fmt.Printf`. {{ if not .HasAssets }}A generated cli declares no `--quiet`: add it as a
`--middleware` whose handler silences `Log` when the project wants one.{{ end }}

## Exit codes

| Code | Const | Meaning |
|---|---|---|
| 0 | `api.ExitOk` | Done |
| 1 | `api.ExitFailure` | A well-formed command failed |
| 2 | `api.ExitUsage` | Bad command line: unknown command or flag, leftover positional, missing required, bad or out-of-range number |

## Docs

- A doc is `docs/<Name>/{doc.md,doc.yaml}`; sub-docs nest as `docs/<Name>/<Sub>/`. Other
  files in a doc dir are assets. Create and delete them with `add-doc` / `remove-doc`.
- Every `docs/**` dir has a parsable `doc.yaml`; a first-level doc names at least one theme of
  `{{.ConfigDir}}/themes.yaml`, a sub-doc names none. A theme no doc names renders no README
  section and is not an error. **(verify)**
- A theme only groups a doc into a section of `README.md`.
- Every entry of `{{.ConfigDir}}/{{.StructureConfFile}}` names a path that exists — a directory
  when it declares `dir: true`, a file otherwise. A path holding `<`, `*` or `?` stands for a
  family, and only its literal head has to exist. Drop the entry when the path goes. **(verify)**
- A generated page is changed at its source, never on the page:
  [PublicApi](../PublicApi/doc.md) from the doc comments of `sandbox/api/` and `sandbox/deps/`,
{{- if .HasCli }}
  [Commands](../Commands/doc.md) from each `command.yaml`,
{{- end }}
  [Structure](../Structure/doc.md) from
  `{{.ConfigDir}}/{{.StructureConfFile}}`, `README.md` from
  `{{.ConfigDir}}/docs/ReadmeHeader.md` and every `doc.yaml`.
- Docs are short, objective and dense: tables, commands, file paths and rules — no prose, no
  narrative, no tutorials, no motivation sections. One page per topic; no sub-doc unless the
  content is a real list of independent items.
{{- if .HasServer }}
- [Routes](../Routes/doc.md) is the exception: it is read by whoever calls the server, a person,
  so it speaks plain words and every page carries requests that run as they are. It is
  generated from each `route.yaml`: a `summary`, a `description` or an `examples` entry is what
  changes a page.
{{- end }}
- Say a rule once, in this page, and link to it. Links are relative to the file that carries
  them: `../X/doc.md` inside `docs/`, `docs/X/doc.md` in `README.md` and `ReadmeHeader.md`.

{{ if .HasExample }}
## Examples

- An example is `examples/<side>/<name>/`, holding exactly one `example.go` under `lib/`{{ if .HasCli }}
  or one `example.sh` under `cli/`{{ end }}. Create and delete them with `add-lib-example` /
  `remove-lib-example`{{ if .HasCli }} and `add-cli-example` / `remove-cli-example`{{ end }},
  never by hand — the same rule as `add-doc` / `remove-doc`.
- An example runs with its own directory as the working directory and writes only inside its own
  `test-dir`, which `run-examples` removes before every run.
- An example ends by copying out of `test-dir` into `assert-dir` the paths it asserts, each keeping
  the place it holds in the tree — `assert-dir` is what `result.yaml` records, and `run-examples`
  removes it before every run too. Copying is not moving: `test-dir` stays whole, for reading.
- An example that copies nothing out fails. Assert the paths the example is about and no more:
  a golden holding the whole project breaks on every unrelated template change.
- `result.yaml` is generated by `run-examples`. Refresh one golden with `update-example <name>`, the
  whole suite with `run-examples --update`, or delete it; never edit one.
- An example's output carries no absolute path other than its own directory, no timestamp and no
  resolved version: those are normalized away or make the golden machine-specific.
{{- if .HasCli }}
- An `example.sh` types the project's `name` exactly as `{{.ConfigDir}}/project.yaml` spells it —
  that name is the alias `run-examples` puts on the PATH, and a case mismatch passes on macOS and
  fails on Linux.
- A `<name>` declared on both sides leaves the same `tree` and exits the same way — so the two
  sides copy the same set into `assert-dir`; `cli-output` is compared per side only.
{{- end }}

Details: [LibExamples](../LibExamples/doc.md){{ if .HasCli }} and
[CliExamples](../CliExamples/doc.md){{ end }}.
{{ end }}
