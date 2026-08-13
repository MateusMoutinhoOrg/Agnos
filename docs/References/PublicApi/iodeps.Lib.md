# `iodeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	ReadFile             func(path string) ([]byte, error)
	WriteFile            func(path string, content []byte) error
	IsDir                func(path string) bool
	IsFile               func(path string) bool
	Exist                func(path string) bool
	CreateDir            func(path string)
	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}
```

## Description

The sandbox's **copy** of the api a filesystem library exposes, declared in `sandbox/contracts/deps/iodeps/` and injected whole as the [`deps.Deps.IoLib`](/docs/References/PublicApi/deps.Deps.md) field. It is the same mechanic as [`verbdeps.Lib`](/docs/References/PublicApi/verbdeps.Lib.md), [`keepdeps.Lib`](/docs/References/PublicApi/keepdeps.Lib.md) and [`embeddeps.Lib`](/docs/References/PublicApi/embeddeps.Lib.md), for the same reason: touching a file is an OS-bound effect, so `os` and `path/filepath` may not appear inside the sandbox. The adapter, which lives outside it, fills every field — see [`standard.New`](/docs/References/PublicApi/standard.New.md).

The financial tracker **does not call it**: every record it keeps is persisted through [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md). It is carried as a standing capability of the template, so a library derived from this repository that must touch the filesystem directly finds the contract already declared and already wired. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md) for why the effect cannot simply be imported.

Paths are whatever the host operating system accepts, resolved by the adapter — unlike `embeddeps.Lib`, which is always slash-separated and rooted at an asset tree. The listing functions return paths that already include the directory they were given, so a result can be passed straight back in.

The predicates report `false` rather than an error: a path that cannot be inspected is neither a directory nor a file, which is the answer the caller wanted either way. Reading and writing report an `error`, because a caller that asked for bytes needs to know why it got none.

## Fields

| Field | Description |
| :--- | :--- |
| `ReadFile func(path string) ([]byte, error)` | Returns the whole content of the file at `path`. The error reports a file that does not exist or could not be read. |
| `WriteFile func(path string, content []byte) error` | Writes `content` to `path`, creating any missing parent directory first and truncating an existing file. |
| `IsDir func(path string) bool` | Reports whether `path` exists and is a directory. |
| `IsFile func(path string) bool` | Reports whether `path` exists and is not a directory. |
| `Exist func(path string) bool` | Reports whether anything exists at `path`, directory or file. |
| `CreateDir func(path string)` | Creates the directory at `path` together with any missing parent. It reports nothing: an existing directory and a newly created one are the same outcome. |
| `ListDirs func(path string) []string` | Returns the directories directly inside `path`, not descending into them. |
| `ListFiles func(path string) []string` | Returns the files directly inside `path`. Directories are not reported. |
| `ListAll func(path string) []string` | Returns every entry directly inside `path`, directories and files alike. |
| `ListDirsRecursively func(path string) []string` | Returns every directory at or below `path`, excluding `path` itself. |
| `ListFilesRecursively func(path string) []string` | Returns every file at or below `path`, at any depth. Directories are never reported. |
| `ListAllRecursively func(path string) []string` | Returns every entry at or below `path`, directories and files alike, excluding `path` itself. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	// The adapter fills IoLib with os/filepath calls and hands it back as one
	// field of the deps contract.
	d := agnosadapter.New("trackerdata")

	// WriteFile creates the missing parent directory on its own.
	if err := d.IoLib.WriteFile("scratch/notes/august.txt", []byte("groceries\n")); err != nil {
		fmt.Println("could not write:", err)
		return
	}

	fmt.Println(d.IoLib.IsFile("scratch/notes/august.txt")) // true
	fmt.Println(d.IoLib.IsDir("scratch/notes"))             // true

	// Listing returns paths that already carry the directory asked for.
	fmt.Println(d.IoLib.ListFilesRecursively("scratch")) // [scratch/notes/august.txt]

	content, err := d.IoLib.ReadFile("scratch/notes/august.txt")
	if err != nil {
		fmt.Println("could not read:", err)
		return
	}
	fmt.Print(string(content)) // groceries
}
```
