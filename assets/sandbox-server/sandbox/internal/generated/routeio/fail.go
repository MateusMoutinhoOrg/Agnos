package routeio

import (
	"{{.Module}}/sandbox/api"
)

// Fail is how an InternalPureHandler refuses a request: it builds the failure
// and returns it as an error, and the handler returns it in turn. A handler is
// handed no route to raise it on, so it is RequestHandler — which holds the
// bound route — that raises what comes back through Raise, reaching the
// project's own Handle* file for that status:
//
//	return routeio.Fail(sandbox, api.StatusUnauthorized, "authorization", "invalid token")
//
// A message left empty is filled by that file's own wording.
func Fail(sandbox *api.Sandbox, status int, field string, message string) error {
	return FailWithCause(sandbox, status, field, message, "")
}

// FailWithCause is Fail carrying what went wrong underneath — an error's text,
// meant for the log, never for the caller.
func FailWithCause(sandbox *api.Sandbox, status int, field string, message string, cause string) error {
	return &api.RouteFailure{
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// Raise is the one way the server layer itself raises a failure on a bound
// route — the dispatch, RequestHandler, a generated ReadBody: it records what
// went wrong on the route and hands it to the project's own handler for that
// status — HandleBadRequest, HandleTooLarge, HandleWrongContentType,
// HandleServerError — through sandbox.Server.Fail.
//
// It goes through the api rather than calling sandbox/internal/server/errors
// because the package that routes to it, sandbox/internal/generated/server/server,
// imports every route package, so no route may import it back. The field on
// the api is what crosses that line, and it is filled by the generated
// sandbox/internal/generated/server/server/new.go.
//
// It returns the error the handler returned — the failure's own message when
// that handler answered it.
func Raise(sandbox *api.Sandbox, route *api.Route, status int, field string, message string) error {
	return RaiseWithCause(sandbox, route, status, field, message, "")
}

// RaiseWithCause is Raise carrying what went wrong underneath — an error's
// text, or the value a handler panicked with. The cause reaches the handler on
// route.Failure and is meant for the log, never for the caller.
func RaiseWithCause(sandbox *api.Sandbox, route *api.Route, status int, field string, message string, cause string) error {
	return RaiseFailure(sandbox, route, &api.RouteFailure{
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	})
}

// RaiseFailure raises one failure already built — what a handler returned
// through Fail — on the bound route it was returned from.
func RaiseFailure(sandbox *api.Sandbox, route *api.Route, failure *api.RouteFailure) error {
	route.Failure = failure

	// A sandbox whose server was never built — a route bound by hand rather
	// than by the dispatch — has no handler to reach, so the failure is
	// still reported rather than swallowed.
	if sandbox.Server.Fail == nil {
		return sandbox.Deps.Std.Errorf("%s", failure.Message)
	}

	return sandbox.Server.Fail(route)
}

// FailureOf is the failure one of the project's Handle* files is answering. A
// file stands for one status and one wording, and passes them here as the
// fallback: a failure that carries its own — a field that would not bind, a
// body the schema rejected — answers with that, because it says something no
// file could say in advance, and one that carries none answers with the file's.
//
// That is the split that makes a Handle* file worth editing. The dispatch
// raises "nothing matched" and "method not allowed" with no message at all, so
// the wording comes from handle_not_found.go and handle_method_not_allowed.go
// and changing it there changes what the server says. It also keeps every
// handler free of a nil check.
func FailureOf(route *api.Route, status int, message string) api.RouteFailure {
	if route.Failure == nil {
		return api.RouteFailure{Status: status, Message: message}
	}

	failure := *route.Failure
	if failure.Status == 0 {
		failure.Status = status
	}
	if failure.Message == "" {
		failure.Message = message
	}
	return failure
}
