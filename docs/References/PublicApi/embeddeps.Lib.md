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

The sandbox's **copy** of the api an embedded-asset library exposes, declared in `sandbox/contracts/deps/embeddeps/` and injected whole as the [`deps.Deps.EmbedDeps`](/docs/References/PublicApi/deps.Deps.md) field. It is the same mechanic as [`verbdeps.Lib`](/docs/References/PublicApi/verbdeps.Lib.md) and [`keepdeps.Lib`](/docs/References/PublicApi/keepdeps.Lib.md), for the same reason: reading a file is an OS-bound effect and compiling one into a binary needs the `//go:embed` directive, so neither may appear inside the sandbox. The adapter, which lives outside it, fills the three fields — see [`standard.New`](/docs/References/PublicApi/standard.New.md).

An asset is a payload better kept as a file than as a Go constant — a template, a long-form document, an image — reached by path at runtime and shipped inside the binary. Adding one is [HandleAssets.md](/docs/Tutorials/HandleAssets.md).

The financial tracker **does not read assets**. Its display text is short and fixed, so it lives in [`sandbox/config`](/docs/References/Structure.md#sandboxconfig) as compile-time constants, where the compiler checks every reference; the asset tree ships empty. This contract is filled anyway, as a standing capability of the template — the mechanic is wired end to end and ready for a derived library to drop files into. See [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md).

The contract is **read-only**: assets ship with the program, and nothing in the library ever writes one back. Writing a file at runtime is [`iodeps.Lib`](/docs/References/PublicApi/iodeps.Lib.md), not this. Paths are slash-separated and relative to the root of the asset tree the adapter serves, so `"templates/report.tmpl"` means the same asset whether the adapter serves it out of the binary, out of a directory on disk, or out of a network store. The root itself is addressed as `"."`.

Because nothing in the library calls it, a hand-built `deps.Deps` can leave this field zero — with the usual caveat that a derived library reading an asset afterwards would panic on the nil `ReadFile`.

## Fields

| Field | Description |
| :--- | :--- |
| `ReadFile func(path string) ([]byte, error)` | Returns the whole content of one asset. The error reports an asset that does not exist or could not be read — a packaging mistake, which the interface reports by path rather than printing nothing. |
| `ListFiles func(path string) ([]string, error)` | Returns the names of the assets directly inside a directory, in lexical order, relative to it. Nested directories are neither descended into nor reported. |
| `ListFilesRecursively func(path string) ([]string, error)` | Returns every asset at or below a directory, in lexical order, as slash-separated paths relative to it — `"templates/report.tmpl"`, not `"report.tmpl"`. Directories are never reported, only the files inside them. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	// The adapter compiles the assets into the binary and hands them back as
	// one field of the deps contract, serving the whole asset tree.
	d := agnosadapter.New("trackerdata")

	// Reads assets/templates/report.tmpl. The asset tree ships empty, so this
	// is what a derived library sees after adding that file.
	template, err := d.EmbedDeps.ReadFile("templates/report.tmpl")
	if err != nil {
		fmt.Println("missing asset:", err)
		return
	}
	fmt.Print(string(template))

	// One level: the files sitting directly in the templates directory,
	// relative to it.
	names, _ := d.EmbedDeps.ListFiles("templates")
	fmt.Println(names[0]) // report.tmpl

	// The whole subtree, relative to what was asked for.
	paths, _ := d.EmbedDeps.ListFilesRecursively(".")
	fmt.Println(paths[0]) // templates/report.tmpl
}
```
