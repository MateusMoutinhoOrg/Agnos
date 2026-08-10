# `api.Lib.AddReceived`

**Type:** Field

## Signature

```go
AddReceived func(category string, description string, amount int64) (Transaction, bool)
```

## Description

Records money **entering** the tracked budget under an existing category and returns the persisted record as an [`api.Transaction`](/docs/LibUsage/References/PublicApi/api.Transaction.md) of kind `Received`, stamped with the injected clock (`l.Deps.Now()`).

It mirrors [`AddSpend`](/docs/LibUsage/References/PublicApi/api.AddSpend.md) in every respect but the direction of the money: the category must already exist, and `amount` is expressed in the smallest currency unit (cents) and must be **positive**. What differs is how [`Balance`](/docs/LibUsage/References/PublicApi/api.Balance.md) counts it — a received transaction's `SignedAmount()` is positive.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `category` | `string` | The name of an existing category. |
| `description` | `string` | The human-readable note to record the transaction with. |
| `amount` | `int64` | The value in the smallest currency unit; must be positive. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Transaction`](/docs/LibUsage/References/PublicApi/api.Transaction.md) | The persisted transaction, or the zero value on failure. |
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
	l := agnoslib.New(agnosadapter.New("trackerdata", "."))
	l.AddCategory("salary")

	transaction, _ := l.AddReceived("salary", "august paycheck", 250000) // 2500.00
	fmt.Println(transaction.SignedAmount())                              // 250000
}
```
