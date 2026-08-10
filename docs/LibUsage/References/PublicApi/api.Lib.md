# `api.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	Deps             deps.Deps
	Sandboxmain      func(args []string) int
	AddCategory      func(name string) (Category, bool)
	GetCategory      func(name string) (Category, bool)
	ListCategories   func() []Category
	AddSpend         func(category string, description string, amount int64) (Transaction, bool)
	AddReceived      func(category string, description string, amount int64) (Transaction, bool)
	ListTransactions func() []Transaction
	Balance          func() int64
}
```

## Description

The library entry point, returned by [`lib.New`](/docs/LibUsage/References/PublicApi/lib.New.md). It is a financial tracker: categories hold spend and received transactions, and every record is persisted through the schema database injected as [`deps.Deps.KeepLib`](/docs/LibUsage/References/PublicApi/keepdeps.Lib.md). It is exposed as a struct of function fields: `lib.New` stores the injected [`deps.Deps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) in `Deps`, then runs the factories in `sandbox/internal/lib/`, each of which fills one function field with a closure reading `Deps` at call time. The same deps are propagated into every [`Category`](/docs/LibUsage/References/PublicApi/api.Category.md) and [`Transaction`](/docs/LibUsage/References/PublicApi/api.Transaction.md) the lib creates. Calling a field reads exactly like calling a method — `l.AddCategory("groceries")`. See [StructContracts.md](/docs/Development/References/StructContracts.md).

Amounts are always expressed in the smallest currency unit (cents) as a positive `int64`; the direction of the money lives in the transaction's `Kind`, and `Balance` sums the signed amounts.

`Sandboxmain` is the same idea taken one step further: the project's whole command-line interface is a field of this struct, so the installed binary in [cmd/main](/cmd/main/) holds no behavior of its own. A Go caller that wants the library rather than the CLI simply never calls it.

`Deps` is exported because the library's own factories read it, but it is **read-only after construction**: the closures already captured the struct they were built over, so reassigning `Deps` here does not change behavior. Patch the `deps.Deps` value before calling `lib.New`.

## Fields

| Field | Description |
| :--- | :--- |
| [`Deps deps.Deps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) | The dependency set injected by `lib.New`; read-only after construction. |
| `Sandboxmain func(args []string) int` | The command-line interface: dispatches a command line over the fields below and returns the exit code. |
| [`AddCategory func(name string) (Category, bool)`](/docs/LibUsage/References/PublicApi/api.AddCategory.md) | Creates a category, or returns the stored one when the name is taken. |
| [`GetCategory func(name string) (Category, bool)`](/docs/LibUsage/References/PublicApi/api.GetCategory.md) | Returns the stored category with the given name, or `false` on a miss. |
| [`ListCategories func() []Category`](/docs/LibUsage/References/PublicApi/api.ListCategories.md) | Returns every stored category, oldest first. |
| [`AddSpend func(category string, description string, amount int64) (Transaction, bool)`](/docs/LibUsage/References/PublicApi/api.AddSpend.md) | Records money leaving the budget under an existing category. |
| [`AddReceived func(category string, description string, amount int64) (Transaction, bool)`](/docs/LibUsage/References/PublicApi/api.AddReceived.md) | Records money entering the budget under an existing category. |
| [`ListTransactions func() []Transaction`](/docs/LibUsage/References/PublicApi/api.ListTransactions.md) | Returns every transaction of every category. |
| [`Balance func() int64`](/docs/LibUsage/References/PublicApi/api.Balance.md) | Sums the signed amounts of every stored transaction. |
