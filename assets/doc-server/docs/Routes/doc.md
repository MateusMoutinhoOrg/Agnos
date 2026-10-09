# Routes
{{ if .RouteDocs }}
Every address this server answers. Open one to see what to send, a request you can run as it is,
and what comes back.

The requests call `localhost:3000`, where `{{.ProjectName}} start-server` listens when that port is free —
it prints the address it took. Change it to wherever your server runs.

## How to read an address

| In the address | Means | For example |
| --- | --- | --- |
| `GET`, `POST`, … | the method to send it with; `ANY` takes every one | `curl -X POST …` |
| `/users` | exactly that text | `/users` |
| `{name}` | a value you choose | `/users/{tenant}` -> `/users/acme` |
| `{name:integer}` | a value of that type: `integer`, `number` or `uuid` | `/articles/{id:integer}` -> `/articles/42` |
| `{*name}` | the rest of the address, one part or more | `/files/{*file}` -> `/files/a/b.png` |
| `*` | anything else, or nothing | `/admin/*` -> `/admin`, `/admin/users` |
| `(a\|b)` | one of these words | `/(en\|pt)` -> `/en` |
{{- range .RouteDocs }}

## {{ .Category }}

| Route | What it does |
| --- | --- |
{{- range .Routes }}
| [`{{ .Method }} {{ .Pattern }}`]({{ .Name }}.md) | {{ .Summary }} |
{{- end }}
{{- end }}

## Postman, Swagger and other tools

[openapi.json](openapi.json) holds every route below as an OpenAPI 3.0 document, and the running
server answers the same file at `/openapi.json`:

- **Postman**: *Import*, then pick the file or paste `http://localhost:3000/openapi.json`. Every
  route lands in a folder named after its section.
- **Swagger UI**, or any tool that reads OpenAPI: open `http://localhost:3000/openapi.json`.

A route that runs in front of others for every method — a check of a token, for example — is
not listed on its own there: what it reads is listed on each route it guards, and a token sent
as `Authorization` is the document's sign-in.

## When something goes wrong

| Status | Means |
| --- | --- |
| `400` | Something you sent is missing or has the wrong type or format |
| `401` | You have to identify yourself first — a token, for example |
| `403` | You are identified, but not allowed to do this |
| `404` | No route answers this address |
| `405` | The address exists, but not for this method — a `GET` where it takes a `POST`, for example |
| `413` | The body is too large |
| `415` | The body is not in the format the route reads — check `Content-Type` |
| `500` | The server failed while answering |

Unless the project changed it, the answer to an error is JSON naming what went wrong and, when
it is one value, which one:

```json
{"error": "required parameter 'authorization' is missing", "field": "authorization"}
```

For developers: each page, and `openapi.json`, is generated on every build from
`sandbox/internal/routes/<name>/route.yaml` ([RouteYaml](../RouteYaml/doc.md#openapi)); hidden
routes are left out. `{{.GeneratorName}} list-routes` prints the routes in the order they run, and
`{{.GeneratorName}} explain-route <METHOD> <path>` which ones a request reaches. The error answers are
the eight files of `sandbox/internal/server/errors/` ([RouteYaml](../RouteYaml/doc.md#failures)).
{{- else }}
No route is declared yet. Run `{{.GeneratorName}} add-route <name> --pattern '/<path>/{id}' --summary "..." --category "..."`,
and every route lands on this page on the next build.
{{- end }}
