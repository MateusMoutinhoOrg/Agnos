package api_set_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// Handle answers POST /api/admin/root/set-backoffice-user. The
// username, email and role of the user whose id the body names are written,
// and their password when one was given, then the user is answered as it now
// stands. A new password ends every session of the user and revokes every API
// token of theirs — the one this request carries included, when a root edits
// their own account; a refused edit — demoting the last root among them — answers a 400
// carrying the reason, and a user that does not exist a 404.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	id := int64(input.Body.Id)
	_, ok := backofficeusers.Find(sandbox, id)
	if !ok {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}

	parsed, ok := backofficeauth.ParseRole(sandbox, input.Body.Role)
	if !ok {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "role", "role must be root or viewer")
	}
	role := int64(parsed)
	message, _, err := backofficeusers.Set(sandbox, *props.User, nil, id, backofficeusers.Fields{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
		Role:     role,
	})
	if err != nil {
		return err
	}
	if message != "" {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "", message)
	}

	user, ok := backofficeusers.Find(sandbox, id)
	if !ok {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.UserResponseJSON(sandbox, user))
}
