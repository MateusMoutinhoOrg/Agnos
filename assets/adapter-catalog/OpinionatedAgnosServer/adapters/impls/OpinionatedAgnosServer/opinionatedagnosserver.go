package opinionatedagnosserver

import (
	"{{.Module}}/sandbox/deps"
	opinionatedagnosserver "{{.Module}}/sandbox/deps/OpinionatedAgnosServer"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// Bind fills deps.Deps.OpinionatedAgnosServer with the agnos server mechanic:
// the request chain of main.go, the binder of run.go, the
// matcher of matches.go, the writers of respond.go and the validators of
// jsonschema.go and formschema.go. Nothing here holds a dep: what a run
// reaches the outside world through comes in its MainProps, and a writer or a
// validator takes the one dep it needs as its first parameter.
func Bind(deps *deps.Deps) {
	deps.OpinionatedAgnosServer = opinionatedagnosserver.Contract{
		Main:           serverMain,
		NewRoute:       newRoute,
		BindRoute:      bindRoute,
		Fail:           fail,
		FailWithCause:  failWithCause,
		FailureOf:      failureOf,
		WriteError:     writeError,
		WriteJSON:      writeJSON,
		WriteText:      writeText,
		Redirect:       redirect,
		ValidateSchema: validateSchema,
		ValidateForm:   validateForm,
		ReadString:     readString,
		ReadInt:        readInt,
		ReadFloat:      readFloat,
		ReadBool:       readBool,
		ReadObject:     readObject,
		ReadItems:      readItems,
		ItemString:     itemString,
		ItemInt:        itemInt,
		ItemFloat:      itemFloat,
		ItemBool:       itemBool,
	}
}

// newRoute returns an empty Route with every slice open.
func newRoute() *opinionatedagnosserver.Route {
	return &opinionatedagnosserver.Route{
		Methods:    []string{},
		Examples:   []string{},
		Paths:      []opinionatedagnosserver.Path{},
		Parameters: []opinionatedagnosserver.Parameter{},
	}
}

// bindRoute copies one declaration into the route a single request runs on,
// with no request, response or failure yet. The dispatch calls it once per
// request, so two requests in flight never share a bound value.
func bindRoute(route *opinionatedagnosserver.Route) *opinionatedagnosserver.Route {
	bound := *route
	bound.Request = serverdeps.Request{}
	bound.Response = serverdeps.Response{}
	bound.Props = nil
	bound.Failure = nil
	return &bound
}

// fail builds the failure a handler refuses a request with.
func fail(status int, field string, message string) error {
	return failWithCause(status, field, message, "")
}

// failWithCause is fail carrying what went wrong underneath.
func failWithCause(status int, field string, message string, cause string) error {
	return &opinionatedagnosserver.RouteFailure{
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// failureOf merges the failure a Handle* file is answering with that file's
// own status and wording: what the failure carries wins, what it leaves empty
// the file fills.
func failureOf(route *opinionatedagnosserver.Route, status int, message string) opinionatedagnosserver.RouteFailure {
	if route.Failure == nil {
		return opinionatedagnosserver.RouteFailure{Status: status, Message: message}
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
