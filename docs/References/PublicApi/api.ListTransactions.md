# `api.Lib.ListTransactions`

**Type:** Field

## Signature

```go
ListTransactions func() []Transaction
```

## Description

Walks every stored category and collects its transactions into one slice of [`api.Transaction`](/docs/References/PublicApi/api.Transaction.md) values, each with the library's deps propagated in. The records are grouped by category, in the order the categories were created, and each category's own transactions come oldest first.

Every call re-reads the injected database, so the result reflects transactions added or removed since the last call. To list a single category's records, use `Category.ListTransactions` on the value returned by [`GetCategory`](/docs/References/PublicApi/api.GetCategory.md).

## Parameters

_None._

## Returns

| Type | Description |
| :--- | :--- |
| `[]api.Transaction` | Every transaction of every category, grouped by category. |

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

	for _, transaction := range l.ListTransactions() {
		fmt.Println(transaction.String())
	}
}
```
