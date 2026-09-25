# RouteYaml

`sandbox/internal/routeslist/<name>/route.yaml` declares one http route. `{{.GeneratorName}} build`
generates two files beside it: `new.go`, the `api.Route` that lands in `Server.Routes` — a 1:1
image of the yaml, built on the generic base of `sandbox/internal/server/route` — and
`entries.go`, the `Entries` struct the route's `InternalPureHandler` is handed, plus the
`ReadBody` its body calls for. The dispatch in `sandbox/internal/server/server/servermain.go` is
generic: it reads every request against those declarations, and nothing about a route is spelled
in Go anywhere else.

| File | Owner |
|---|---|
| `route.yaml` | the editors below |
| `new.go`, `entries.go` | every build |
| `InternalPureHandler.go` | written once by `{{.GeneratorName}} add-route`, then the project's |

Grow the file with the editors of [Workflow](../Workflow/doc.md#change-the-route-surface) — one
per place it holds something — not by hand: they re-render it with keys in alphabetical order
and drop comments.

| Section | Editors |
|---|---|
| route-level keys | `set-route` |
| `paths` | `add-path` / `set-path` / `remove-path` |
| `parameters` | `add-parameter` / `set-parameter` / `remove-parameter` |
| `body` | `set-body` |
| `body.json-schema` | `add-body-field` / `set-body-field` / `remove-body-field` / `import-body` |

`{{.GeneratorName}} show-route <route>` prints the whole of it as a tree, which is what the
file reads as once the schema is more than a few keys deep.

```yaml
methods: [POST]
priority: 0
response-type: application/json
paths:
  - id: Route
    start: 0
    end: -1
    trigger: { type: prefix, value: /users/ }
  - id: Tenant
    start: 1
    end: 1
parameters:
  - id: Page
    key: page
    type: number
    fonts: [query]
    default: "1"
category: Users
help: Create a user under a tenant
body:
  type: json
  required: true
  max-bytes: 1048576
  content-type: application/json
  json-schema:
    type: object
    required: [email]
    properties:
      email: { type: string, format: email, maxLength: 254 }
```

## Route keys

| Key | Effect |
|---|---|
| `methods` | Required, never empty: `GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS` |
| `priority` | **Required.** The rung this route runs on when several match one request: lowest first, never negative |
| `response-type` | **Required.** The `Content-Type` set on the response before the handler runs; the handler may set another |
| `paths` | The slices of the request path it reads. Required and never empty |
| `parameters` | The values it reads from the query string and the headers |
| `category`, `help`, `long-description`, `examples`, `hidden` | As in [EntriesYaml](../EntriesYaml/doc.md#command-keys); feeds [Routes](../Routes/doc.md) |
| `body` | The request body, one object rather than a sequence |

`method`, `headers`, `params` and the `identifier` / `name` spelling of `paths` are the older
declaration; `verify` names each one it finds.

## Path keys

A path reads the request segments from `start` to `end` — both inclusive, `-1` the last one — as
`/` followed by them joined by `/`, and binds that text to `Entries.<id>`. The root path, with no
segment at all, reads as `/` for a path of `start: 0, end: -1`.

| Key | Effect |
|---|---|
| `id` | The `Entries` field the slice binds to: an exported Go name, unique across `paths` and `parameters`, never `FullRoute` or `Body` |
| `start` | Index of the first segment. Default `0` |
| `end` | Index of the last segment, `-1` the last one. Default `-1` |
| `trigger` | `{type, value}`: the slice has to match it for the route to run at all |
| `description` | One-line help text |

A request with no such slice — fewer segments than `start` — is not for the route. Every path
has to find its slice and match its trigger, in any order.

| `trigger.type` | Matches a slice that |
|---|---|
| `equal` | is exactly `value` |
| `prefix` | begins with `value` |
| `suffix` | ends with `value` |
| `regex` | `value`, a regular expression, matches |

```yaml
paths:
  - { id: Route, start: 0, end: 0, trigger: { type: equal, value: /static } }
  - { id: Item, start: 1, end: -1 }
```

```
GET /static          404, no segment 1
GET /static/a        Entries.Item = "/a"
GET /static/a/b.png  Entries.Item = "/a/b.png"
```

## Parameter keys

| Key | Effect |
|---|---|
| `id` | The `Entries` field the value binds to, as for a path |
| `key` | The query key or header name it is read under; a header is matched without regard to case. Default `id` |
| `type` | `string` (default, `string`), `number` (`float64`), `boolean` (`bool`, `true`/`1`/`false`/`0`), `datetime` (`string`, RFC 3339), `string-array` (`[]string`: every occurrence of a query key, a header's comma-separated values) |
| `fonts` | Where it is read from, in order: `query`, `header`. The first that brings a value wins |
| `required` | No font brings it: `400` through `HandleBadRequest`. Never on a `boolean`, never with `default` |
| `default` | The literal bound when no font brings it |
| `trigger` | As on a path: the value has to match for the route to run at all |
| `description`, `examples` | One-line help text, whole requests |

A parameter that fails its `trigger` does not fail the route — this route is not the one for the
request, so the chain moves on and, if nothing else answers, `HandleNotFound` does. That is the
difference from `required`, which says the route **is** the one and the request is malformed.

```yaml
methods: [GET]
priority: 5
response-type: application/json
paths:
  - { id: Route, start: 0, end: -1, trigger: { type: prefix, value: /admin } }
parameters:
  - { id: Authorization, key: authorization, type: string, fonts: [header], trigger: { type: prefix, value: Bearer } }
```

```
GET /admin                                  404, this route never ran
GET /admin  Authorization: Basic abc        404, this route never ran
GET /admin  Authorization: Bearer abc       this route
```

## Entries and `InternalPureHandler`

`entries.go` declares one struct, every field tagged with the id `RequestHandler` fills it by:

```go
type Entries struct {
	FullRoute string   `id:"FullRoute"` // the whole request path, always
	Route     string   `id:"Route"`     // one per path
	Tenant    string   `id:"Tenant"`
	Page      float64  `id:"Page"`      // one per parameter, in its type
}
```

The hand-written half is one function, which `verify` holds to this signature:

```go
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error
```

`route` is a copy of the declaration, made per request by `api.BindRoute`, so two requests in
flight never share a value. The generic `RequestHandler` builds `Entries` and calls the handler
through `Deps.Reflectdeps`, since every route's `Entries` is a type of its own.

## Body keys

| Key | Effect |
|---|---|
| `type` | `none` (default, no `ReadBody` at all), `raw` (`[]byte`), `text` (`string`), `json` |
| `required` | An absent or empty body is `400` |
| `max-bytes` | A longer body is `413`. Default `1048576` |
| `content-type` | A divergent one is `415`. Default `application/json` for `type: json` |
| `json-schema` | A subset of JSON Schema, only with `type: json` |

Supported schema keywords: `type` (`object`/`array`/`string`/`integer`/`number`/`boolean`/
`null`), `properties`, `required`, `additionalProperties`, `items`, `enum`, `const`, `minimum`,
`maximum`, `exclusiveMinimum`, `exclusiveMaximum`, `minLength`, `maxLength`, `pattern`,
`minItems`, `maxItems`, `uniqueItems`, `format` (`email`, `uuid`, `date-time`, `uri`),
`nullable`. Anything else (`$ref`, `oneOf`, `allOf`, `anyOf`, `patternProperties`) fails the
build rather than being ignored.

## Generated `ReadBody`

`ReadBody(sandbox *api.Sandbox, route *api.Route)` is generated into the route's own
`entries.go`, returning what its `body.type` declares:

| `body.type` | Returns |
|---|---|
| `none` | none is generated |
| `raw` | `([]byte, error)` |
| `text` | `(string, error)` |
| `json` with an object `json-schema` | `(Body, error)` |
| `json` without one | `(*serializables.SerializibleObject, error)` |

Every variant does, in order: `Request.ReadBody(MaxBodyBytes)` (`413`), the `required` check
(`400`) and — for `json` — `routeio.ValidateSchema` against `BodySchema` (`400` on the first
violation, its field path in the response's `field`). A nested object becomes `Body<Path>`; an
object inside an array becomes `Body<Path>Item`.

A non-nil error means the request has already been answered, by whichever `Handle*` file of
`sandbox/internal/server/errors/` the failure belongs to, so the handler only has to return it.

## The chain

`ServerMain` hands every request to one dispatch, which walks `Server.Routes` in `priority`
order — lowest rung first, by name within one — and runs every route whose `IsActionable` says
the request is for it: the method is one of `methods`, and every path and every parameter
trigger matches.

A route answers the request by **setting a status**. A route that writes no status has declined,
and the next one runs; that is the whole of what makes a route a middleware, and nothing else
distinguishes one:

```go
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	sandbox.Deps.Std.Log("%s\n", entries.FullRoute)
	return nil // no status: the next route of the chain runs
}
```

```yaml
# routeslist/logger/route.yaml — runs first, answers nothing
methods: [GET]
priority: 0
response-type: text/plain
paths:
  - { id: Route, start: 0, end: -1, trigger: { type: prefix, value: / } }

# routeslist/hello/route.yaml — runs after it, answers
methods: [GET]
priority: 5
response-type: text/plain
paths:
  - { id: Route, start: 0, end: -1, trigger: { type: equal, value: /hello } }
```

A middleware in front of everything sits on a lower rung than the routes it guards — `0` for
it, `5` for them.

## Failures

Nothing the dispatch does writes a response. Every way a request can end without a route
answering it is handed to one of six files of `sandbox/internal/server/errors/`, each written once by
`{{.GeneratorName}} server-init` and never regenerated:

| Situation | File | Status |
|---|---|---|
| no route matched, or every matching route declined | `handle_not_found.go` | 404 |
| the path matched under another method | `handle_method_not_allowed.go` | 405 |
| a parameter that will not convert, a missing `required` one, a body the schema rejected | `handle_bad_request.go` | 400 |
| `Content-Length` or the body itself above `max-bytes` | `handle_too_large.go` | 413 |
| `content-type` differs | `handle_wrong_content_type.go` | 415 |
| a handler returned an error without answering, or panicked | `handle_server_error.go` | 500 |

Each holds one function, writes the response itself and returns what it could not answer:

```go
func HandleNotFound(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	failure := routeio.FailureOf(route, api.StatusNotFound, "route not found")
	return routeio.WriteError(sandbox, response, failure.Status, failure.Field, failure.Message)
}
```

`routeio.FailureOf` is what the failure says, falling back to what the file says: the two the
dispatch raises with nothing to add — nothing matched, method not allowed — carry no message, so
the wording is the one spelled in that file and changing it there changes what the server says.
The failures that know something the file could not — which field would not bind, and why —
carry their own.

Raise a failure with `routeio.Fail`, from anywhere:

```go
return routeio.Fail(sandbox, route, api.StatusFailure, "", "not authorized")
```

It reaches the right file through `sandbox.Server.Fail`, which is a field on the api rather than
a call, because a route package may not import `sandbox/internal/server/server` — that package
imports every route. `routeio.WriteError` is the writer underneath, and the default body every one of
them produces is `{"error": "...", "field": "..."}`, logged on `deps.Std.Log` as it is written.

A `Handle*` file answers a failure and never raises one: `routeio.Fail` from inside one comes
back to it.
