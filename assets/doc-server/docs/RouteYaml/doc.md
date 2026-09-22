# RouteYaml

`sandbox/internal/routes/<name>/route.yaml` declares one http route. `{{.GeneratorName}} build` generates
`new.go` from it — the `api.Route` that lands in `Server.Routes`, plus the `ReadBody` its body
calls for. The dispatch in `sandbox/internal/server/servermain.go` is generic: it reads every
request against those declarations, and nothing about a route is spelled in Go anywhere else.

Grow the file with the editors of [Workflow](../Workflow/doc.md#change-the-route-surface) — one
per place it holds something — not by hand: they re-render it with keys in alphabetical order
and drop comments.

| Section | Editors |
|---|---|
| route-level keys | `set-route` |
| `paths` | `add-segment` / `set-segment` / `remove-segment` |
| `headers` | `add-header` / `set-header` / `remove-header` |
| `params` | `add-param` / `set-param` / `remove-param` |
| `body` | `set-body` |
| `body.json-schema` | `add-body-field` / `set-body-field` / `remove-body-field` / `import-body` |

`{{.GeneratorName}} show-route <route>` prints the whole of it as a tree, which is what the
file reads as once the schema is more than a few keys deep.

```yaml
method: POST
paths:
  - identifier: "/users"
  - name: tenant
    type: string
    required: true
  - identifier: "/create"
category: Users
help: Create a user under a tenant
params:
  - name: page
    type: int
    default: "1"
    min: 1
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
| `paths` | The URL's segments, in order. Required and never empty |
| `method` | `GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS`. Default `GET` |
| `priority` | The rung this route runs on when several match one request: lowest first, never negative. Default `0`, and written back only when it is not |
| `category`, `help`, `long-description`, `examples`, `hidden` | As in [EntriesYaml](../EntriesYaml/doc.md#command-keys); feeds [Routes](../Routes/doc.md) |
| `headers`, `params` | Sequences of fields read from the request headers / the query string |
| `body` | The request body, one object rather than a sequence |

## Segment keys

There is no `path` key: the URL is the concatenation of `paths`, and a segment is one of three
kinds.

| Key | Effect |
|---|---|
| `identifier` | A literal segment, **always starting with `/`** (`/users`, never `users`). `/` alone is the root; no other holds an inner or trailing slash. Excludes `starts-with-identifier`, `name` and every field key |
| `starts-with-identifier` | The same, loosened to a prefix: the path only has to **begin** with it, and every segment after is left unmatched. It may spell more than one segment (`/api/v1`), `/` alone matches every request, and it is always the last entry of `paths` |
| `name` | A captured segment: it matches anything and is bound under this name, already converted |
| `type` | `string` (default), `boolean`, `int`, `float` |
| `description`, `examples`, `min`, `max` | As in [EntriesYaml](../EntriesYaml/doc.md#field-keys) |
| `required` | Always `true` on a capture — `false` is a `verify` violation, and `default` is refused |
| `array` | The capture takes **every segment left** in the path, read back with `GetStrings`/`GetInts`/`GetFloats`. Only on the last entry of `paths`, and its name is declared nowhere else |

The first literal segment is the trigger that names the route. Run order is `priority` first,
lowest rung to highest; within one rung it is specificity, not directory order: most
`identifier`s first, then the longest `identifier`s, then the routes of fixed length before the
ones taking the rest of the path, then the pattern alphabetically — without which a route on
`/` would swallow one on `/home`.

A route ending in an `array` capture matches **one or more** remaining segments, never zero:
the capture is required like any other, so `/static` does not reach `/static/{rest...}`. A
`starts-with-identifier` is the other way of taking the rest of the path, and matches zero
remaining segments as happily as ten — a route declares one or the other, never both.

```yaml
method: GET
paths:
  - identifier: "/static"
  - name: rest
    type: string
    required: true
    array: true
```

```
GET /static          404
GET /static/a        GetStrings("rest") = ["a"]
GET /static/a/b.png  GetStrings("rest") = ["a", "b.png"]
```

## Field keys

`headers` and `params` take the field keys of [EntriesYaml](../EntriesYaml/doc.md#field-keys)
(`name`, `description`, `examples`, `type`, `default`, `required`, `array`, `min`, `max`), with
two differences: `name` **is** the external spelling — the header name, matched without regard
to case, or the query key — and `array: true` is refused in `headers` (in `params` it collects
every occurrence of the key; in `paths` it takes the rest of the path).

One `name` may be declared in more than one place. The value is bound once and filled by the
first origin, in `paths` → `headers` → `params` order, that brings a value; any one of them
satisfies `required`, and all of them must agree on `type`.

A header or a query parameter may also carry a **match condition**, which is what puts it into
what the route matches on rather than only into what it binds:

| Key | Effect |
|---|---|
| `identifier` | The route runs only when the request brings exactly this value under that name |
| `starts-with-identifier` | The route runs only when the value the request brings begins with this |

A field declares one or the other, never both. A request that fails a condition does not fail
the route — it means this route is not the one for it, so the chain moves on and, if nothing
else answers, `HandleNotFound` does. That is the difference from `required`, which says the
route **is** the one and the request is malformed, and answers `400` through `HandleBadRequest`.

```yaml
method: GET
priority: 5
paths:
  - identifier: "/admin"
headers:
  - name: authorization
    type: string
    starts-with-identifier: "Bearer"
```

```
GET /admin                                  404, this route never ran
GET /admin  Authorization: Basic abc        404, this route never ran
GET /admin  Authorization: Bearer abc       this route
```

## Reading the values

A handler reads a bound value by the name its declaration gives it, never off a field:

| Reader | Returns |
|---|---|
| `route.GetString(name)`, `GetInt`, `GetFloat`, `GetBool` | the first value bound under that name, the zero value when none was |
| `route.GetStrings(name)`, `GetInts`, `GetFloats` | every value bound under it, in order — for an `array` capture or param |
| `route.GetItem(name)` | the raw `[]any` behind them |

`route` is a copy of the declaration, made per request by `api.BindRoute`, so two requests in
flight never share a value.

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

`ReadBody(sandbox *api.Sandbox, route *api.Route)` is generated into the route's own `new.go`,
returning what its `body.type` declares:

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
`sandbox/internal/server/` the failure belongs to, so the handler only has to return it.

## The chain

`ServerMain` hands every request to one dispatch, which slices the path and collects **every**
route that matches it — the method, every path identifier, and every header and query-parameter
condition. Those routes then run in `priority` order, lowest rung first.

A route answers the request by **setting a status**. A route that writes no status has declined,
and the next one runs; that is the whole of what makes a route a middleware, and nothing else
distinguishes one:

```go
func RouteHandler(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	request := routeio.RequestOf(route)
	sandbox.Deps.Std.Log("%s %s\n", request.GetMethod(), request.GetPath())
	return nil // no status: the next route of the chain runs
}
```

```yaml
# routes/logger/route.yaml — runs first, answers nothing
method: GET
paths:
  - starts-with-identifier: "/"

# routes/hello/route.yaml — runs after it, answers
method: GET
priority: 5
paths:
  - identifier: "/hello"
```

Two routes on the same rung are ordered by specificity, so a middleware in front of everything
sits on a lower rung than the routes it guards — `0` for it, `5` for them.

## Failures

Nothing the dispatch does writes a response. Every way a request can end without a route
answering it is handed to one of six files of `sandbox/internal/server/`, each written once by
`{{.GeneratorName}} server-init` and never regenerated:

| Situation | File | Status |
|---|---|---|
| no route matched, or every matching route declined | `handle_not_found.go` | 404 |
| the path matched under another method | `handle_method_not_allowed.go` | 405 |
| invalid header/param, missing `required`, out of `min`/`max`, body the schema rejected | `handle_bad_request.go` | 400 |
| `Content-Length` or the body itself above `max-bytes` | `handle_too_large.go` | 413 |
| `content-type` differs | `handle_wrong_content_type.go` | 415 |
| a handler returned an error without answering, or panicked | `handle_server_error.go` | 500 |

Each holds one function with the route handler's own signature, writes the response itself and
returns what it could not answer:

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
a call, because a route package may not import `sandbox/internal/server` — that package imports
every route. `routeio.WriteError` is the writer underneath, and the default body every one of
them produces is `{"error": "...", "field": "..."}`, logged on `deps.Std.Log` as it is written.

A `Handle*` file answers a failure and never raises one: `routeio.Fail` from inside one comes
back to it.
