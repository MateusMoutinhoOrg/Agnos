# RouteYaml

`sandbox/internal/routeslist/<name>/route.yaml` declares one http route. `{{.GeneratorName}} build`
generates two files beside it: `new.go`, the `api.Route` that lands in `Server.Routes` — a 1:1
image of the yaml, built on the generic base of `sandbox/internal/generated/server/route` — and
`entries.go`, the `Entries` struct the route's `InternalPureHandler` is handed, plus the
`ReadBody` its body calls for. The dispatch in `sandbox/internal/generated/server/server/servermain.go` is
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
| route-level keys | `set-route` (`--before` / `--after` place it one rung from another route) |
| `paths` | `add-path` / `set-path` / `remove-path` |
| `parameters` | `add-parameter` / `set-parameter` / `remove-parameter` |
| `body` | `set-body` |
| `body.json-schema` | `add-body-field` / `set-body-field` / `remove-body-field` / `import-body` |

`{{.GeneratorName}} show-route <route>` prints the whole of it as a tree, with its place in the
chain. Three more readers write nothing: `list-routes` (every route in run order),
`explain-route <METHOD> <path> [--header k=v] [--cookie k=v]` (which routes one request reaches,
and why each other one is skipped — no server needed). `rename-route` and `rebalance-routes`
(every route `--step` rungs apart, in the order they run now) are the two whole-route editors.

`add-route --pattern` declares the paths from one url shape instead of `--trigger`:

| Pattern piece | Becomes |
|---|---|
| literal segments in a row | one path with an `equal` trigger, named after them (`/get-article` -> `GetArticle`) |
| `{name}` | one segment, `Entries.Name string` |
| `{name:integer\|number\|uuid}` | one segment of that type |
| `{*rest}`, last only | the rest of the path, `start: i, end: -1` |
| no `{*…}` | `segments: <n>` on the route: the segment count is fixed |

```
{{.GeneratorName}} add-route get-article --pattern '/get-article/{article:integer}'
GET /get-article/42      Entries.Article = 42
GET /get-article/abc     not this route (the type is part of the match)
GET /get-article/42/x    not this route (segments: 2)
```

`add-route --middleware` writes `methods: [ANY]`, a `prefix` trigger on `/` unless `--trigger`
says otherwise, `priority: 10` and a stub that answers nothing.

```yaml
methods: [POST]
priority: 100
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
| `methods` | Required, never empty: `GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS`, or `ANY` alone (`*` on the command line) |
| `priority` | **Required.** The rung this route runs on when several match one request: lowest first, never negative. `add-route` writes `100`, `10` for a `--middleware` |
| `response-type` | **Required.** The `Content-Type` set on the response before the handler runs; the handler may set another |
| `segments` | The segment count the request path has to have, `≥ 1`. Absent: any count |
| `phase` | `before` (default, omitted): a rung of the chain. `after`: runs once the request has been answered, see [The chain](#the-chain) |
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
| `type` | `string` (default, omitted), `integer` (`int`), `number` (`float64`), `uuid` (`string`). Anything but `string` needs `start == end`; a segment that does not convert is a **non-match**, never a `400` |
| `trigger` | `{type, value, negate, ignore-case}`: the slice has to match it for the route to run at all |
| `description` | One-line help text |

A request with no such slice — fewer segments than `start` — is not for the route. Every path
has to find its slice, convert to its type and match its trigger, in any order. A path of one
segment (`start == end`) binds the segment itself (`"42"`); a range binds the slice with its
leading slash (`"/a/b.png"`). Triggers always compare the slice with its leading slash.

| `trigger.type` | On a path slice, matches one that | On a parameter value |
|---|---|---|
| `equal` | is exactly `value` | same |
| `prefix` | is `value` or continues it with a new segment: `/admin` matches `/admin`, `/admin/users`, never `/administrator`; `/` matches all | begins with `value` |
| `text-prefix` | begins with `value`, whatever follows | begins with `value` |
| `suffix` | ends with `value` | same |
| `regex` | `value`, a regular expression, matches | same |

`negate: true` inverts the result; `ignore-case: true` compares without regard to case. The
command line also takes `starts-with`, `ends-with`, `exact`/`equals` and `matches`, and writes the
canonical name. Two paths over the same slice are an AND — "`/admin` but not `/admin/login`":

```yaml
paths:
  - { id: Route, start: 0, end: -1, trigger: { type: prefix, value: /admin } }
  - { id: NotLogin, start: 0, end: -1, trigger: { type: equal, value: /admin/login, negate: true } }
```

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
| `type` | `string` (default, `string`), `integer` (`int`), `number` (`float64`), `boolean` (`bool`, `true`/`1`/`false`/`0`), `datetime` (`string`, RFC 3339), `string-array` (`[]string`: every occurrence of a query key, a header's comma-separated values), `integer-array` (`[]int`, read the same way) |
| `fonts` | Where it is read from, in order: `query`, `header`, `cookie`. The first that brings a value wins |
| `required` | No font brings it: `400` through `HandleBadRequest`. Never on a `boolean`, never with `default` |
| `default` | The literal bound when no font brings it |
| `trigger` | As on a path: the value has to match for the route to run at all |
| `description`, `examples` | One-line help text, whole requests |

A parameter that fails its `trigger` does not fail the route — this route is not the one for the
request, so the chain moves on and, if nothing else answers, `HandleNotFound` does. That is the
difference from `required`, which says the route **is** the one and the request is malformed.

```yaml
methods: [GET]
priority: 100
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
	FullRoute string `id:"FullRoute"` // the whole request path, always
	Route     string `id:"Route"`     // one per path, in its type
	Tenant    string `id:"Tenant"`
	Page      int    `id:"Page"`      // one per parameter, in its type
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
| `type` | `none` (default, no `ReadBody` at all), `raw` (`[]byte`), `text` (`string`), `json`, `form` (`map[string][]string`, `application/x-www-form-urlencoded`) |
| `required` | An absent or empty body is `400` |
| `max-bytes` | A longer body is `413`. Default `1048576` |
| `content-type` | A divergent one is `415`. Default `application/json` for `json`, `application/x-www-form-urlencoded` for `form` |
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
| `form` | `(map[string][]string, error)`; a body that does not parse is `400` |
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
order — lowest rung first, by name within one — and runs every route of the `before` phase whose
`IsActionable` says the request is for it: the method is one of `methods` (or they are `ANY`),
the path has `segments` segments when the route declares a count, and every path and parameter
trigger matches.

A route **answers** the request by setting a status or by writing a byte, which sends a `200`.
A route that does neither has declined, and the next one runs; that is the whole of what makes a
route a middleware. `SetHeader` alone answers nothing, which is how a middleware adds a header to
whatever answers after it.

```go
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	token := routeio.RequestOf(route).GetHeader("Authorization")
	if token == "" {
		return routeio.Fail(sandbox, route, api.StatusUnauthorized, "authorization", "")
	}
	routeio.SetLocal(route, "user", token) // the routes after it read it
	return nil                            // nothing answered: the next route runs
}
```

```bash
{{.GeneratorName}} add-route admin-guard --middleware --trigger /admin --before admin
```

Every route of one request shares `route.Locals`: `routeio.SetLocal(route, key, value)` stores,
`routeio.GetLocal[T](route, key)` reads what a route earlier in the chain stored.

When no route answers:

| Situation | Answer |
|---|---|
| a `HEAD` nothing declares `HEAD` for | the chain runs again as a `GET`; the body is dropped |
| a route with explicit `methods` matched the path under another method, and no route with explicit `methods` ran | `405` — an `ANY` route running does not hide it |
| anything else | `404` |

Once the request is answered — by a route or by a failure — every route of the `after` phase
the request is for runs, lowest rung first. Its response is frozen (a status, a header or a byte
written there is logged and dropped), `routeio.AnsweredStatus(route)` is the status it went out
with, and a panic in one is logged. It is where an access log or a metric goes:

```bash
{{.GeneratorName}} add-route access-log --middleware --phase after
```

`routeio.WriteJSON`, `routeio.WriteText` and `routeio.Redirect` answer in one call; the
`Response` also carries `AddHeader` (a second `Set-Cookie`) and `GetHeader`, the `Request`
`GetHeaders`, `GetHost` and `GetCookie`.

## Failures

Nothing the dispatch does writes a response. Every way a request can end without a route
answering it is handed to one of eight files of `sandbox/internal/server/errors/`, each written
once by `{{.GeneratorName}} build` and never regenerated:

| Situation | File | Status |
|---|---|---|
| no route matched, or every matching route declined | `handle_not_found.go` | 404 |
| the path matched under another method | `handle_method_not_allowed.go` | 405 |
| a parameter that will not convert, a missing `required` one, a body the schema rejected | `handle_bad_request.go` | 400 |
| a route raised `api.StatusUnauthorized` — a guard, most often | `handle_unauthorized.go` | 401 |
| a route raised `api.StatusForbidden` | `handle_forbidden.go` | 403 |
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
a call, because a route package may not import `sandbox/internal/generated/server/server` — that package
imports every route. `routeio.WriteError` is the writer underneath, and the default body every one of
them produces is `{"error": "...", "field": "..."}`, logged on `deps.Std.Log` as it is written.

A `Handle*` file answers a failure and never raises one: `routeio.Fail` from inside one comes
back to it.
