package snapshots

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// A snapshot can also be built by hand, one file at a time, over as many
// requests as it takes: CreateOpen records it empty under StatusOpen, AddFile
// stores one file in it, AddBlob stores a content alone and AddReference names
// one already stored at a path of it, and Close turns it StatusReady once
// every content it names is stored. Until then it is neither downloaded whole
// nor restored, and a server that stops leaves it open: RecoverInterrupted
// marks StatusCreating alone. ListFiles and ReadFile read any snapshot one file
// at a time.
//
// Adding a path the snapshot already holds replaces that file. Every pair
// stays until Close, and of the ones at one path the pair of the highest id —
// the last added, as an id is never reused — is the one listed, read and
// kept; Close removes the others.

// shaPattern is the sha a content is stored under: the SHA-256 of its bytes,
// in lowercase hex.
const shaPattern = `^[0-9a-f]{64}$`

// CreateOpen records a new empty snapshot under StatusOpen. It is named name,
// its surrounding spaces trimmed; an empty one names it manual followed by the
// current instant. The outcome is OutcomeOk, OutcomeInvalidName,
// OutcomeNameTaken or OutcomeBusy, as for StartCreate; nothing is recorded
// unless OutcomeOk.
func CreateOpen(sandbox *api.Sandbox, name string) (backup.SnapshotRecord, string, error) {
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
	defer release(sandbox)

	snapshot, outcome, err := record(sandbox, backup.New(sandbox), name, "manual", StatusOpen)
	if err != nil || outcome != OutcomeOk {
		return backup.SnapshotRecord{}, outcome, err
	}
	sandbox.Deps.StdDeps.Logf("snapshot %s opened\n", snapshot.Name)
	return snapshot, OutcomeOk, nil
}

// AddFile stores content as the file at path — slash-separated, relative to
// DataDir, a leading slash dropped — of the open snapshot of id, replacing the
// one it held there. The outcome is OutcomeOk, OutcomeInvalidPath for a path
// SafePath refuses, OutcomeNotFound, OutcomeNotOpen, or OutcomeBusy while
// another job runs; nothing is stored unless OutcomeOk.
func AddFile(sandbox *api.Sandbox, id int64, path string, content []byte) (backup.SnapshotContentRecord, string, error) {
	path = cleanPath(sandbox, path)
	if !SafePath(sandbox, path) {
		return backup.SnapshotContentRecord{}, OutcomeInvalidPath, nil
	}
	if !acquire(sandbox) {
		return backup.SnapshotContentRecord{}, OutcomeBusy, nil
	}
	defer release(sandbox)

	db := backup.New(sandbox)
	if _, outcome := opened(sandbox, db, id); outcome != OutcomeOk {
		return backup.SnapshotContentRecord{}, outcome, nil
	}
	sha, err := storeBlob(sandbox, db, content)
	if err != nil {
		return backup.SnapshotContentRecord{}, "", err
	}
	file, err := db.AddSnapshotContent(id, backup.SnapshotContentInput{Path: path, Sha: sha})
	if err != nil {
		return backup.SnapshotContentRecord{}, "", err
	}
	return file, OutcomeOk, nil
}

// AddBlob stores content once, under its SHA-256, and answers that sha: the
// one AddReference names it by. The outcome is OutcomeOk, or OutcomeBusy while
// another job runs — an optimize running beside it would remove it at once.
// A blob no snapshot names yet is one the next optimize removes all the same.
func AddBlob(sandbox *api.Sandbox, content []byte) (string, string, error) {
	if !acquire(sandbox) {
		return "", OutcomeBusy, nil
	}
	defer release(sandbox)

	sha, err := storeBlob(sandbox, backup.New(sandbox), content)
	if err != nil {
		return "", "", err
	}
	return sha, OutcomeOk, nil
}

// AddReference names the content stored under sha as the file at path of the
// open snapshot of id, replacing the one it held there. sha is read in
// lowercase, as every content is stored: optimize keeps a blob only when a pair
// names its sha exactly. The outcome is OutcomeOk, OutcomeInvalidPath,
// OutcomeInvalidSha, OutcomeNotFound, OutcomeNotOpen, OutcomeBlobMissing when
// no content is stored under sha — never sent, or removed by an optimize before
// a pair named it — or OutcomeBusy; nothing is stored unless OutcomeOk.
func AddReference(sandbox *api.Sandbox, id int64, path string, sha string) (backup.SnapshotContentRecord, string, error) {
	strings := sandbox.Deps.StringsDeps
	path = cleanPath(sandbox, path)
	if !SafePath(sandbox, path) {
		return backup.SnapshotContentRecord{}, OutcomeInvalidPath, nil
	}
	sha = strings.ToLower(strings.TrimSpace(sha))
	valid, err := strings.MatchPattern(shaPattern, sha)
	if err != nil {
		return backup.SnapshotContentRecord{}, "", err
	}
	if !valid {
		return backup.SnapshotContentRecord{}, OutcomeInvalidSha, nil
	}
	if !acquire(sandbox) {
		return backup.SnapshotContentRecord{}, OutcomeBusy, nil
	}
	defer release(sandbox)

	db := backup.New(sandbox)
	if _, outcome := opened(sandbox, db, id); outcome != OutcomeOk {
		return backup.SnapshotContentRecord{}, outcome, nil
	}
	stored, err := backup.HasBlob(sandbox, db, sha)
	if err != nil {
		return backup.SnapshotContentRecord{}, "", err
	}
	if !stored {
		return backup.SnapshotContentRecord{}, OutcomeBlobMissing, nil
	}
	file, err := db.AddSnapshotContent(id, backup.SnapshotContentInput{Path: path, Sha: sha})
	if err != nil {
		return backup.SnapshotContentRecord{}, "", err
	}
	return file, OutcomeOk, nil
}

// Close turns the open snapshot of id StatusReady, once every content its
// files name is stored, and removes the pairs a later one at the same path
// replaced. The outcome is OutcomeOk, OutcomeNotFound, OutcomeNotOpen,
// OutcomeEmpty for a snapshot holding no file, OutcomeConflict with the path
// of a file another one needs as a folder as missing, OutcomeBlobMissing with the
// path whose content is not stored as missing, or OutcomeBusy; the snapshot
// stays open unless OutcomeOk.
func Close(sandbox *api.Sandbox, id int64) (snapshot backup.SnapshotRecord, outcome string, missing string, err error) {
	if !acquire(sandbox) {
		return snapshot, OutcomeBusy, "", nil
	}
	defer release(sandbox)

	db := backup.New(sandbox)
	snapshot, outcome = opened(sandbox, db, id)
	if outcome != OutcomeOk {
		return snapshot, outcome, "", nil
	}
	contents, err := db.ListSnapshotContents(id)
	if err != nil {
		return snapshot, "", "", err
	}
	kept, replaced := latest(sandbox, contents)
	if len(kept) == 0 {
		return snapshot, OutcomeEmpty, "", nil
	}
	paths := []string{}
	for _, file := range kept {
		paths = append(paths, file.Path)
	}
	if clash := Conflict(sandbox, paths); clash != "" {
		return snapshot, OutcomeConflict, clash, nil
	}
	for _, file := range kept {
		stored, err := backup.HasBlob(sandbox, db, file.Sha)
		if err != nil {
			return snapshot, "", "", err
		}
		if !stored {
			return snapshot, OutcomeBlobMissing, file.Path, nil
		}
	}

	for _, file := range replaced {
		if err := backup.RemoveSnapshotContent(sandbox, db, id, file.Id); err != nil {
			return snapshot, "", "", err
		}
	}
	if err := db.SetSnapshotStatus(id, StatusReady); err != nil {
		return snapshot, "", "", err
	}
	snapshot.Status = StatusReady
	sandbox.Deps.StdDeps.Logf("snapshot %s closed: %d files\n", snapshot.Name, len(kept))
	return snapshot, OutcomeOk, "", nil
}

// ListFiles is the snapshot of id and every file it holds whose path starts
// with prefix, in path order: every one when prefix, its surrounding spaces
// and a leading slash trimmed, is empty. A snapshot of any status is listed —
// an open one as it stands, a path added twice once, as last added. The
// outcome is OutcomeOk or OutcomeNotFound.
func ListFiles(sandbox *api.Sandbox, id int64, prefix string) (backup.SnapshotRecord, []backup.SnapshotContentRecord, string, error) {
	strings := sandbox.Deps.StringsDeps
	prefix = cleanPath(sandbox, strings.TrimSpace(prefix))
	db := backup.New(sandbox)
	snapshot, ok := db.FindSnapshotById(id)
	if !ok {
		return snapshot, nil, OutcomeNotFound, nil
	}
	contents, err := db.ListSnapshotContents(id)
	if err != nil {
		return snapshot, nil, "", err
	}
	kept, _ := latest(sandbox, contents)
	listed := []backup.SnapshotContentRecord{}
	for _, file := range kept {
		if strings.HasPrefix(file.Path, prefix) {
			listed = append(listed, file)
		}
	}
	return snapshot, listed, OutcomeOk, nil
}

// ReadFile is the file at path — a leading slash dropped — of the snapshot of
// id, of any status, and its content. The outcome is OutcomeOk,
// OutcomeNotFound, OutcomeFileNotFound, or OutcomeBlobMissing when its content
// is not stored.
func ReadFile(sandbox *api.Sandbox, id int64, path string) (backup.SnapshotContentRecord, []byte, string, error) {
	path = cleanPath(sandbox, path)
	db := backup.New(sandbox)
	if _, ok := db.FindSnapshotById(id); !ok {
		return backup.SnapshotContentRecord{}, nil, OutcomeNotFound, nil
	}
	contents, err := db.ListSnapshotContents(id)
	if err != nil {
		return backup.SnapshotContentRecord{}, nil, "", err
	}
	found := false
	file := backup.SnapshotContentRecord{}
	for _, content := range contents {
		if content.Path == path && (!found || content.Id > file.Id) {
			file = content
			found = true
		}
	}
	if !found {
		return file, nil, OutcomeFileNotFound, nil
	}
	blob, ok := db.FindBlobBySha(file.Sha)
	if !ok {
		return file, nil, OutcomeBlobMissing, nil
	}
	return file, blob.Value, OutcomeOk, nil
}

// opened is the snapshot of id and whether files may be added to it:
// OutcomeOk when it is StatusOpen, OutcomeNotFound or OutcomeNotOpen
// otherwise.
func opened(sandbox *api.Sandbox, db *backup.Backup, id int64) (backup.SnapshotRecord, string) {
	snapshot, ok := db.FindSnapshotById(id)
	if !ok {
		return snapshot, OutcomeNotFound
	}
	if snapshot.Status != StatusOpen {
		return snapshot, OutcomeNotOpen
	}
	return snapshot, OutcomeOk
}

// latest splits contents into the pair each path is held by — of the ones at
// that path, the one of the highest id — in path order, and every pair a
// later one at its path replaced.
func latest(sandbox *api.Sandbox, contents []backup.SnapshotContentRecord) (kept []backup.SnapshotContentRecord, replaced []backup.SnapshotContentRecord) {
	at := map[string]int{}
	kept = []backup.SnapshotContentRecord{}
	replaced = []backup.SnapshotContentRecord{}
	for _, content := range contents {
		index, seen := at[content.Path]
		if !seen {
			at[content.Path] = len(kept)
			kept = append(kept, content)
			continue
		}
		if content.Id > kept[index].Id {
			replaced = append(replaced, kept[index])
			kept[index] = content
		} else {
			replaced = append(replaced, content)
		}
	}
	sandbox.Deps.SortDeps.SliceStable(kept, func(i int, j int) bool {
		return kept[i].Path < kept[j].Path
	})
	return kept, replaced
}

// cleanPath is path, a file of a snapshot as an address or a body hands it, as
// a snapshot stores it: without a leading slash.
func cleanPath(sandbox *api.Sandbox, path string) string {
	return sandbox.Deps.StringsDeps.TrimPrefix(path, "/")
}
