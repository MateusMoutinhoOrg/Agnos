package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// HandleInternalServerError answers a request no route could carry out: a handler
// that returned an error without answering, one that panicked, or a failure
// raised with a status no other file of this package stands for.
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
func HandleInternalServerError(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	failure := sandbox.Deps.OpinionatedAgnosServer.FailureOf(route, api.StatusInternalServerError, "the route could not answer this request")

	// Cause is what went wrong underneath — an error's text, or the value a
	// handler panicked with. It is for whoever runs the server, never for
	// whoever called it, so it is logged and never written.
	if failure.Cause != "" {
		sandbox.Deps.StdDeps.Eprintf("route %s failed: %s\n", route.Name, failure.Cause)
	}

	sandbox.Deps.StdDeps.Logf("route error %d %s %s \n", failure.Status, failure.Field, failure.Message)
	return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, response, failure.Status, failure.Field, failure.Message)
}
