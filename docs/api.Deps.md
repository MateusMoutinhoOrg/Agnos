# `api.Lib.Deps` / `api.Category.Deps` / `api.Transaction.Deps`

**Type:** Field

## Signature

```go
Deps deps.Deps
```

## Description

The injected dependency set the struct was built with. [`lib.New`](/docs/lib.New.md) writes it onto [`api.Lib`](/docs/api.Lib.md), which propagates the same value onto every [`api.Category`](/docs/api.Category.md) it creates, and each category propagates it again onto every [`api.Transaction`](/docs/api.Transaction.md) it hands back. Every other function field on those structs is a closure that reads this field at call time — it is how a dependency injected once reaches the whole object graph. See [DepsMechanic.md](/docs/DepsMechanic.md).

The field is exported because the library's own factories, which live in another package, must read it. It is **not** a customization point: the closures captured the struct the factories ran over, so assigning to `Deps` on a struct you already hold changes nothing. To replace a behavior, patch the [`deps.Deps`](/docs/deps.Deps.md) value **before** passing it to `lib.New`.

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/deps.Deps.md) | The filled dependency contract; read-only after construction. |

## Examples

### Patch Before Injection, Not After

```go
package main

import (
	"fmt"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	myDeps := agnosadapter.New("trackerdata")

	// Correct: replace the clock while it is still a plain deps value.
	frozen := time.Unix(0, 0)
	myDeps.Now = func() time.Time { return frozen }

	l := agnoslib.New(myDeps)

	// Wrong: the factories already captured l — this assignment is inert.
	l.Deps.Now = time.Now

	l.AddCategory("groceries")
	transaction, _ := l.AddSpend("groceries", "coffee beans", 1290)

	// The frozen clock is still the one in effect.
	fmt.Println(transaction.OccurredAt.Equal(frozen)) // true
}
```
