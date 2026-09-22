package routeio

import (
	"{{.Module}}/sandbox/api"
)

// Fail is the one way any part of the server layer raises a failure: it records
// what went wrong on the route and hands it to the project's own handler for
// that status — HandleBadRequest, HandleTooLarge, HandleWrongContentType,
// HandleServerError — through sandbox.Server.Fail.
//
// It goes through the api rather than calling sandbox/internal/server directly
// because that package imports every route package, so no route may import it
// back. The field on the api is what crosses that line, and it is filled by the
// generated sandbox/internal/server/new.go.
//
// It returns the error the handler returned — the failure's own message when
// that handler answered it — so a RouteHandler ends on one line:
//
//	return routeio.Fail(sandbox, route, api.StatusFailure, "", "not authorized")
func Fail(sandbox *api.Sandbox, route *api.Route, status int, field string, message string) error {
	return FailWithCause(sandbox, route, status, field, message, "")
}

// FailWithCause is Fail carrying what went wrong underneath — an error's text,
// or the value a handler panicked with. The cause reaches the handler on
// route.Failure and is meant for the log, never for the caller.
func FailWithCause(sandbox *api.Sandbox, route *api.Route, status int, field string, message string, cause string) error {
	route.Failure = &api.RouteFailure{
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	}

	// A sandbox whose server was never built — a route bound by hand rather
	// than by the dispatch — has no handler to reach, so the failure is
	// still reported rather than swallowed.
	if sandbox.Server.Fail == nil {
		return sandbox.Deps.Std.Errorf("%s", message)
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
