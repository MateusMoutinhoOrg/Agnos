# Manage Categories

## Description
Covers creating the categories transactions are tracked under, reading them back, and removing one. Installing and initializing the lib is covered by [LibInitialization.md](/docs/LibInitialization.md); recording the transactions themselves is covered by [TrackTransactions.md](/docs/TrackTransactions.md).

---

## Workflow
1. Initialize the lib as shown in [LibInitialization.md](/docs/LibInitialization.md):
   ```go
   deps := agnosadapter.New("trackerdata")
   l := agnoslib.New(deps)
   ```
2. Create a category with `AddCategory`. It returns an `api.Category` and a boolean that is `false` when the name is empty or the record could not be written:
   ```go
   groceries, ok := l.AddCategory("groceries")
   if !ok {
       println("could not create the category")
       return
   }
   ```
3. Call `AddCategory` again with the same name whenever you need the category. Creation is idempotent — the stored category is returned instead of a failure — so a program can declare the categories it needs on every run:
   ```go
   l.AddCategory("salary")
   l.AddCategory("salary") // same category, no duplicate
   ```
4. Look a stored category up by name with `GetCategory`, which does **not** create anything. Branch on the boolean: a struct return has no nil to compare against:
   ```go
   if salary, ok := l.GetCategory("salary"); ok {
       println(salary.String())
   } else {
       println("salary: no such category")
   }
   ```
5. List every stored category with `ListCategories`. Each value is read fresh from the database, so its `Balance` and `ListTransactions` are current:
   ```go
   for _, category := range l.ListCategories() {
       println(category.String())
   }
   ```
6. Remove a category with `Remove`. Every transaction stored under it goes with it:
   ```go
   if holidays, ok := l.GetCategory("holidays"); ok {
       holidays.Remove()
   }
   ```
7. Run the program:
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
    deps := agnosadapter.New("trackerdata")
    l := agnoslib.New(deps)

    groceries, ok := l.AddCategory("groceries")
    if !ok {
        println("could not create the category")
        return
    }
    println(groceries.Name)

    // Creation is idempotent — no duplicate is stored
    l.AddCategory("salary")
    l.AddCategory("salary")

    if salary, ok := l.GetCategory("salary"); ok {
        println(salary.String())
    } else {
        println("salary: no such category")
    }

    for _, category := range l.ListCategories() {
        println(category.String())
    }

    // Removing a category removes its transactions too
    if holidays, ok := l.GetCategory("holidays"); ok {
        holidays.Remove()
    }
}
```
