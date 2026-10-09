package api_download_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /api/admin/root/download-backup/{id}
// with the zip archive of the snapshot. A snapshot that does not exist
// answers 404, and one that is not ready 400.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, archive, outcome, err := snapshots.Export(sandbox, int64(input.Id))
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeNotReady:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", "only a ready snapshot can be downloaded")
	}
	return backofficesnapshots.WriteArchive(sandbox, response, snapshot, archive)
}
