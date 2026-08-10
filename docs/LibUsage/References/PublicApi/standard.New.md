# `standard.New`

**Type:** Function

## Signature

```go
func New(basePath string, embedDir string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) backed by the standard adapter: the real wall clock for `Now`, the process's standard output for `Printf`, the embedded Verb argv-parser library — initialized over the process's own command line, `os.Args[1:]` — for [`VerbLib`](/docs/LibUsage/References/PublicApi/verbdeps.Lib.md), the embedded Keep schema-database library — wired to Keep's own filesystem adapter, which writes one file per key under `basePath` — for [`KeepLib`](/docs/LibUsage/References/PublicApi/keepdeps.Lib.md), and the project's assets, compiled into the binary and rooted at `embedDir`, for [`EmbedDeps`](/docs/LibUsage/References/PublicApi/embeddeps.Lib.md). Categories and transactions therefore survive across runs of the process, and the interface's text travels inside the binary with no files beside it. Handing `api.Lib.Sandboxmain` the same `os.Args[1:]` this adapter wired the parser over is what keeps the interface and the parser in agreement about the command line. Initializing Verb and Keep here, outside the sandbox, is what keeps the sandbox free of the embedded libraries. The factory returns the **contract struct**, never the concrete `StandardAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the adapter's state travels with the deps. For all shipped adapters, see [Adapters.md](/docs/LibUsage/References/Adapters.md).

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `basePath` | `string` | Directory the embedded Keep database writes its records under. |
| `embedDir` | `string` | Directory inside the compiled-in assets every asset path is resolved against. Pass `"."` for the whole asset tree — the usage screen, the version, and every message of the interface. Pass a subdirectory to root the library at part of the tree. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	d := agnosadapter.New("trackerdata", ".")
	l := agnoslib.New(d)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)

	// Still there on the next run — the records are on disk.
	fmt.Println(l.Balance())
}
```
