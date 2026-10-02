package home

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// InternalPureHandler answers GET /admin/home with backoffice/home.html,
// rendered for the user the authentication middleware put on props.User.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}
	return backofficerender.Home(sandbox, response, props.User, backofficeauth.SessionSeconds/60)
}
