# Routes
{{ if .RouteDocs }}
Every route this server answers, generated from
`sandbox/internal/routes/<name>/route.yaml` ([RouteYaml](../RouteYaml/doc.md)) on each build.
Hidden routes are not listed.

A path is matched segment by segment, most specific route first. A path nothing matches is
`404`; one matched under another method is `405`. Every field below is bound, converted and
range-checked before the handler runs — a failure there is `400`, never the handler's call.
{{- range .RouteDocs }}

## {{ .Category }}
{{- range .Routes }}

### `{{ .Method }} {{ .Pattern }}`

{{ .Help }}
{{- with .LongDescription }}

{{ . }}
{{- end }}
{{- if .Fields }}

| Field | In | Type | Default | Description |
| --- | --- | --- | --- | --- |
{{- range .Fields }}
| `{{ .Key }}` | {{ .In }} | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- with .Body }}

Body: {{ . }}
{{- end }}
{{- if .Examples }}

```bash
{{- range .Examples }}
{{ . }}
{{- end }}
```
{{- end }}
{{- end }}
{{- end }}
{{- else }}
No route is declared yet. Run `{{.GeneratorName}} add-route <name> --trigger /<path> --help "..." --category "..."`,
and every route lands on this page on the next build.
{{- end }}

Statuses and who answers each one are in [RouteYaml](../RouteYaml/doc.md#dispatch).
