# `{{ .Database.Type }}`

`sandbox/internal/databases/{{ .Database.Package }}/`, keys under `{{ .Database.Prefix }}`. Build one with
`{{ .Database.Package }}.New(sandbox)` — it touches no key, so building one is free.
{{- range .Database.Tables }}

## `{{ .Name }}`

| Field | Type | Required | Target |
| --- | --- | --- | --- |
{{- range .Fields }}
| `{{ .Name }}` | {{ .Type }} | {{ .Required }} | {{ .Target }} |
{{- end }}
{{- end }}
{{- if .Database.Methods }}

## Methods

| Method | What it does |
| --- | --- |
{{- range .Database.Methods }}
| `{{ .Signature }}` | {{ .Help }} |
{{- end }}
{{- end }}

A query this page does not list goes in `methods_custom.go`, hand-written beside these and
rewritten by no build.

[every database](doc.md)
