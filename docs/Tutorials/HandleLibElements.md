# Handle Library Elements

## Description
Covers adding new functions to the library's internal logic and its public API. Public functions are exposed to external callers (like the CLI or other consumers) via the `Lib` contract and documented in the public API. Internal functions are used exclusively inside `sandbox/` and are simpler to add. 

### Rules
- `sandbox/` is a closed sandbox: library code must never import [adapters/](/adapters/), [examples/](/examples/), a third-party module, or an OS-bound standard-library package (`os`, `net`, `syscall`, …) — reach every such effect through `Deps`. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- Adding a directory or file to [sandbox/](/sandbox/) requires updating [Structure.md](/docs/References/Structure.md).

---

## AddLibFunction

Internal library functions reside in packages within [sandbox/lib/](/sandbox/lib/) (e.g., `store`, `category`, `transaction`) and are meant for internal domain logic. They are not exposed to external consumers.

### Workflow
1. Create or open the relevant domain package in [sandbox/lib/](/sandbox/lib/) (e.g., [sandbox/lib/store/store.go](/sandbox/lib/store/store.go)).
2. Write the standard Go function. If it requires dependencies, pass them as arguments (typically `deps.Deps` or specific values):
   ```go
   func CalculateTotal(d deps.Deps, category string) int64 {
       // logic...
   }
   ```
3. If a new file or package was created, register it in [Structure.md](/docs/References/Structure.md).
4. If the function needs a dependency that is not yet in the contract, add it following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).

---

## AddPublicLibFunction

Public library functions form the API exposed to the CLI and external consumers. They must be declared in the API contract, constructed via a factory, and documented.

### Rules
- A public function is only usable once its factory's return value is assigned in the library's `New(d deps.Deps)` constructor inside [sandbox/lib/new.go](/sandbox/lib/new.go).
- One factory per function, named `<Function>Factory`, returning one closure.
- Dependencies are reached as `l.Deps.<Field>(...)` **inside** the closure, never captured at factory time.
- Every public-facing entry must be listed in [PublicApi.md](/docs/References/PublicApi.md) with a detail page under [docs/References/PublicApi/](/docs/References/PublicApi/) named `<pkg>.<Symbol>.md`.

### Workflow
1. Declare the function as a field of the `Lib` struct in [sandbox/contracts/api/api.go](/sandbox/contracts/api/api.go):
   ```go
   type Lib struct {
       Deps           deps.Deps
       // ... existing functions
       HasCategory    func(name string) bool // new function
   }
   ```
2. Write its factory in a new or existing file under [sandbox/lib/publicfunctions/](/sandbox/lib/publicfunctions/). The factory must return the closure:
   ```go
   package publicfunctions

   import (
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/lib/store"
   )

   // HasCategoryFactory returns the closure that fills api.Lib.HasCategory.
   func HasCategoryFactory(l *api.Lib) func(name string) bool {
       return func(name string) bool {
           _, ok := store.FindCategory(l.Deps, name)
           return ok
       }
   }
   ```
3. Assign the factory's return value in the package's `New` constructor inside [sandbox/lib/new.go](/sandbox/lib/new.go) — without this line the field stays nil and panics when called:
   ```go
   func New(d deps.Deps) api.Lib {
       l := api.Lib{Deps: d}
       // ... existing assignments
       l.HasCategory = publicfunctions.HasCategoryFactory(&l) // register the new function
       return l
   }
   ```
4. If a new file was created, register it in [Structure.md](/docs/References/Structure.md).
5. Build the project and call the new field once to confirm it is not nil.
6. Expose it in the public API:
   - Add the function to [PublicApi.md](/docs/References/PublicApi.md) with a one-line description.
   - Create its detail page under [docs/References/PublicApi/](/docs/References/PublicApi/), named `api.HasCategory.md`, following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md).
   - Link the new detail page from its entry in [PublicApi.md](/docs/References/PublicApi.md).
   - Register the detail page in [Structure.md](/docs/References/Structure.md).
