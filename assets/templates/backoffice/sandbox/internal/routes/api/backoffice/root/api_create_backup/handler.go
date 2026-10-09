package api_create_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/create-backup: a
// snapshot of every database, named after the body's name or, without one,
// after the current instant, is recorded as creating and answered 202 at
// once, while its files are stored in the background; list-backups
// tells when it is ready. A name it cannot be given is answered 400, one
// another snapshot holds 409, and another backup job running 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, outcome, err := snapshots.StartCreate(sandbox, input.Body.Name)
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
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, backofficesnapshots.StatusAccepted, backofficeapi.SnapshotCreatingJSON(sandbox, snapshot))
}
