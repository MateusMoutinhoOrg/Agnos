# Handle Library Elements

## Description
Covers adding new elements — functions and objects — to the library's public API: declare the struct/field in [sandbox/contracts/api/api.go](../sandbox/contracts/api/api.go), write the factory under [sandbox/internal/](../sandbox/internal/), and register it in the package's `New` constructor. Assumes the mechanics in [StructContracts.md](/docs/StructContracts.md). The CLI command calling the new element is a separate goal — [HandleCliCommands.md](/docs/HandleCliCommands.md); publishing it is [ExposePublicApi.md](/docs/ExposePublicApi.md).

### Rules
- A function or object field is only usable once its factory's return value is assigned from the package's `New(d deps.Deps, …)` constructor, which doubles as the factory aggregate — an unassigned field stays nil and panics on first call. The compiler does not catch this.
- One factory per field, named `<Field>Factory`, returning one closure.
- Dependencies are reached as `l.Deps.<Field>(...)` or `b.Deps.<Field>(...)` **inside** the closure, never captured at factory time — that is what keeps the injected value authoritative.
- `sandbox/` is a closed sandbox: library code must never import [adapters/](../adapters/), [examples/libraryExamples/](../examples/libraryExamples/), a third-party module, or an OS-bound standard-library package (`os`, `net`, `syscall`, …) — reach every such effect through `Deps`. See [SandboxIsolation.md](/docs/SandboxIsolation.md).
- Adding a directory or file to [sandbox/internal/](../sandbox/internal/) requires updating [Structure.md](/docs/Structure.md).

---

## Add a Library Function

### Workflow
1. Declare the function as a field of the `Lib` struct in [sandbox/contracts/api/api.go](../sandbox/contracts/api/api.go):
   ```go
   type Lib struct {
       Deps           deps.Deps
       AddCategory    func(name string) (Category, bool)
       GetCategory    func(name string) (Category, bool)
       HasCategory    func(name string) bool // new function
   }
   ```
2. Write its factory in a new or existing file in [sandbox/internal/lib/](../sandbox/internal/lib/), with the identical signature, returning the closure:
   ```go
   // HasCategoryFactory returns the closure that fills api.Lib.HasCategory,
   // reporting whether a category is stored under that name.
   func HasCategoryFactory(l *api.Lib) func(name string) bool {
       return func(name string) bool {
           _, ok := l.GetCategory(name)
           return ok
       }
   }
   ```
   > Calling another field from inside a closure (`l.GetCategory` above) is fine: by the time `HasCategory` runs, `New` has already filled every field.
3. Assign the factory's return value in the package's `New` constructor — without this line the field stays nil and the function panics when called:
   ```go
   func New(d deps.Deps) api.Lib {
       l := api.Lib{Deps: d}
       l.AddCategory = AddCategoryFactory(&l)
       l.GetCategory = GetCategoryFactory(&l)
       l.HasCategory = HasCategoryFactory(&l) // register the new function
       return l
   }
   ```
4. If the function needs a dependency that is not yet in the contract, add it following [HandleDependencies.md](/docs/HandleDependencies.md).
5. Expose the function following [ExposePublicApi.md](/docs/ExposePublicApi.md).
6. If a new file was created, register it in [Structure.md](/docs/Structure.md).
7. If the function needs a runnable demonstration, add one following [HandleSamples.md](/docs/HandleSamples.md).
8. Build the project and call the new field once to confirm it is not nil.

---

## Add a Library Object

### Rules
- The object **is** its api struct. There is no internal mirror type: [sandbox/internal/](../sandbox/internal/) holds only the factories and the constructor.
- An object that needs dependencies declares a `Deps deps.Deps` field, filled by its `New` constructor from the parent lib's `l.Deps`. Its factories read that field inside their closures.
- Every api field must be exported: the factories fill them from another package, and consumers read them.

### Workflow
1. Declare the object's struct in [sandbox/contracts/api/api.go](../sandbox/contracts/api/api.go).
   ```go
   type Budget struct {
       Deps     deps.Deps
       Category string
       Limit    int64
       Exceeded func() bool
   }
   ```
2. Create the object's package and file, both named after the object (e.g. `sandbox/internal/budget/budget.go`) — no `internal_` prefix, the `internal/` parent already marks it private:
   ```go
   package budget

   import (
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/store"
   )

   // ExceededFactory returns the closure that fills api.Budget.Exceeded...
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

   // New builds an api.Budget, propagating the library's Deps into it...
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
3. Declare the constructor as a field of the `Lib` api struct in [sandbox/contracts/api/api.go](../sandbox/contracts/api/api.go), returning the object's api struct:
   ```go
   type Lib struct {
       NewBudget func(category string, limit int64) Budget
   }
   ```
4. Write the constructor's factory in [sandbox/internal/lib/](../sandbox/internal/lib/), propagating `l.Deps` into the new object:
   ```go
   // NewBudgetFactory returns the closure that fills api.Lib.NewBudget...
   func NewBudgetFactory(l *api.Lib) func(category string, limit int64) api.Budget {
       return func(category string, limit int64) api.Budget {
           return budget.New(l.Deps, category, limit)
       }
   }
   ```
5. Assign `NewBudgetFactory`'s return value in the lib package's `New` constructor (Step 3 of "Add a Library Function").
6. Expose the object, its constructor, and its fields following [ExposePublicApi.md](/docs/ExposePublicApi.md).
7. Register the new directory and file in [Structure.md](/docs/Structure.md).
