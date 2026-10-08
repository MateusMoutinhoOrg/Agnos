package add_backoffice_user_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers POST /admin/root/add-backoffice-user. A user the
// form describes well is added and the browser is sent to the list; anything
// else answers the form again, filled with what was sent but the password,
// under a 400 with the reason above it.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	fields := backofficeusers.Fields{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
		Role:     int64(input.Body.Role),
	}
	_, message, err := backofficeusers.Add(sandbox, fields)
	if err != nil {
		return err
	}
	if message != "" {
		return backofficerender.RenderAddUserPage(sandbox, response, api.StatusBadRequest, props.User, fields, message)
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, backofficeusers.NoticeAdded))
}
