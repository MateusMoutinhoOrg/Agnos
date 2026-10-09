package api_list_backups

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /api/admin/root/list-backups: every
// snapshot whose name starts with the prefix query value — every one without
// it — newest first, as {"snapshots": [...], "busy": ...}, busy telling
// whether a snapshot or a restore is running.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listed, err := snapshots.List(sandbox, input.Prefix)
	if err != nil {
		return err
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.SnapshotListJSON(sandbox, listed, snapshots.Busy(sandbox)))
}
