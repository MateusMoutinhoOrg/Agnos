# RouteYaml

`sandbox/internal/routes/<name>/route.yaml` declares one http route. `{{.GeneratorName}} build` generates
`entries.go` (the `Entries` struct and its `ReadBody`) and a match/handle pair in
`sandbox/internal/server/servermain.go` from it. Grow it with the editors of
[Workflow](../Workflow/doc.md#change-the-route-surface) — one per place this file holds
something — not by hand: they re-render it with keys in alphabetical order and drop comments.

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
| `name` | A captured segment: it matches anything and becomes an `Entries` field, already converted |
| `type` | `string` (default), `boolean`, `int`, `float` |
| `description`, `examples`, `min`, `max` | As in [EntriesYaml](../EntriesYaml/doc.md#field-keys) |
| `required` | Always `true` on a capture — `false` is a `verify` violation, and `default` is refused |
| `array` | The capture takes **every segment left** in the path, as a `[]T` field. Only on the last entry of `paths`, and its name is declared nowhere else |

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
GET /static/a        Rest = ["a"]
GET /static/a/b.png  Rest = ["a", "b.png"]
```

## Field keys

`headers` and `params` take the field keys of [EntriesYaml](../EntriesYaml/doc.md#field-keys)
(`name`, `description`, `examples`, `type`, `default`, `required`, `array`, `min`, `max`), with
two differences: `name` **is** the external spelling — the header name, matched without regard
to case, or the query key — and `array: true` is refused in `headers` (in `params` it collects
every occurrence of the key; in `paths` it takes the rest of the path).

One `name` may be declared in more than one place. The `Entries` field is written once and
filled by the first origin, in `paths` → `headers` → `params` order, that brings a value; any
one of them satisfies `required`, and all of them must agree on `type`.

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

| `body.type` | Signature |
|---|---|
| `none` | none is generated |
| `raw` | `([]byte, int)` |
| `text` | `(string, int)` |
| `json` with an object `json-schema` | `(Body, int)` |
| `json` without one | `(*serializables.SerializibleObject, int)` |

Every variant does, in order: `Request.ReadBody(MaxBodyBytes)` (`413`), the `required` check
(`400`) and — for `json` — `routeio.ValidateSchema` against `EntriesSchema` (`400` on the first
violation, its field path in the response's `field`). A second call returns the cached value.
A nested object becomes `Body<Path>`; an object inside an array becomes `Body<Path>Item`.

## Dispatch

`ServerMain` hands every request to one dispatch, which slices the path and tests each route in
match order. Everything but the body is settled before the handler runs.

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
