# Workflow

Every change this project takes and the command that makes it. `agnos` owns every generated
file; what stays hand-written is listed in [GeneratedFiles](../GeneratedFiles/doc.md), and the
rules each recipe holds to are in [Rules](../Rules/doc.md).

## The loop

```bash
agnos build      # verify + regenerate every generated file + go mod tidy + compile
agnos verify     # the schema check alone, writes nothing
```

`build` is the only thing that regenerates{{ if .HasCli }} the dispatch,{{ end }} the wiring,
`README.md` and `docs/`, so run it after every hand edit. It is idempotent: a second run leaves
the tree unchanged. Every command below takes `--path <dir>` (default `.`) and `-q`, and runs
`build` for you.
{{- if .HasAssets }}

This project carries its own `assets/` template tree, so an installed `agnos` would rewrite it
to that older binary's shape. Build it with a binary compiled from this tree instead.
{{- end }}

No recipe below asks for a Go file to be created by hand except the cases listed under
[Hand-written code](#hand-written-code).
{{ if .HasCli }}
## Change the command surface

```bash
agnos add-command <name> --help "one line" --category "Core"
agnos add-flag <name> --command <cmd> --type string --description "..." [--default . | --required]
agnos add-arg  <name> --command <cmd> --type int --min 1 --description "..."
agnos set-command <cmd> --long-description "..." --example "<cmd> --flag v" --identifier <alias>
agnos remove-flag <name> --command <cmd>
agnos remove-arg  <name> --command <cmd>
agnos remove-command <cmd>
```

`add-command` writes `sandbox/internal/commands/<name>/entries.yaml` (the declaration) and a
stub `handler.go` (yours), then generates `entries.go` and the dispatch arm. Every key these
editors write is in [EntriesYaml](../EntriesYaml/doc.md); never edit `entries.yaml` by hand.

Then write `handler.go` — the whole hand-written half of a command:

```go
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	if err := something(deps, entries.Path); err != nil {
		deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	deps.Std.Printf("%s\n", result)
	return api.ExitOk
}
```

`Entries` arrives typed, defaulted and range-checked: bad input already exited 2 before the
handler ran. [Commands](../Commands/doc.md) documents the command on the next build.
{{- else }}
## Add the CLI layer

```bash
agnos cli-init     # sandbox/internal/cli, cmd/main, the help and version commands, argvdeps + std
```

From there `agnos add-command <name> --help "..." --category "..."` declares a command and
`agnos add-flag` / `add-arg` its fields. `agnos cli-purge` removes the layer again.
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
3. `sandbox/binds/<x>.go` — `func <X>Bind(deps *deps.Deps, sandbox *api.Sandbox)`, assigning
   each field of `sandbox.<X>`. One binds file per api file, functions only.

## Add a dependency

Everything the sandbox is not allowed to do itself — filesystem, clock, network, subprocess —
arrives through `deps.Deps`. Install a ready-made one:

```bash
agnos dep-list                 # every installable contract
agnos dep-install <dep>        # sandbox/deps/<dep>/ + adapters/libs/<lib>/ + the go.mod require
agnos dep-remove <dep>
```

[DepList](../DepList/doc.md) is the catalogue. For one of your own, write the two halves and
`build` picks them up from the directory listing:

1. `sandbox/deps/<x>/<x>.go` — `type Lib struct { ... }` of function fields, stdlib imports only.
2. `adapters/libs/<x>/<x>.go` — `func Bind(deps *deps.Deps) { deps.<X> = <x>.Lib{...} }`, any
   import allowed.

Reach it as `deps.<X>` from anywhere inside `sandbox/`.
{{- if not .HasDeps }}

This project has no `sandbox/deps/` yet: `agnos deps-init` creates it (`deps-purge` removes it).
{{- end }}

## Add a doc

```bash
agnos add-doc <Name> --theme <id> --description "one line"    # themes: {{.ConfigDir}}/themes.yaml
agnos add-doc <Name>/<Sub> --description "one line"           # sub-doc, no theme
agnos remove-doc <Name>
```

Write `docs/<Name>/doc.md`; `README.md`'s index, and the parent `Index.md` of a sub-doc, are
regenerated. Describe any new path worth naming in `{{.ConfigDir}}/{{.StructureConfFile}}` — that
file is what renders [Structure](../Structure/doc.md).

## Add an example

```bash
{{ if .HasCli }}agnos add-cli-example <name>       # examples/cli/<name>/example.sh
{{ end }}agnos add-lib-example <name>       # examples/lib/<name>/example.go
agnos exec-test                    # run them all, check each against its golden
agnos exec-test --only <name> --update
{{ if .HasCli }}agnos remove-cli-example <name>
{{ end }}agnos remove-lib-example <name>
```

Write the example itself; `result.yaml` is written by the first `exec-test` and refreshed with
`--update`. Details in [LibExamples](../LibExamples/doc.md){{ if .HasCli }} and
[CliExamples](../CliExamples/doc.md){{ end }}.

## Hand-written code

| File | Written when |
| --- | --- |
{{- if .HasCli }}
| `sandbox/internal/commands/<name>/handler.go` | a command does something |
{{- end }}
| `sandbox/internal/<pkg>/*.go` | logic worth reusing |
| `sandbox/api/<x>.go` + `sandbox/binds/<x>.go` | a new api surface |
| `sandbox/deps/<x>/<x>.go` + `adapters/libs/<x>/<x>.go` | a new dependency |

Everything else is regenerated over. Two more files are yours: `{{.ConfigDir}}/docs/ReadmeHeader.md`
is the whole of `README.md` above the documentation index, and `LICENSE` is pasted verbatim into
its License section — put whatever license you want there.

## Ship
{{ if .HasCli }}
```bash
agnos compile --target all   # cross-compile ./cmd/main into release/
agnos publish                # build, compile, then a gh release
```

`go build -o release/{{.Name}} ./cmd/main` is the plain local binary.
`publish` names the release after `version` in `{{.ConfigDir}}/project.yaml`; bump it there
first. `compile` targets: `linux86`, `linuxarm64`, `linuxi32`, `mac86`, `macarm64`,
`windows86`, `windowsi32`, or `all`.
{{- else }}
This project has no `cmd/main` to compile: it ships as the Go module other programs import
(see [LibUsage](../LibUsage/doc.md)). Bump `version` in `{{.ConfigDir}}/project.yaml` and tag
the repository; `agnos cli-init` adds a binary if you want one.
{{- end }}
