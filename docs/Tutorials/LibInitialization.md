# Library Initialization

## Description
Covers installing the library and initializing it with the standard adapter in a new program. Creating categories after initialization is covered by [ManageCategories.md](/docs/Tutorials/ManageCategories.md), and recording transactions by [TrackTransactions.md](/docs/Tutorials/TrackTransactions.md). For other ways to build the dependencies, see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

---

## Workflow
1. Install the lib:
   ```bash
   go get github.com/MateusMoutinhoOrg/Agnos-Cli@v0.0.7
   ```
2. Create a file called `main.go` with the following code:
   ```go
   package main

   // 1. Import the standard adapter and the lib
   import (
       agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
       agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
   )

   func main() {
       // 2. Create deps via an adapter (the "opinionated" part)
       deps := agnosadapter.New("trackerdata")

       // 3. Inject deps into the pure library — a financial tracker
       l := agnoslib.New(deps)

       // 4. Use the library — it never knows which adapter is behind the scenes
       l.AddCategory("groceries")
       l.AddSpend("groceries", "weekly shopping", 8450) // 84.50
       println(l.Balance())
   }
   ```
3. Run the code:
   ```bash
   go run main.go
   ```
