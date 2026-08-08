# Track Transactions

## Description
Covers recording spend and received transactions, listing them, and reading a balance. Creating the categories they are recorded under is covered by [ManageCategories.md](/docs/ManageCategories.md); installing and initializing the lib is covered by [LibInitialization.md](/docs/LibInitialization.md).

### Rules
- Amounts are expressed in the **smallest currency unit** (cents) and must be **positive**: `8450` is `84.50`. The direction of the money is the transaction's `Kind`, never the sign of its amount.
- A transaction is always recorded under an existing category, so create it first — see [ManageCategories.md](/docs/ManageCategories.md).

---

## Workflow
1. Initialize the lib and create the category to record under:
   ```go
   l := agnoslib.New(agnosadapter.New("trackerdata"))
   groceries, _ := l.AddCategory("groceries")
   ```
2. Record money leaving the budget with `AddSpend`. It returns the persisted `api.Transaction` and a boolean that is `false` when the amount is not positive or the record could not be written:
   ```go
   transaction, ok := groceries.AddSpend("weekly shopping", 8450) // 84.50
   if !ok {
       println("nothing recorded")
       return
   }
   println(transaction.String())
   ```
3. Record money entering the budget with `AddReceived`, which follows the same rules:
   ```go
   l.AddCategory("salary")
   salary, _ := l.GetCategory("salary")
   salary.AddReceived("august paycheck", 250000) // 2500.00
   ```
4. When the category is already stored, record straight from the lib instead — `l.AddSpend` and `l.AddReceived` take the category name as their first argument and report `false` when it is unknown:
   ```go
   if _, ok := l.AddSpend("holidays", "flight", 42000); !ok {
       println("holidays: no such category, nothing recorded")
   }
   ```
5. List the records. `Category.ListTransactions` returns one category's transactions, `l.ListTransactions` returns every category's:
   ```go
   for _, transaction := range l.ListTransactions() {
       println(transaction.String())
   }
   ```
6. Read a balance. Both `Category.Balance` and `l.Balance` sum signed amounts, so received money counts up and spent money counts down:
   ```go
   println(groceries.Balance()) // -8450
   println(l.Balance())         // 241550
   ```
7. Delete a transaction recorded by mistake with `Remove`, which reports whether the record was found and deleted:
   ```go
   mistake, _ := groceries.AddSpend("wrong amount", 99900)
   mistake.Remove()
   ```
8. Run the program:
   ```bash
   go run main.go
   ```

---

## Full Code
```go
package main

import (
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
    l := agnoslib.New(agnosadapter.New("trackerdata"))

    groceries, _ := l.AddCategory("groceries")

    transaction, ok := groceries.AddSpend("weekly shopping", 8450) // 84.50
    if !ok {
        println("nothing recorded")
        return
    }
    println(transaction.String())

    l.AddCategory("salary")
    salary, _ := l.GetCategory("salary")
    salary.AddReceived("august paycheck", 250000) // 2500.00

    // The category must exist to record straight from the lib
    if _, ok := l.AddSpend("holidays", "flight", 42000); !ok {
        println("holidays: no such category, nothing recorded")
    }

    for _, transaction := range l.ListTransactions() {
        println(transaction.String())
    }

    println(groceries.Balance())
    println(l.Balance())

    // A transaction recorded by mistake removes itself
    mistake, _ := groceries.AddSpend("wrong amount", 99900)
    mistake.Remove()
}
```
