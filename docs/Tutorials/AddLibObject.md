# Add a Library Object

## Description
Covers adding an object created by the library: its struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), the factories that fill its function fields under [sandbox/internal/](../../sandbox/internal/), and the `New` constructor that propagates the deps into it.

### Rules
- The object **is** its api struct. There is no internal mirror type: [sandbox/internal/](../../sandbox/internal/) holds only the factories and the constructor.
- An object that needs dependencies declares a `Deps deps.Deps` field, filled by its `New` constructor from the parent lib's `l.Deps`. Its factories read that field inside their closures.
- Every field factory's return value must be assigned in the package's `New(d deps.Deps, …) api.<Object>` constructor, which doubles as the factory aggregate — an unassigned field stays nil and panics on first call.
- Every api field must be exported: the factories fill them from another package, and consumers read them.
- Adding a directory or file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Declare the object's struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go). This walkthrough adds a `Budget`, a monthly cap set on a category:
   ```go
   type Budget struct {
       Deps     deps.Deps
       Category string
       Limit    int64
       Exceeded func() bool
   }
   ```
2. Create the object's package and file, both named after the object (e.g. `sandbox/internal/budget/budget.go`) — no `internal_` prefix, the `internal/` parent already marks it private — holding its field factories and its `New` constructor:
   ```go
   package budget

   import (
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/store"
   )

   // ExceededFactory returns the closure that fills api.Budget.Exceeded,
   // comparing the category's live balance against the budget's limit.
   func ExceededFactory(b *api.Budget) func() bool {
       return func() bool {
           record, ok := store.FindCategory(b.Deps, b.Category)
           if !ok {
               return false
           }
           spent := int64(0)
           for _, transaction := range record.ListAll(store.TransactionsField) {
               spent += store.Number(transaction, store.AmountField)
           }
           return spent > b.Limit
       }
   }

   // New builds an api.Budget, propagating the library's Deps into it, and
   // runs every budget factory over it, assigning each return value into its
   // matching function field.
   func New(d deps.Deps, category string, limit int64) api.Budget {
       b := api.Budget{
           Deps:     d,
           Category: category,
           Limit:    limit,
       }
       b.Exceeded = ExceededFactory(&b)
       return b
   }
   ```
3. Declare the constructor as a field of the `Lib` api struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), returning the object's api struct:
   ```go
   type Lib struct {
       NewBudget func(category string, limit int64) Budget
   }
   ```
4. Write the constructor's factory in [sandbox/internal/lib/](../../sandbox/internal/lib/), propagating `l.Deps` into the new object:
   ```go
   // NewBudgetFactory returns the closure that fills api.Lib.NewBudget,
   // creating a Budget with the library's injected dependencies wired in.
   func NewBudgetFactory(l *api.Lib) func(category string, limit int64) api.Budget {
       return func(category string, limit int64) api.Budget {
           return budget.New(l.Deps, category, limit)
       }
   }
   ```
5. Assign `NewBudgetFactory`'s return value in the lib package's `New` constructor, as described in [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
6. If a factory needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
7. Expose the object, its constructor, and its fields following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
8. Register the new directory and file in [Structure.md](/docs/References/Structure.md).
9. If the object needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
10. Build the project, then call the new field once to confirm it is not nil — a missing assignment in `New` compiles cleanly:
    ```bash
    go build ./...
    ```
