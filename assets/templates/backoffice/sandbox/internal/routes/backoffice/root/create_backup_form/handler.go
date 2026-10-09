package create_backup_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /admin/root/create-backup: a snapshot
// of every database, named after the form's name or, left empty, after the
// current instant, is recorded as creating and the browser is sent to the
// list at once, while its files are stored in the background. A name it
// cannot be given or another snapshot holds, or another backup job running,
// refuses it with a notice.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	_, outcome, err := snapshots.StartCreate(sandbox, input.Body.Name)
	if err != nil {
		return err
	}
	notice := outcome
	if outcome == snapshots.OutcomeOk {
		notice = backofficesnapshots.NoticeCreating
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficesnapshots.ListLocation(sandbox, notice))
}
