# `api.Transaction`

**Type:** Struct

## Definition

```go
const (
	Spend = iota
	Received
)

type Transaction struct {
	Deps         deps.Deps
	Id           int64
	Category     string
	Reference    string
	Description  string
	Amount       int64
	Kind         int
	OccurredAt   time.Time
	SignedAmount func() int64
	Remove       func() bool
	String       func() string
}
```

## Description

A single spend or received record handed back by the library, already persisted in the injected database. `Deps` is the dependency set propagated from the parent [`api.Lib`](/docs/PublicApi/api.Lib.md); the plain data fields are read from the stored record when the transaction is built; `SignedAmount`, `Remove` and `String` are function fields filled by factories in `sandbox/internal/transaction/`. A `Transaction` is always constructed by [`AddSpend`](/docs/PublicApi/api.AddSpend.md), [`AddReceived`](/docs/PublicApi/api.AddReceived.md), [`ListTransactions`](/docs/PublicApi/api.ListTransactions.md), or the matching fields of [`api.Category`](/docs/PublicApi/api.Category.md), which propagate the deps in.

`Amount` is always **positive** and expressed in the smallest currency unit (cents); the direction of the money is `Kind`, one of `Spend` or `Received`. Use `SignedAmount()` to get the value with its sign applied, so a list of transactions can be summed directly.

`Reference` is the record's unique key inside its category. The injected database offers unique string keys and integers, so a description — which is not unique — travels inside the reference, composed by the library from a sequence number and the description itself. `Description` is that value read back out.

## Fields

| Field | Description |
| :--- | :--- |
| [`Deps deps.Deps`](/docs/PublicApi/api.Deps.md) | The dependency set propagated from the lib; read-only after construction. |
| `Id int64` | The record's permanent identifier inside its category. |
| `Category string` | The name of the category the transaction belongs to. |
| `Reference string` | The record's unique key inside its category. |
| `Description string` | The human-readable note the transaction was recorded with. |
| `Amount int64` | The value in the smallest currency unit, always positive. |
| `Kind int` | `Spend` or `Received`. |
| `OccurredAt time.Time` | When the transaction was recorded, stamped from the injected clock. |
| `SignedAmount func() int64` | `Amount` negated for a spend, unchanged for a received. |
| `Remove func() bool` | Deletes the transaction from its category. |
| `String func() string` | Renders the transaction as one line. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func main() {
	l := agnoslib.New(agnosadapter.New("trackerdata"))

	l.AddCategory("groceries")
	transaction, _ := l.AddSpend("groceries", "coffee beans", 1290)

	fmt.Println(transaction.Kind == agnostypes.Spend) // true
	fmt.Println(transaction.SignedAmount())           // -1290
	fmt.Println(transaction.Remove())                 // true
}
```
