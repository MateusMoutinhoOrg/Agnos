# `iodeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	ReadFile  func(path string) ([]byte, error)
	WriteFile func(path string, content []byte) error

	IsDir  func(path string) bool
	IsFile func(path string) bool
	Exist  func(path string) bool

	CreateDir func(path string)
	RemoveDir func(path string)

	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}
```

## Description

The filesystem library injected whole as `Deps.Iodeps` — the sandbox's copy of what `os` and `path/filepath` provide. Paths are whatever the host accepts, resolved by the adapter; the listing functions report paths that already include the directory they were given, so a result can be passed straight back in. The predicates report `false` rather than an error: a path that cannot be stat'd is not a directory and is not a file, which is the answer the caller wanted either way. `CreateDir` and `RemoveDir` report nothing: already-there and just-created, missing and just-removed, are the same outcome.

Inside Agnos, no action calls it directly: every action goes through [SmartIO](/docs/References/SmartIO.md), which wraps it with a transaction and a root and calls it only at the boundary. Installed by the `iodeps` dep.

## Fields

| Field | Description |
| :--- | :--- |
| `ReadFile(path)` | The whole content; errors when missing or unreadable. |
| `WriteFile(path, content)` | Creates missing parents, truncates an existing file. |
| `IsDir`, `IsFile`, `Exist` | Predicates; never an error. |
| `CreateDir(path)` | `MkdirAll`. |
| `RemoveDir(path)` | Removes a directory with its children, or a file. |
| `ListDirs`, `ListFiles`, `ListAll` | Direct children of `path`, by kind. |
| `List*Recursively` | Every entry at or below `path`, excluding `path` itself. |

## Examples

```go
// A handler in a generated project writing a report through the injected filesystem.
if err := deps.Iodeps.WriteFile(entries.Output, report); err != nil {
	deps.Std.Error("could not write %s: %v\n", entries.Output, err)
	return api.ExitFailure
}
for _, file := range deps.Iodeps.ListFilesRecursively("sandbox/internal/commands") {
	deps.Std.Log("%s\n", file)
}
```
