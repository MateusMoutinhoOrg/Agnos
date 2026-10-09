package api_create_empty_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/create-empty-backup: an
// empty snapshot, named after the body's name or, without one, after the
// current instant, is recorded as open and answered 201 as {"snapshot": ...}.
// Files are then added to it one by one — add-backup-file,
// add-backup-reference — until close-backup turns it ready.
// A name it cannot be given is answered 400, one another snapshot holds 409,
// and another backup job running 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, outcome, err := snapshots.CreateOpen(sandbox, input.Body.Name)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeInvalidName:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "name", backofficesnapshots.InvalidNameMessage)
	case snapshots.OutcomeNameTaken:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "name", backofficesnapshots.NameTakenMessage)
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusCreated, backofficeapi.SnapshotResponseJSON(sandbox, snapshot))
}
