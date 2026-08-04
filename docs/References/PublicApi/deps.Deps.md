# `deps.Deps`

**Type:** Struct

## Definition

```go
type Deps struct {
	Now     func() time.Time
	VerbLib verbdeps.Lib
	KeepLib keepdeps.Lib
}
```

## Description

The dependency contract every adapter must fill. Each field is one injectable behavior the financial tracker needs: `Now` is the clock (injecting it fixes the timestamp a category or transaction is stamped with, which is what makes the tracker testable), and `KeepLib` is the schema database every category and transaction is persisted in (an adapter can back it with the filesystem or with any other backend). A filled `Deps` is built by an adapter — see [`standard.New`](/docs/References/PublicApi/standard.New.md) — and passed to [`lib.New`](/docs/References/PublicApi/lib.New.md).

`VerbLib` and `KeepLib` are the exceptions to "every field is a function": the dependency is itself a library built with this pattern, so it arrives as one plain struct field — [`verbdeps.Lib`](/docs/References/PublicApi/verbdeps.Lib.md), [`keepdeps.Lib`](/docs/References/PublicApi/keepdeps.Lib.md) — with no getter around it. The sandbox never imports the embedded Verb or Keep libraries; it declares a copy of each api in `sandbox/contracts/deps/verbdeps/` and `sandbox/contracts/deps/keepdeps/`, and the adapter, which lives outside the sandbox, initializes the real library and assigns its fields onto that copy.

Because it is a struct and not an interface, a value returned by an adapter can be patched field by field before injection, and a custom contract needs no type declaration at all. The trade-off: the compiler cannot detect a field you forgot to fill — it stays nil and panics on first call. See [StructContracts.md](/docs/Explanations/StructContracts.md) and [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Now func() time.Time` | Returns the current time, used to stamp categories and transactions as they are created. |
| `VerbLib verbdeps.Lib` | The embedded Verb argv parser, already initialized by the adapter over the argument vector that adapter chose. |
| `KeepLib keepdeps.Lib` | The embedded Keep schema database every category and transaction is stored in, already wired by the adapter to the storage backend that adapter chose. |

## Examples

```go
package main

import (
	"fmt"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// Start from an adapter — it is what fills KeepLib, the one dependency
	// the tracker cannot work without — then patch the single behavior this
	// program wants to control.
	d := agnosadapter.New("trackerdata")

	frozen := time.Unix(0, 0)
	d.Now = func() time.Time { return frozen }

	l := agnoslib.New(d)

	l.AddCategory("groceries")
	transaction, _ := l.AddSpend("groceries", "coffee beans", 1290)

	fmt.Println(transaction.OccurredAt.Equal(frozen)) // true
}
```
