package snapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// A restore never leaves DataDir half written. Before it removes anything it
// takes a snapshot of what DataDir holds, named pre-restore, and records it
// StatusRollback: that status is the journal of the restore. A restore that
// fails puts the pre-restore snapshot back at once; one the process stopped in
// the middle of — an interrupt, a kill, a crash — is put back by the next
// start-server, through StartRecover, before the server answers anything else.
// Only once every file of the target is written does the pre-restore snapshot
// turn StatusReady. While a restore or a roll-back runs, every request is
// answered 503 (Restoring), so nothing reads a half-written database nor writes
// into one about to be replaced.

// KeptPreRestore is how many pre-restore snapshots a restore that finished
// keeps, the newest ones: each holds every database whole, so without a
// bound every restore would add one more for good. The older ones are
// removed; their contents are freed by the next optimize.
const KeptPreRestore = 5

// preRestorePrefix starts the name of every snapshot a restore takes first.
const preRestorePrefix = "pre-restore"

// BackofficeDir is the directory of DataDir the backoffice's own users,
// sessions and API tokens live in, its key-prefix. A restore leaves it as it
// is unless asked to include it: putting an older copy back would bring back
// revoked tokens, closed sessions, removed users and changed passwords.
const BackofficeDir = "backofficedb"

// restoring holds a token while a restore or a roll-back writes DataDir.
var restoring = make(chan struct{}, 1)

// Restoring tells whether a restore or a roll-back is writing DataDir: every
// request is refused until it ends.
func Restoring(sandbox *api.Sandbox) bool {
	return len(restoring) > 0
}

// StartRestore starts putting the snapshot of id back over DataDir and
// answers at once, leaving the work to a goroutine that runs after the caller
// returns: OutcomeOk when it started, OutcomeNotFound, OutcomeNotReady or
// OutcomeBusy when it did not. BackofficeDir is restored only when
// includeBackoffice is true. How the restore ended is written on stderr.
func StartRestore(sandbox *api.Sandbox, id int64, includeBackoffice bool) (string, error) {
	db := backup.New(sandbox)
	if !acquire(sandbox) {
		return OutcomeBusy, nil
	}
	// Read under the job token: a remove that ended before it was taken is
	// seen, and none can start until the restore ends.
	target, ok := db.FindSnapshotById(id)
	if !ok {
		release(sandbox)
		return OutcomeNotFound, nil
	}
	if target.Status != StatusReady {
		release(sandbox)
		return OutcomeNotReady, nil
	}
	restoring <- struct{}{}
	go func() {
		defer release(sandbox)
		defer func() { <-restoring }()
		if err := restore(sandbox, db, target, includeBackoffice); err != nil {
			sandbox.Deps.StdDeps.Eprintf("restore of snapshot %s failed: %s\n", target.Name, err.Error())
			return
		}
		sandbox.Deps.StdDeps.Logf("snapshot %s restored\n", target.Name)
	}()
	return OutcomeOk, nil
}

// restore puts target back over DataDir: every database directory but
// BackupDir — and BackofficeDir, unless includeBackoffice — is removed, and
// every file of target under the others written in its place. Every file is
// checked first — its path, its content, no path a folder of another — so a
// restore that cannot finish fails before anything is removed. One that fails
// while writing puts the pre-restore snapshot back. The caller holds the job
// token and the restoring one.
func restore(sandbox *api.Sandbox, db *backup.Backup, target backup.SnapshotRecord, includeBackoffice bool) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = sandbox.Deps.StdDeps.Errorf("%v", recovered)
		}
	}()

	all, err := db.ListSnapshotContents(target.Id)
	if err != nil {
		return err
	}
	contents := scoped(sandbox, all, includeBackoffice)
	if len(contents) == 0 {
		return sandbox.Deps.StdDeps.Errorf("the snapshot holds no file to restore, nothing was restored")
	}
	if err := checkContents(sandbox, db, contents); err != nil {
		return sandbox.Deps.StdDeps.Errorf("%s, nothing was restored", err.Error())
	}

	safety, err := create(sandbox, db, preRestorePrefix)
	if err != nil {
		return sandbox.Deps.StdDeps.Errorf("the safety snapshot failed, nothing was restored: %s", err.Error())
	}
	if err := db.SetSnapshotStatus(safety.Id, StatusRollback); err != nil {
		return sandbox.Deps.StdDeps.Errorf("the safety snapshot could not be recorded, nothing was restored: %s", err.Error())
	}
	sandbox.Deps.StdDeps.Logf("restore of snapshot %s: what %s held is kept as snapshot %s\n", target.Name, DataDir(sandbox), safety.Name)

	if err := write(sandbox, db, contents, includeBackoffice); err != nil {
		if backErr := rollBack(sandbox, db, safety); backErr != nil {
			return sandbox.Deps.StdDeps.Errorf("%s; putting %s back failed too: %s — the next start-server tries again", err.Error(), safety.Name, backErr.Error())
		}
		return sandbox.Deps.StdDeps.Errorf("%s; %s was put back as it was, from snapshot %s", err.Error(), DataDir(sandbox), safety.Name)
	}
	if err := db.SetSnapshotStatus(safety.Id, StatusReady); err != nil {
		return err
	}
	return prunePreRestore(sandbox, db)
}

// prunePreRestore removes every ready pre-restore snapshot but the
// KeptPreRestore newest.
func prunePreRestore(sandbox *api.Sandbox, db *backup.Backup) error {
	listed, err := List(sandbox, preRestorePrefix+"-")
	if err != nil {
		return err
	}
	kept := 0
	for _, snapshot := range listed {
		if snapshot.Status != StatusReady {
			continue
		}
		kept++
		if kept <= KeptPreRestore {
			continue
		}
		if err := db.RemoveSnapshot(snapshot.Id); err != nil {
			return err
		}
		sandbox.Deps.StdDeps.Logf("snapshot %s removed: only the %d newest pre-restore snapshots are kept\n", snapshot.Name, KeptPreRestore)
	}
	return nil
}

// scoped is every one of contents a restore writes: every one but those under
// BackofficeDir, unless includeBackoffice.
func scoped(sandbox *api.Sandbox, contents []backup.SnapshotContentRecord, includeBackoffice bool) []backup.SnapshotContentRecord {
	if includeBackoffice {
		return contents
	}
	kept := []backup.SnapshotContentRecord{}
	for _, content := range contents {
		if !sandbox.Deps.StringsDeps.HasPrefix(content.Path, BackofficeDir+"/") {
			kept = append(kept, content)
		}
	}
	return kept
}

// checkContents answers why contents cannot be written over DataDir, nil when
// they can: every path SafePath accepts, every content stored, and no path a
// folder another one needs.
func checkContents(sandbox *api.Sandbox, db *backup.Backup, contents []backup.SnapshotContentRecord) error {
	paths := []string{}
	for _, content := range contents {
		if !SafePath(sandbox, content.Path) {
			return sandbox.Deps.StdDeps.Errorf("the path %q is not one a snapshot may hold", content.Path)
		}
		stored, err := backup.HasBlob(sandbox, db, content.Sha)
		if err != nil {
			return err
		}
		if !stored {
			return sandbox.Deps.StdDeps.Errorf("the content of %s is missing", content.Path)
		}
		paths = append(paths, content.Path)
	}
	if clash := Conflict(sandbox, paths); clash != "" {
		return sandbox.Deps.StdDeps.Errorf("the path %q is both a file and a folder of another one", clash)
	}
	return nil
}

// Conflict is a path of paths that another one needs as a folder — data/a
// beside data/a/b, which no filesystem holds both of — or "" when there is
// none.
func Conflict(sandbox *api.Sandbox, paths []string) string {
	strings := sandbox.Deps.StringsDeps
	files := map[string]bool{}
	for _, path := range paths {
		files[path] = true
	}
	for _, path := range paths {
		for cut := strings.LastIndex(path, "/"); cut > 0; cut = strings.LastIndex(path[:cut], "/") {
			if files[path[:cut]] {
				return path[:cut]
			}
		}
	}
	return ""
}

// write removes every database directory of DataDir but BackupDir — and
// BackofficeDir, unless includeBackoffice — and writes every one of contents
// in their place.
func write(sandbox *api.Sandbox, db *backup.Backup, contents []backup.SnapshotContentRecord, includeBackoffice bool) error {
	io := sandbox.Deps.IoDeps
	root := DataDir(sandbox)
	kept := map[string]bool{io.Join(root, BackupDir): true}
	if !includeBackoffice {
		kept[io.Join(root, BackofficeDir)] = true
	}
	for _, dir := range io.ListDirs(root) {
		if !kept[dir] {
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

// rollBack puts safety — the pre-restore snapshot of a restore that did not
// finish — back over DataDir, BackofficeDir included, and turns it
// StatusReady once it is. It stays StatusRollback when that fails, so the
// next StartRecover tries again.
func rollBack(sandbox *api.Sandbox, db *backup.Backup, safety backup.SnapshotRecord) error {
	contents, err := db.ListSnapshotContents(safety.Id)
	if err != nil {
		return err
	}
	if err := write(sandbox, db, contents, true); err != nil {
		return err
	}
	return db.SetSnapshotStatus(safety.Id, StatusReady)
}

// RecoverRestores puts back every snapshot left StatusRollback by a restore
// that did not finish — the process stopped in the middle, or its roll-back
// failed — and answers the ones it put back. The caller holds the job token
// and the restoring one.
func RecoverRestores(sandbox *api.Sandbox) ([]string, error) {
	db := backup.New(sandbox)
	pending, err := db.ListSnapshots(backup.SnapshotFilter{StatusEquals: StatusRollback})
	if err != nil {
		return nil, err
	}
	sandbox.Deps.SortDeps.SliceStable(pending, func(i int, j int) bool {
		return pending[i].Id > pending[j].Id
	})
	names := []string{}
	for index, safety := range pending {
		// The newest one is what DataDir held before the last restore began;
		// any older one was already superseded by it.
		if index > 0 {
			if err := db.SetSnapshotStatus(safety.Id, StatusReady); err != nil {
				return names, err
			}
			continue
		}
		if err := rollBack(sandbox, db, safety); err != nil {
			return names, err
		}
		names = append(names, safety.Name)
	}
	return names, nil
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
