# Handle Dependencies

## Description
Explains how the library receives its dependencies — the `Deps` contract in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go), how an injected value propagates through the object graph, and how to add a requirement to the contract. Third stage of the Development learning path: it assumes [SandboxIsolation.md](/docs/Development/References/SandboxIsolation.md) and [StructContracts.md](/docs/Development/References/StructContracts.md). Using the dependency from library code is covered by [HandleLibElements.md](/docs/Development/Tutorials/HandleLibElements.md).

---

## The Contract

`sandbox/contracts/deps` declares what the library needs; `sandbox/contracts/api` declares what it hands back. Nothing else crosses the boundary, and `lib.New(deps.Deps) api.Lib` is the single wiring point.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now     func() time.Time
	Printf  func(format string, a ...any) (int, error)
	VerbLib verbdeps.Lib
	KeepLib keepdeps.Lib
}
```

`Deps` is the *only* door in the sandbox wall: since nothing under `sandbox/` may import an adapter, a third-party module, or an OS-bound stdlib package, every effect the library performs has to be a field on this struct.

---

## Propagation

`lib.New` stores the `Deps` on the `api.Lib` struct and runs the factories over it; each closure reads `l.Deps` when the field is *called*, not when the factory ran. Every object the lib creates receives the same `Deps`, passed into the object package's `New` constructor:

```go
// sandbox/internal/lib/lib.go — inside GetCategoryFactory's closure
c := category.New(l.Deps, record) // the object stores Deps on its own api struct
```

So a dependency injected once is reachable from anywhere in the object graph:

```
standard.New() ──▶ deps.Deps ──▶ lib.New(deps) ──▶ api.Lib
                                                     │ GetCategory()
                                                     ▼
                                               api.Category
                                                     │ ListTransactions()
                                                     ▼
                                             api.Transaction
```

---

## Using an Adapter

The simplest setup — an adapter builds a ready-to-use `deps.Deps`:

```go
myDeps := agnosadapter.New("trackerdata") // fills every field
l := agnoslib.New(myDeps)
```

---

## Overwriting a Single Behavior

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

> **Careful:** patch the `deps.Deps` value **before** calling `lib.New`. The factories close over the `api.Lib` they ran on, so assigning to `l.Deps.Now` afterwards changes nothing — see [StructContracts.md](/docs/Development/References/StructContracts.md#what-it-costs).

---

## Writing Custom Deps

For complete control, build the `deps.Deps` as a struct literal — no type to declare, no method set to satisfy:

```go
myDeps := agnosdeps.Deps{
	Now:     func() time.Time { return time.Unix(0, 0) },
	KeepLib: agnoskeepdeps.Lib{NewDatabase: myOwnDatabase},
	// Printf and VerbLib left zero: this program never calls Sandboxmain
}
l := agnoslib.New(myDeps)
```

> **Careful:** the compiler cannot tell you a field is missing — an unfilled field panics on first call. In practice, start from an adapter and patch what you need: `KeepLib` is a whole database api, and no program should reimplement it just to change the clock.

---

## Injecting a Whole Library

Two fields are not behaviors but other libraries built with this same pattern: [`VerbLib`](/docs/LibUsage/References/PublicApi/verbdeps.Lib.md), the Verb argv parser, and [`KeepLib`](/docs/LibUsage/References/PublicApi/keepdeps.Lib.md), the Keep schema database. Because a contract is a struct of function fields, the whole library fits in one plain field — no getter, no bridging type.

The sandbox cannot import Verb or Keep, so it declares a copy of each api in `sandbox/contracts/deps/verbdeps/` and `keepdeps/`. The adapter initializes the real library and assigns its fields onto that copy — the factory returns a **value** instead of a closure, the one shape variation the [factory pattern](/docs/Development/References/RULES.md#factory-pattern) allows:

```go
// adapters/standard/standard.go
func VerbLibFactory(s *StandardAdapter) verbdeps.Lib {
	inner := verblib.New(s.args)
	return verbdeps.Lib{
		Args:      inner.Args,
		IsPresent: inner.IsPresent,
		// ... one assignment per field
	}
}
```

When the embedded library hands back structs of its own — Keep's `NewDatabase` returns a `KeepDatabase`, whose `GetSchema` returns a `SchemaInstance` — assigning straight across would leak the embedded types into the sandbox. The factory wraps each such function in a closure that converts what it returns, recursing as far as the api tree goes:

```go
// adapters/standard/standard.go
func KeepLibFactory(s *StandardAdapter) keepdeps.Lib {
	inner := keeplib.New(keepadapter.NewWithBase(s.keepBasePath))
	return keepdeps.Lib{
		NewDatabase: func(props keepdeps.Props) keepdeps.KeepDatabase {
			return fromKeepDatabase(inner.NewDatabase(toKeepProps(props)))
		},
	}
}
```

That conversion is the price of the wall, paid once per adapter, outside the sandbox. `bootstrap/` demonstrates the same mechanic from the other side: there, *this* library is the one being embedded — see [Bootstrap.md](/docs/LibUsage/References/Bootstrap.md) and [StructContracts.md](/docs/Development/References/StructContracts.md#consuming-a-library-that-uses-this-pattern).

---

## Add a Dependency

### Rules
- A requirement is a **function field** declaring behavior, never a concrete implementation.
- A new field must be filled by **every** adapter in [adapters/](/adapters/) (and `bootstrap/adapters/`) in the same commit. The compiler will **not** catch a missing one: the field stays nil and panics on first call.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/Development/References/Specs.md).

### Workflow
1. Add the field to `Deps` in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Now     func() time.Time
       Printf  func(format string, a ...any) (n int, err error)
       VerbLib verbdeps.Lib
       KeepLib keepdeps.Lib
       Uuid    func() string // new requirement
   }
   ```
2. On every adapter, write a `<Field>Factory` returning the closure and assign it from that adapter's `New`, following the adapter specification located in [Specs.md](/docs/Development/References/Specs.md):
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
       adapter.Deps.Uuid = UuidFactory(adapter) // the new field
       return adapter.Deps
   }
   ```
3. Grep every adapter's `New` for the new assignment — this step replaces the compiler check:
   ```bash
   grep -rn "Factory(adapter)" adapters/ bootstrap/
   ```
4. Use the dependency from the library through `l.Deps.<Field>(...)`, following [HandleLibElements.md](/docs/Development/Tutorials/HandleLibElements.md).
5. Build and run a sample — an unfilled field surfaces at runtime, not at build time:
   ```bash
   go build ./...
   go run ./examples/libraryExamples/TrackSpendSample/TrackSpendSample.go
   ```
