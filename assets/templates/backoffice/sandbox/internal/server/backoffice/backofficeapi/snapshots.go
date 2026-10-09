package backofficeapi

import (
	"{{.Module}}/sandbox/api"
	serializabledeps "{{.Module}}/sandbox/deps/serializabledeps"
	"{{.Module}}/sandbox/internal/databases/backup"
)

// SnapshotJSON is the JSON object a snapshot is answered as: its id, name,
// the Unix instant it was taken at as data, and its status. Its files are
// listed apart, by SnapshotFilesJSON, and their contents never travel in
// JSON: they are the archive and the files the download routes answer.
func SnapshotJSON(sandbox *api.Sandbox, snapshot backup.SnapshotRecord) *serializabledeps.SerializableObject {
	object := sandbox.Deps.SerializableDeps.CreateObject()
	object.AddItemToObject("id", snapshot.Id)
	object.AddItemToObject("name", snapshot.Name)
	object.AddItemToObject("data", snapshot.Data)
	object.AddItemToObject("status", snapshot.Status)
	return object
}

// SnapshotResponseJSON is {"snapshot": SnapshotJSON(snapshot)}.
func SnapshotResponseJSON(sandbox *api.Sandbox, snapshot backup.SnapshotRecord) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("snapshot", SnapshotJSON(sandbox, snapshot))
	return document
}

// SnapshotListJSON is every snapshot, newest first, and whether a backup job
// is running.
func SnapshotListJSON(sandbox *api.Sandbox, listed []backup.SnapshotRecord, busy bool) *serializabledeps.SerializableObject {
	items := sandbox.Deps.SerializableDeps.CreateArray()
	for _, snapshot := range listed {
		items.AddItemToArray(SnapshotJSON(sandbox, snapshot))
	}

	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("snapshots", items)
	document.AddItemToObject("busy", busy)
	return document
}

// SnapshotCreatingJSON is {"status": "creating", "snapshot": SnapshotJSON(snapshot)}:
// a snapshot recorded and still being filled.
func SnapshotCreatingJSON(sandbox *api.Sandbox, snapshot backup.SnapshotRecord) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("status", "creating")
	document.AddItemToObject("snapshot", SnapshotJSON(sandbox, snapshot))
	return document
}

// JobJSON is {"status": status}: a job started and not waited for —
// "restoring", "optimizing".
func JobJSON(sandbox *api.Sandbox, status string) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("status", status)
	return document
}

// SnapshotFileJSON is the JSON object one file of a snapshot is answered as:
// its path below the --database folder and the sha of its content.
func SnapshotFileJSON(sandbox *api.Sandbox, file backup.SnapshotContentRecord) *serializabledeps.SerializableObject {
	object := sandbox.Deps.SerializableDeps.CreateObject()
	object.AddItemToObject("path", file.Path)
	object.AddItemToObject("sha", file.Sha)
	return object
}

// SnapshotFileResponseJSON is {"file": SnapshotFileJSON(file)}.
func SnapshotFileResponseJSON(sandbox *api.Sandbox, file backup.SnapshotContentRecord) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("file", SnapshotFileJSON(sandbox, file))
	return document
}

// SnapshotFilesJSON is the snapshot and every one of its files listed, in
// path order.
func SnapshotFilesJSON(sandbox *api.Sandbox, snapshot backup.SnapshotRecord, files []backup.SnapshotContentRecord) *serializabledeps.SerializableObject {
	items := sandbox.Deps.SerializableDeps.CreateArray()
	for _, file := range files {
		items.AddItemToArray(SnapshotFileJSON(sandbox, file))
	}

	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("snapshot", SnapshotJSON(sandbox, snapshot))
	document.AddItemToObject("files", items)
	return document
}

// BlobResponseJSON is {"blob": {"sha": sha, "size": size}}: a content stored,
// the sha a snapshot names it by and how many bytes it holds.
func BlobResponseJSON(sandbox *api.Sandbox, sha string, size int) *serializabledeps.SerializableObject {
	blob := sandbox.Deps.SerializableDeps.CreateObject()
	blob.AddItemToObject("sha", sha)
	blob.AddItemToObject("size", size)

	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("blob", blob)
	return document
}
