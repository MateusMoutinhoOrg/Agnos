# SmartIO

## Description
Explains `sandbox/internal/smartio` — the transactional filesystem every action reads and writes through, rooted at the project directory named by `--path`. It is what lets an action compose several changes and flush them once, and what makes `--path` an absolute boundary. The actions using it are described in [BuildPipeline](/docs/BuildPipeline/doc.md); writing an action against it is [AddAction](/docs/AddAction/doc.md).

---

## One Root, Project-Relative Paths

`smartio.New(deps, path, projectName)` returns a `*SmartIO` whose `Root` is the `--path` value, normalized so `""`, `"."` and `"./"` all mean "no prefix". Every path an action hands it afterwards is **project-relative** — `"go.mod"`, `"sandbox/api"`, `"sandbox/internal/commands/greet/entries.yaml"` — and never `path + "/..."`:

```go
io := smartio.New(deps, "./my-tool", config.ProjectName)
content, err := io.ReadFile("go.mod")        // reads ./my-tool/go.mod
io.CreateDir("sandbox/deps")                 // will create ./my-tool/sandbox/deps
```

`Root` is joined on by `rootedPath` only at the boundary calls into `deps.Iodeps` — reads, existence checks, listings, and the writes inside `Persist` — and stripped back off listing results by `unrootedPath`. `rootedPath` is idempotent, so a path already under `Root` is left alone. Nothing an action does can address a file outside the project.

`New` also loads `AgnosConfig/ignore.yaml` and `AgnosConfig/paths.yaml` from the project: the first hides matching paths from every listing, the second rewrites them.

---

## The Transaction

Writes never hit the disk when they are made. `SmartIO` keeps three pending sets:

| Call | Pending effect |
|------|----------------|
| `WriteFile(path, content)` | Buffers the content in `Transactions[path]`. **Refuses** to overwrite a file that exists on disk or in the transaction. |
| `WriteFileOverwrite(path, content)` | Buffers the content, replacing whatever is there. Every generated file is written this way. |
| `CreateDir(path)` | Appends to `PendingCreateDirs`. |
| `RemoveDir(path)` | Appends to `PendingRemoveDirs` — a file or a directory, with its children. |

`Persist()` flushes them in one fixed order: pending removals, then pending directory creations, then buffered file writes. A removal and a re-creation of the same path in one transaction therefore ends with the path present.

---

## Transaction-Aware Reads

Every read answers as if the transaction had already been applied, which is what lets a follow-up step inside the same action see an earlier step's work before anything is persisted:

| Call | Sees |
|------|------|
| `ReadFile` | Buffered content first, then the disk. |
| `Exist`, `IsFile`, `IsDir` | Pending creations as existing, pending removals as gone. |
| `ListDirs`, `ListFiles`, `ListAll` and their `Recursively` variants | The disk, with `ignore.yaml` applied and `paths.yaml` rewrites made. |

The listings are the one place the transaction is **not** consulted — they come from `deps.Iodeps`. This is why every action that runs `build` as a follow-up persists first: `build`'s collectors list `sandbox/deps/`, `adapters/libs/` and `sandbox/internal/commands/` from disk and would not see files still pending. `start` is the one exception — it renders its group and calls `BuildInternal` inside the same transaction, which works because `BuildInternal` reads `project.yaml` through `ReadFile`, not through a listing.

---

## Why Actions Compose

Because an open `*SmartIO` is a value, one action can call another's `*Internal` function and both write into the same transaction; the outer action persists once. `start` calls `BuildInternal` this way, and every editor action loads `entries.yaml`, mutates the parsed struct, saves it, persists, and then runs `build` as a fresh transaction. The two-file shape every action has — a public entry that owns the transaction and an internal function that does not — exists to make this composition possible; see the Action specification in [Specs](/docs/Specs/doc.md).
