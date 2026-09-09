# ServerUsage

The server layer mirrors the cli layer file for file: `sandbox/internal/routes/<name>/` is to a
route what `sandbox/internal/commands/<name>/` is to a command, and `route.yaml` is to it what
`entries.yaml` is to a command.

| Concept | CLI | Server |
|---|---|---|
| External input contract | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` |
| Generated dispatch | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/servermain.go` |
| Declared unit | `commands/<name>/entries.yaml` | `routes/<name>/route.yaml` |
| Hand-written half | `handler.go` -> `CommandHandler` | `handler.go` -> `RouteHandler` |
| Install / remove | `cli-init` / `cli-purge` | `server-init` / `server-purge` |

## Bring it up

```bash
{{.Name}} server-init          # serverdeps, the server layer, the health route, start-server
{{.Name}} start-server         # listens on :8080
{{.Name}} start-server --addr :3000 --read-timeout-ms 30000
curl localhost:8080/health
```

`server-init` installs the CLI layer first when the project has none — a server needs a command
that starts it. `server-purge` removes the server layer again and leaves the CLI in place.

A Go caller skips the command entirely:

```go
sandbox := sandbox.New(&deps)
err := sandbox.Server.Serve(api.ServeProps{Addr: ":8080", ReadTimeoutMs: 10000, WriteTimeoutMs: 10000})
```

`Serve` blocks until the server stops.

## Declare a route

```bash
{{.Name}} add-route create-user --trigger /users --method POST --help "Create a user" --category Users
{{.Name}} add-segment tenant --route create-user
{{.Name}} add-header authorization --route create-user --required
{{.Name}} add-param page --route create-user --type int --default 1 --min 1
{{.Name}} set-body create-user --type json --required
{{.Name}} add-body-field email --route create-user --format email --required
{{.Name}} remove-param page --route create-user
{{.Name}} remove-route create-user
```

`add-route` writes `route.yaml` (the declaration) and a stub `handler.go` (yours); `build`
generates `entries.go` and the dispatch arm. One editor per place the declaration holds
something — `add-segment`, `add-header`, `add-param`, `set-body`, `add-body-field`, `set-route`,
each with its `remove-` inverse — so every key of [RouteYaml](../RouteYaml/doc.md) is reachable
from the command line and `route.yaml` is never edited by hand.

`add-segment` takes `--identifier /users` for a literal segment, or a name for a capture; an
identifier is normalized to start with `/`, and an inner or trailing slash is refused. With
`--array` the capture takes every segment left in the path into a `[]T` field
(`{{.Name}} add-segment rest --route static --array` matches `/static/a/b.png`), which only the
last segment of a route may do.
`add-body-field` takes a dotted path (`address.city`), creating the intervening objects in the
`json-schema`.

## Write the handler

```go
func RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int {
	// the body has not been read yet — refuse early if you can
	if !isAuthorized(deps, entries.Authorization) {
		return routeio.WriteError(deps, response, api.StatusFailure, "", "not authorized")
	}
	body, status := entries.ReadBody(deps, response)
	if status != api.StatusOk {
		return status
	}
	return writeJson(deps, response, api.StatusCreated, createUser(deps, entries.Tenant, body))
}
```

`Entries` arrives bound, converted and range-checked; a bad request was already answered `400`
before the handler ran. The handler returns the status it answered with, and propagates the one
`ReadBody` gives it. [Routes](../Routes/doc.md) documents the route on the next build.
