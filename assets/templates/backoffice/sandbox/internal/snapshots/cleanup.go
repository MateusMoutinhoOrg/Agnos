package snapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// Remove deletes the snapshot of id with every one of its (path, sha) pairs,
// answering OutcomeOk, OutcomeNotFound, OutcomeNotReady for a StatusRollback
// one — the only copy of what DataDir held before a restore that has not
// finished — or OutcomeBusy while a job runs — one may be creating, restoring
// or uploading that very snapshot. Its blobs stay:
// another snapshot may hold the same contents, and StartOptimize is what
// removes the ones none holds anymore.
func Remove(sandbox *api.Sandbox, id int64) (string, error) {
	if !acquire(sandbox) {
		return OutcomeBusy, nil
	}
	defer release(sandbox)

	db := backup.New(sandbox)
	snapshot, ok := db.FindSnapshotById(id)
	if !ok {
		return OutcomeNotFound, nil
	}
	if snapshot.Status == StatusRollback {
		return OutcomeNotReady, nil
	}
	if err := db.RemoveSnapshot(id); err != nil {
		return "", err
	}
	sandbox.Deps.StdDeps.Logf("snapshot %s removed\n", snapshot.Name)
	return OutcomeOk, nil
}

// StartOptimize starts removing every blob no snapshot holds anymore and
// answers at once, leaving the work to a goroutine that runs after the caller
// returns; started is false while another job runs. How many blobs it removed
// is written on stderr.
func StartOptimize(sandbox *api.Sandbox) (started bool) {
	if !acquire(sandbox) {
		return false
	}
	go func() {
		defer release(sandbox)
		removed, kept, err := optimize(sandbox, backup.New(sandbox))
		if err != nil {
			sandbox.Deps.StdDeps.Eprintf("optimize of the backups failed after removing %d blobs: %s\n", removed, err.Error())
			return
		}
		sandbox.Deps.StdDeps.Logf("optimize of the backups: %d blobs no snapshot held removed, %d kept\n", removed, kept)
	}()
	return true
}

// optimize first repairs the backup database — deleting what a job stopped
// in the middle of a write left behind — then removes every blob no (path,
// sha) pair of any snapshot names — a failed snapshot's included, so what it
// holds stays whole until it is removed — and answers how many it removed
// and how many it kept. The caller holds the job token, so no snapshot gains
// a pair while it runs.
func optimize(sandbox *api.Sandbox, db *backup.Backup) (removed int, kept int, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = sandbox.Deps.StdDeps.Errorf("%v", recovered)
		}
	}()

	if err := backup.RepairBackup(sandbox, db); err != nil {
		return 0, 0, err
	}

	held := map[string]bool{}
	listed, err := db.ListSnapshots(backup.SnapshotFilter{})
	if err != nil {
		return 0, 0, err
	}
	for _, snapshot := range listed {
		contents, err := db.ListSnapshotContents(snapshot.Id)
		if err != nil {
			return 0, 0, err
		}
		for _, content := range contents {
			held[content.Sha] = true
		}
	}

	blobs, err := backup.ListBlobShas(sandbox, db)
	if err != nil {
		return 0, 0, err
	}
	for _, blob := range blobs {
		if held[blob.Sha] {
			kept++
			continue
		}
		if err := db.RemoveBlob(blob.Id); err != nil {
			return removed, kept, err
		}
		removed++
	}
	return removed, kept, nil
}
