# Public API

## Description
Index of every public-facing entry of the library, logically grouped by their role in the system. Callers hold **structs of function fields** declared in `sandbox/contracts/api` and `sandbox/contracts/deps`; the **factories** that fill those fields live in `sandbox/` and are unreachable from outside `sandbox/`. See [StructContracts.md](/docs/References/StructContracts.md).

---

## Entry Points

### [lib.New](/docs/References/PublicApi/lib.New.md)
Injects a `deps.Deps` into the library and returns an `api.Lib`.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Creates a `deps.Deps` using the standard library adapter (real clock, a Keep database on disk, and the assets compiled into the binary).

---

## Core Interface

### [api.Lib](/docs/References/PublicApi/api.Lib.md)
The library entry point — a financial tracker persisting categories and transactions in the injected database. Returned by `lib.New`; exposes all library functions, and the command-line interface itself, as fields.

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

---

## Data Models

### [api.Category](/docs/References/PublicApi/api.Category.md)
One bucket transactions are tracked under, with its dependencies already wired into every field it exposes.

### [api.Transaction](/docs/References/PublicApi/api.Transaction.md)
A single spend or received record handed back by the library, with its dependencies already wired in.

### [api.Lib.Deps / api.Category.Deps / api.Transaction.Deps](/docs/References/PublicApi/api.Deps.md)
The injected dependency set the struct was built with; read-only after construction.

---

## Dependency Contracts

### [deps.Deps](/docs/References/PublicApi/deps.Deps.md)
The dependency contract every adapter must fill: the clock, the writer the interface reports through, the embedded Verb argv parser, the embedded Keep schema database, the embedded assets, the filesystem, and the HTTP client.

### [verbdeps.Lib](/docs/References/PublicApi/verbdeps.Lib.md)
The sandbox's copy of the embedded Verb argv-parser library's api, injected whole as the `deps.Deps.VerbLib` field.

### [keepdeps.Lib](/docs/References/PublicApi/keepdeps.Lib.md)
The sandbox's copy of the embedded Keep schema-database library's api, injected whole as the `deps.Deps.KeepLib` field.

---

## Standing Capabilities

The three contracts below are declared and filled like any other dependency, but the financial tracker never calls them. They ship as capabilities a library derived from this template gets already wired — see [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md).

### [embeddeps.Lib](/docs/References/PublicApi/embeddeps.Lib.md)
The sandbox's copy of an embedded-asset library's api, injected whole as the `deps.Deps.EmbedDeps` field: reads the files shipped under `/assets/`.

### [iodeps.Lib](/docs/References/PublicApi/iodeps.Lib.md)
The sandbox's copy of a filesystem library's api, injected whole as the `deps.Deps.IoLib` field: reads, writes, and lists paths on disk.

### [requestdeps.Request](/docs/References/PublicApi/requestdeps.Request.md)
One HTTP request under construction, handed back by the `deps.Deps.NewRequest` field already bound to a url.

### [requestdeps.Response](/docs/References/PublicApi/requestdeps.Response.md)
One HTTP response, handed back by `requestdeps.Request.Fetch` with its body still open.
