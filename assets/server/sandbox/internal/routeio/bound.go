package routeio

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// api.Route carries the request it was bound from and the response being
// written as `any`: sandbox/api may name no type of sandbox/deps, so the two
// readers below are where the server layer puts the names back on. Every
// generated NewRoute and every generated ReadBody reads them through these, and
// a hand-written handler is handed the response already read back.
//
// A route nothing bound — one taken straight off sandbox.Routes rather than
// minted per request by the dispatch — reads back as the zero value, whose
// function fields are nil.

// RequestOf is the http request the route was bound from.
func RequestOf(route *api.Route) serverdeps.Request {
	request, _ := route.Request.(serverdeps.Request)
	return request
}

// ResponseOf is the http response the route is answering on.
func ResponseOf(route *api.Route) serverdeps.Response {
	response, _ := route.Response.(serverdeps.Response)
	return response
}
