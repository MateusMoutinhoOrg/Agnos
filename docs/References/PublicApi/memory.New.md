# `memory.New`

**Type:** Function

## Signature

```go
func New(args []string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the memory adapter: the real wall clock for `Now`, the embedded Verb argv-parser library — initialized over the `args` the caller passes in rather than over the process's command line — for [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), and the embedded Keep schema-database library — wired to Keep's own native adapter, so its records live in memory — for [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md). Nothing is persisted: the tracked categories and transactions do not survive the process, which makes it the fastest choice for tests and for programs that start from an empty tracker every run. The factory returns the **contract struct**, never the concrete `MemoryAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the adapter's state travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `args` | `[]string` | The argument vector the embedded Verb library parses. Pass `os.Args[1:]` to read the real command line, or a fixed slice in a test. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"
	"os"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/memory"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	d := agnosadapter.New(os.Args[1:])
	l := agnoslib.New(d)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)

	fmt.Println(l.Balance()) // -8450 — gone when the process exits
}
```
