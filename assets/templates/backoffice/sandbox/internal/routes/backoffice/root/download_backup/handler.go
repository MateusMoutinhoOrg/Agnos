package download_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /admin/root/download-backup/{id} with
// the zip archive of the snapshot, saved under its name. A snapshot that does
// not exist or is not ready sends the browser back to the list with a notice.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, archive, outcome, err := snapshots.Export(sandbox, int64(input.Id))
	if err != nil {
		return err
	}
	if outcome != snapshots.OutcomeOk {
		return sandbox.Deps.OpinionatedAgnosServer.Redirect(*response, api.StatusSeeOther, backofficesnapshots.ListLocation(sandbox, outcome))
	}
	return backofficesnapshots.WriteArchive(sandbox, response, snapshot, archive)
}
