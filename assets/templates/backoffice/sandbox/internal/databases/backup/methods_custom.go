package backup

import (
	api "{{.Module}}/sandbox/api"
)

// BlobSha is one blob record without its value: its permanent id and its sha.
type BlobSha struct {
	Id  int64
	Sha string
}

// RepairBackup runs the store's Repair over the blob and the snapshot
// collections, the content nested under every snapshot included: what a
// process stopped in the middle of a write left behind — a blob's bytes among
// it — is deleted, and every invariant of both is restored. No generated
// method reaches Repair, so it is written here.
func RepairBackup(sandbox *api.Sandbox, db *Backup) error {
	for _, name := range []string{"blob", "snapshot"} {
		collection, err := sandbox.Deps.OpinionatedAgnosDatabase.Collection(db.database, name)
		if err != nil {
			return err
		}
		if failure := collection.Repair(); failure != nil {
			return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
		}
	}
	return nil
}

// ListBlobShas reads the id and the sha of every blob record, never its
// value. ListBlobs would load every stored content to answer the same, so a
// sweep over every blob — the one that removes the blobs no snapshot holds —
// is written here.
func ListBlobShas(sandbox *api.Sandbox, db *Backup) ([]BlobSha, error) {
	collection, err := sandbox.Deps.OpinionatedAgnosDatabase.Collection(db.database, "blob")
	if err != nil {
		return nil, err
	}
	records, failure := collection.ListAll()
	if failure != nil {
		return nil, sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	listed := []BlobSha{}
	for _, record := range records {
		sha, err := sandbox.Deps.OpinionatedAgnosDatabase.ReadString(record, "sha")
		if err != nil {
			return nil, err
		}
		listed = append(listed, BlobSha{Id: record.ID, Sha: sha})
	}
	return listed, nil
}

// HasBlob tells whether a blob record holds sha, never reading its value:
// FindBlobBySha would load the whole stored content to answer the same.
func HasBlob(sandbox *api.Sandbox, db *Backup, sha string) (bool, error) {
	collection, err := sandbox.Deps.OpinionatedAgnosDatabase.Collection(db.database, "blob")
	if err != nil {
		return false, err
	}
	_, ok, failure := collection.FindByKey("sha", sha)
	if failure != nil {
		return false, sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	return ok, nil
}

// RemoveSnapshotContent deletes one content record nested under one snapshot
// record. The generated methods add and list the content of a snapshot, never
// remove one of it, so it is written here. A content record already gone is
// not an error.
func RemoveSnapshotContent(sandbox *api.Sandbox, db *Backup, snapshotId int64, contentId int64) error {
	collection, err := sandbox.Deps.OpinionatedAgnosDatabase.Collection(db.database, "snapshot")
	if err != nil {
		return err
	}
	parent, ok, failure := collection.FindByID(snapshotId)
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	if !ok {
		return sandbox.Deps.StdDeps.Errorf("snapshot %d not found", snapshotId)
	}
	nested, failure := parent.Nested("content")
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	record, ok, failure := nested.FindByID(contentId)
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	if !ok {
		return nil
	}
	return sandbox.Deps.OpinionatedAgnosDatabase.Fail(record.Remove())
}
