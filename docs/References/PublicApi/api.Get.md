# `api.Lib.Get`

**Type:** Field

## Signature

```go
Get func(key string) (Entry, bool)
```

## Description

Loads the record stored under `key` via `l.Deps.Load`, builds an [`api.Entry`](/docs/References/PublicApi/api.Entry.md) with the library's deps propagated into its `IsExpired`, and returns it. Returns the **zero `Entry`** and `false` when the key is absent **or** when the entry has already expired according to the injected clock — always branch on the `bool`, since a struct return has no nil to compare against.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `key` | `string` | The key to look up. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Entry`](/docs/References/PublicApi/api.Entry.md) | The live entry, or the zero value on a miss or expiry. |
| `bool` | `true` when a live entry was found, `false` otherwise. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	l := agnoslib.New(agnosadapter.New("cache.json"))
	l.Set("greeting", "hello", 60)

	if entry, ok := l.Get("greeting"); ok {
		fmt.Println(entry.Value) // hello
	}
}
```
