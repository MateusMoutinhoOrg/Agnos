# Add an Adapter

## Description
Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](../../adapters/).

### Rules
- Each adapter lives in its own directory under [adapters/](../../adapters/) and uses a package named after that directory.
- The adapter is a struct carrying a `Deps deps.Deps` field, filled by one **factory** per field of the contract — the same factory pattern `sandbox/internal/` uses, and a binding rule of the project. See [RULES.md](/docs/References/RULES.md#factory-pattern).
- Fields are never filled by binding methods of the adapter. Methods may exist only as unexported helpers a closure calls.
- A single `New(...) deps.Deps` constructor calls **every** field factory, assigns its return value into the matching field, and returns the `deps.Deps` contract struct, never the concrete adapter type.
- Filling every field is the author's job — an unassigned field compiles and panics on first call. See [StructContracts.md](/docs/Explanations/StructContracts.md).
- An adapter lives outside the sandbox and is the only place OS-bound and third-party code is allowed. It may import [sandbox/contracts/deps](../../sandbox/contracts/deps/), but never [sandbox/internal/](../../sandbox/internal/) — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- The adapter file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

---

## Workflow
1. Create the adapter directory and its file, both named after the adapter (e.g., `adapters/memory/memory.go`).
2. Declare the package and the adapter struct — the **carrier**, leading with the `Deps` field its factories fill, followed by its configuration and state:
   ```go
   package frozen

   import (
       "time"

       "github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
   )

   // FrozenAdapter fills deps.Deps with a fixed clock, so expiry
   // is deterministic. Records are kept in a map.
   type FrozenAdapter struct {
       // Deps is the contract this adapter fills; its factories assign into it.
       Deps  deps.Deps
       now   time.Time
       store map[string]record
   }

   type record struct {
       value         string
       expiresAtUnix int64
   }
   ```
3. Write one `<Field>Factory` per field of the `Deps` contract, each returning a single closure that reads the adapter's state through the pointer:
   ```go
   // NowFactory returns the closure that fills deps.Deps.Now, returning the
   // adapter's fixed clock.
   func NowFactory(f *FrozenAdapter) func() time.Time {
       return func() time.Time { return f.now }
   }

   // LoadFactory returns the closure that fills deps.Deps.Load, fetching a
   // record from the map.
   func LoadFactory(f *FrozenAdapter) func(key string) (string, int64, bool) {
       return func(key string) (string, int64, bool) {
           r, ok := f.store[key]
           return r.value, r.expiresAtUnix, ok
       }
   }

   // StoreFactory returns the closure that fills deps.Deps.Store, writing a
   // record into the map.
   func StoreFactory(f *FrozenAdapter) func(key, value string, expiresAtUnix int64) {
       return func(key, value string, expiresAtUnix int64) {
           f.store[key] = record{value: value, expiresAtUnix: expiresAtUnix}
       }
   }
   ```
   Reading `f.now` and `f.store` inside the closure — instead of capturing them when the factory runs — is what carries the adapter's live state into the library.
4. Expose the `New` constructor: build the adapter instance, run every field factory over it, assign each return value into its matching field, and return its `Deps`:
   ```go
   // New creates a deps.Deps whose clock is frozen at the given time.
   func New(now time.Time) deps.Deps {
       adapter := &FrozenAdapter{now: now, store: make(map[string]record)}
       adapter.Deps.Now = NowFactory(adapter)
       adapter.Deps.Load = LoadFactory(adapter)
       adapter.Deps.Store = StoreFactory(adapter)
       return adapter.Deps
   }
   ```
5. Compare the assignments in your `New` against `sandbox/contracts/deps/deps.go` field by field. A missing field will **not** fail the build.
6. Register the new directory and file in [Structure.md](/docs/References/Structure.md), and add a row for the adapter in [Adapters.md](/docs/References/Adapters.md).
7. If the adapter is public-facing, expose its `New` factory following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
8. If the adapter needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Build the project and exercise the adapter:
   ```bash
   go build ./...
   ```
