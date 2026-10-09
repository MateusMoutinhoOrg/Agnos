package api_close_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/close-backup: the open
// snapshot whose id the body names turns ready — downloadable and restorable,
// taking no file anymore — once every content its files name is stored, and
// is answered 200 as {"snapshot": ...}. A snapshot that does not exist answers
// 404, one that is not open, holds no file or a file at a path others need as
// a folder 400, and one naming a content that
// is not stored 400 with the path of that file; it stays open. Another backup
// job running answers 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, outcome, missing, err := snapshots.Close(sandbox, int64(input.Body.Id))
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeNotOpen:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", backofficesnapshots.NotOpenMessage)
	case snapshots.OutcomeEmpty:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", backofficesnapshots.EmptyMessage)
	case snapshots.OutcomeConflict:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", "the snapshot holds a file at "+missing+" and others inside it as a folder: no restore could write both, so remove one of them first")
	case snapshots.OutcomeBlobMissing:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", "the content of "+missing+" is not stored: add that file again, then close the snapshot")
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.SnapshotResponseJSON(sandbox, snapshot))
}
