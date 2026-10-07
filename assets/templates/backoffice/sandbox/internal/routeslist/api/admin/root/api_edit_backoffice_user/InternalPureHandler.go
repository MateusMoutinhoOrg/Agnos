package api_edit_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/root/edit-backoffice-user. The
// username, email and role of the user whose id the body names are written,
// and their password when one was given, then the user is answered as it now
// stands. A new password ends every session of the user and revokes every API
// token of theirs — the one this request carries included, when a root edits
// their own account; a refused edit — demoting the last root among them — answers a 400
// carrying the reason, and a user that does not exist a 404.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	id := int64(entries.Body.Id)
	_, ok := backofficeusers.Find(sandbox, id)
	if !ok {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}

	role, err := backofficeapi.Role(sandbox, entries.Body.Role)
	if err != nil {
		return err
	}
	message, _, err := backofficeusers.Update(sandbox, *props.User, nil, id, backofficeusers.Fields{
		Username: entries.Body.Username,
		Email:    entries.Body.Email,
		Password: entries.Body.Password,
		Role:     role,
	})
	if err != nil {
		return err
	}
	if message != "" {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusBadRequest, "", message)
	}

	user, ok := backofficeusers.Find(sandbox, id)
	if !ok {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}
	return sandbox.Deps.OpinatedAgnosServer.WriteJSON(sandbox.Deps.Serializables, *response, api.StatusOk, backofficeapi.UserDocument(sandbox, user))
}
