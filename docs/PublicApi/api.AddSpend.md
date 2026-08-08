# `api.Lib.AddSpend`

**Type:** Field

## Signature

```go
AddSpend func(category string, description string, amount int64) (Transaction, bool)
```

## Description

Records money **leaving** the tracked budget under an existing category and returns the persisted record as an [`api.Transaction`](/docs/PublicApi/api.Transaction.md) of kind `Spend`, stamped with the injected clock (`l.Deps.Now()`).

It is the lib-level shorthand for `l.GetCategory(category)` followed by `Category.AddSpend(...)`: the category must already exist, so declare it with [`AddCategory`](/docs/PublicApi/api.AddCategory.md) first. `amount` is expressed in the smallest currency unit (cents) and must be **positive** — the direction of the money is the transaction's `Kind`, not the sign of its amount.

Descriptions need not be unique: the library composes the record's unique key from a sequence number and the description, so two transactions can carry the same note.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `category` | `string` | The name of an existing category. |
| `description` | `string` | The human-readable note to record the transaction with. |
| `amount` | `int64` | The value in the smallest currency unit; must be positive. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Transaction`](/docs/PublicApi/api.Transaction.md) | The persisted transaction, or the zero value on failure. |
| `bool` | `false` when the category is unknown, `amount` is not positive, or the record could not be written. |

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
	l.AddCategory("groceries")

	transaction, ok := l.AddSpend("groceries", "weekly shopping", 8450) // 84.50
	if !ok {
		fmt.Println("nothing recorded")
		return
	}
	fmt.Println(transaction.String())
}
```
