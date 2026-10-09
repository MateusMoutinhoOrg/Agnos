package api_add_backup_file

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/add-backup-file/{id}/{path}:
// the body is stored as the file at path, below the --database folder, of the open snapshot
// of id, replacing the one it held there, and answered 201 as
// {"file": {path, sha}}. A path a snapshot cannot hold is answered 400, a
// snapshot that does not exist 404, one that is not open 400, and another
// backup job running 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	file, outcome, err := snapshots.AddFile(sandbox, int64(input.Id), input.Path, input.Body)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeInvalidPath:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "path", backofficesnapshots.InvalidPathMessage)
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeNotOpen:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", backofficesnapshots.NotOpenMessage)
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusCreated, backofficeapi.SnapshotFileResponseJSON(sandbox, file))
}
