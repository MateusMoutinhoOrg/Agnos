# {{ .Title }}
{{- with .Source }}

{{ . }}
{{- end }}
{{- range .Files }}
{{- with .Doc }}

{{ . }}
{{- end }}
{{- template "declarations" . }}
{{- end }}

[every contract](doc.md)
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

## `{{ .Name }}`
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
