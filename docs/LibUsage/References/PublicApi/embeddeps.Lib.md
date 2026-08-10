# `embeddeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	ReadFile             func(path string) ([]byte, error)
	ListFiles            func(path string) ([]string, error)
	ListFilesRecursively func(path string) ([]string, error)
}
```

## Description

The sandbox's **copy** of the api an embedded-asset library exposes, declared in `sandbox/contracts/deps/embeddeps/` and injected whole as the [`deps.Deps.EmbedDeps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) field. It is the same mechanic as [`verbdeps.Lib`](/docs/LibUsage/References/PublicApi/verbdeps.Lib.md) and [`keepdeps.Lib`](/docs/LibUsage/References/PublicApi/keepdeps.Lib.md), for the same reason: reading a file is an OS-bound effect and compiling one into a binary needs the `//go:embed` directive, so neither may appear inside the sandbox. The adapter, which lives outside it, fills the three fields — see [`standard.New`](/docs/LibUsage/References/PublicApi/standard.New.md).

Assets are how the library holds no display text of its own. Every word the command-line interface prints — the usage screen, the version, and each message — is a file under [`/assets/`](/docs/Development/References/Structure.md#assets) that `Sandboxmain` reads by path. Changing the interface's wording, or translating it, is editing those files; no Go changes and no recompilation of the sandbox are involved. The mechanic is explained in [EmbeddedAssets.md](/docs/LibUsage/References/EmbeddedAssets.md).

The contract is **read-only**: assets ship with the program, and nothing in the library ever writes one back. Paths are slash-separated and relative to the asset root the adapter was pointed at (`embedDir` in `standard.New`), so `"version.txt"` means the same asset whether the adapter serves it out of the binary, out of a directory on disk, or out of a network store. The root itself is addressed as `"."`.

Only `Sandboxmain` reads assets. A program that calls the library functions directly never touches one, so a hand-built `deps.Deps` can leave this field zero — with the usual caveat that calling into the interface afterwards would then panic on the nil `ReadFile`.

## Fields

| Field | Description |
| :--- | :--- |
| `ReadFile func(path string) ([]byte, error)` | Returns the whole content of one asset. The error reports an asset that does not exist or could not be read — a packaging mistake, which the interface reports by path rather than printing nothing. |
| `ListFiles func(path string) ([]string, error)` | Returns the names of the assets directly inside a directory, in lexical order, relative to it. Nested directories are neither descended into nor reported. |
| `ListFilesRecursively func(path string) ([]string, error)` | Returns every asset at or below a directory, in lexical order, as slash-separated paths relative to it — `"messages/no-command.txt"`, not `"no-command.txt"`. Directories are never reported, only the files inside them. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	// The adapter compiles the assets into the binary and hands them back as
	// one field of the deps contract, rooted at the directory it was given.
	d := agnosadapter.New("trackerdata", ".")

	version, err := d.EmbedDeps.ReadFile("version.txt")
	if err != nil {
		fmt.Println("missing asset:", err)
		return
	}
	fmt.Print(string(version)) // v0.1.0

	// One level: the files sitting directly in the messages directory.
	lines, _ := d.EmbedDeps.ListFiles("messages")
	fmt.Println(lines[0]) // amount-invalid.txt

	// The whole subtree, relative to what was asked for.
	texts, _ := d.EmbedDeps.ListFilesRecursively(".")
	fmt.Println(len(texts), texts[1]) // 21 messages/amount-invalid.txt
}
```
