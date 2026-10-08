package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// HandleMethodNotAllowed answers a request whose path a route declares,
// under a method none of them does.
//
// It is a route handler like any other — same signature, it writes the response
// itself, it returns what it could not answer — and it is **yours**: written
// once by `{{.GeneratorName}} server-init` and never regenerated, so whatever
// you put here is what your server says. The default is the same JSON shape
// every other failure carries, {"error": "...", "field": "..."}.
//
// What went wrong is on `route.Failure`, read through Deps.OpinionatedAgnosServer.FailureOf so a
// route carrying none still answers something. The route itself is bound when a
// declared route raised the failure and bare when none did, so
// route.Request reads the request either way.
//
// Answer a failure here; never raise one: a failure raised from here comes back to this file.
func HandleMethodNotAllowed(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	failure := sandbox.Deps.OpinionatedAgnosServer.FailureOf(route, api.StatusMethodNotAllowed, "method not allowed")

	sandbox.Deps.StdDeps.Logf("route error %d %s %s \n", failure.Status, failure.Field, failure.Message)
	return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, response, failure.Status, failure.Field, failure.Message)
}
