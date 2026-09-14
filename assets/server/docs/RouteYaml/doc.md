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
| `paths` | `add-segment` / `remove-segment` |
| `headers` | `add-header` / `remove-header` |
| `params` | `add-param` / `remove-param` |
| `body` | `set-body` |
| `body.json-schema` | `add-body-field` / `remove-body-field` |

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
| `category`, `help`, `long-description`, `examples`, `hidden` | As in [EntriesYaml](../EntriesYaml/doc.md#command-keys); feeds [Routes](../Routes/doc.md) |
| `headers`, `params` | Sequences of fields read from the request headers / the query string |
| `body` | The request body, one object rather than a sequence |

## Segment keys

There is no `path` key: the URL is the concatenation of `paths`, and a segment is one of two
kinds.

| Key | Effect |
|---|---|
| `identifier` | A literal segment, **always starting with `/`** (`/users`, never `users`). `/` alone is the root; no other holds an inner or trailing slash. Excludes `name` and every field key |
| `name` | A captured segment: it matches anything and is bound under this name, already converted |
| `type` | `string` (default), `boolean`, `int`, `float` |
| `description`, `examples`, `min`, `max` | As in [EntriesYaml](../EntriesYaml/doc.md#field-keys) |
| `required` | Always `true` on a capture — `false` is a `verify` violation, and `default` is refused |
| `array` | The capture takes **every segment left** in the path, read back with `GetStrings`/`GetInts`/`GetFloats`. Only on the last entry of `paths`, and its name is declared nowhere else |

The first `identifier` is the trigger that names the route. Match order is by specificity, not
by directory: most `identifier`s first, then the longest `identifier`s, then the routes of
fixed length before the ones taking the rest of the path, then the pattern alphabetically —
without which a route on `/` would swallow one on `/home`.

A route ending in an `array` capture matches **one or more** remaining segments, never zero:
the capture is required like any other, so `/static` does not reach `/static/{rest...}`.

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

`ReadBody(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response)` is generated
into the route's own `new.go`, returning what its `body.type` declares:

| `body.type` | Returns |
|---|---|
| `none` | none is generated |
| `raw` | `([]byte, int)` |
| `text` | `(string, int)` |
| `json` with an object `json-schema` | `(Body, int)` |
| `json` without one | `(*serializables.SerializibleObject, int)` |

Every variant does, in order: `Request.ReadBody(MaxBodyBytes)` (`413`), the `required` check
(`400`) and — for `json` — `routeio.ValidateSchema` against `BodySchema` (`400` on the first
violation, its field path in the response's `field`). A nested object becomes `Body<Path>`; an
object inside an array becomes `Body<Path>Item`.

## Dispatch

`ServerMain` hands every request to one dispatch, which slices the path and tests each route of
`Server.Routes` in match order. Everything but the body is settled before the handler runs.

| Situation | Status | Answered by |
|---|---|---|
| no route matched the path | 404 | dispatch |
| path matched, method differs | 405 | dispatch |
| `content-type` differs | 415 | dispatch |
| `Content-Length` above `max-bytes` | 413 | dispatch |
| invalid header/param, missing `required`, out of `min`/`max` | 400 | dispatch |
| body over `max-bytes` while reading | 413 | `ReadBody` |
| body absent with `required: true`, invalid json, schema rejected | 400 | `ReadBody` |
| handler returned `0`, or panicked | 500 | dispatch |

A failure is written by `routeio.WriteError` as `{"error": "...", "field": "..."}` and logged on
`deps.Std.Log`. A `RouteHandler` returns the status it answered with; the only way it returns
`400`/`413`/`415` is by propagating one from `ReadBody`.
