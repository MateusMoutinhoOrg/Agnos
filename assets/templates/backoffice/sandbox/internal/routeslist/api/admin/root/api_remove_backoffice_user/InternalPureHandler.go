package api_remove_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers POST /api/admin/root/remove-backoffice-user: the
// user whose id the body names is removed, with every session of it, so its
// tokens are refused from here on. A root removing its own account is refused
// with a 403, and a user that does not exist answers a 404.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	if props.User == nil {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "", "no authenticated user")
	}

	notice, err := backofficeusers.Remove(sandbox, *props.User, int64(entries.Body.Id))
	if err != nil {
		return err
	}
	switch notice {
	case backofficeusers.NoticeSelf:
		return routeio.Fail(sandbox, api.StatusForbidden, "id", "you cannot remove your own account")
	case backofficeusers.NoticeNotFound:
		return routeio.Fail(sandbox, api.StatusNotFound, "id", "that user does not exist")
	}
	return routeio.WriteJSON(sandbox, *response, api.StatusOk, backofficeapi.Ok(sandbox))
}
