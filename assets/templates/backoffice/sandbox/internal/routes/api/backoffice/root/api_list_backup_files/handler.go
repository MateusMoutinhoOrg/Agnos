package api_list_backup_files

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapi"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /api/admin/root/list-backup-files/{id}
// with {"snapshot": ..., "files": [{path, sha}]}: every file of the snapshot
// of id, or the ones whose path below the --database folder starts with the prefix query, in
// path order. A snapshot of any status is listed — an open one as it stands,
// a path added twice once, as last added. A snapshot that does not exist
// answers 404.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	snapshot, files, outcome, err := snapshots.ListFiles(sandbox, int64(input.Id), input.Prefix)
	if err != nil {
		return err
	}
	if outcome == snapshots.OutcomeNotFound {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	}
	return sandbox.Deps.OpinionatedAgnosServer.WriteJSON(sandbox.Deps.SerializableDeps, *response, api.StatusOK, backofficeapi.SnapshotFilesJSON(sandbox, snapshot, files))
}
