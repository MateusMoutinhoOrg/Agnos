# `api.Entry.IsExpired`

**Type:** Field

## Signature

```go
IsExpired func() bool
```

## Description

Reports whether the injected clock (`e.Deps.Now()`) has passed the entry's [`ExpiresAt`](/docs/References/PublicApi/api.ExpiresAt.md). Unlike `Value` and `ExpiresAt`, this is a function field, not data: the factory that filled it assigned a closure over the entry, so it re-reads the clock on every call rather than freezing a verdict at construction time. Because the clock comes from the deps, expiry can be tested deterministically by supplying a custom `Deps` to `lib.New`.

Note that [`api.Lib.Get`](/docs/References/PublicApi/api.Get.md) already filters out expired entries, so an entry obtained from `Get` reports `false` — this field is what makes that filtering possible and is available for entries held and checked over time.

## Parameters

_None._

## Returns

| Type | Description |
| :--- | :--- |
| `bool` | `true` when the entry has expired, `false` while it is still valid. |

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
	fmt.Println(entry.IsExpired()) // false
}
```
