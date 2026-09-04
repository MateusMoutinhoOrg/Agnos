# LibUsage

`{{.Name}}` is a Go module before it is anything else: every feature lives in `sandbox/`
and is reachable from any Go program that imports it.

```bash
go get {{.Module}}@latest
```

## Wiring
{{if .HasDeps}}
`sandbox/` performs no OS effects of its own — filesystem, clock, stdout, processes all
arrive through a `deps.Deps` struct. `adapters/availables/standard` builds the ready-made
assembly, and `sandbox.New` turns it into the API object.

```go
package main

import (
	"{{.Module}}/adapters/availables/standard"
	"{{.Module}}/sandbox"
)

func main() {
	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox

	_ = lib
}
```
{{else}}
`sandbox.New` takes no arguments and returns the API object.

```go
package main

import (
	"{{.Module}}/sandbox"
)

func main() {
	lib := sandbox.New() // *api.Sandbox

	_ = lib
}
```
{{end}}
## What the sandbox exposes

`*api.Sandbox` is a flat struct, one field per contract declared in `sandbox/api/`.
Everything callable from Go is behind one of them.

| Field | Type |
| --- | --- |
{{- range .Constructors }}
| `lib.{{ . }}` | `api.{{ . }}` |
{{- end }}

Read `sandbox/api/` for the exact function signatures and props structs: it is pure
contract — no logic, no imports — so it is the shortest description of the whole surface.
{{if .HasDeps}}
## Custom deps

Every sub-contract is a struct of function fields, so any of them can be swapped for a
test double, an in-memory implementation or an instrumented wrapper. Patch fields **before**
`sandbox.New(&deps)`: the binders capture the pointer.

```go
deps := standard.New()

var out bytes.Buffer
deps.Std.Printf = func(f string, a ...any) (int, error) {
	return fmt.Fprintf(&out, f, a...)
}

lib := sandbox.New(&deps)
```

The contracts available to patch:

| Field | Contract package |
| --- | --- |
{{- range .DepsLibs }}
| `deps.{{ .Title }}` | `sandbox/deps/{{ .Name }}` |
{{- end }}

Each one is filled by a matching implementation under `adapters/libs/`, every package
exposing the same `Bind(deps *deps.Deps)` entry point:

| Adapter lib | Binder |
| --- | --- |
{{- range .AdapterLibs }}
| `adapters/libs/{{ .Name }}` | `{{ .Name }}.Bind(&deps)` |
{{- end }}

Starting from `standard.New()` is the safe default: an unfilled field is a nil func that
panics on first call. For a permanent mix, write your own
`adapters/availables/<name>/new.go` binding only the libs you want — `standard/new.go` is
regenerated on every build, while other directories under `availables/` are left alone.
{{end}}
## Rules that hold for callers

| Rule | Why it matters to you |
| --- | --- |
| `sandbox/api` is pure contract | Import it freely; it pulls in nothing else. |{{if .HasDeps}}
| `sandbox/` never touches the OS | Anything the library does to the outside world goes through a field you can replace. |
| `deps.Deps` field names are mechanical | A field is the title-cased contract directory name, so an added contract never renames an existing one. |{{end}}
| Generated files are overwritten | Build your own entry point in a package of your own instead of editing regenerated ones. |
