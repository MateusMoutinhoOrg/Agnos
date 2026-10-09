package api_restore_backup

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/restore-backup: the
// snapshot whose id the body names starts being put back over every database
// — after a pre-restore snapshot of them; the backoffice's own users and
// tokens only when include-backoffice is true — and 202 is answered at once. A
// snapshot that does not exist answers 404, one that is not ready 400, and
// another backup job running 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	outcome, err := snapshots.StartRestore(sandbox, int64(input.Body.Id), input.Body.IncludeBackoffice)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeNotReady:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusBadRequest, "id", "only a ready snapshot can be restored")
	case snapshots.OutcomeBusy:
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, backofficesnapshots.StatusAccepted, backofficeapi.JobJSON(sandbox, "restoring"))
}
