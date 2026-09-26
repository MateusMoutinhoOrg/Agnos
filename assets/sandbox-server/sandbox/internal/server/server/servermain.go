package server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeio"
)

// ServerMain opens the port through sandbox.Deps.Serverdeps — which routes
// nothing — and hands every request to dispatch, whatever its method or path.
// The first request to stop the process — an interrupt, a termination —
// shuts it down gracefully through sandbox.Deps.Signaldeps: no new request is
// taken, and the ones in flight get ShutdownTimeoutMs to finish, after which
// Listen returns nil.
// Nothing here is generated per route: every route is one declaration built by
// its own NewRoute and collected by sandbox/internal/server/server/new.go, so
// this file is the same in every project.
func ServerMain(sandbox *api.Sandbox, props api.ServeProps) error {
	server := sandbox.Deps.Serverdeps.NewServer(serverdeps.ServerProps{
		Addr:              props.Addr,
		ReadTimeoutMs:     props.ReadTimeoutMs,
		WriteTimeoutMs:    props.WriteTimeoutMs,
		ShutdownTimeoutMs: props.ShutdownTimeoutMs,
		Handler: func(request serverdeps.Request, response serverdeps.Response) {
			dispatch(sandbox, request, response)
		},
	})

	sandbox.Deps.Signaldeps.OnInterrupt(func() {
		sandbox.Deps.Std.Log("server shutting down \n")
		if err := server.Shutdown(); err != nil {
			sandbox.Deps.Std.Error("server shutdown: %s \n", err.Error())
		}
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
// one runs in turn until one of them **answers** — sets a status, or writes a
// byte, which sends a 200. A handler that does neither has declined, so the
// next route runs — that is the whole of what makes a middleware a
// middleware. A handler that returns an error without answering ends the
// chain too, through HandleServerError. Every route of one request shares one
// Locals, which is how a middleware hands what it learned to the routes after
// it.
//
// Nothing here writes a response itself. Every way a request can end without a
// route answering it — nothing matched, matched under another method, a value
// that will not bind, a panic — is raised through routeio.Fail and answered by
// one of the project's own Handle* files. Once it is answered, whatever
// answered it, the routes of the `after` phase run.
func dispatch(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response) {
	tracked, status := routeio.Tracked(response)
	locals := map[string]any{}

	answer(sandbox, request, tracked, status, locals)
	runAfter(sandbox, request, response, status(), locals)
}

// answer runs the chain for one request and, when no route answered it,
// raises the failure that says why. A HEAD request nothing declares HEAD for
// is run again as the GET it asks the headers of: net/http drops the body a
// HEAD answer writes.
func answer(sandbox *api.Sandbox, request serverdeps.Request, tracked serverdeps.Response, status func() int, locals map[string]any) {
	defer recoverRoute(sandbox, request, tracked, status, locals)

	ran, method_mismatch := runChain(sandbox, request, tracked, status, locals)
	if status() != 0 {
		return
	}

	if !ran && request.GetMethod() == "HEAD" {
		as_get := request
		as_get.GetMethod = func() string { return "GET" }
		ran, _ = runChain(sandbox, as_get, tracked, status, locals)
		if status() != 0 {
			return
		}
	}

	// A path some route answers under another method is the one failure the
	// chain can tell apart from a path nothing knows at all — and only when
	// no route ran, since a chain that ran and wrote nothing is the project
	// declining to answer, not the method being wrong.
	if !ran && method_mismatch {
		failRequest(sandbox, request, tracked, locals, api.StatusMethodNotAllowed)
		return
	}

	failRequest(sandbox, request, tracked, locals, api.StatusNotFound)
}

// runChain runs every route of the `before` phase the request is for, until
// one answers or fails. It reports whether a route declaring its methods ran,
// and whether a route matched the path under another method. A route on ANY —
// a middleware, most often — says nothing about which methods the path
// takes, so it running does not keep a 405 from being told apart.
func runChain(sandbox *api.Sandbox, request serverdeps.Request, tracked serverdeps.Response, status func() int, locals map[string]any) (bool, bool) {
	method_mismatch := false
	ran := false

	for _, declared := range sandbox.Server.Routes {
		if declared.After {
			continue
		}

		bound := api.BindRoute(declared)
		bound.Request = request
		bound.Response = tracked
		bound.Locals = locals

		if !bound.IsActionable(bound) {
			if bound.MatchesPath(bound) && !accepts(bound, request.GetMethod()) {
				method_mismatch = true
			}
			continue
		}

		if !acceptsAny(bound) {
			ran = true
		}
		err := bound.RequestHandler(bound)

		// The answer is what ends the chain, so a handler that answered
		// *and* returned something has answered: what it returned is
		// reported and goes no further.
		if status() != 0 {
			if err != nil {
				sandbox.Deps.Std.Log("route %s: %s \n", bound.Name, err.Error())
			}
			return ran, method_mismatch
		}
		if err != nil {
			routeio.FailWithCause(sandbox, bound, api.StatusFailure, "", "", err.Error())
			return ran, method_mismatch
		}
	}

	return ran, method_mismatch
}

// runAfter runs every route of the `after` phase the request is for, once it
// has been answered. They read the status it was answered with through
// routeio.AnsweredStatus and cannot change the answer: their response is
// frozen, and a panic in one is reported and goes no further.
func runAfter(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, status int, locals map[string]any) {
	routeio.SetAnsweredStatus(locals, status)
	frozen := routeio.Frozen(sandbox, response)

	for _, declared := range sandbox.Server.Routes {
		if !declared.After {
			continue
		}

		bound := api.BindRoute(declared)
		bound.Request = request
		bound.Response = frozen
		bound.Locals = locals

		if !bound.IsActionable(bound) {
			continue
		}
		runAfterRoute(sandbox, bound)
	}
}

// runAfterRoute runs one `after` route, reporting what it returned or panicked
// with instead of answering it: the request already has its answer.
func runAfterRoute(sandbox *api.Sandbox, bound *api.Route) {
	defer func() {
		if failure := recover(); failure != nil {
			sandbox.Deps.Std.Error("after route %s panicked: %v\n", bound.Name, failure)
		}
	}()

	if err := bound.RequestHandler(bound); err != nil {
		sandbox.Deps.Std.Log("after route %s: %s \n", bound.Name, err.Error())
	}
}

// acceptsAny reports a route declared for every method.
func acceptsAny(route *api.Route) bool {
	return len(route.AcceptMethods) == 1 && route.AcceptMethods[0] == api.AnyMethod
}

// accepts reports whether a request method is one of the route's.
func accepts(route *api.Route, method string) bool {
	for _, accepted := range route.AcceptMethods {
		if accepted == method || accepted == api.AnyMethod {
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
func failRequest(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, locals map[string]any, status int) {
	route := api.NewRoute()
	route.Request = request
	route.Response = response
	route.Locals = locals

	routeio.Fail(sandbox, route, status, "", "")
}

// recoverRoute turns a panicking handler into one answered request, so a single
// bad route cannot take the process down with it. A handler that panicked after
// answering has already answered: the panic is reported and nothing is written
// over it.
func recoverRoute(sandbox *api.Sandbox, request serverdeps.Request, response serverdeps.Response, status func() int, locals map[string]any) {
	failure := recover()
	if failure == nil {
		return
	}

	sandbox.Deps.Std.Error("route panicked: %v\n", failure)
	if status() != 0 {
		return
	}

	route := api.NewRoute()
	route.Request = request
	route.Response = response
	route.Locals = locals

	routeio.FailWithCause(sandbox, route, api.StatusFailure, "", "",
		sandbox.Deps.Std.Sprintf("%v", failure))
}
