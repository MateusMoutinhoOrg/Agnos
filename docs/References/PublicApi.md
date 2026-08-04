# Public API

## Description
Index of every public-facing entry of the library. Callers hold **structs of function fields** declared in `sandbox/contracts/api` and `sandbox/contracts/deps`; the **factories** that fill those fields live in `sandbox/internal/` and are unreachable from outside `sandbox/`. See [StructContracts.md](/docs/Explanations/StructContracts.md).

---

## Structs

### [api.Lib](/docs/References/PublicApi/api.Lib.md)
The library entry point — a financial tracker persisting categories and transactions in the injected database. Returned by `lib.New`; exposes all library functions, and the command-line interface itself, as fields.

### [api.Category](/docs/References/PublicApi/api.Category.md)
One bucket transactions are tracked under, with its dependencies already wired into every field it exposes.

### [api.Transaction](/docs/References/PublicApi/api.Transaction.md)
A single spend or received record handed back by the library, with its dependencies already wired in.

### [deps.Deps](/docs/References/PublicApi/deps.Deps.md)
The dependency contract every adapter must fill: the clock, the writer the interface reports through, the embedded Verb argv parser, and the embedded Keep schema database.

### [verbdeps.Lib](/docs/References/PublicApi/verbdeps.Lib.md)
The sandbox's copy of the embedded Verb argv-parser library's api, injected whole as the `deps.Deps.VerbLib` field.

### [keepdeps.Lib](/docs/References/PublicApi/keepdeps.Lib.md)
The sandbox's copy of the embedded Keep schema-database library's api, injected whole as the `deps.Deps.KeepLib` field.

---

## Functions

### [lib.New](/docs/References/PublicApi/lib.New.md)
Injects a `deps.Deps` into the library and returns an `api.Lib`.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Creates a `deps.Deps` using the standard library adapter (real clock + a Keep database on disk).

---

## Fields

The fields of [`api.Category`](/docs/References/PublicApi/api.Category.md) and [`api.Transaction`](/docs/References/PublicApi/api.Transaction.md) are documented on those structs' own pages; the entries below are the library functions `api.Lib` exposes.

### [api.Lib.Deps / api.Category.Deps / api.Transaction.Deps](/docs/References/PublicApi/api.Deps.md)
The injected dependency set the struct was built with; read-only after construction.

### [api.Lib.Sandboxmain](/docs/References/PublicApi/api.Sandboxmain.md)
Runs the whole command-line interface over an argument vector and returns the process exit code.

### [api.Lib.AddCategory](/docs/References/PublicApi/api.AddCategory.md)
Creates a category, or returns the stored one when the name is already taken.

### [api.Lib.GetCategory](/docs/References/PublicApi/api.GetCategory.md)
Returns the stored category with the given name, or `false` on a miss.

### [api.Lib.ListCategories](/docs/References/PublicApi/api.ListCategories.md)
Returns every stored category, oldest first.

### [api.Lib.AddSpend](/docs/References/PublicApi/api.AddSpend.md)
Records money leaving the budget under an existing category.

### [api.Lib.AddReceived](/docs/References/PublicApi/api.AddReceived.md)
Records money entering the budget under an existing category.

### [api.Lib.ListTransactions](/docs/References/PublicApi/api.ListTransactions.md)
Returns every transaction of every category.

### [api.Lib.Balance](/docs/References/PublicApi/api.Balance.md)
Sums the signed amounts of every stored transaction.
