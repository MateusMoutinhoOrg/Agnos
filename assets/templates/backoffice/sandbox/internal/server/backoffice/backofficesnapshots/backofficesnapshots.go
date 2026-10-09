package backofficesnapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// The backups of the backoffice: the pages and the api through which a root
// takes, downloads, uploads, restores and removes the snapshots of
// sandbox/internal/snapshots, and removes the blobs none of them holds.
// Every one of them is under /admin/root or /api/admin/root, so only a root
// reaches it: a restore replaces every database, and an archive carries
// every password hash and session.

// ListPath is the page every action on a snapshot sends the browser back to.
const ListPath = "/admin/root/list-backups"

// The notices an action hands the list page through ListPath?notice=<code>.
// Each one is a fixed word the page words itself, so nothing the client sent
// is ever shown back. NoticeBusy, NoticeNotFound, NoticeNotReady,
// NoticeInvalidName and NoticeNameTaken spell the snapshots outcome of the
// same meaning.
const (
	// NoticeCreating follows a snapshot started.
	NoticeCreating = "creating"
	// NoticeRestoring follows a restore started.
	NoticeRestoring = "restoring"
	// NoticeUploaded follows an archive uploaded.
	NoticeUploaded = "uploaded"
	// NoticeRemoved follows a snapshot removed.
	NoticeRemoved = "removed"
	// NoticeOptimizing follows an optimize of the backups started.
	NoticeOptimizing = "optimizing"
	// NoticeBusy follows a job refused while another one runs.
	NoticeBusy = "busy"
	// NoticeNotFound follows an action on a snapshot that does not exist.
	NoticeNotFound = "not-found"
	// NoticeNotReady follows a download or a restore of a snapshot that is
	// not ready.
	NoticeNotReady = "not-ready"
	// NoticeInvalidName follows a snapshot refused for a name it cannot be
	// given.
	NoticeInvalidName = "invalid-name"
	// NoticeNameTaken follows a snapshot refused for a name another one
	// holds.
	NoticeNameTaken = "name-taken"
)

// BusyMessage is what the api answers, under a 409, to a job refused while
// another one runs.
const BusyMessage = "another backup job is running (a snapshot, a restore, an upload, a removal, an optimize or a write to a snapshot built by hand): try again once it ends"

// InvalidNameMessage is what the api answers, under a 400, to a snapshot
// name the snapshots package refuses.
const InvalidNameMessage = "a snapshot name is 1 to 100 letters, digits, '.', '_' or '-', starting with a letter or a digit"

// NameTakenMessage is what the api answers, under a 409, to a snapshot name
// another one holds.
const NameTakenMessage = "another snapshot already has that name"

// NotOpenMessage is what the api answers, under a 400, to a file added to or
// a close of a snapshot that is not open.
const NotOpenMessage = "only an open snapshot, one create-empty-backup recorded and close-backup has not closed yet, takes files or closes"

// InvalidPathMessage is what the api answers, under a 400, to a file path a
// snapshot cannot hold.
const InvalidPathMessage = "a file path is relative to the --database folder (data by default), inside a database folder other than backup, as in backofficedb/backoffice-user/1/values/username: no empty, '.' or '..' segment, no '\\' nor ':'"

// InvalidShaMessage is what the api answers, under a 400, to a sha that is
// not one a content is stored under.
const InvalidShaMessage = "a sha is the 64 hex digits of the SHA-256 of the content, as add-backup-blob answers it"

// BlobMissingMessage is what the api answers, under a 404, to a sha no
// stored content has.
const BlobMissingMessage = "no content is stored under that sha: send it to add-backup-blob first (an optimize removes the ones no snapshot names yet)"

// EmptyMessage is what the api answers, under a 400, to a close of a snapshot
// holding no file.
const EmptyMessage = "the snapshot holds no file: restoring it would empty every database, so it is not closed"

// StatusAccepted is the status of a job the api started and does not wait
// for; api has no constant of its own for it.
const StatusAccepted = 202

// ListLocation is the list page showing the notice code notice.
func ListLocation(sandbox *api.Sandbox, notice string) string {
	return ListPath + "?notice=" + notice
}

// WriteArchive answers archive, the zip of snapshot, as a file to save under
// the snapshot's name. A snapshot name holds no character a header or a file
// name has to escape: snapshots names every one of them.
func WriteArchive(sandbox *api.Sandbox, response *serverdeps.Response, snapshot backup.SnapshotRecord, archive []byte) error {
	response.SetHeader("Content-Type", "application/zip")
	response.SetHeader("Content-Disposition", "attachment; filename=\""+snapshot.Name+".zip\"")
	response.SetStatus(api.StatusOK)
	return response.Write(archive)
}

// WriteFile answers content, the file at path of a snapshot, as a file to
// save under the last segment of path. A '"', a '\' or a control character
// in it is answered as '_', so the header holds one quoted file name.
func WriteFile(sandbox *api.Sandbox, response *serverdeps.Response, path string, content []byte) error {
	strings := sandbox.Deps.StringsDeps
	name := []rune(path[strings.LastIndex(path, "/")+1:])
	for index, char := range name {
		if char == '"' || char == '\\' || char < ' ' || char == 0x7f {
			name[index] = '_'
		}
	}
	response.SetHeader("Content-Type", "application/octet-stream")
	response.SetHeader("Content-Disposition", "attachment; filename=\""+string(name)+"\"")
	response.SetStatus(api.StatusOK)
	return response.Write(content)
}
