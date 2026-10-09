package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/databases/backup"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficesnapshots"
	"{{.Module}}/sandbox/internal/snapshots"
)

// BackupSnapshotsPage is what backoffice/backup_snapshots.html is rendered
// with.
type BackupSnapshotsPage struct {
	Viewer Viewer
	// Notice is shown above the list, its Text "" for none.
	Notice Notice
	// Prefix is the name prefix the list is narrowed to, "" for none.
	Prefix    string
	Snapshots []BackupSnapshotRow
	// Busy disables what would start a second job, while one runs.
	Busy bool
	// Refresh reloads the page every few seconds, while a job runs or a
	// snapshot is still being created.
	Refresh bool
	// MaxUploadMegabytes is the largest archive the upload form takes.
	MaxUploadMegabytes int
}

// BackupSnapshotRow is one snapshot of the list.
type BackupSnapshotRow struct {
	Id   string
	Name string
	// Date is the instant the snapshot was taken at.
	Date       string
	Status     string
	IsReady    bool
	IsCreating bool
	IsFailed   bool
}

// backupSnapshotNoticeOf is how the list page words a backofficesnapshots
// notice code; an unknown code shows nothing.
func backupSnapshotNoticeOf(sandbox *api.Sandbox, code string) Notice {
	switch code {
	case backofficesnapshots.NoticeCreating:
		return Notice{Text: "Snapshot started. It is listed as creating until every database is copied, and this page refreshes until then.", Kind: "ok"}
	case backofficesnapshots.NoticeRestoring:
		return Notice{Text: "Restore started. A pre-restore snapshot of the current data is taken first, then every database but the backoffice users and tokens — unless you ticked the box — is replaced, and the site answers 503 until it ends. If it fails, the pre-restore snapshot is put back.", Kind: "ok"}
	case backofficesnapshots.NoticeUploaded:
		return Notice{Text: "Snapshot uploaded. It can be downloaded or restored now.", Kind: "ok"}
	case backofficesnapshots.NoticeRemoved:
		return Notice{Text: "Snapshot deleted. The contents only it held stay stored until you optimize the backups size.", Kind: "ok"}
	case backofficesnapshots.NoticeOptimizing:
		return Notice{Text: "Optimizing started. The stored contents no snapshot holds anymore are being removed in the background.", Kind: "ok"}
	case backofficesnapshots.NoticeBusy:
		return Notice{Text: "Another backup job is running: a snapshot, a restore, an upload, a deletion, an optimization or a write to a snapshot built by hand. Try again once it ends.", Kind: "error"}
	case backofficesnapshots.NoticeNotFound:
		return Notice{Text: "That snapshot no longer exists.", Kind: "error"}
	case backofficesnapshots.NoticeNotReady:
		return Notice{Text: "That snapshot is not ready: only a finished snapshot can be downloaded or restored, and the safety copy of a restore that has not finished cannot be deleted.", Kind: "error"}
	case backofficesnapshots.NoticeInvalidName:
		return Notice{Text: "That name cannot be used: a snapshot name is 1 to 100 letters, digits, '.', '_' or '-', starting with a letter or a digit.", Kind: "error"}
	case backofficesnapshots.NoticeNameTaken:
		return Notice{Text: "Another snapshot already has that name. Pick another one, or delete that snapshot first.", Kind: "error"}
	}
	return Notice{}
}

// backupSnapshotRow is the row of snapshot on the list page.
func backupSnapshotRow(sandbox *api.Sandbox, snapshot backup.SnapshotRecord) BackupSnapshotRow {
	return BackupSnapshotRow{
		Id:         sandbox.Deps.StringsDeps.FormatInt(snapshot.Id, 10),
		Name:       snapshot.Name,
		Date:       instantOf(sandbox, snapshot.Data),
		Status:     snapshot.Status,
		IsReady:    snapshot.Status == snapshots.StatusReady,
		IsCreating: snapshot.Status == snapshots.StatusCreating,
		IsFailed:   snapshot.Status == snapshots.StatusFailed,
	}
}

// RenderBackupSnapshotsPage answers, under status, the snapshot list page
// for user: listed, newest first, narrowed to the names starting with prefix,
// with the notice code notice above it, and what starts a job disabled while
// busy.
func RenderBackupSnapshotsPage(sandbox *api.Sandbox, response *serverdeps.Response, status int, user *backoffice_db.BackofficeUserRecord, listed []backup.SnapshotRecord, busy bool, prefix string, notice string) error {
	rows := []BackupSnapshotRow{}
	refresh := busy
	for _, snapshot := range listed {
		row := backupSnapshotRow(sandbox, snapshot)
		refresh = refresh || row.IsCreating
		rows = append(rows, row)
	}

	return RenderHTML(sandbox, response, status, "backoffice/backup_snapshots.html", BackupSnapshotsPage{
		Viewer:             viewerOf(sandbox, user),
		Notice:             backupSnapshotNoticeOf(sandbox, notice),
		Prefix:             sandbox.Deps.StringsDeps.TrimSpace(prefix),
		Snapshots:          rows,
		Busy:               busy,
		Refresh:            refresh,
		MaxUploadMegabytes: snapshots.MaxArchiveBytes >> 20,
	})
}
