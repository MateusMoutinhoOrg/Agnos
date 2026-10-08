package set_backoffice_user_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
)

// Handle answers POST /admin/root/set-backoffice-user/{id}. The
// user's username, email and role are written, and their password when one was
// given, then the browser is sent to the list. A new password ends every
// session of the user and revokes their API tokens — the session of this
// request survives when a root edits their own account — and the list says
// so. A refused edit answers the form
// again, filled with what was sent but the password, under a 400 with the
// reason above it. A user that does not exist sends the browser back to the
// list.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	id := int64(input.Id)
	_, ok := backofficeusers.Find(sandbox, id)
	if !ok {
		return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, backofficeusers.NoticeNotFound))
	}

	fields := backofficeusers.Fields{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
		Role:     int64(input.Body.Role),
	}
	message, notice, err := backofficeusers.Set(sandbox, *props.User, props.Session, id, fields)
	if err != nil {
		return err
	}
	if message != "" {
		return backofficerender.RenderSetUserPage(sandbox, response, api.StatusBadRequest, props.User, id, fields, message)
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficeusers.ListLocation(sandbox, notice))
}
