# `{{ .Route.Method }} {{ .Route.Pattern }}`

{{ .Route.Help }}
{{- with .Route.LongDescription }}

{{ . }}
{{- end }}
{{- if .Route.Fields }}

| Entries | Read from | In | Type | Default | Description |
| --- | --- | --- | --- | --- | --- |
{{- range .Route.Fields }}
| `{{ .Id }}` | {{ .Key }} | {{ .In }} | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- with .Route.Body }}

Body: {{ . }}
{{- end }}
{{- if .Route.Examples }}

```bash
{{- range .Route.Examples }}
{{ . }}
{{- end }}
```
{{- end }}

`sandbox/internal/routeslist/{{ .Route.Name }}/` · {{ .Category }} · [every route](doc.md) · [RouteYaml](../RouteYaml/doc.md)
