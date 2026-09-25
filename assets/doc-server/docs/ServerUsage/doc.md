# ServerUsage

The server layer mirrors the cli layer file for file: `sandbox/internal/routeslist/<name>/` is to
a route what `sandbox/internal/commands/<name>/` is to a command, and `route.yaml` is to it what
`entries.yaml` is to a command.

| Concept | CLI | Server |
|---|---|---|
| External input contract | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` |
| Dispatch | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/server/servermain.go` |
| Declared unit | `commands/<name>/entries.yaml` | `routeslist/<name>/route.yaml` |
| Generated declaration | `new.go` -> `NewCommand` | `new.go` -> `NewRoute`, `entries.go` -> `Entries` |
| Generic matcher and binder | the dispatch | `sandbox/internal/server/route/` (`IsActionable`, `RequestHandler`) |
| Surface on the sandbox | `Cli.Commands` | `Server.Routes` |
| Built by | `sandbox/internal/cli/new.go` | `sandbox/internal/server/server/new.go` |
| Hand-written half | `handler.go` -> `CommandHandler` | `InternalPureHandler.go` -> `InternalPureHandler` |
| Answer to bad input | the dispatch, exit 2 | `sandbox/internal/server/errors/handle_*.go`, yours |
| Install / remove | `cli-init` / `cli-purge` | `server-init` / `server-purge` |

## Bring it up

```bash
{{.GeneratorName}} server-init  # serverdeps, the server layer, the health route, start-server
{{.Name}} start-server  # listens on :8080
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
{{.GeneratorName}} add-route create-user --trigger /users/ --trigger-type prefix --method POST --help "Create a user" --category Users
{{.GeneratorName}} add-route logger --trigger / --trigger-type prefix --priority 0 --response-type text/plain --help "Logs every request" --category Server
{{.GeneratorName}} add-path tenant --route create-user --start 1 --end 1
{{.GeneratorName}} add-parameter authorization --route create-user --font header --required
{{.GeneratorName}} add-parameter page --route create-user --type number --default 1
{{.GeneratorName}} set-body create-user --type json --required
{{.GeneratorName}} add-body-field email --route create-user --format email --required
{{.GeneratorName}} set-parameter page --route create-user --font query --font header
{{.GeneratorName}} show-route create-user
{{.GeneratorName}} remove-parameter page --route create-user
{{.GeneratorName}} remove-route create-user
```

`add-route` writes `route.yaml` (the declaration — `priority` and `response-type` always
included, and one path, `Route`, reading the whole request path against `--trigger`) and a stub
`InternalPureHandler.go` (yours); `build` generates `new.go`, the `api.Route` that lands in
`Server.Routes`, and `entries.go`, the `Entries` the handler is handed. One editor per place the
declaration holds something — `add-path`, `add-parameter`, `set-body`, `add-body-field`,
`set-route`, each with its `remove-` inverse — so every key of [RouteYaml](../RouteYaml/doc.md)
is reachable from the command line and `route.yaml` is never edited by hand.

Each `add-` has a `set-` beside it — `set-path`, `set-parameter`, `set-body-field` — which edits
the declaration that is there instead of replacing it: the keys given are written over the ones
already declared, `--clear <key>` takes one off, `--rename` changes the name it answers to, and
the result goes through the same constructor the `add-` side calls.

`add-path` reads a slice of the request path, `--start` to `--end` (`-1` the last segment), into
`Entries.<Id>`; with `--trigger` (and `--trigger-type equal|prefix|suffix|regex`) the route only
runs when the slice matches. `add-parameter` reads one value from the `--font`s given, in order;
its `--trigger` is a condition on the **value**, and the route runs only when it holds
([RouteYaml](../RouteYaml/doc.md#parameter-keys)).
`add-body-field` takes a dotted path (`address.city`), creating the intervening objects in the
`json-schema`. `import-body` declares a whole payload at once from an example of it:

```bash
{{.GeneratorName}} import-body create-user --file payload.json --required --infer-format
{{.GeneratorName}} import-body create-user --json '{"email":"a@b.co","age":30,"tags":["x"]}'
```

It infers a type per key, the objects and lists around them, `--required` for every key the
example carries, and — with `--infer-format` — the `email`, `uuid`, `date-time` and `uri` a
string spells. A property already declared is never written over; `--replace` starts the schema
over instead. What it infers is a starting point, and every bound after that is
`set-body-field`'s.

`show-route <route>` prints the declaration as a tree — the request line, the paths, the
parameters and the body schema property by property, with the keywords declared on each. It
writes nothing and runs no build.

## Write the handler

```go
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	// the body has not been read yet — refuse early if you can
	if !isAuthorized(sandbox, entries.Authorization) {
		return routeio.Fail(sandbox, route, api.StatusFailure, "", "not authorized")
	}
	body, err := ReadBody(sandbox, route)
	if err != nil {
		return err
	}
	response.SetStatus(api.StatusCreated)
	response.Write(payload(sandbox, createUser(sandbox, entries.Tenant, body)))
	return nil
}
```

`entries` arrives bound and converted, and the response already carries the route's
`response-type`; a bad request was already answered `400` before the handler ran. Every value is
a field of `Entries`, named by its id ([RouteYaml](../RouteYaml/doc.md#entries-and-internalpurehandler)).

**Setting a status is what answers the request.** A handler that writes none has declined, and
the next route matching this request runs — that is the whole of what a middleware is. Returning
an error means "I could not answer this", and hands it to `handle_server_error.go`; returning
`nil` means "done" or "not mine", which the written status tells apart.
[Routes](../Routes/doc.md) documents the route on the next build, and
[RouteYaml](../RouteYaml/doc.md#the-chain) has the chain in full.

## Answer the failures

`server-init` writes six more files into `sandbox/internal/server/errors/`, one per way a request can
end without a route answering it — `handle_not_found.go`, `handle_method_not_allowed.go`,
`handle_bad_request.go`, `handle_too_large.go`, `handle_wrong_content_type.go` and
`handle_server_error.go`. Each is yours: written
once, never regenerated. Editing what your server says when nothing matches is editing
`handle_not_found.go` and nothing else. The table and the shape are in
[RouteYaml](../RouteYaml/doc.md#failures).

A Go caller reads the same surface without a socket: `Server.Routes` is every declared route,
in run order, and `api.BindRoute` copies one into the route a single request runs on.
