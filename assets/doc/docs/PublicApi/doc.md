# PublicApi

Every exported symbol of `{{.Module}}`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns{{if .HasDeps}}, `sandbox/deps/`
the contracts an adapter fills and a caller may replace{{end}}. Each description is the
doc comment of the declaration itself — change the comment, run `build`, and the page
follows.

One page per contract: the tables below say which page declares a symbol, so open that page
rather than reading the whole surface.

## Entry points

| Symbol | Signature |
| --- | --- |
{{- if .HasDeps }}
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/availables/standard`) |
{{- else }}
| `sandbox.New` | `func() *api.Sandbox` |
{{- end }}

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

## The sandbox api

| Page | Declares |
| --- | --- |
{{- range .PublicApi }}
| [`{{ .Path }}`]({{ .Page }}) | {{ .Symbols }} |
{{- end }}
{{- if .DepsApi }}

## Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Sandbox` struct, filled by `adapters/libs/<name>.Bind(&deps)`.

| Page | Declares |
| --- | --- |
{{- range .DepsApi }}
| [`deps.{{ .Title }}`]({{ .Page }}) | {{ .Symbols }} |
{{- end }}
{{- end }}
