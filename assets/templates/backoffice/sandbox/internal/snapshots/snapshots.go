package snapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// A snapshot is every file of every database under DataDir — the folder
// --database names, the backup database itself, under BackupDir, aside — as
// it was at one instant, kept in
// the backup database: the content of each file once in the blob table,
// indexed by its SHA-256, and the snapshot a list of (path, sha) pairs under
// its content. A file that did not change between two snapshots is stored
// once, whatever the number of snapshots holding it.
//
// Creating one, restoring one and removing the blobs no snapshot holds are
// slow, so a route only starts them: the work runs on a goroutine of its own
// once the route has answered. A snapshot is recorded the moment it is
// started, under StatusCreating, and turns StatusReady or StatusFailed when
// its job ends. Every job — those three, an upload, a removal and every write
// of manual.go — runs one at a time.

// DataDir is the directory every database of the project lives under, one
// directory per database: sandbox.Config.DatabaseDir, which --database sets on
// every command line and which is data by default.
func DataDir(sandbox *api.Sandbox) string {
	return sandbox.Config.DatabaseDir
}

// BackupDir is the directory of DataDir the backup database lives in, its
// key-prefix. It is never put in a snapshot, and a restore never touches it.
const BackupDir = "backup"

// The status of a snapshot.
const (
	// StatusCreating is a snapshot whose files are still being stored.
	StatusCreating = "creating"
	// StatusReady is a snapshot holding every file it was taken with.
	StatusReady = "ready"
	// StatusFailed is a snapshot whose job failed or was interrupted: it
	// holds part of its files at most, so it is never restored nor
	// downloaded.
	StatusFailed = "failed"
	// StatusOpen is a snapshot built by hand, still taking files one at a
	// time until Close turns it StatusReady. It is never restored nor
	// downloaded whole, and a server that stops leaves it open.
	StatusOpen = "open"
	// StatusRollback is the pre-restore snapshot of a restore that has not
	// finished: what DataDir held before it began. A restore that fails, or
	// the next start after one the process stopped in the middle of, puts it
	// back and turns it StatusReady. It is never restored, downloaded nor
	// removed by hand while it holds this status.
	StatusRollback = "rollback"
)

// The outcomes StartRestore, Export, Import and Remove answer with.
const (
	// OutcomeOk is a restore started, an archive built or stored, or a
	// snapshot removed.
	OutcomeOk = "ok"
	// OutcomeBusy is a job refused because another one runs.
	OutcomeBusy = "busy"
	// OutcomeRefused is an archive that is not one a snapshot can be read
	// from.
	OutcomeRefused = "refused"
	// OutcomeNotFound is an id no snapshot holds.
	OutcomeNotFound = "not-found"
	// OutcomeNotReady is a snapshot that is not StatusReady.
	OutcomeNotReady = "not-ready"
	// OutcomeInvalidName is a name a snapshot cannot be given: one
	// namePattern refuses.
	OutcomeInvalidName = "invalid-name"
	// OutcomeNameTaken is a name another snapshot already holds.
	OutcomeNameTaken = "name-taken"
	// OutcomeNotOpen is a file added to, or a close of, a snapshot that is
	// not StatusOpen.
	OutcomeNotOpen = "not-open"
	// OutcomeInvalidPath is a file path SafePath refuses.
	OutcomeInvalidPath = "invalid-path"
	// OutcomeInvalidSha is a sha that is not the lowercase hex SHA-256 a
	// content is stored under.
	OutcomeInvalidSha = "invalid-sha"
	// OutcomeEmpty is a close of a snapshot holding no file: restoring it
	// would empty every database.
	OutcomeEmpty = "empty"
	// OutcomeFileNotFound is a path the snapshot holds no file at.
	OutcomeFileNotFound = "file-not-found"
	// OutcomeBlobMissing is a sha no stored content has.
	OutcomeBlobMissing = "blob-missing"
	// OutcomeConflict is a close of a snapshot holding a file at a path
	// another of its files needs as a folder: no restore could write both.
	OutcomeConflict = "conflict"
)

// MaxArchiveBytes is the largest archive an upload may send.
const MaxArchiveBytes = 256 << 20

// MaxUnpackedBytes is the most an uploaded archive may hold once unpacked.
const MaxUnpackedBytes = 1 << 30

// NameLayout is how the instant a snapshot was taken at is spelled in its
// name, after its prefix.
const NameLayout = "20060102-150405"

// keepLockSuffix ends the lock files the store leaves beside a key while it
// writes it.
const keepLockSuffix = ".keeplock"

// jobs holds a token while a job runs, so a second one is refused rather than
// run over the first: a create and a restore read and write every database
// of DataDir, and an optimize removes every blob no (path, sha) pair names —
// one a create or an upload is about to name included.
var jobs = make(chan struct{}, 1)

// acquire takes the job token and reports whether it was free.
func acquire(sandbox *api.Sandbox) bool {
	select {
	case jobs <- struct{}{}:
		return true
	default:
		return false
	}
}

// release gives the job token back.
func release(sandbox *api.Sandbox) {
	<-jobs
}

// Busy tells whether a job is running.
func Busy(sandbox *api.Sandbox) bool {
	return len(jobs) > 0
}

// creating holds a token while a snapshot reads DataDir file by file.
var creating = make(chan struct{}, 1)

// Creating tells whether a snapshot is reading DataDir: a request that may
// write a database is refused until it ends, so every file of the snapshot is
// taken at the same instant.
func Creating(sandbox *api.Sandbox) bool {
	return len(creating) > 0
}

// nowSeconds is the current instant in Unix seconds.
func nowSeconds(sandbox *api.Sandbox) int64 {
	return sandbox.Deps.StdDeps.Now() / 1_000_000_000
}

// StartCreate records a new snapshot under StatusCreating and answers it at
// once, leaving its files to a goroutine that runs after the caller returns.
// The snapshot is named name, its surrounding spaces trimmed; an empty one
// names it snapshot followed by the current instant. The outcome is
// OutcomeOk, OutcomeInvalidName for a name namePattern refuses,
// OutcomeNameTaken for one another snapshot holds, or OutcomeBusy while
// another job runs; nothing is recorded unless OutcomeOk.
func StartCreate(sandbox *api.Sandbox, name string) (snapshot backup.SnapshotRecord, outcome string, err error) {
	name, valid, err := cleanName(sandbox, name)
	if err != nil {
		return backup.SnapshotRecord{}, "", err
	}
	if !valid {
		return backup.SnapshotRecord{}, OutcomeInvalidName, nil
	}
	if !acquire(sandbox) {
		return backup.SnapshotRecord{}, OutcomeBusy, nil
	}
	db := backup.New(sandbox)
	snapshot, outcome, err = record(sandbox, db, name, "snapshot", StatusCreating)
	if err != nil || outcome != OutcomeOk {
		release(sandbox)
		return backup.SnapshotRecord{}, outcome, err
	}
	creating <- struct{}{}
	go func() {
		defer release(sandbox)
		defer func() { <-creating }()
		finish(sandbox, db, snapshot)
	}()
	return snapshot, OutcomeOk, nil
}

// ValidName tells whether name is one a snapshot may be given: namePattern,
// so it is safe in a file name and a header.
func ValidName(sandbox *api.Sandbox, name string) (bool, error) {
	return sandbox.Deps.StringsDeps.MatchPattern(namePattern, name)
}

// cleanName is name with its surrounding spaces trimmed, and whether a
// snapshot may be given it: an empty one may, standing for one after the
// current instant.
func cleanName(sandbox *api.Sandbox, name string) (string, bool, error) {
	name = sandbox.Deps.StringsDeps.TrimSpace(name)
	if name == "" {
		return name, true, nil
	}
	valid, err := ValidName(sandbox, name)
	return name, valid, err
}

// record records a new snapshot under status, named name — cleanName
// already passed it — or, when name is empty, prefix followed by the current
// instant. The outcome is OutcomeOk, or OutcomeNameTaken for a name another
// snapshot holds. The caller holds the job token.
func record(sandbox *api.Sandbox, db *backup.Backup, name string, prefix string, status string) (backup.SnapshotRecord, string, error) {
	if name == "" {
		snapshot, err := begin(sandbox, db, prefix, status)
		if err != nil {
			return backup.SnapshotRecord{}, "", err
		}
		return snapshot, OutcomeOk, nil
	}
	if _, taken := db.FindSnapshotByName(name); taken {
		return backup.SnapshotRecord{}, OutcomeNameTaken, nil
	}
	snapshot, err := db.AddSnapshot(backup.SnapshotInput{Name: name, Data: nowSeconds(sandbox), Status: status})
	if err != nil {
		return backup.SnapshotRecord{}, "", err
	}
	return snapshot, OutcomeOk, nil
}

// create takes one snapshot named after prefix and waits for it. The caller
// holds the job token.
func create(sandbox *api.Sandbox, db *backup.Backup, prefix string) (backup.SnapshotRecord, error) {
	snapshot, err := begin(sandbox, db, prefix, StatusCreating)
	if err != nil {
		return snapshot, err
	}
	if err := finish(sandbox, db, snapshot); err != nil {
		return snapshot, err
	}
	snapshot.Status = StatusReady
	return snapshot, nil
}

// begin records a new snapshot under status, named prefix followed by the
// current instant.
func begin(sandbox *api.Sandbox, db *backup.Backup, prefix string, status string) (backup.SnapshotRecord, error) {
	now := nowSeconds(sandbox)
	name := uniqueName(sandbox, db, prefix+"-"+sandbox.Deps.TimeDeps.FormatUnix(now, NameLayout))
	return db.AddSnapshot(backup.SnapshotInput{Name: name, Data: now, Status: status})
}

// finish stores every file of DataDir under snapshot and records how that
// ended: StatusReady, or StatusFailed with the cause on stderr. A panic is
// caught and fails the snapshot, so a job never takes the server down.
func finish(sandbox *api.Sandbox, db *backup.Backup, snapshot backup.SnapshotRecord) (err error) {
	files := 0
	defer func() {
		if recovered := recover(); recovered != nil {
			err = sandbox.Deps.StdDeps.Errorf("%v", recovered)
		}
		status := StatusReady
		if err != nil {
			status = StatusFailed
			sandbox.Deps.StdDeps.Eprintf("snapshot %s failed: %s\n", snapshot.Name, err.Error())
		} else {
			sandbox.Deps.StdDeps.Logf("snapshot %s ready: %d files\n", snapshot.Name, files)
		}
		if setErr := db.SetSnapshotStatus(snapshot.Id, status); setErr != nil {
			sandbox.Deps.StdDeps.Eprintf("snapshot %s: recording it %s failed: %s\n", snapshot.Name, status, setErr.Error())
			if err == nil {
				err = setErr
			}
		}
	}()
	files, err = fill(sandbox, db, snapshot.Id)
	return err
}

// fill stores every file collect finds under the snapshot of id, one file at
// a time, and answers how many it stored. A file removed between the listing
// and its reading is left out: the store removes one whenever a record goes.
func fill(sandbox *api.Sandbox, db *backup.Backup, id int64) (int, error) {
	io := sandbox.Deps.IoDeps
	files := 0
	for _, file := range collect(sandbox) {
		content, err := io.ReadFile(file.host)
		if err != nil {
			if !io.Exists(file.host) {
				continue
			}
			return files, err
		}
		sha, err := storeBlob(sandbox, db, content)
		if err != nil {
			return files, err
		}
		if _, err := db.AddSnapshotContent(id, backup.SnapshotContentInput{Path: file.path, Sha: sha}); err != nil {
			return files, err
		}
		files++
	}
	return files, nil
}

// collected is one file a snapshot holds: the host path it is read at, and
// the path it is stored under, relative to DataDir.
type collected struct {
	host string
	path string
}

// collect is every file a snapshot holds: every file under every database
// directory of DataDir but BackupDir, in lexical order, leaving out the
// temporary and lock files the store keeps while it writes.
func collect(sandbox *api.Sandbox) []collected {
	io := sandbox.Deps.IoDeps
	root := DataDir(sandbox)
	backupDir := io.Join(root, BackupDir)
	files := []collected{}
	for _, dir := range io.ListDirs(root) {
		if dir == backupDir {
			continue
		}
		for _, path := range io.ListFilesRecursively(dir) {
			if kept(sandbox, path) {
				files = append(files, collected{host: path, path: relative(sandbox, dir, path)})
			}
		}
	}
	return files
}

// kept tells whether the file at path belongs in a snapshot: the store names
// its temporary files with a leading "." and its locks with keepLockSuffix,
// which no file of a record ever carries.
func kept(sandbox *api.Sandbox, path string) bool {
	strings := sandbox.Deps.StringsDeps
	slashed := strings.ReplaceAll(path, "\\", "/")
	name := slashed[strings.LastIndex(slashed, "/")+1:]
	return !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, keepLockSuffix)
}

// relative is path, the host path of a file under dir — a database directory
// io.ListDirs found in DataDir — as a snapshot stores it: slash-separated and
// relative to DataDir, the name of dir first. It is read off dir rather than
// off DataDir, so a --database spelled ./data, absolute or with a trailing
// slash stores the same paths.
func relative(sandbox *api.Sandbox, dir string, path string) string {
	strings := sandbox.Deps.StringsDeps
	slashedDir := strings.ReplaceAll(dir, "\\", "/")
	slashed := strings.ReplaceAll(path, "\\", "/")
	name := slashedDir[strings.LastIndex(slashedDir, "/")+1:]
	return name + "/" + strings.TrimPrefix(slashed, slashedDir+"/")
}

// storeBlob stores content once, under its SHA-256, and answers that sha.
// Content already stored is not stored again; an insert that lost a race
// against another one storing the same content reads as stored.
func storeBlob(sandbox *api.Sandbox, db *backup.Backup, content []byte) (string, error) {
	sha := sandbox.Deps.HashDeps.Sha256Hex(content)
	if _, ok := db.FindBlobBySha(sha); ok {
		return sha, nil
	}
	if _, err := db.AddBlob(backup.BlobInput{Sha: sha, Value: content}); err != nil {
		if _, ok := db.FindBlobBySha(sha); ok {
			return sha, nil
		}
		return "", err
	}
	return sha, nil
}

// uniqueName is base when no snapshot holds that name, base-2, base-3 and on
// otherwise: the first one free.
func uniqueName(sandbox *api.Sandbox, db *backup.Backup, base string) string {
	name := base
	for count := 2; ; count++ {
		if _, ok := db.FindSnapshotByName(name); !ok {
			return name
		}
		name = base + "-" + sandbox.Deps.StringsDeps.FormatInt(int64(count), 10)
	}
}

// List is every snapshot whose name starts with prefix, newest first: every
// one when prefix, its surrounding spaces trimmed, is empty. The match is
// case-sensitive, as names are.
func List(sandbox *api.Sandbox, prefix string) ([]backup.SnapshotRecord, error) {
	prefix = sandbox.Deps.StringsDeps.TrimSpace(prefix)
	listed, err := backup.New(sandbox).ListSnapshots(backup.SnapshotFilter{NameStartsWith: prefix})
	if err != nil {
		return nil, err
	}
	sandbox.Deps.SortDeps.SliceStable(listed, func(i int, j int) bool {
		if listed[i].Data != listed[j].Data {
			return listed[i].Data > listed[j].Data
		}
		return listed[i].Id > listed[j].Id
	})
	return listed, nil
}

// StartRecover runs RecoverRestores and RecoverInterrupted on a goroutine
// holding the job token, so no snapshot or restore starts before what a
// stopped process left is dealt with, and writes on stderr what it did. While
// a restore the last run did not finish is put back, every request is refused
// with a 503 (Restoring). It does not wait: a process killed while writing
// leaves its write lease behind, and the first write after it waits until
// that lease expires — up to a minute the server would otherwise spend not
// listening.
func StartRecover(sandbox *api.Sandbox) {
	if !acquire(sandbox) {
		return
	}
	restoring <- struct{}{}
	go func() {
		defer release(sandbox)
		rolledBack, err := RecoverRestores(sandbox)
		<-restoring
		if err != nil {
			sandbox.Deps.StdDeps.Eprintf("warning: a restore the last run did not finish could not be undone: %s; it is tried again at the next start\n", err.Error())
		}
		for _, name := range rolledBack {
			sandbox.Deps.StdDeps.Eprintf("warning: the last run stopped in the middle of a restore; %s was put back as it was before it, from snapshot %s: run the restore again\n", DataDir(sandbox), name)
		}
		interrupted, err := RecoverInterrupted(sandbox)
		if err != nil {
			sandbox.Deps.StdDeps.Eprintf("warning: the snapshots left creating by the last run could not be marked failed: %s\n", err.Error())
		} else if interrupted > 0 {
			sandbox.Deps.StdDeps.Eprintf("warning: %d snapshot(s) were still being created when the last run stopped, and are marked failed\n", interrupted)
		}
	}()
}

// RecoverInterrupted marks StatusFailed every snapshot left StatusCreating by
// a process that stopped before its job ended, and answers how many it
// marked. The caller holds the job token, so none of them is being created.
func RecoverInterrupted(sandbox *api.Sandbox) (int, error) {
	db := backup.New(sandbox)
	stuck, err := db.ListSnapshots(backup.SnapshotFilter{StatusEquals: StatusCreating})
	if err != nil {
		return 0, err
	}
	for _, snapshot := range stuck {
		if err := db.SetSnapshotStatus(snapshot.Id, StatusFailed); err != nil {
			return 0, err
		}
	}
	return len(stuck), nil
}
