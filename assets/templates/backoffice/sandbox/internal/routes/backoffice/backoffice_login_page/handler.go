package backoffice_login_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers GET /admin/login: the sign-in form, under a 200,
// whether or not the browser holds a session. The form posts to the
// backoffice-login route at the same path.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	return backofficerender.RenderLoginPage(sandbox, response, api.StatusOK, "", "")
}
