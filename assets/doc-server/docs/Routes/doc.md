# Routes
{{ if .RouteDocs }}
Every route this server answers, one page each, generated from
`sandbox/internal/routeslist/<name>/route.yaml` ([RouteYaml](../RouteYaml/doc.md)) on each build —
open the one you need rather than this whole page. Hidden routes are not listed.

Routes run lowest `priority` first, the `after` phase once the request is answered; the segment
count, every path (its type and its trigger) and every parameter trigger of a route have to match
for it to run. A path nothing answers is `404`; one matched under another method is `405`.
`{{.GeneratorName}} list-routes` is the chain in run order, and
`{{.GeneratorName}} explain-route <METHOD> <path>` says which routes one request reaches.
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
No route is declared yet. Run `{{.GeneratorName}} add-route <name> --pattern '/<path>/{id}' --help "..." --category "..."`,
and every route lands on this page on the next build.
{{- end }}

Statuses and who answers each one are in [RouteYaml](../RouteYaml/doc.md#failures).
