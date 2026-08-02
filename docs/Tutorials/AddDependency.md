# Add a Dependency

## Description
Covers adding a requirement to the `Deps` contract in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go) and filling it in every existing adapter.

### Rules
- `Deps` is a struct of function fields — a requirement is a **function field**, declaring behavior and never a concrete implementation.
- A new field must be filled by **every** adapter in [adapters/](../../adapters/) in the same commit. The compiler will **not** catch an adapter that misses it: the field stays nil and panics on the first call, so this check is on you.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

---

## Workflow
1. Add the field to the `Deps` struct in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Now    func() time.Time
       Load   func(key string) (value string, expiresAtUnix int64, ok bool)
       Store  func(key string, value string, expiresAtUnix int64)
       Delete func(key string) // new requirement
   }
   ```
2. Write a `<Field>Factory` for the new field on every adapter under [adapters/](../../adapters/), returning the closure, and assign its return value from that adapter's `New`, following the adapter specification located in [Specs.md](/docs/References/Specs.md):
   ```go
   // DeleteFactory returns the closure that fills the new deps.Deps.Delete
   // requirement, removing a record from the store.
   func DeleteFactory(s *StandardAdapter) func(key string) {
       return func(key string) {
           s.mu.Lock()
           defer s.mu.Unlock()
           delete(s.store, key)
       }
   }

   func New(filePath string) deps.Deps {
       adapter := &StandardAdapter{filePath: filePath}
       adapter.Deps.Now = NowFactory(adapter)
       adapter.Deps.Load = LoadFactory(adapter)
       adapter.Deps.Store = StoreFactory(adapter)
       adapter.Deps.Delete = DeleteFactory(adapter) // assign the new field's factory
       return adapter.Deps
   }
   ```
3. Grep every adapter's `New` for the new factory assignment to be sure none was missed — this is the step that replaces the compiler check:
   ```bash
   grep -rn "Factory(adapter)" adapters/ bootstrap/
   ```
4. Use the dependency from the library through `l.Deps.<Field>(...)`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
5. If the requirement changes how dependencies behave for consumers, update [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).
6. Build the project and run a sample — an unfilled field surfaces at runtime, not at build time:
   ```bash
   go build ./...
   go run ./examples/KvCacheSample/KvCacheSample.go
   ```
