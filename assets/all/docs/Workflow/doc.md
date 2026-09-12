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
{{ if .HasCli }}
## Change the command surface

```bash
{{.GeneratorName}} add-command <name> --help "one line" --category "Core"
{{.GeneratorName}} add-flag <name> --command <cmd> --type string --description "..." [--default . | --required]
{{.GeneratorName}} add-arg  <name> --command <cmd> --type int --min 1 --description "..."
{{.GeneratorName}} set-command <cmd> --long-description "..." --example "<cmd> --flag v" --identifier <alias>
{{.GeneratorName}} remove-flag <name> --command <cmd>
{{.GeneratorName}} remove-arg  <name> --command <cmd>
{{.GeneratorName}} remove-command <cmd>
```

`add-command` writes `sandbox/internal/commands/<name>/entries.yaml` (the declaration) and a
stub `handler.go` (yours), then generates `entries.go` and the dispatch arm. Every key these
editors write is in [EntriesYaml](../EntriesYaml/doc.md); never edit `entries.yaml` by hand.

Then write `handler.go` — the whole hand-written half of a command:

```go
func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	if err := something(sandbox, entries.Path); err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	sandbox.Deps.Std.Printf("%s\n", result)
	return api.ExitOk
}
```

`Entries` arrives typed, defaulted and range-checked: bad input already exited 2 before the
handler ran. [Commands](../Commands/doc.md) documents the command on the next build.
{{- else }}
## Add the CLI layer

```bash
{{.GeneratorName}} cli-init     # sandbox/internal/cli, cmd/main, the help and version commands, argvdeps + std
```

From there `{{.GeneratorName}} add-command <name> --help "..." --category "..."` declares a command and
`{{.GeneratorName}} add-flag` / `add-arg` its fields. `{{.GeneratorName}} cli-purge` removes the layer again.
{{- end }}

{{ if .HasServer }}
## Change the route surface

```bash
{{.GeneratorName}} add-route <name> --trigger /<path> --method POST --help "one line" --category "Users"
{{.GeneratorName}} set-route <route> --method PUT --example "curl localhost:8080/users"
{{.GeneratorName}} add-segment <name> --route <route>            # a capture; --identifier /users for a literal
{{.GeneratorName}} add-segment <name> --route <route> --array    # the last one, taking the rest of the path
{{.GeneratorName}} add-header <name> --route <route> --required
{{.GeneratorName}} add-param <name> --route <route> --type int --default 1 --min 1
{{.GeneratorName}} set-body <route> --type json --required --max-bytes 2097152
{{.GeneratorName}} add-body-field <dotted.name> --route <route> --format email --required
{{.GeneratorName}} remove-segment <name> --route <route>         # and remove-header / remove-param /
{{.GeneratorName}} remove-body-field <dotted.name> --route <route>
{{.GeneratorName}} remove-route <route>
```

`add-route` writes `sandbox/internal/routes/<name>/route.yaml` (the declaration) and a stub
`handler.go` (yours), then generates `entries.go` and the match/handle pair of the dispatch.
One editor per place the declaration holds something, so every key of
[RouteYaml](../RouteYaml/doc.md) is reachable from the command line and `route.yaml` is never
edited by hand. `add-body-field` takes a dotted path (`address.city`) and creates the objects
it passes through; `set-body` covers the envelope around the schema — how the body is read,
whether it is required, its size limit and its content-type.

Then write `handler.go` — the whole hand-written half of a route:

```go
func RouteHandler(sandbox *api.Sandbox, entries *Entries, response serverdeps.Response) int {
	body, status := entries.ReadBody(sandbox, response)
	if status != api.StatusOk {
		return status
	}
	return writeJson(sandbox, response, api.StatusCreated, create(sandbox, entries.Tenant, body))
}
```

`Entries` arrives bound, converted and range-checked: a bad request was already answered `400`
before the handler ran. The body is the exception — it is read only when `ReadBody` asks for
it. [Routes](../Routes/doc.md) documents the route on the next build, and
[ServerUsage](../ServerUsage/doc.md) is the whole recipe.
{{- else }}
## Add the server layer

```bash
{{.GeneratorName}} server-init      # serverdeps, sandbox/internal/server, the health route, start-server
{{.Name}} start-server  # listens on :8080
```

From there `{{.GeneratorName}} add-route <name> --trigger /<path> --help "..." --category "..."` declares a
route and `{{.GeneratorName}} add-segment` / `add-header` / `add-param` / `add-body-field` its fields. A
project with no CLI gets one first: a server needs a command that starts it.
`{{.GeneratorName}} server-purge` removes the layer again.
{{- end }}

{{ if .HasFront }}
## Change the page surface

```bash
{{.GeneratorName}} add-page <name> --trigger /<path> --title "One Line"
{{.GeneratorName}} remove-page <name>                             # the route and the html both
```

`add-page` writes the route that answers the page (`route.yaml` + a `handler.go` rendering
through `pageio`) and `assets/frontend/pages/<name>.html`, then generates `entries.go` and the
match/handle pair of the dispatch. A page **is** a route, so every editor of
[RouteYaml](../RouteYaml/doc.md) works on its declaration and [Routes](../Routes/doc.md)
documents it.

Then write the html, naming assets rather than hardcoding links:

```html
{{"{{ dirref \"styles\" }}"}}
<h1>{{"{{ .Title }}"}}</h1>
```

Every `{{"{{ .Field }}"}}` of the page is one exported field of the `pageVars` struct in its
handler, so a new variable is a compile error until it is declared.
[FrontUsage](../FrontUsage/doc.md) is the whole recipe.
{{- else }}
## Add the front layer

```bash
{{.GeneratorName}} front-init                  # pageio, the static route, assets/frontend/
{{.GeneratorName}} add-page home --trigger /   # a page answering GET /
{{.Name}} start-server             # serves it
```

From there `{{.GeneratorName}} add-page <name>` declares a page and `remove-page` drops it,
html included. A project with no server layer gets one first: a page is answered over http.
`{{.GeneratorName}} front-purge` removes the layer again, leaving `assets/frontend/` alone.
{{- end }}
## Add reusable logic

`sandbox/internal/<pkg>/`, one directory per concern, imported by whatever needs it. No
declaration, no generated counterpart — write the package and run `build`.

## Add a surface to the sandbox api

The api is what a Go caller gets back from `sandbox.New` (see [LibUsage](../LibUsage/doc.md)).
Three hand-written files, then `build` regenerates `sandbox/api/sandbox.go` and `sandbox/new.go`
around them:

1. `sandbox/api/<x>.go` — the contract: `type <X> struct { ... }` of function fields, named
   after the file, every declaration doc-commented (those comments render
   [PublicApi](../PublicApi/doc.md)). It becomes the `api.Sandbox` field `<X>`.
2. `sandbox/internal/<x>/` — the implementation.
3. `sandbox/binds/<x>.go` — `func <X>Bind(sandbox *api.Sandbox)`, assigning
   each field of `sandbox.<X>`. One binds file per api file, functions only.

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

1. `sandbox/deps/<x>/<x>.go` — `type Sandbox struct { ... }` of function fields, no import at all.
2. `adapters/libs/<x>/<x>.go` — `func Bind(deps *deps.Deps) { deps.<X> = <x>.Sandbox{...} }`, any
   import allowed, beside an `adapter.yaml` saying `dep: <x>`.

Then bind it: add `<x>` to `adapters/availables/standard/available.yaml`, or let
`{{.GeneratorName}} add-dep` do both for a dep of the catalogue. Reach it as `sandbox.Deps.<X>`
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
{{.GeneratorName}} exec-test                    # run them all, check each against its golden
{{.GeneratorName}} exec-test --only <name>      # one example, both sides
{{.GeneratorName}} update-test <name>           # rewrite that one golden with what it produces now
{{.GeneratorName}} exec-test --update           # rewrite every golden at once
{{ if .HasCli }}{{.GeneratorName}} remove-cli-example <name>
{{ end }}{{.GeneratorName}} remove-lib-example <name>
```

Write the example itself, ending with the copy out of `TestDir` into `AssertDir` that says what
it asserts: `result.yaml` records `AssertDir`, and an example that copies nothing out fails.
The golden is written by the first `exec-test` and refreshed with `update-test <name>`, which
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
| `assets/frontend/pages/<page>.html`, `assets/frontend/static/**` | a page looks like something |
{{- end }}
| `sandbox/internal/<pkg>/*.go` | logic worth reusing |
| `sandbox/api/<x>.go` + `sandbox/binds/<x>.go` | a new api surface |
| `sandbox/deps/<x>/<x>.go` + `adapters/libs/<x>/<x>.go` + its `adapter.yaml` | a new dependency |

Everything else is regenerated over. Two more files are yours: `{{.ConfigDir}}/docs/ReadmeHeader.md`
is the whole of `README.md` above the documentation index, and `LICENSE` is pasted verbatim into
its License section — put whatever license you want there.

## Ship
{{ if .HasCli }}
```bash
{{.GeneratorName}} compile --target all   # cross-compile ./cmd/main into release/
{{.GeneratorName}} publish                # build, compile, then a gh release
```

`go build -o release/{{.Name}} ./cmd/main` is the plain local binary.
`publish` names the release after `version` in `{{.ConfigDir}}/project.yaml`; bump it there
first. `compile` targets: `linux86`, `linuxarm64`, `linuxi32`, `mac86`, `macarm64`,
`windows86`, `windowsi32`, or `all`.
{{- else }}
This project has no `cmd/main` to compile: it ships as the Go module other programs import
(see [LibUsage](../LibUsage/doc.md)). Bump `version` in `{{.ConfigDir}}/project.yaml` and tag
the repository; `{{.GeneratorName}} cli-init` adds a binary if you want one.
{{- end }}
