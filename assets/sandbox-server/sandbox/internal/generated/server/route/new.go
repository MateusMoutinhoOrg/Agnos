package route

import (
	"{{.Module}}/sandbox/api"
)

// NewRoute returns the base every route of sandbox/internal/routeslist/ is
// built on: an empty declaration whose IsActionable, MatchesPath and
// RequestHandler are the generic ones of this package, closed over the
// sandbox. A route's generated NewRoute fills its declaration and its
// InternalPurehandler on top of what this returns, so every route matches and
// binds a request the same way.
func NewRoute(sandbox *api.Sandbox) *api.Route {
	self := api.NewRoute()

	self.IsActionable = func(bound *api.Route) bool {
		return IsActionable(sandbox, bound)
	}
	self.MatchesPath = func(bound *api.Route) bool {
		return MatchesPath(sandbox, bound)
	}
	self.RequestHandler = func(bound *api.Route) error {
		return RequestHandler(sandbox, bound)
	}

	return self
}
