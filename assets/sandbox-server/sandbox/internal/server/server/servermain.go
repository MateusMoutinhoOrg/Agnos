package server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeio"
)

// ServerMain opens the port through sandbox.Deps.Serverdeps — which routes
// nothing — and hands every request to dispatch, whatever its method or path.
// Nothing here is generated per route: every route is one declaration built by
// its own NewRoute and collected by sandbox/internal/server/server/new.go, so
// this file is the same in every project.
func ServerMain(sandbox *api.Sandbox, props api.ServeProps) error {
	server := sandbox.Deps.Serverdeps.NewServer(serverdeps.ServerProps{
		Addr:           props.Addr,
		ReadTimeoutMs:  props.ReadTimeoutMs,
		WriteTimeoutMs: props.WriteTimeoutMs,
		Handler: func(request serverdeps.Request, response serverdeps.Response) {
			dispatch(sandbox, request, response)
		},
	})

	sandbox.Deps.Std.Log("server listening on %s \n", props.Addr)
	return server.Listen()
}

// dispatch is the whole routing layer: it runs every route of Server.Routes
// the request is for, in the order the collector put them — lowest `priority`
// first. Whether a route is for the request is its own IsActionable's to say;
// binding and running it is its own RequestHandler's.
//
// More than one route may be for one request, which is what a chain is: each
// one runs in turn until one of them **sets a status**. A handler that writes
// no status has not answered, so the next route runs — that is the whole of
// what makes a middleware a middleware. A handler that returns an error
// without answering ends the chain too, through HandleServerError.
//
// Nothing here writes a response itself. Every way a request can end without a
// route answering it — nothing matched, matched under another method, a value
// that will not bind, a panic — is raised through routeio.Fail and answered by
// one of the project's own Handle* files.
func dispatch(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response) {
	tracked, answered := routeio.Tracked(response)
	defer recoverRoute(sandbox, request, tracked, answered)

	method_mismatch := false
	ran := false

	for _, declared := range sandbox.Server.Routes {
		bound := api.BindRoute(declared)
		bound.Request = request
		bound.Response = tracked

		if !bound.IsActionable(bound) {
			if bound.MatchesPath(bound) && !accepts(bound, request.GetMethod()) {
				method_mismatch = true
			}
			continue
		}

		ran = true
		err := bound.RequestHandler(bound)

		// The status is what ends the chain, so a handler that answered
		// *and* returned something has answered: what it returned is
		// reported and goes no further.
		if answered() {
			if err != nil {
				sandbox.Deps.Std.Log("route %s: %s \n", bound.Name, err.Error())
			}
			return
		}
		if err != nil {
			routeio.FailWithCause(sandbox, bound, api.StatusFailure, "", "", err.Error())
			return
		}
	}

	// A path some route answers under another method is the one failure the
	// chain can tell apart from a path nothing knows at all — and only when
	// no route ran, since a chain that ran and wrote nothing is the project
	// declining to answer, not the method being wrong.
	if !ran && method_mismatch {
		failRequest(sandbox, request, tracked, api.StatusMethodNotAllowed)
		return
	}

	failRequest(sandbox, request, tracked, api.StatusNotFound)
}

// accepts reports whether a request method is one of the route's.
func accepts(route *api.Route, method string) bool {
	for _, accepted := range route.AcceptMethods {
		if accepted == method {
			return true
		}
	}
	return false
}

// failRequest raises a failure that belongs to no route — nothing matched the
// path, or nothing matched it under this method. The handler still gets an
// api.Route, carrying the request and the response and nothing else, so every
// Handle* file reads the same whichever failure brought it there.
//
// It carries no message on purpose. There is nothing to say about these two
// beyond the status, so the wording is the project's: routeio.FailureOf fills
// in the one the Handle* file spells, which is what makes editing that file
// change what the server says.
func failRequest(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, status int) {
	route := api.NewRoute()
	route.Request = request
	route.Response = response

	routeio.Fail(sandbox, route, status, "", "")
}

// recoverRoute turns a panicking handler into one answered request, so a single
// bad route cannot take the process down with it. A handler that panicked after
// writing its status has already answered: the panic is reported and nothing is
// written over it.
func recoverRoute(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, answered func() bool) {
	failure := recover()
	if failure == nil {
		return
	}

	sandbox.Deps.Std.Error("route panicked: %v\n", failure)
	if answered() {
		return
	}

	route := api.NewRoute()
	route.Request = request
	route.Response = response

	routeio.FailWithCause(sandbox, route, api.StatusFailure, "", "",
		sandbox.Deps.Std.Sprintf("%v", failure))
}
