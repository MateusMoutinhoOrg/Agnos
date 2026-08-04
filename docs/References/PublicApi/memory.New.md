# `memory.New`

**Type:** Function

## Signature

```go
func New() deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the memory adapter: an in-memory map (guarded by a mutex) for storage and the real wall clock for `Now`. Nothing is persisted — the store vanishes when the process exits, which makes it the fastest choice for ephemeral caches and tests. The factory returns the **contract struct**, never the concrete `MemoryAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the map travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

## Parameters

_None._

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/memory"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	d := agnosadapter.New()
	l := agnoslib.New(d)

	l.Set("greeting", "hello world", 60)
	if entry, ok := l.Get("greeting"); ok {
		fmt.Println(entry.Value) // hello world — gone when the process exits
	}
}
```
