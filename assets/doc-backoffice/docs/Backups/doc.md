# Backups

Snapshots of every database of the `--database` folder (`data` by default), taken, downloaded,
uploaded, restored and deleted from the [Backoffice](../Backoffice/doc.md) by a root — or built by
hand over the JSON api, one file per request, for a backup too long for one job or one archive.
Kept in the `backup` database ([Databases](../Databases/backup.md)), under `./data/backup` by
default, gitignored. `{{.GeneratorName}} backoffice-init` wrote every file of it once; each one is
the project's from then on.

## Model

| Table | Holds |
|---|---|
| `blob` | one file content, `sha` its SHA-256 (indexed), `value` its bytes; stored once whatever the number of snapshots holding it |
| `snapshot` | `name` (indexed), `data` (Unix seconds it was taken at), `status`, and `content`: one `(path, sha)` per file, `path` relative to the `--database` folder |

| `status` | Means |
|---|---|
| `creating` | its files are being stored |
| `ready` | holds every file; the only status that downloads or restores |
| `failed` | its job failed or the server stopped during it: `start-server` marks those in the background, create and restore answering `busy` until it is done (up to a minute after a crash, while the store's write lease expires) |
| `open` | built by hand: takes files until it is closed, then `ready`; a server that stops leaves it `open` |

## Where it lives

| Path | What |
|---|---|
| `sandbox/internal/snapshots/` | take (`StartCreate`), restore (`StartRestore`), `Export`, `Import`, `Remove`, `StartOptimize`, the one-job lock; `DataDir` is `sandbox.Config.DatabaseDir` |
| `sandbox/internal/snapshots/manual.go` | by hand: `CreateOpen`, `AddFile`, `AddBlob`, `AddReference`, `Close`, `ListFiles`, `ReadFile` |
| `sandbox/internal/databases/backup/methods_custom.go` | `RepairBackup`; `ListBlobShas` and `HasBlob`, which never read a blob's value; `RemoveSnapshotContent`: one `(path, sha)` of a snapshot |
| `sandbox/internal/server/backoffice/backofficesnapshots/` | list path, notices, messages, the zip and the file download writers |
| `sandbox/internal/server/backoffice/backofficerender/backup_snapshots.go`, `assets/backoffice/backup_snapshots.html` | the list page; its upload form is sent by `assets/front/backoffice/backoffice.js` |
| `sandbox/internal/server/backoffice/backofficeapi/snapshots.go` | the JSON documents |
| `sandbox/internal/routes/backoffice/root/*backup*`, `sandbox/internal/routes/api/backoffice/root/api_*backup*` | the routes |
| `sandbox/internal/commands/middleware/backoffice_start_server/handler.go` | `snapshots.StartRecover`: marks `failed` what the last run left `creating` |
| `sandbox/deps/archivedeps/`, `adapters/impls/ziparchive/` | zip and unzip, over `archive/zip` |
| `sandbox/deps/iodeps/`, `adapters/impls/osio/` | the filesystem the `--database` folder is read and written through |

## Routes

Root only: `backoffice-root-guard` and `backoffice-api-root-guard` answer `403` to anyone else.

| Page | Does |
|---|---|
| GET `/admin/root/list-backups?prefix=` | every snapshot, newest first, or the ones whose name starts with `prefix`; reloads every 5 s while a job runs or one is `creating` |
| POST `/admin/root/create-backup` | form field `name` (optional): starts one, `303` to the list at once |
| POST `/admin/root/restore-backup/{id}` | starts restoring one, `303` to the list at once |
| GET `/admin/root/download-backup/{id}` | the zip, `attachment; filename="<name>.zip"` |
| POST `/admin/root/upload-backup` | the zip as the whole body (`application/zip`), sent by `backoffice.js`: `201`, `400` or `409` |
| POST `/admin/root/remove-backup/{id}` | deletes one, `303` to the list |
| POST `/admin/root/optimize-backup-storage` | starts removing every blob no snapshot holds, `303` to the list at once |

| JSON (`Authorization: Bearer <token>`) | Answers |
|---|---|
| GET `/api/admin/root/list-backups?prefix=` | `200 {snapshots: [{id, name, data, status}], busy}`, narrowed to the names starting with `prefix` when given |
| POST `/api/admin/root/create-backup` `{name?}` | `202 {status: "creating", snapshot}`; `400` invalid name, `409` name taken or while a job runs |
| POST `/api/admin/root/restore-backup` `{id}` | `202 {status: "restoring"}`; `404`, `400` not ready, `409` |
| GET `/api/admin/root/download-backup/{id}` | the zip; `404`, `400` not ready |
| POST `/api/admin/root/upload-backup` | `201 {snapshot}`; `400` with why, `409` |
| POST `/api/admin/root/remove-backup` `{id}` | `200 {status: "ok"}`; `404`, `409` |
| POST `/api/admin/root/optimize-backup-storage` | `202 {status: "optimizing"}`; `409` |

By hand — `path` is relative to the `--database` folder, as in `backofficedb/backoffice-user/1/values/username`:

| JSON (`Authorization: Bearer <token>`) | Answers |
|---|---|
| POST `/api/admin/root/create-empty-backup` `{name?}` | `201 {snapshot}`, `open`; `400` invalid name, `409` name taken or while a job runs |
| POST `/api/admin/root/add-backup-file/{id}/{path}`, the file as the body (`application/octet-stream`) | `201 {file: {path, sha}}`; `400` invalid path or not `open`, `404`, `409` |
| POST `/api/admin/root/add-backup-blob`, a content as the body (`application/octet-stream`) | `201 {blob: {sha, size}}`; `409` |
| POST `/api/admin/root/add-backup-reference` `{id, path, sha}` | `201 {file: {path, sha}}`; `400` invalid path or sha or not `open`, `404` snapshot or sha not stored (`field: sha`), `409` |
| POST `/api/admin/root/close-backup` `{id}` | `200 {snapshot}`, `ready`; `400` not `open`, holding no file, or naming a content not stored (its path in the message), `404`, `409` |
| GET `/api/admin/root/list-backup-files/{id}?prefix=` | `200 {snapshot, files: [{path, sha}]}`, in path order, narrowed to the paths starting with `prefix` when given; any status; `404` |
| GET `/api/admin/root/download-backup-file/{id}/{path}` | the file, `attachment; filename="<last segment>"`; any status; `404` snapshot, path or content |

Every request with its curl is in [Routes](../Routes/doc.md).

## Rules

- Names: a create may name its snapshot; the name is trimmed and must match
  `^[A-Za-z0-9][A-Za-z0-9._-]{0,99}$` (`400`, or the `invalid-name` notice), and one another snapshot
  holds is refused (`409`, or `name-taken`), never renamed. Without one it is `snapshot-<ts>`; a
  restore's safety copy is `pre-restore-<ts>`, an upload without a usable name `upload-<ts>`, a
  snapshot built by hand `manual-<ts>`, `<ts>` being `YYYYMMDD-hhmmss`.
- Search: `prefix` keeps the snapshots whose name starts with it, case-sensitive, its surrounding
  spaces trimmed; empty keeps every one. Naming snapshots by a scheme (`daily-…`, `before-…`) is
  what makes a prefix find a family of them.
- A snapshot is every file of every folder of the `--database` folder but `backup`; the store's
  temporary (`.…`) and `*.keeplock` files are left out. `{{.ProjectName}} start-server --database var/app`
  snapshots and restores `./var/app`.
- Create, restore and optimize answer at once and run on a goroutine after the route returns. Every
  job — those three, an upload, a delete and every write by hand — runs one at a time: a second one
  is refused (`busy`, `409`). An upload is checked before, so a bad archive is a `400` even then.
- Delete removes the snapshot and its `(path, sha)` list, never a blob: another snapshot may hold
  the same content. Optimize first runs the store's `Repair` over `blob` and `snapshot`, then
  removes every blob no snapshot names — a `failed` one's included — and writes how many it
  removed on stderr.
- Restore: first a `pre-restore-<ts>` snapshot of the current data, then every blob of the target is
  checked, and only then is every folder of the `--database` folder but `backup` replaced.
  Sessions and API tokens are restored too, so the root restoring may be signed out.
- Archive: `snapshot.json` (`{name, data, files}`) + `data/<db>/...`, whatever `--database` names:
  unzipped at the root of a project running with the default folder, it restores by hand.
- Upload: at most 256 MiB, 1 GiB unpacked; every entry under `data/<db>/`, never `data/backup`, an
  empty, `.` or `..` segment, a `\` or a `:`, nor the same path twice. The name comes from
  `snapshot.json` when it is a valid one, else `upload-<ts>`; a name taken gets `-2`, `-3`….
- By hand: `create-empty-backup` records an `open` snapshot; `add-backup-file` stores one file in
  it, or `add-backup-blob` stores a content and `add-backup-reference` names it at a path;
  `close-backup` turns it `ready` once every sha it names is stored. Only an `open` snapshot takes
  files, and only a `ready` one downloads whole or restores; list and download one file read any
  status, so an interrupted backup lists what it already holds and resumes from there. Every
  write is a job (`409` while another runs); list and download one file are not.
- A path added twice replaces the file: the one added last is the one listed and downloaded, and
  the close removes the others. A close refuses a snapshot holding no file — restoring it would
  empty every database — and one naming a sha not stored, which stays `open`.
- A sha is the 64 hex digits of the SHA-256, read in lowercase. A blob no snapshot names is removed
  by the next optimize. A `404` on a reference is how to skip a content already stored: reference
  it first, send it only when it is missing.
- A large snapshot may need `start-server --read-timeout-ms` / `--write-timeout-ms` above 10 s.
- A project whose backoffice was installed before backups: `{{.GeneratorName}} backoffice-init`
  again writes every missing file, never one already there — the nav link, the upload script of
  `backoffice.js` and the `StartRecover` call of `backoffice-start-server` are added by hand.
