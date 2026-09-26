package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
)

// HandleNotFound answers every request no route of the chain answered — one
// that matched nothing at all, and one every matching route declined by
// writing no status.
//
// It is a route handler like any other — same signature, it writes the response
// itself, it returns what it could not answer — and it is **yours**: written
// once by `{{.GeneratorName}} server-init` and never regenerated, so whatever
// you put here is what your server says. The default is the same JSON shape
// every other failure carries, {"error": "...", "field": "..."}.
//
// What went wrong is on `route.Failure`, read through routeio.FailureOf so a
// route carrying none still answers something. The route itself is bound when a
// declared route raised the failure and bare when none did, so
// routeio.RequestOf(route) reads the request either way.
//
// Answer a failure here; never raise one. routeio.Fail comes back to this file.
func HandleNotFound(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	failure := routeio.FailureOf(route, api.StatusNotFound, "route not found")

	return routeio.WriteError(sandbox, response, failure.Status, failure.Field, failure.Message)
}
