# `api.Entry.Value`

**Type:** Field

## Signature

```go
Value string
```

## Description

The cached value the entry was created with. A plain data field: it is written by the entry constructor when [`api.Lib.Get`](/docs/References/PublicApi/api.Get.md) builds the `Entry`, before any factory runs, and reading it touches no dependency.

## Type

| Type | Description |
| :--- | :--- |
| `string` | The stored value. |

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
	fmt.Println(entry.Value) // v
}
```
