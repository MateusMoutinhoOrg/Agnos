package api_add_backup_reference

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/add-backup-reference: the
// content add-backup-blob stored under the body's sha is named as the file at
// its path, below the --database folder, of the open snapshot of its id, replacing the one it
// held there, and answered 201 as {"file": {path, sha}}. A path a snapshot
// cannot hold or a sha that is not one is answered 400, a snapshot that does
// not exist 404, one that is not open 400, a sha no content is stored under
// 404 — send it to add-backup-blob, then reference it again — and another
// backup job running 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	file, outcome, err := snapshots.AddReference(sandbox, int64(input.Body.Id), input.Body.Path, input.Body.Sha)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeInvalidPath:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "path", backofficesnapshots.InvalidPathMessage)
	case snapshots.OutcomeInvalidSha:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "sha", backofficesnapshots.InvalidShaMessage)
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeNotOpen:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", backofficesnapshots.NotOpenMessage)
	case snapshots.OutcomeBlobMissing:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "sha", backofficesnapshots.BlobMissingMessage)
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusCreated, backofficeapi.SnapshotFileResponseJSON(sandbox, file))
}
