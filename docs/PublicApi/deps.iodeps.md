# `deps.Iodeps`

`sandbox/deps/iodeps`

## `Sandbox`

Sandbox is the filesystem library injected whole as the Deps.IoLib field. Paths are whatever the host operating system accepts, resolved by the adapter — unlike embeddeps.Sandbox, which is always slash-separated and rooted at an asset tree. The listing functions report paths that already include the directory they were given, so a result can be passed straight back in. The predicates report false rather than an error: a path that cannot be stat'd is not a directory and is not a file, which is the answer the caller wanted either way.

| Field | Type | Description |
| --- | --- | --- |
| `ReadFile` | `func(path string) ([]byte, error)` | ReadFile returns the whole content of the file at path. The error reports a file that does not exist or could not be read. |
| `WriteFile` | `func(path string, content []byte) error` | WriteFile writes content to path, creating any missing parent directory first and truncating an existing file. The error reports a directory or a file that could not be written. |
| `IsDir` | `func(path string) bool` | IsDir reports whether path exists and is a directory. |
| `IsFile` | `func(path string) bool` | IsFile reports whether path exists and is not a directory. |
| `Exist` | `func(path string) bool` | Exist reports whether anything exists at path, directory or file. |
| `CreateDir` | `func(path string)` | CreateDir creates the directory at path together with any missing parent. It reports nothing: a directory that already exists and a directory just created are the same outcome to the caller. |
| `RemoveDir` | `func(path string)` | RemoveDir removes the directory or file at path and any children it contains. It reports nothing: a missing path and a path just removed are the same outcome. |
| `ListDirs` | `func(path string) []string` | ListDirs returns the directories directly inside path. Nested directories are not descended into. |
| `ListFiles` | `func(path string) []string` | ListFiles returns the files directly inside path. Directories are not reported. |
| `ListAll` | `func(path string) []string` | ListAll returns every entry directly inside path, directories and files alike. |
| `ListDirsRecursively` | `func(path string) []string` | ListDirsRecursively returns every directory at or below path, excluding path itself. |
| `ListFilesRecursively` | `func(path string) []string` | ListFilesRecursively returns every file at or below path, at any depth. Directories are never reported. |
| `ListAllRecursively` | `func(path string) []string` | ListAllRecursively returns every entry at or below path, directories and files alike, excluding path itself. |
| `Join` | `func(elements ...string) string` | Join joins the given path elements with the separator the host operating system uses, cleaning the result. It is the only way the sandbox can build a host path, which is always separator-dependent — unlike the project-relative paths it otherwise passes, which are always slash-separated. |
| `Dir` | `func(path string) string` | Dir returns path without its last element, the directory holding it. |
| `UserHomeDir` | `func() (string, error)` | UserHomeDir returns the home directory of the user running the process. The error reports a home directory that could not be determined. |

[every contract](doc.md)
