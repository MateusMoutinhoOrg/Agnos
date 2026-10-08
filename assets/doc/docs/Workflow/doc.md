# Workflow

Every change this project takes and the command that makes it. `{{.GeneratorName}}` owns every generated
file; what stays hand-written is listed in [GeneratedFiles](../GeneratedFiles/doc.md), and the
rules each recipe holds to are in [Rules](../Rules/doc.md).

## The loop

```bash
{{.GeneratorName}} build      # verify + regenerate every generated file + go mod tidy + compile
{{.GeneratorName}} verify     # the schema check alone, writes nothing
```

`build` is the only thing that regenerates{{ if .HasCli }} the dispatch,{{ end }} the wiring,
`README.md` and `docs/`, so run it after every hand edit. It is idempotent: a second run leaves
the tree unchanged. Every command below takes `--path <dir>` (default `.`) and `-q`, and runs
`build` for you.
{{- if .HasAssets }}

This project carries its own `assets/` template tree, so an installed `{{.GeneratorName}}` would rewrite it
to that older binary's shape. Build it with a binary compiled from this tree instead.
{{- end }}

No recipe below asks for a Go file to be created by hand except the cases listed under
[Hand-written code](#hand-written-code).

## Choose what {{.GeneratorName}} generates

```bash
{{.GeneratorName}} list-extensions             # every generation mechanic and whether it is on
{{.GeneratorName}} enable-extension readme     # start generating README.md again
{{.GeneratorName}} disable-extension doc       # stop generating docs/, keep what is there
```

`{{.ConfigDir}}/extensions.yaml` is what `build` reads to decide what to render. Turning a
mechanic off stops the generation and removes nothing: the files stay, and they are yours to
edit. Every key is in [Extensions](../Extensions/doc.md).
{{ if .HasCli }}
## Change the command surface

```bash
{{.GeneratorName}} add-command <name> --summary "one line" [--category "Core"] [--pattern 'route add {name}']
{{.GeneratorName}} add-command <name> --middleware --summary "..."      # runs in front of every command line
{{.GeneratorName}} add-arg  <name> --command <cmd> [--type integer] [--required] [--start 1 --end -1]
{{.GeneratorName}} add-flag <name> --command <cmd> [--key --out --key -o] [--type integer --min 1] [--enum a --enum b]
{{.GeneratorName}} set-command <cmd> --description "..." --example "<cmd> --flag v" --identifier <alias>
{{.GeneratorName}} set-arg <name> --command <cmd> ... / set-flag <name> --command <cmd> ...
{{.GeneratorName}} remove-arg <name> --command <cmd> / remove-flag <name> --command <cmd> / remove-command <cmd>
{{.GeneratorName}} list-commands / show-command <cmd> / explain-command -- <argv…>
```

`add-command` writes `sandbox/internal/commands/[<--dir>/]<name>/command.yaml` (the declaration) and a
stub `handler.go` (yours), then generates `new.go` — the `api.Command` that joins
`Cli.Commands` — and `input.go`, the `Input` it is handed. Every key these editors write is
in [CommandYaml](../CommandYaml/doc.md); never edit `command.yaml` by hand.

Then write `handler.go` — the whole hand-written half of a command:

```go
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	result, err := something(sandbox, input.Name)
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.Printf("%s\n", result)
	return nil
}
```

Every value arrives typed, defaulted and checked: bad input was answered with exit 2 before the
handler ran. Printing through `response` answers the line with exit 0; a returned error fails it
with exit 1, even after a print. A command that prints nothing has still run — only a
`--middleware` that answers nothing hands the line to the next command of the chain. [Commands](../Commands/doc.md) documents the command on the
next build.
{{- else }}
## Add the CLI layer

```bash
{{.GeneratorName}} cli-init     # sandbox/internal/generated/cli, cmd/main, help, version and help-flag, stddeps + argvdeps + stringsdeps + OpinionatedAgnosCli
```

From there `{{.GeneratorName}} add-command <name> --summary "..." --category "..."` declares a command and
`{{.GeneratorName}} add-flag` / `add-arg` its fields. `{{.GeneratorName}} cli-purge` removes the layer again.
{{- end }}

{{ if .HasServer }}
## Change the route surface

```bash
{{.GeneratorName}} add-route <name> --pattern '/users/{id:integer}' --method POST --summary "one line" --category "Users"
{{.GeneratorName}} add-route <name> --trigger /admin --trigger-type prefix   # /admin and under, never /administrator
{{.GeneratorName}} add-route <name> --middleware --trigger /admin --before <route>
{{.GeneratorName}} set-route <route> --method PUT --response-type text/plain --example "curl localhost:8080/users"
{{.GeneratorName}} add-path <name> --route <route> --start 1 --end 1 --type integer   # one slice of the path
{{.GeneratorName}} add-path <name> --route <route> --start 0 --end 0 --trigger /v1
{{.GeneratorName}} add-parameter <name> --route <route> --source header --required
{{.GeneratorName}} add-parameter <name> --route <route> --type integer --default 1
{{.GeneratorName}} set-body <route> --type json --required --max-bytes 2097152
{{.GeneratorName}} add-body-field <dotted.name> --route <route> --format email --required
{{.GeneratorName}} import-body <route> --file payload.json --required --infer-format
{{.GeneratorName}} set-path <name> --route <route> --end -1                  # and set-parameter
{{.GeneratorName}} set-body-field <dotted.name> --route <route> --max 130 --clear format
{{.GeneratorName}} show-route <route>                                      # the whole declaration as a tree
{{.GeneratorName}} list-routes                                             # the chain, in run order
{{.GeneratorName}} explain-route GET /admin/users --header authorization=x   # which routes one request reaches
{{.GeneratorName}} rename-route <route> <name>
{{.GeneratorName}} rebalance-routes --step 10                              # room between the rungs again
{{.GeneratorName}} remove-parameter <name> --route <route>                 # and remove-path
{{.GeneratorName}} remove-body-field <dotted.name> --route <route>
{{.GeneratorName}} remove-route <route>
```

`add-route` writes `sandbox/internal/routes/[<--dir>/]<name>/route.yaml` (the declaration, `priority`
and `response-type` always included — `100` for a route, `10` for a `--middleware`) and a stub
`handler.go` (yours), then
generates `new.go` — the `api.Route` that lands in `Server.Routes`, a 1:1 image of the yaml —
and `input.go` — the `Input` the handler is handed.
One editor per place the declaration holds something, so every key of
[RouteYaml](../RouteYaml/doc.md) is reachable from the command line and `route.yaml` is never
edited by hand. `add-body-field` takes a dotted path (`address.city`) and creates the objects
it passes through — into the `json-schema` of a json body, or the flat `form-schema` of a form
one; `set-body` covers the envelope around the schema — how the body is read, whether it is
required, its size limit and its content-type.

Each `add-` has a `set-` beside it, so a key that was forgotten is added to the declaration
that is there instead of removing it and declaring it again: the keys given are written over
the ones already declared, `--clear <key>` takes one off, and the result goes through the same
constructor the `add-` side calls. `import-body` is `add-body-field` run once per key of an
example payload — a document pasted with `--json` or read with `--file`, inferring a type per
key, the objects and lists around them and, with `--infer-format`, the four formats a string
may spell; it never writes over a property already declared, and `--replace` starts the schema
over. `show-route` prints the whole declaration as a tree, `list-routes` the chain, and
`explain-route` which routes one request reaches and why the others are skipped — the three
write nothing, and `explain-route` is the first step when a route does not run.

Then write `handler.go` — the whole hand-written half of a route:

```go
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	response.SetStatus(api.StatusCreated)
	response.Write(payload(sandbox, create(sandbox, props.User, input.Tenant, input.Body)))
	return nil
}
```

`entries` arrives bound and converted — one field per path and per parameter, named by its id,
and the body on `Body` — so a bad request was already answered `400` before the handler ran.

Setting a status or writing a byte is what answers the request. Several routes may match one
request; they run in `priority` order and stop at the first one that answers, so a handler that
does neither has declined and the next one runs — that is the whole of what a middleware is, and
`props` — the request's `routeprops.RouteProps`, typed in `sandbox/internal/routeprops/project.go` — carries what it
learned to the routes after it. A handler refuses a request by returning `sandbox.Deps.OpinionatedAgnosServer.Fail`. What no route answers is answered
by the eight `sandbox/internal/server/errors/handle_*.go`, which `build` writes once and no build
rewrites: they are where a 404, a 405, a 401 or a 500 is worded.
[Routes](../Routes/doc.md) documents the route on the next build, and
[ServerUsage](../ServerUsage/doc.md) is the whole recipe.
{{- else }}
## Add the server layer

```bash
{{.GeneratorName}} server-init      # serverdeps, signaldeps, sandbox/internal/server, the health route, start-server
{{ if .HasAssets }}<name>{{ else }}{{.ProjectName}}{{ end }} start-server  # listens on the first free port of 3000..4000
```

From there `{{.GeneratorName}} add-route <name> --pattern '/<path>/{id}'` declares a
route and `{{.GeneratorName}} add-path` / `add-parameter` / `add-body-field` what it reads. A
project with no CLI gets one first: a server needs a command that starts it.
`{{.GeneratorName}} server-purge` removes the layer again.
{{- end }}

{{ if .HasFront }}
## Change the page surface

```bash
{{.GeneratorName}} add-page <name> --title "One Line"   # assets/front/<name>.html, answered on /<name>
{{.GeneratorName}} remove-page <name>                   # deletes the html
```

A page is a file of `assets/front/`, served as it is by the `front` route: write it by
hand, scaffold it with `add-page`, or point a bundler's output there. Data comes from api
routes the page's js calls. [FrontUsage](../FrontUsage/doc.md) is the whole recipe.
{{- else }}
## Add the front layer

```bash
{{.GeneratorName}} front-init      # the OpinionatedAgnosFront lib, the front route, assets/front/{index,404}.html
{{ if .HasAssets }}<name>{{ else }}{{.ProjectName}}{{ end }} start-server  # serves every file of assets/front
```

From there any file under `assets/front/` is served; `{{.GeneratorName}} add-page <name>`
scaffolds an html one and `remove-page` deletes it. A project with no server layer gets one
first: the front is answered over http. `{{.GeneratorName}} front-purge` removes the layer
again, leaving `assets/front/` alone.
{{- end }}
{{ if .HasDatabase }}
## Change the database surface

```bash
{{.GeneratorName}} add-database app-database --key-prefix app
{{.GeneratorName}} add-table url --database app-database
{{.GeneratorName}} add-table-field alias --database app-database --table url --type key --required
{{.GeneratorName}} add-table-field visits --database app-database --table url --type object
{{.GeneratorName}} add-table-field agent --database app-database --table url --parent visits
{{.GeneratorName}} show-database app-database                  # read the declaration back
```

`set-table-field` and the `remove-` half of each pair are the inverses. Every command rewrites
`sandbox/internal/databases/<db>/database.yaml` and runs `build`, which regenerates `api.go`,
`new.go` and `methods.go` from it — the records, the insert structs, the filter and the body
of every method.

Then call it from wherever needs it:

```go
db := app_database.New(sandbox)
url, err := db.AddUrl(app_database.UrlNew{Alias: "gh", Link: "https://github.com"})
found, ok := db.FindUrlByAlias("gh")
```

A query the declaration cannot describe goes in `methods_custom.go` beside them, hand-written
and rewritten by no build. [Databases](../Databases/doc.md) is the whole recipe.
{{- else }}
## Add the database layer

```bash
{{.GeneratorName}} database-init                     # the store contract, the OpinionatedAgnosDatabase lib, the mechanic on
{{.GeneratorName}} add-database app-database         # the first database
{{.GeneratorName}} add-table url --database app-database
```

From there `add-table-field` declares what a table holds and every method it generates is
written for you. `{{.GeneratorName}} database-purge` removes the layer again.
{{- end }}
{{ if .HasBackoffice }}
## Run the backoffice

```bash
export {{.SecretEnv}}=$(openssl rand -hex 32)   # optional: unset, one is generated per run
{{ if .HasAssets }}<name>{{ else }}{{.ProjectName}}{{ end }} add-backoffice-user --username admin --email admin@example.com --role root
{{ if .HasAssets }}<name>{{ else }}{{.ProjectName}}{{ end }} start-server --insecure-http   # then /admin/login
```

Every page, route and package of it is the project's, written once: change it by hand.
[Backoffice](../Backoffice/doc.md) is the whole map.
{{- else }}
## Add the backoffice

```bash
{{.GeneratorName}} backoffice-init   # /admin pages, /api/admin, users, API tokens, backoffice-db
```

It installs the server, front and database layers it is missing, and writes every file once.
`start-server` then reads the session secret from `{{.SecretEnv}}`, or generates one per run when
it is unset. `{{.GeneratorName}} backoffice-purge` removes it again.
{{- end }}
## Add reusable logic

`sandbox/internal/<pkg>/`, one directory per concern, imported by whatever needs it. No
declaration, no generated counterpart — write the package and run `build`.

## Add a surface to the sandbox api

The api is what a Go caller gets back from `sandbox.New` (see [LibUsage](../LibUsage/doc.md)).
Two hand-written places, then `build` regenerates `sandbox/api/sandbox.go` and
`sandbox/new.go` around them:

1. `sandbox/api/<x>.go` — the contract: `type <X> struct { ... }` of function fields, named
   after the file, every declaration doc-commented (those comments render
   [PublicApi](../PublicApi/doc.md)). It becomes the `api.Sandbox` field `<X>`.
2. `sandbox/internal/<x>/new.go` — `func New<X>(sandbox *api.Sandbox) api.<X>`, assigning each
   field of the contract, with the implementation beside it.

`build` then writes `sandbox/constructors/<x>/constructor.go` — `sandbox.<X> =
<x>.New<X>(sandbox)` — **once**, and `sandbox/new.go` calls it. From there the constructor is
yours: wrap the implementation, decorate the contract, or build a different one entirely.

## Construct a field yourself

`sandbox/new.go` is one `<x>.Constructor(&self)` per directory of `sandbox/constructors/`, so
adding a directory is adding a call. Write `sandbox/constructors/<x>/constructor.go` with
`func Constructor(sandbox *api.Sandbox)` in `package <x>`, run `build`, and it is wired — the
same way a generated one is, and with no generated file to fight over. Editing a constructor
`build` wrote earlier works for the same reason: nothing rewrites it.

## Add a dependency

Everything the sandbox is not allowed to do itself — filesystem, clock, network, subprocess —
arrives through `sandbox.Deps`. Install a ready-made one:

```bash
{{.GeneratorName}} list-deps                 # every installable contract
{{.GeneratorName}} add-dep <dep>             # sandbox/deps/<dep>/ + its default adapter + the go.mod require
{{.GeneratorName}} add-dep <dep> --adapter <adapter>
{{.GeneratorName}} remove-dep <dep> [--with-adapters]
```

[DepList](../DepList/doc.md) is the catalogue. For one of your own, write the two halves:

1. `sandbox/deps/<x>deps/<x>deps.go` — `type Contract struct { ... }` of function fields, no import at all.
2. `adapters/impls/<impl><x>/<impl><x>.go` — `func Bind(deps *deps.Deps) { deps.<X>Deps = <x>deps.Contract{...} }`,
   any import allowed, beside an `adapter.yaml` saying `dep: <x>deps`.

Then bind it: add `<impl><x>` to `adapters/bindings/standard/binding.yaml`, or let
`{{.GeneratorName}} add-dep` do both for a dep of the catalogue. Reach it as `sandbox.Deps.<X>Deps`
from anywhere inside `sandbox/`.

One contract may have several adapters — see [Adapters](../Adapters/doc.md).
{{- if not .HasDeps }}

This project has no `sandbox/deps/` yet: `{{.GeneratorName}} deps-init` creates it (`deps-purge` removes it).
{{- end }}

## Add a doc

```bash
{{.GeneratorName}} add-doc <Name> --theme <id> --description "one line"    # themes: {{.ConfigDir}}/themes.yaml
{{.GeneratorName}} add-doc <Name>/<Sub> --description "one line"           # sub-doc, no theme
{{.GeneratorName}} remove-doc <Name>
```

Write `docs/<Name>/doc.md`; `README.md`'s index, and the parent `Index.md` of a sub-doc, are
regenerated. Describe any new path worth naming in `{{.ConfigDir}}/{{.StructureConfFile}}` — that
file is what renders [Structure](../Structure/doc.md).

## Add an example

```bash
{{ if .HasCli }}{{.GeneratorName}} add-cli-example <name>       # examples/cli/<name>/example.sh
{{ end }}{{.GeneratorName}} add-lib-example <name>       # examples/lib/<name>/example.go
{{.GeneratorName}} run-examples                    # run them all, check each against its golden
{{.GeneratorName}} run-examples --only <name>      # one example, both sides
{{.GeneratorName}} update-example <name>           # rewrite that one golden with what it produces now
{{.GeneratorName}} run-examples --update           # rewrite every golden at once
{{ if .HasCli }}{{.GeneratorName}} remove-cli-example <name>
{{ end }}{{.GeneratorName}} remove-lib-example <name>
```

Write the example itself, ending with the copy out of `test-dir` into `assert-dir` that says what
it asserts: `result.yaml` records `assert-dir`, and an example that copies nothing out fails.
The golden is written by the first `run-examples` and refreshed with `update-example <name>`, which
prints what it changes before writing. Details in [LibExamples](../LibExamples/doc.md){{ if .HasCli }} and
[CliExamples](../CliExamples/doc.md){{ end }}.

## Hand-written code

| File | Written when |
| --- | --- |
{{- if .HasCli }}
| `sandbox/internal/commands/<name>/handler.go` | a command does something |
{{- end }}
{{- if .HasServer }}
| `sandbox/internal/routes/<name>/handler.go` | a route answers something |
{{- end }}
{{- if .HasFront }}
| `assets/front/**` | the site looks like something |
{{- end }}
{{- if .HasBackoffice }}
| `sandbox/internal/server/backoffice/**`, `assets/backoffice/*.html` | the backoffice behaves or looks otherwise |
{{- end }}
| `sandbox/internal/<pkg>/*.go` (never under `generated/`) | logic worth reusing |
| `sandbox/api/<x>.go` + `sandbox/internal/<x>/new.go` | a new api surface |
| `sandbox/constructors/<x>/constructor.go` | how a field of the `Sandbox` is built |
| `sandbox/deps/<x>/<x>.go` + `adapters/impls/<x>/<x>.go` + its `adapter.yaml` | a new dependency |

Everything else is regenerated over. Two more files are yours: `{{.ConfigDir}}/docs/ReadmeHeader.md`
is the whole of `README.md` above the documentation index, and `LICENSE` is pasted verbatim into
its License section — put whatever license you want there.

A project built before the `OpinionatedAgnos<X>` libs keeps hand-written files written against
the generated packages they replaced. `add-dep` the lib of every mechanic that is on (`verify`
names the missing ones); the next `build` removes `sandbox/internal/generated/{cliio,trigger,routeio,frontio,databaseio}`,
`cli/command`, `server/route`, `main.go` and `main.go`; then `verify` names every
hand-written import of them with its replacement — `cliio.Fail(sandbox, …)` is
`sandbox.Deps.OpinionatedAgnosCli.Fail(…)`, `routeio.RequestOf(route)` is `route.Request`, a `Handle*`
logs and then calls `OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, …)`, and
`start-server` calls `sandbox.Server.Serve` rather than `server.Main`.

## Ship
{{ if .HasCli }}
```bash
{{.GeneratorName}} compile --target all   # cross-compile ./cmd/main into release/
{{.GeneratorName}} publish                # build, compile, then a gh release
```

`go build -o release/{{.ProjectName}} ./cmd/main` is the plain local binary.
`publish` names the release after `version` in `{{.ConfigDir}}/project.yaml`; bump it there
first. `compile` targets: `linux86`, `linuxarm64`, `linuxi32`, `mac86`, `macarm64`,
`windows86`, `windowsi32`, or `all`.
{{- else }}
This project has no `cmd/main` to compile: it ships as the Go module other programs import
(see [LibUsage](../LibUsage/doc.md)). Bump `version` in `{{.ConfigDir}}/project.yaml` and tag
the repository; `{{.GeneratorName}} cli-init` adds a binary if you want one.
{{- end }}
