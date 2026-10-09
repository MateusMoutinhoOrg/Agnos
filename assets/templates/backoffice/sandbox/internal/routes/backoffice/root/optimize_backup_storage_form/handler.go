package optimize_backup_storage_form

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /admin/root/optimize-backup-storage: every
// stored content no snapshot holds anymore starts being removed in the
// background, and the browser is sent to the list at once. Another backup job
// running refuses it with a notice.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	notice := backofficesnapshots.NoticeOptimizing
	if !snapshots.StartOptimize(sandbox) {
		notice = backofficesnapshots.NoticeBusy
	}
	return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficesnapshots.ListLocation(sandbox, notice))
}
