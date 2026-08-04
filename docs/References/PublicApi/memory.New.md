# `memory.New`

**Type:** Function

## Signature

```go
func New(args []string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the memory adapter: an in-memory map (guarded by a mutex) for storage, the real wall clock for `Now`, and the embedded Verb argv-parser library — initialized over the `args` the caller passes in rather than over the process's command line — for [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md). Nothing is persisted — the store vanishes when the process exits, which makes it the fastest choice for ephemeral caches and tests. The factory returns the **contract struct**, never the concrete `MemoryAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the map travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

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

	l.Set("greeting", "hello world", 60)
	if entry, ok := l.Get("greeting"); ok {
		fmt.Println(entry.Value) // hello world — gone when the process exits
	}
}
```
