# Add a Library Function

## Description
Covers adding a function to the library: declaring it as a field of the `Lib` struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), writing the factory that fills it in [sandbox/internal/lib/](../../sandbox/internal/lib/), and registering that factory in the package's `New` constructor.

### Rules
- A function is only usable once its factory's return value is assigned from the package's `New(d deps.Deps) api.Lib` constructor, which doubles as the factory aggregate — an unassigned field stays nil and panics on first call. The compiler does not catch this.
- One factory per field, named `<Field>Factory`, taking a single `*api.Lib` and returning one closure.
- Dependencies are reached as `l.Deps.<Field>(...)` **inside** the closure, never captured at factory time — that is what keeps the injected value authoritative.
- `sandbox/` is a closed sandbox: library code must never import [adapters/](../../adapters/), [examples/](../../examples/), a third-party module, or an OS-bound standard-library package (`os`, `net`, `syscall`, …) — reach every such effect through `l.Deps`. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- Adding a file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Declare the function as a field of the `Lib` struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go):
   ```go
   type Lib struct {
       Deps           deps.Deps
       AddCategory    func(name string) (Category, bool)
       GetCategory    func(name string) (Category, bool)
       HasCategory    func(name string) bool // new function
   }
   ```
2. Write its factory in a new or existing file in [sandbox/internal/lib/](../../sandbox/internal/lib/), with the identical signature, returning the closure:
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
4. If the function needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. If the function returns a new object, create it following [AddLibObject.md](/docs/Tutorials/AddLibObject.md) and return the object's `api` struct.
6. Expose the function following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
7. If a new file was created, register it in [Structure.md](/docs/References/Structure.md).
8. If the function needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Build the project — a signature mismatch fails here, but a forgotten assignment in `New` does not:
   ```bash
   go build ./...
   ```
10. Call the new field once — from a sample or a test — to confirm it is not nil. That is the only check that catches a missing assignment.
