package routeio

import (
	"{{.Module}}/sandbox/api"
)

// answeredStatusKey is the entry of Locals the dispatch puts the status a
// request was answered with under, before the routes of the `after` phase run.
const answeredStatusKey = "routeio.answered-status"

// SetLocal stores one value under key in the Locals of the request the route
// is bound to, for every route of the chain that runs after it:
//
//	routeio.SetLocal(route, "user", user)
func SetLocal(route *api.Route, key string, value any) {
	if route.Locals == nil {
		route.Locals = map[string]any{}
	}
	route.Locals[key] = value
}

// GetLocal reads back what a route earlier in the chain stored under key. It
// reports false when nothing was stored there, or something of another type:
//
//	user, ok := routeio.GetLocal[auth.User](route, "user")
func GetLocal[T any](route *api.Route, key string) (T, bool) {
	var zero T
	raw, has := route.Locals[key]
	if !has {
		return zero, false
	}
	value, is := raw.(T)
	if !is {
		return zero, false
	}
	return value, true
}

// AnsweredStatus is the status the request was answered with, as a route of
// the `after` phase reads it; 0 before the chain has answered.
func AnsweredStatus(route *api.Route) int {
	status, _ := GetLocal[int](route, answeredStatusKey)
	return status
}

// SetAnsweredStatus records the status the request was answered with. It is
// the dispatch's to call, once, before the routes of the `after` phase run.
func SetAnsweredStatus(locals map[string]any, status int) {
	locals[answeredStatusKey] = status
}
