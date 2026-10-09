package snapshots

import (
	"{{.Module}}/sandbox/api"
	archivedeps "{{.Module}}/sandbox/deps/archivedeps"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// A snapshot travels as a zip archive: ManifestPath, naming it, and every one
// of its files under ArchiveDataDir at the path it holds below DataDir — so
// the archive unzipped at the root of a project running with the default
// --database puts the databases back by hand. The archive's folder is data/
// whatever --database names: an archive taken from one folder restores into
// another.

// ManifestPath is the archive entry naming the snapshot: a JSON object with
// its name, the Unix instant it was taken at as data, and how many files it
// holds.
const ManifestPath = "snapshot.json"

// ArchiveDataDir is the folder of the archive every file of the snapshot
// sits under.
const ArchiveDataDir = "data/"

// namePattern is a snapshot name a root may give StartCreate or an upload
// may keep from its manifest: one that is safe in a file name and a header.
const namePattern = `^[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`

// Export packs the snapshot of id into its archive. The outcome is OutcomeOk,
// OutcomeNotFound or OutcomeNotReady; the archive is nil unless OutcomeOk.
func Export(sandbox *api.Sandbox, id int64) (backup.SnapshotRecord, []byte, string, error) {
	db := backup.New(sandbox)
	snapshot, ok := db.FindSnapshotById(id)
	if !ok {
		return snapshot, nil, OutcomeNotFound, nil
	}
	if snapshot.Status != StatusReady {
		return snapshot, nil, OutcomeNotReady, nil
	}
	contents, err := db.ListSnapshotContents(id)
	if err != nil {
		return snapshot, nil, "", err
	}

	files := []archivedeps.File{ {Path: ManifestPath, Content: manifest(sandbox, snapshot, len(contents))} }
	for _, content := range contents {
		blob, ok := db.FindBlobBySha(content.Sha)
		if !ok {
			return snapshot, nil, "", sandbox.Deps.StdDeps.Errorf("snapshot %s: the content of %s is missing", snapshot.Name, content.Path)
		}
		files = append(files, archivedeps.File{Path: ArchiveDataDir + content.Path, Content: blob.Value})
	}
	archive, err := sandbox.Deps.ArchiveDeps.Zip(files)
	if err != nil {
		return snapshot, nil, "", err
	}
	return snapshot, archive, OutcomeOk, nil
}

// manifest is the ManifestPath entry of the archive of snapshot, which holds
// files files.
func manifest(sandbox *api.Sandbox, snapshot backup.SnapshotRecord, files int) []byte {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("name", snapshot.Name)
	document.AddItemToObject("data", snapshot.Data)
	document.AddItemToObject("files", files)
	return []byte(sandbox.Deps.SerializableDeps.SerializeToJson(document))
}

// Import stores the snapshot an archive Export built — or one written by
// hand the same way — as a new StatusReady snapshot, and answers it with
// OutcomeOk. The name and the instant come from ManifestPath when it carries
// usable ones, and the name is made unique; without them the snapshot is
// named upload followed by the current instant.
//
// The archive is checked first, and an archive that is not one a snapshot
// can be read from answers OutcomeRefused with why in refused: one that is
// not a zip, unpacks past MaxUnpackedBytes, holds an entry outside
// ArchiveDataDir, a path SafePath refuses, the same path twice, a file at a
// path other entries need as a folder, or no file at all. Storing it is a job — an optimize running beside it could remove a
// blob it is about to name — so while another one runs it answers
// OutcomeBusy and stores nothing.
func Import(sandbox *api.Sandbox, archive []byte) (snapshot backup.SnapshotRecord, outcome string, refused string, err error) {
	files, name, data, refused, err := unpack(sandbox, archive)
	if err != nil {
		return snapshot, "", "", err
	}
	if refused != "" {
		return snapshot, OutcomeRefused, refused, nil
	}
	if !acquire(sandbox) {
		return snapshot, OutcomeBusy, "", nil
	}
	defer release(sandbox)

	db := backup.New(sandbox)
	snapshot, err = db.AddSnapshot(backup.SnapshotInput{Name: uniqueName(sandbox, db, name), Data: data, Status: StatusCreating})
	if err != nil {
		return snapshot, "", "", err
	}
	if err := store(sandbox, db, snapshot.Id, files); err != nil {
		if removeErr := db.RemoveSnapshot(snapshot.Id); removeErr != nil {
			sandbox.Deps.StdDeps.Eprintf("snapshot %s: removing it after a failed upload failed: %s\n", snapshot.Name, removeErr.Error())
		}
		return backup.SnapshotRecord{}, "", "", err
	}
	if err := db.SetSnapshotStatus(snapshot.Id, StatusReady); err != nil {
		return snapshot, "", "", err
	}
	snapshot.Status = StatusReady
	sandbox.Deps.StdDeps.Logf("snapshot %s uploaded: %d files\n", snapshot.Name, len(files))
	return snapshot, OutcomeOk, "", nil
}

// unpack reads archive into the files a snapshot of it holds, each at its
// path below DataDir, and the name and the instant to store it under — the
// manifest's, or upload followed by now. refused is why the archive cannot
// be one, "" when it can.
func unpack(sandbox *api.Sandbox, archive []byte) (files []archivedeps.File, name string, data int64, refused string, err error) {
	strings := sandbox.Deps.StringsDeps
	unpacked, err := sandbox.Deps.ArchiveDeps.Unzip(archive, MaxUnpackedBytes)
	if err != nil {
		return nil, "", 0, "not a zip archive a snapshot can be read from: " + err.Error(), nil
	}

	data = nowSeconds(sandbox)
	files = []archivedeps.File{}
	seen := map[string]bool{}
	for _, file := range unpacked {
		if file.Path == ManifestPath {
			named, taken, ok := readManifest(sandbox, file.Content)
			if !ok {
				return nil, "", 0, ManifestPath + " is not a JSON object", nil
			}
			name = named
			if taken > 0 {
				data = taken
			}
			continue
		}
		path := strings.TrimPrefix(file.Path, ArchiveDataDir)
		if !strings.HasPrefix(file.Path, ArchiveDataDir) || !SafePath(sandbox, path) {
			return nil, "", 0, "the entry " + strings.Quote(file.Path) + " is not a file a snapshot may hold: every file sits in a database folder under " + ArchiveDataDir + ", never " + ArchiveDataDir + BackupDir, nil
		}
		if seen[path] {
			return nil, "", 0, "the entry " + strings.Quote(file.Path) + " is in the archive twice", nil
		}
		seen[path] = true
		files = append(files, archivedeps.File{Path: path, Content: file.Content})
	}
	if len(files) == 0 {
		return nil, "", 0, "the archive holds no file under " + ArchiveDataDir, nil
	}
	paths := []string{}
	for _, file := range files {
		paths = append(paths, file.Path)
	}
	if clash := Conflict(sandbox, paths); clash != "" {
		return nil, "", 0, "the entry " + strings.Quote(ArchiveDataDir+clash) + " is a file, and other entries sit inside it as a folder: no restore could write both", nil
	}

	valid, err := ValidName(sandbox, name)
	if err != nil {
		return nil, "", 0, "", err
	}
	if !valid {
		name = "upload-" + sandbox.Deps.TimeDeps.FormatUnix(nowSeconds(sandbox), NameLayout)
	}
	return files, name, data, "", nil
}

// store stores every one of files under the snapshot of id.
func store(sandbox *api.Sandbox, db *backup.Backup, id int64, files []archivedeps.File) error {
	for _, file := range files {
		sha, err := storeBlob(sandbox, db, file.Content)
		if err != nil {
			return err
		}
		if _, err := db.AddSnapshotContent(id, backup.SnapshotContentInput{Path: file.Path, Sha: sha}); err != nil {
			return err
		}
	}
	return nil
}

// readManifest reads the name and the instant out of a ManifestPath entry:
// "" and 0 for either one it does not carry as a string and an integer. ok
// is false when content is not a JSON object.
func readManifest(sandbox *api.Sandbox, content []byte) (name string, data int64, ok bool) {
	document, err := sandbox.Deps.SerializableDeps.ParseJson(string(content))
	if err != nil || !document.IsObject() {
		return "", 0, false
	}
	if item, err := document.GetObjectItem("name"); err == nil && item != nil && item.IsString() {
		name, _ = item.GetString()
	}
	if item, err := document.GetObjectItem("data"); err == nil && item != nil && item.IsInt() {
		data, _ = item.GetInt()
	}
	return name, data, true
}
