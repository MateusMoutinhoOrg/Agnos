package api_remove_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// Handle answers POST /api/admin/root/remove-backoffice-user: the
// user whose id the body names is removed, with every session of it, so its
// tokens are refused from here on. A root removing its own account is refused
// with a 403, and a user that does not exist answers a 404.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	notice, err := backofficeusers.Remove(sandbox, *props.User, int64(input.Body.Id))
	if err != nil {
		return err
	}
	switch notice {
	case backofficeusers.NoticeSelf:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusForbidden, "id", "you cannot remove your own account")
	case backofficeusers.NoticeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that user does not exist")
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.OkJSON(sandbox))
}
