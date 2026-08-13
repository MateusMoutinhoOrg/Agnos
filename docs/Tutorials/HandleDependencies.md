# Handle Dependencies

## Description
Explains how the library receives its dependencies — the `Deps` contract in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go), how an injected value propagates through the object graph, and how to add a requirement to the contract. Third stage of the Development learning path: it assumes [SandboxIsolation.md](/docs/References/SandboxIsolation.md) and [StructContracts.md](/docs/References/StructContracts.md). Using the dependency from library code is covered by [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).

---

## Find Dependencies Functions you can use

`sandbox/contracts/deps` declares what the library needs; `sandbox/contracts/api` declares what it hands back. Nothing else crosses the boundary, and `lib.New(deps.Deps) api.Lib` is the single wiring point.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now       func() time.Time
	Printf    func(format string, a ...any) (int, error)
	VerbLib   verbdeps.Lib
	KeepLib   keepdeps.Lib
	EmbedDeps embeddeps.Lib
}
```

`Deps` is the *only* door in the sandbox wall: since nothing under `sandbox/` may import an adapter, a third-party module, or an OS-bound stdlib package, every effect the library performs has to be a field on this struct.

`lib.New` stores the `Deps` on the `api.Lib` struct and runs the factories over it; each closure reads `l.Deps` when the field is *called*, not when the factory ran. Every object the lib creates receives the same `Deps`, passed into the object package's `New` constructor. So a dependency injected once is reachable from anywhere in the object graph.

---

## Add New Dependencie

### Rules
- A requirement is a **function field** declaring behavior, never a concrete implementation.
- A new field must be filled by **every** adapter in [adapters/](/adapters/) in the same commit. The compiler will **not** catch a missing one: the field stays nil and panics on first call.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

### Workflow
1. Add the field to `Deps` in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Now       func() time.Time
       Printf    func(format string, a ...any) (n int, err error)
       VerbLib   verbdeps.Lib
       KeepLib   keepdeps.Lib
       EmbedDeps embeddeps.Lib
       Uuid      func() string // new requirement
   }
   ```
2. On every adapter, write a `<Field>Factory` returning the closure and assign it from that adapter's `New`, following the adapter specification located in [Specs.md](/docs/References/Specs.md):
   ```go
   // UuidFactory returns the closure that fills deps.Deps.Uuid,
   // handing back a fresh random identifier.
   func UuidFactory(s *StandardAdapter) func() string {
       return func() string { return uuid.NewString() }
   }

   func New(basePath string) deps.Deps {
       adapter := &StandardAdapter{args: os.Args[1:], keepBasePath: basePath}
       adapter.Deps.Now = NowFactory(adapter)
       adapter.Deps.VerbLib = VerbLibFactory(adapter)
       adapter.Deps.KeepLib = KeepLibFactory(adapter)
       adapter.Deps.EmbedDeps = EmbedDepsFactory(adapter)
       adapter.Deps.Uuid = UuidFactory(adapter) // the new field
       return adapter.Deps
   }
   ```
3. Grep every adapter's `New` for the new assignment — this step replaces the compiler check:
   ```bash
   grep -rn "Factory(adapter)" adapters/
   ```
4. Use the dependency from the library through `l.Deps.<Field>(...)`, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).
5. Build and run a sample — an unfilled field surfaces at runtime, not at build time.

Three fields are not behaviors but other libraries built with this same pattern: `VerbLib`, `KeepLib`, and `EmbedDeps`. Because a contract is a struct of function fields, the whole library fits in one plain field. The sandbox cannot import Verb or Keep, nor Go's `embed` machinery, so it declares a copy of each api in `sandbox/contracts/deps/verbdeps/`, `keepdeps/` and `embeddeps/`. The adapter initializes the real library and assigns its fields onto that copy — the factory returns a **value** instead of a closure.

---

## Overwrinting a adapter function

Take the `deps.Deps` an adapter returns and reassign the field you want; every other field keeps the adapter's implementation:

```go
myDeps := agnosadapter.New("trackerdata")

// Replace only the clock — KeepLib stays as the adapter built it
now := time.Unix(0, 0)
myDeps.Now = func() time.Time { return now }

l := agnoslib.New(myDeps)
l.AddCategory("groceries")

// Moving the captured variable moves the clock the lib sees
now = time.Unix(120, 0)
transaction, _ := l.AddSpend("groceries", "weekly shopping", 8450)
println(transaction.OccurredAt.Unix()) // 120
```

> **Careful:** patch the `deps.Deps` value **before** calling `lib.New`. The factories close over the `api.Lib` they ran on, so assigning to `l.Deps.Now` afterwards changes nothing — see [StructContracts.md](/docs/References/StructContracts.md#what-it-costs).

---

## Creating a adapter in repo

Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](/adapters/). Assumes the mechanics in [StructContracts.md](/docs/References/StructContracts.md); the shipped adapters are listed in [Adapters.md](/docs/References/Adapters.md).

### Rules
- Each adapter lives in its own directory under [adapters/](/adapters/) and uses a package named after that directory.
- The adapter is a struct carrying a `Deps deps.Deps` field, filled by one **factory** per field of the contract — the same factory pattern `sandbox/` uses, and a binding rule of the project. See [RULES.md](/docs/References/RULES.md#factory-pattern).
- Fields are never filled by binding methods of the adapter. Methods may exist only as unexported helpers a closure calls.
- A single `New(...) deps.Deps` constructor calls **every** field factory, assigns its return value into the matching field, and returns the `deps.Deps` contract struct, never the concrete adapter type.
- Filling every field is the author's job — an unassigned field compiles and panics on first call. See [StructContracts.md](/docs/References/StructContracts.md).
- An adapter lives outside the sandbox and is the only place OS-bound and third-party code is allowed. It may import [sandbox/contracts/deps](/sandbox/contracts/deps/), but never [sandbox/](/sandbox/) — see [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- The adapter file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

### Workflow
1. Create the adapter directory and its file, both named after the adapter (e.g., `adapters/frozen/frozen.go`).
2. Declare the package and the adapter struct — the **carrier**, leading with the `Deps` field its factories fill, followed by its configuration and state:
   ```go
   package frozen

   import (
       "time"

       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
   )

   // FrozenAdapter fills deps.Deps with a fixed clock, so every category
   // and transaction is stamped with a known time. The database and the argv
   // parser come from the embedded libraries, as in every other adapter.
   type FrozenAdapter struct {
       // Deps is the contract this adapter fills; its factories assign into it.
       Deps deps.Deps
       now  time.Time
       args []string
   }
   ```
3. Write one `<Field>Factory` per field of the `Deps` contract, each returning a single closure that reads the adapter's state through the pointer:
   ```go
   // NowFactory returns the closure that fills deps.Deps.Now, returning the
   // adapter's fixed clock.
   func NowFactory(f *FrozenAdapter) func() time.Time {
       return func() time.Time { return f.now }
   }

   // KeepLibFactory returns the value that fills deps.Deps.KeepLib: the
   // embedded Keep schema-database library, wired to Keep's in-memory
   // adapter and copied onto the sandbox's local keepdeps.Lib. A field that
   // is not a function has its factory return a value, not a closure.
   func KeepLibFactory(f *FrozenAdapter) keepdeps.Lib {
       inner := keeplib.New(keepadapter.New())
       return keepdeps.Lib{
           NewDatabase: func(props keepdeps.Props) keepdeps.KeepDatabase {
               return fromKeepDatabase(inner.NewDatabase(toKeepProps(props)))
           },
       }
   }
   ```
   Reading `f.now` inside the closure — instead of capturing it when the factory runs — is what carries the adapter's live state into the library.
4. Give a factory a file of its own when its body brings conversion helpers with it, named after the field it fills — `adapters/standard/embed.go` holds `EmbedDepsFactory` and the asset walk behind it. The `Factories` specification applies the same way there.
5. Expose the `New` constructor: build the adapter instance, run every field factory over it, assign each return value into its matching field, and return its `Deps`:
   ```go
   // New creates a deps.Deps whose clock is frozen at the given time.
   func New(now time.Time) deps.Deps {
       adapter := &FrozenAdapter{now: now}
       adapter.Deps.Now = NowFactory(adapter)
       adapter.Deps.VerbLib = VerbLibFactory(adapter)
       adapter.Deps.KeepLib = KeepLibFactory(adapter)
       adapter.Deps.EmbedDeps = EmbedDepsFactory(adapter)
       return adapter.Deps
   }
   ```
6. Compare the assignments in your `New` against `sandbox/contracts/deps/deps.go` field by field. A missing field will **not** fail the build.
7. Register the new directory and file in [Structure.md](/docs/References/Structure.md), and add a row for the adapter in [Adapters.md](/docs/References/Adapters.md).
8. If the adapter is public-facing, expose its `New` factory following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).
9. If the adapter needs a runnable demonstration, add one following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md).
10. Build the project and exercise the adapter:
   ```bash
   go build ./...
   ```

---

## Creating a adapter in your project

For complete control, build the `deps.Deps` as a struct literal — no type to declare, no method set to satisfy:

```go
myDeps := agnosdeps.Deps{
	Now:     func() time.Time { return time.Unix(0, 0) },
	KeepLib: agnoskeepdeps.Lib{NewDatabase: myOwnDatabase},
	// Printf, VerbLib and EmbedDeps left zero: this program never calls
	// Sandboxmain, and only Sandboxmain prints or reads an asset
}
l := agnoslib.New(myDeps)
```

> **Careful:** the compiler cannot tell you a field is missing — an unfilled field panics on first call. In practice, start from an adapter and patch what you need: `KeepLib` is a whole database api, and no program should reimplement it just to change the clock.

If your project requires setting up everything from scratch, simply provide a factory or a function literal for each required field in `deps.Deps`.
