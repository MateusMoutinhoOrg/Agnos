package restore_backup_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /admin/root/restore-backup/{id}: the
// snapshot starts being put back over every database — after a pre-restore
// snapshot of them — and the browser is sent to the list at once, with the
// outcome. Only a ready snapshot is restored, and never while another job
// runs.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	outcome, err := snapshots.StartRestore(sandbox, int64(input.Id))
	if err != nil {
		return err
	}
	notice := outcome
	if outcome == snapshots.OutcomeOk {
		notice = backofficesnapshots.NoticeRestoring
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficesnapshots.ListLocation(sandbox, notice))
}
