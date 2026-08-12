# `api.Lib.Balance`

**Type:** Field

## Signature

```go
Balance func() int64
```

## Description

Sums the signed amounts of every stored transaction: received money counts up, spent money counts down. The result is expressed in the smallest currency unit (cents), and is negative when the tracked budget has spent more than it received.

It is built on [`ListTransactions`](/docs/References/PublicApi/api.ListTransactions.md) and each transaction's `SignedAmount()`, so it re-reads the injected database on every call and never serves a stale total. For a single category's total, use `Category.Balance` on the value returned by [`GetCategory`](/docs/References/PublicApi/api.GetCategory.md).

## Parameters

_None._

## Returns

| Type | Description |
| :--- | :--- |
| `int64` | The signed total in the smallest currency unit. |

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

	l.AddCategory("salary")
	l.AddReceived("salary", "august paycheck", 250000)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)

	fmt.Println(l.Balance()) // 241550
}
```
