package snapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// StartRestore starts putting the snapshot of id back over DataDir and
// answers at once, leaving the work to a goroutine that runs after the caller
// returns: OutcomeOk when it started, OutcomeNotFound, OutcomeNotReady or
// OutcomeBusy when it did not. How the restore ended is written on stderr.
func StartRestore(sandbox *api.Sandbox, id int64) (string, error) {
	db := backup.New(sandbox)
	target, ok := db.FindSnapshotById(id)
	if !ok {
		return OutcomeNotFound, nil
	}
	if target.Status != StatusReady {
		return OutcomeNotReady, nil
	}
	if !acquire(sandbox) {
		return OutcomeBusy, nil
	}
	go func() {
		defer release(sandbox)
		if err := restore(sandbox, db, target); err != nil {
			sandbox.Deps.StdDeps.Eprintf("restore of snapshot %s failed: %s\n", target.Name, err.Error())
			return
		}
		sandbox.Deps.StdDeps.Logf("snapshot %s restored\n", target.Name)
	}()
	return OutcomeOk, nil
}

// restore puts target back over DataDir: every database directory but
// BackupDir is removed, and every file of target written in its place. It
// first takes a snapshot of what DataDir holds now, named pre-restore, and
// checks that every file of target can be read back, so a restore that
// cannot finish fails before anything is removed. The caller holds the job
// token.
func restore(sandbox *api.Sandbox, db *backup.Backup, target backup.SnapshotRecord) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = sandbox.Deps.StdDeps.Errorf("%v", recovered)
		}
	}()

	safety, err := create(sandbox, db, "pre-restore")
	if err != nil {
		return sandbox.Deps.StdDeps.Errorf("the safety snapshot failed, nothing was restored: %s", err.Error())
	}
	sandbox.Deps.StdDeps.Logf("restore of snapshot %s: what %s held is kept as snapshot %s\n", target.Name, DataDir(sandbox), safety.Name)

	contents, err := db.ListSnapshotContents(target.Id)
	if err != nil {
		return err
	}
	for _, content := range contents {
		if !SafePath(sandbox, content.Path) {
			return sandbox.Deps.StdDeps.Errorf("the path %q is not one a snapshot may hold, nothing was restored", content.Path)
		}
		if _, ok := db.FindBlobBySha(content.Sha); !ok {
			return sandbox.Deps.StdDeps.Errorf("the content of %s is missing, nothing was restored", content.Path)
		}
	}

	io := sandbox.Deps.IoDeps
	root := DataDir(sandbox)
	backupDir := io.Join(root, BackupDir)
	for _, dir := range io.ListDirs(root) {
		if dir != backupDir {
			io.RemoveDir(dir)
		}
	}
	for _, content := range contents {
		blob, ok := db.FindBlobBySha(content.Sha)
		if !ok {
			return sandbox.Deps.StdDeps.Errorf("the content of %s went missing while restoring", content.Path)
		}
		segments := append([]string{root}, sandbox.Deps.StringsDeps.Split(content.Path, "/")...)
		if err := io.WriteFile(io.Join(segments...), blob.Value); err != nil {
			return err
		}
	}
	return nil
}

// SafePath tells whether path, slash-separated and relative to DataDir, is
// one a snapshot may hold: a file inside a database directory other than
// BackupDir, with no segment that is empty, "." or "..", nor one holding a
// backslash or a colon — so that written under DataDir it lands nowhere else.
func SafePath(sandbox *api.Sandbox, path string) bool {
	strings := sandbox.Deps.StringsDeps
	segments := strings.Split(path, "/")
	if len(segments) < 2 || strings.ToLower(segments[0]) == BackupDir {
		return false
	}
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." || strings.ContainsAny(segment, "\\:\x00") {
			return false
		}
	}
	return true
}
