# `standard.New`

**Type:** Function

## Signature

```go
func New(basePath string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the standard adapter: the real wall clock for `Now`, the embedded Verb argv-parser library — initialized over the process's own command line, `os.Args[1:]` — for [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), and the embedded Keep schema-database library — wired to Keep's own filesystem adapter, which writes one file per key under `basePath` — for [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md). Categories and transactions therefore survive across runs of the process. Initializing Verb and Keep here, outside the sandbox, is what keeps the sandbox free of the embedded libraries. The factory returns the **contract struct**, never the concrete `StandardAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the adapter's state travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `basePath` | `string` | Directory the embedded Keep database writes its records under. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	d := agnosadapter.New("trackerdata")
	l := agnoslib.New(d)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)

	// Still there on the next run — the records are on disk.
	fmt.Println(l.Balance())
}
```
