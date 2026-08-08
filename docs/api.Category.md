# `api.Category`

**Type:** Struct

## Definition

```go
type Category struct {
	Deps             deps.Deps
	Id               int64
	Name             string
	CreatedAt        time.Time
	AddSpend         func(description string, amount int64) (Transaction, bool)
	AddReceived      func(description string, amount int64) (Transaction, bool)
	ListTransactions func() []Transaction
	Balance          func() int64
	Remove           func() bool
	String           func() string
}
```

## Description

One bucket transactions are tracked under — "groceries", "salary" — already persisted in the injected database. `Deps` is the dependency set propagated from the parent [`api.Lib`](/docs/api.Lib.md); `Id`, `Name` and `CreatedAt` are plain data read from the stored record; the function fields are filled by factories in `sandbox/internal/category/`, whose closures re-read the stored record through `Deps` on every call. A `Category` is always constructed by [`AddCategory`](/docs/api.AddCategory.md), [`GetCategory`](/docs/api.GetCategory.md), or [`ListCategories`](/docs/api.ListCategories.md), which propagate the deps in. The factories behind it live inside the closed sandbox and are not reachable by callers.

Because every function field re-reads the database, a `Category` value never goes stale: transactions added through another value of the same category show up in `ListTransactions` and `Balance`.

## Fields

| Field | Description |
| :--- | :--- |
| [`Deps deps.Deps`](/docs/api.Deps.md) | The dependency set propagated from the lib; read-only after construction. |
| `Id int64` | The record's permanent identifier. |
| `Name string` | The category's unique name. |
| `CreatedAt time.Time` | When the category was created, stamped from the injected clock. |
| `AddSpend func(description string, amount int64) (Transaction, bool)` | Records money leaving the budget under this category; `false` when `amount` is not positive or the write failed. |
| `AddReceived func(description string, amount int64) (Transaction, bool)` | Records money entering the budget, with the same rules as `AddSpend`. |
| `ListTransactions func() []Transaction` | Every transaction stored under this category, oldest first. |
| `Balance func() int64` | The sum of this category's signed amounts. |
| `Remove func() bool` | Deletes the category and every transaction under it. |
| `String func() string` | Renders the category as one line: name, balance, transaction count. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	l := agnoslib.New(agnosadapter.New("trackerdata"))

	groceries, _ := l.AddCategory("groceries")
	groceries.AddSpend("weekly shopping", 8450) // 84.50

	fmt.Println(groceries.String()) // groceries  -84.50  1 transactions
	fmt.Println(groceries.Balance())
}
```
