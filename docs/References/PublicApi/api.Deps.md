# `api.Lib.Deps` / `api.Entry.Deps`

**Type:** Field

## Signature

```go
Deps deps.Deps
```

## Description

The injected dependency set the struct was built with. [`lib.New`](/docs/References/PublicApi/lib.New.md) writes it onto [`api.Lib`](/docs/References/PublicApi/api.Lib.md), and [`api.Lib.Get`](/docs/References/PublicApi/api.Get.md) propagates the same value onto every [`api.Entry`](/docs/References/PublicApi/api.Entry.md) it creates. Every other function field on those structs is a closure that reads this field at call time — it is how a dependency injected once reaches the whole object graph. See [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

The field is exported because the library's own factories, which live in another package, must read it. It is **not** a customization point: the closures captured the struct the factories ran over, so assigning to `Deps` on a struct you already hold changes nothing. To replace a behavior, patch the [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) value **before** passing it to `lib.New`.

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | The filled dependency contract; read-only after construction. |

## Examples

### Patch Before Injection, Not After

```go
package main

import (
	"fmt"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/memory"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	myDeps := agnosadapter.New()

	// Correct: replace the clock while it is still a plain deps value.
	frozen := time.Unix(0, 0)
	myDeps.Now = func() time.Time { return frozen }

	l := agnoslib.New(myDeps)

	// Wrong: the factories already captured l — this assignment is inert.
	l.Deps.Now = time.Now

	l.Set("k", "v", 60)
	frozen = time.Unix(120, 0) // past the TTL

	_, ok := l.Get("k")
	fmt.Println(ok) // false — the frozen clock is still the one in effect
}
```
