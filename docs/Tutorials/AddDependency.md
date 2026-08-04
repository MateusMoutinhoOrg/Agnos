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
       Now     func() time.Time
       VerbLib verbdeps.Lib
       KeepLib keepdeps.Lib
       Uuid    func() string // new requirement
   }
   ```
2. Write a `<Field>Factory` for the new field on every adapter under [adapters/](../../adapters/), returning the closure, and assign its return value from that adapter's `New`, following the adapter specification located in [Specs.md](/docs/References/Specs.md):
   ```go
   // UuidFactory returns the closure that fills the new deps.Deps.Uuid
   // requirement, handing back a fresh random identifier.
   func UuidFactory(s *StandardAdapter) func() string {
       return func() string {
           return uuid.NewString()
       }
   }

   func New(basePath string) deps.Deps {
       adapter := &StandardAdapter{args: os.Args[1:], keepBasePath: basePath}
       adapter.Deps.Now = NowFactory(adapter)
       adapter.Deps.VerbLib = VerbLibFactory(adapter)
       adapter.Deps.KeepLib = KeepLibFactory(adapter)
       adapter.Deps.Uuid = UuidFactory(adapter) // assign the new field's factory
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
   go run ./examples/TrackSpendSample/TrackSpendSample.go
   ```
