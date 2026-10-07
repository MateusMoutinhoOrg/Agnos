package api_add_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/root/add-backoffice-user. A user
// the body describes well is added and answered under a 201; anything else is
// refused with a 400 carrying the reason.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	role, err := backofficeapi.Role(sandbox, entries.Body.Role)
	if err != nil {
		return err
	}
	user, message, err := backofficeusers.Add(sandbox, backofficeusers.Fields{
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
	return sandbox.Deps.OpinatedAgnosServer.WriteJSON(sandbox.Deps.Serializables, *response, api.StatusCreated, backofficeapi.UserDocument(sandbox, user))
}
