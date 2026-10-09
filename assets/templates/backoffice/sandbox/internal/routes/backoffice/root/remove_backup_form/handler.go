package remove_backup_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /admin/root/remove-backup/{id}: the
// snapshot is deleted with its list of files, and the browser is sent to the
// list with the outcome. The contents only it held stay stored until the
// backups size is optimized. Another backup job running refuses it with a
// notice: it may be creating, restoring or uploading that very snapshot.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	outcome, err := snapshots.Remove(sandbox, int64(input.Id))
	if err != nil {
		return err
	}
	notice := outcome
	if outcome == snapshots.OutcomeOk {
		notice = backofficesnapshots.NoticeRemoved
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficesnapshots.ListLocation(sandbox, notice))
}
