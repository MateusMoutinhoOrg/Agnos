# `{{ .Route.Method }} {{ .Route.Pattern }}`
{{- with .Route.Help }}

{{ . }}
{{- end }}
{{- with .Route.LongDescription }}

{{ . }}
{{- end }}

## Try it
{{- range .Route.Requests }}
{{- with .Title }}

{{ . }}
{{- end }}

```bash
{{ .Command }}
```
{{- end }}
{{- if .Route.Examples }}

More examples:

```bash
{{- range .Route.Examples }}
{{ . }}
{{- end }}
```
{{- end }}
{{- if .Route.Address }}

## In the address

| Part | What goes there | Example | Description |
| --- | --- | --- | --- |
{{- range .Route.Address }}
| {{ .Name }} | {{ .Type }} | {{ .Example }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Route.Parameters }}

## Query string, headers and cookies

| Name | Sent in | What goes there | Required | Example | Description |
| --- | --- | --- | --- | --- | --- |
{{- range .Route.Parameters }}
| {{ .Name }} | {{ .Where }} | {{ .Type }} | {{ .Required }} | {{ .Example }} | {{ .Description }}{{ if .From }}{{ if .Description }} — {{ end }}read by [`{{ .From }}`]({{ .FromPage }}), which runs first{{ end }} |
{{- end }}
{{- end }}
{{- with .Route.Body }}

## Body

{{ .Summary }}
{{- if .Fields }}

| Field | What goes there | Required | Rules |
| --- | --- | --- | --- |
{{- range .Fields }}
| {{ .Name }} | {{ .Type }} | {{ .Required }} | {{ .Rules }} |
{{- end }}
{{- end }}
{{- with .Sample }}

Example:

```{{ $.Route.Body.SampleLang }}
{{ . }}
```
{{- end }}
{{- end }}

## What comes back

| Status | Means |
| --- | --- |
{{- range .Route.Statuses }}
| {{ .Code }} | {{ .Meaning }} |
{{- end }}

Any route may also answer `404`, `405` or `500`: see [when something goes wrong](doc.md#when-something-goes-wrong).
{{- if .Route.Middlewares }}

## Runs first

These routes run before this one, on the same request. Any of them may refuse it — with `401`
or `403`, for example — or let it through to this route.

| Route | When |
| --- | --- |
{{- range .Route.Middlewares }}
| [`{{ .Name }}`]({{ .Page }}) | {{ with .Condition }}{{ . }}{{ else }}always{{ end }} |
{{- end }}
{{- end }}

---

For developers: `sandbox/internal/routeslist/{{ .Route.Name }}/` · {{ .Category }} · [every route](doc.md) · [RouteYaml](../RouteYaml/doc.md)
