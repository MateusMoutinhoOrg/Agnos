# Get Started with the Library

## Description
Installs the module and runs the smallest program that uses it: an adapter builds the `Deps`, `lib.New` turns them into an `api.Lib`, and the tracker records its first transactions. Choosing and patching an adapter is covered by [LibInitialization.md](/docs/Tutorials/LibInitialization.md); every public symbol is listed in [PublicApi.md](/docs/References/PublicApi.md).

### Rules
- Consumers alias every import with the `agnos` prefix — see the Import Aliases rule in [RULES.md](/docs/References/RULES.md#import-aliases).
- Money is an `int64` in the smallest currency unit: `8450` is `84.50`.
- `Deps` is read-only after construction — patch the value before calling `agnoslib.New`, never the returned `api.Lib`.

---

## Workflow

1. Install the module:

   ```bash
   go get github.com/MateusMoutinhoOrg/Agnos-Cli@v0.0.3
   ```

2. Create a `main.go` file:

   ```go
   package main

   import (
       agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
       agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
   )

   func main() {
       // 1. Create deps via an adapter (the "opinionated" part:
       //    real clock + standard output + a schema database on disk)
       deps := agnosadapter.New("trackerdata")

       // 2. Inject deps into the pure library — a financial tracker
       l := agnoslib.New(deps)

       // 3. Use the library — it never knows which adapter is behind the scenes.
       //    Amounts are in the smallest currency unit, so 8450 is 84.50.
       l.AddCategory("groceries")
       l.AddSpend("groceries", "weekly shopping", 8450)

       l.AddCategory("salary")
       l.AddReceived("salary", "august paycheck", 250000)

       println(l.Balance()) // 241550
   }
   ```

3. Run it:

   ```bash
   go run main.go
   ```

4. Continue with [ManageCategories.md](/docs/Tutorials/ManageCategories.md) and [TrackTransactions.md](/docs/Tutorials/TrackTransactions.md) for the full set of operations.

## Full Code

```go
package main

import (
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
    deps := agnosadapter.New("trackerdata")
    l := agnoslib.New(deps)

    l.AddCategory("groceries")
    l.AddSpend("groceries", "weekly shopping", 8450)

    l.AddCategory("salary")
    l.AddReceived("salary", "august paycheck", 250000)

    println(l.Balance()) // 241550
}
```
