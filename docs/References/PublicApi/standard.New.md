# `standard.New`

**Type:** Function

## Signature

```go
func New(filePath string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the standard adapter, which fills the contract using the Go standard library only: a single JSON file configured via `filePath` for storage and the real wall clock for `Now`. Values survive across runs of the process. The factory returns the **contract struct**, never the concrete `StandardAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the adapter's state travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `filePath` | `string` | Path to the JSON file where records should be persisted. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	d := agnosadapter.New("cache.json")
	l := agnoslib.New(d)

	l.Set("greeting", "hello world", 60)
	if entry, ok := l.Get("greeting"); ok {
		fmt.Println(entry.Value) // hello world — still there on the next run
	}
}
```
