# PublicApi

Every exported symbol of `{{.Module}}`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns{{if .HasDeps}}, `sandbox/deps/`
the contracts an adapter fills and a caller may replace{{end}}. Each description below is the
doc comment of the declaration itself — change the comment, run `build`, and this page
follows.

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

# The sandbox api
{{- range .PublicApi }}

## `{{ .Path }}`
{{- with .Doc }}

{{ . }}
{{- end }}
{{- template "declarations" . }}
{{- end }}
{{- if .DepsApi }}

# Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Lib` struct, filled by `adapters/libs/<name>.Bind(&deps)`.
{{- range .DepsApi }}

## `deps.{{ .Title }}`

`sandbox/deps/{{ .Name }}`
{{- range .Files }}
{{- with .Doc }}

{{ . }}
{{- end }}
{{- template "declarations" . }}
{{- end }}
{{- end }}
{{- end }}
{{- define "declarations" }}
{{- if .Constants }}

| Constant | Value | Description |
| --- | --- | --- |
{{- range .Constants }}
| `{{ .Name }}` | {{ if .Value }}`{{ .Value }}`{{ end }} | {{ .Doc }} |
{{- end }}
{{- end }}
{{- if .Variables }}

| Variable | Type | Description |
| --- | --- | --- |
{{- range .Variables }}
| `{{ .Name }}` | {{ if .Type }}`{{ .Type }}`{{ end }} | {{ .Doc }} |
{{- end }}
{{- end }}
{{- range .Types }}

### `{{ .Name }}`
{{- with .Doc }}

{{ . }}
{{- end }}
{{- if .Fields }}
{{- if .FieldsDocumented }}

| Field | Type | Description |
| --- | --- | --- |
{{- range .Fields }}
| `{{ .Name }}` | `{{ .Type }}` | {{ .Doc }} |
{{- end }}
{{- else }}

| Field | Type |
| --- | --- |
{{- range .Fields }}
| `{{ .Name }}` | `{{ .Type }}` |
{{- end }}
{{- end }}
{{- end }}
{{- if .Methods }}
{{- if .MethodsDocumented }}

| Method | Description |
| --- | --- |
{{- range .Methods }}
| `{{ .Signature }}` | {{ .Doc }} |
{{- end }}
{{- else }}

| Method |
| --- |
{{- range .Methods }}
| `{{ .Signature }}` |
{{- end }}
{{- end }}
{{- end }}
{{- if .Underlying }}

`type {{ .Name }} {{ .Underlying }}`
{{- end }}
{{- end }}
{{- if .Functions }}

| Function | Description |
| --- | --- |
{{- range .Functions }}
| `{{ .Signature }}` | {{ .Doc }} |
{{- end }}
{{- end }}
{{- end }}
