# Routes
{{ if .RouteDocs }}
Every route this server answers, one page each, generated from
`sandbox/internal/routeslist/<name>/route.yaml` ([RouteYaml](../RouteYaml/doc.md)) on each build —
open the one you need rather than this whole page. Hidden routes are not listed.

Routes run lowest `priority` first; every path and every parameter trigger of a route has to
match for it to run. A path nothing matches is `404`; one matched under another method is `405`.
Every parameter of a route is bound and converted before the handler runs — a failure there is
`400`, never the handler's call.
{{- range .RouteDocs }}

## {{ .Category }}

| Route | Answers | Package |
| --- | --- | --- |
{{- range .Routes }}
| [`{{ .Method }} {{ .Pattern }}`]({{ .Name }}.md) | {{ .Help }} | `{{ .Name }}` |
{{- end }}
{{- end }}
{{- else }}
No route is declared yet. Run `{{.GeneratorName}} add-route <name> --trigger /<path> --help "..." --category "..."`,
and every route lands on this page on the next build.
{{- end }}

Statuses and who answers each one are in [RouteYaml](../RouteYaml/doc.md#dispatch).
