package list_backups_page

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficerender"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /admin/root/list-backups, open to a
// root only, with backoffice/backup_snapshots.html: every snapshot of the
// databases whose name starts with the prefix query value — every one
// without it — newest first, with the controls that take, upload, download
// and restore one, and the outcome of the last action above them.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	listed, err := snapshots.List(sandbox, input.Prefix)
	if err != nil {
		return err
	}
	return backofficerender.RenderBackupSnapshotsPage(sandbox, response, api.StatusOK, props.User, listed, snapshots.Busy(sandbox), input.Prefix, input.Notice)
}
