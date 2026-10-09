package api_download_backup_file

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle answers GET /api/admin/root/download-backup-file/{id}/{path}
// with the bytes of the file at path, below the --database folder, of the snapshot of id, of
// any status. A snapshot that does not exist, a path it holds no file at, and
// a file whose content is not stored answer 404.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if props.User == nil {
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "", "no authenticated user")
	}

	file, content, outcome, err := snapshots.ReadFile(sandbox, int64(input.Id), input.Path)
	if err != nil {
		return err
	}
	switch outcome {
	case snapshots.OutcomeNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "id", "that snapshot does not exist")
	case snapshots.OutcomeFileNotFound:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "path", "that snapshot holds no file at that path")
	case snapshots.OutcomeBlobMissing:
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusNotFound, "path", "the content of that file is not stored")
	}
	return backofficesnapshots.WriteFile(sandbox, response, file.Path, content)
}
