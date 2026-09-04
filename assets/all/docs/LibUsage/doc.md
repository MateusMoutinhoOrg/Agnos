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

[PublicApi](../PublicApi/doc.md) lists every one of them — signatures, props structs{{if .HasDeps}} and
dependency contracts{{end}} — generated from `sandbox/api/` itself on every build.
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
The rules a caller can count on are in [Rules](../Rules/doc.md#callers).
