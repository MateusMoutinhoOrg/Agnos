# `api.Entry.ExpiresAt`

**Type:** Field

## Signature

```go
ExpiresAt time.Time
```

## Description

The moment the entry stops being valid, reconstructed from the unix timestamp it was stored with when [`api.Lib.Get`](/docs/References/PublicApi/api.Get.md) builds the `Entry`. Compare it against the injected clock — or call [`IsExpired`](/docs/References/PublicApi/api.IsExpired.md) — to decide whether the value is still live.

## Type

| Type | Description |
| :--- | :--- |
| `time.Time` | The instant after which the entry is expired. |

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
	l.Set("k", "v", 60)

	entry, _ := l.Get("k")
	fmt.Println(entry.ExpiresAt) // ~60s from now
}
```
