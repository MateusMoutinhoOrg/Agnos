package api_upload_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/upload-backup, the
// zip archive a download built as the whole body: it is stored as a new
// ready snapshot, answered 201 as {"snapshot": ...}. An archive that is not
// one a download built is answered 400 with why, and one sent while another
// backup job runs is answered 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, outcome, refused, err := snapshots.Import(sandbox, input.Body)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeRefused:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "body", refused)
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusCreated, backofficeapi.SnapshotResponseJSON(sandbox, snapshot))
}
