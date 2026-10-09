package api_add_backup_blob

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers POST /api/admin/root/add-backup-blob: the body is stored
// once, under its SHA-256, and answered 201 as {"blob": {sha, size}} — the
// sha add-backup-reference names it by. A content no snapshot names
// is removed by the next optimize-backup-storage, so it is referenced soon
// after. Another backup job running is answered 409.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	sha, outcome, err := snapshots.AddBlob(sandbox, input.Body)
	if err != nil {
		return err
	}
	if outcome == snapshots.OutcomeBusy {
		return sandbox.Deps.OpinionatedAgnosServer.WriteError(sandbox.Deps.SerializableDeps, *response, api.StatusConflict, "", backofficesnapshots.BusyMessage)
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusCreated, backofficeapi.BlobResponseJSON(sandbox, sha, len(input.Body)))
}
