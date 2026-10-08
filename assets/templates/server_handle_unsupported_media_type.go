package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// HandleWrongContentType answers a request whose Content-Type is not the
// one its route declares.
//
// It is a route handler like any other — same signature, it writes the response
// itself, it returns what it could not answer — and it is **yours**: written
// once by `{{.GeneratorName}} server-init` and never regenerated, so whatever
// you put here is what your server says. The default is the same JSON shape
// every other failure carries, {"error": "...", "field": "..."}.
//
// What went wrong is on `route.Failure`, read through Deps.OpinatedAgnosServer.FailureOf so a
// route carrying none still answers something. The route itself is bound when a
// declared route raised the failure and bare when none did, so
// route.Request reads the request either way.
//
// Answer a failure here; never raise one: a failure raised from here comes back to this file.
func HandleWrongContentType(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	failure := sandbox.Deps.OpinatedAgnosServer.FailureOf(route, api.StatusUnsupportedMedia, "this route does not accept that content type")

	sandbox.Deps.Std.Log("route error %d %s %s \n", failure.Status, failure.Field, failure.Message)
	return sandbox.Deps.OpinatedAgnosServer.WriteError(sandbox.Deps.Serializables, response, failure.Status, failure.Field, failure.Message)
}
