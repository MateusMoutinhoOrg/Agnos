# Handle Dependencies

## Description
Explains how the library receives its dependencies and covers adding a requirement to the `Deps` contract in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go) and filling it in every existing adapter.

`Deps` is a **struct of function fields**: the library declares the behaviors it needs as fields, and the caller supplies a value with every field filled — a pre-built adapter, an adapter with one field replaced, or a struct literal written by hand.

---

## Add a Dependency

### Rules
- `Deps` is a struct of function fields — a requirement is a **function field**, declaring behavior and never a concrete implementation.
- A new field must be filled by **every** adapter in [adapters/](../../adapters/) in the same commit. The compiler will **not** catch an adapter that misses it: the field stays nil and panics on the first call, so this check is on you.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/Specs.md).

### Workflow
1. Add the field to the `Deps` struct in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Now     func() time.Time
       Printf  func(format string, a ...any) (n int, err error)
       VerbLib verbdeps.Lib
       KeepLib keepdeps.Lib
       Uuid    func() string // new requirement
   }
   ```
2. Write a `<Field>Factory` for the new field on every adapter under [adapters/](../../adapters/), returning the closure, and assign its return value from that adapter's `New`, following the adapter specification located in [Specs.md](/docs/Specs.md):
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
4. Use the dependency from the library through `l.Deps.<Field>(...)`, following [HandleLibElements.md](/docs/HandleLibElements.md).
5. Build the project and run a sample — an unfilled field surfaces at runtime, not at build time:
   ```bash
   go build ./...
   go run ./libraryExamples/TrackSpendSample/TrackSpendSample.go
   ```

---

## The Contract

`sandbox/contracts/deps` declares what the library needs; `sandbox/contracts/api` declares what it hands back. Nothing else crosses the boundary.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now     func() time.Time
	Printf  func(format string, a ...any) (int, error)
	VerbLib verbdeps.Lib
	KeepLib keepdeps.Lib
}

// sandbox/contracts/api/api.go — what the library returns
type Lib struct {
	Deps        deps.Deps
	Sandboxmain func(args []string) int
	AddCategory func(name string) (Category, bool)
	AddSpend    func(category string, description string, amount int64) (Transaction, bool)
	Balance     func() int64
	// ... one field per library function
}
```

`lib.New(deps.Deps) api.Lib` is the single wiring point. Because both sides are plain structs, `sandbox` never imports an adapter, and callers never import `sandbox/internal/`. For why the contracts are structs rather than interfaces, see [StructContracts.md](/docs/StructContracts.md).

`Deps` is the *only* door in the sandbox wall: since nothing under `sandbox/` may import an adapter, a third-party module, or an OS-bound standard-library package, every effect the library performs has to be a field on this struct. That constraint is what this page's mechanic exists to serve — see [SandboxIsolation.md](/docs/SandboxIsolation.md).

---

## Using an Adapter

The simplest setup: an adapter builds a ready-to-use `deps.Deps`.

```go
package main

import (
	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// The standard adapter fills every field of deps.Deps
	myDeps := agnosadapter.New("trackerdata")
	l := agnoslib.New(myDeps)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)
	println(l.Balance()) // -8450
}
```

---

## Overwriting a Single Behavior

To keep an adapter but change one behavior, take the `deps.Deps` it returns and assign the field you want to replace. Every other field keeps the adapter's implementation — no embedding, no wrapper type.

This is exactly how you make a timestamp deterministic: keep the standard database, but replace `Now` so the clock is under your control.

```go
package main

import (
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Get the default implementation from an adapter
	myDeps := agnosadapter.New("trackerdata")

	// 2. Replace only the clock — KeepLib stays as the adapter built it
	now := time.Unix(0, 0)
	myDeps.Now = func() time.Time { return now }

	// 3. Inject — the lib sees a normal deps.Deps
	l := agnoslib.New(myDeps)
	l.AddCategory("groceries")

	// 4. Move the clock, and the next record is stamped with the new value
	now = time.Unix(120, 0)
	transaction, _ := l.AddSpend("groceries", "weekly shopping", 8450)
	println(transaction.OccurredAt.Unix()) // 120
}
```

> **Careful:** patch the `deps.Deps` value **before** calling `lib.New`. The library's factories close over the `api.Lib` they were run over, so assigning to `l.Deps.Now` on the returned struct changes nothing. Swapping the *variable* a field's closure captured — as `now` is swapped in step 4 above — works, because that happens inside the function the adapter's field already points at.

---

## Writing Custom Deps

For complete control, build the `deps.Deps` yourself. No adapter is involved: it is a struct literal, so there is no type to declare and no method set to satisfy.

```go
package main

import (
	"time"

	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
	agnosdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	agnoskeepdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
)

func main() {
	// 1. Build your own implementation. Now is a one-line closure; KeepLib
	//    is an embedded library, so filling it by hand means supplying the
	//    whole api — see "Injecting a Whole Library" below.
	frozen := time.Unix(0, 0)

	myDeps := agnosdeps.Deps{
		Now: func() time.Time { return frozen },
		KeepLib: agnoskeepdeps.Lib{
			NewDatabase: myOwnDatabase, // your implementation of the api
		},
		// Printf and VerbLib are left zero here: this program calls the
		// library functions directly, never Sandboxmain, so nothing
		// parses arguments or prints.
	}

	// 2. Inject it into the library
	l := agnoslib.New(myDeps)

	// 3. Use the library normally
	l.AddCategory("groceries")
	println(l.Balance())
}
```

> **Careful:** the compiler cannot tell you a field is missing. A `Deps` built by hand with an unfilled field holds a nil function that panics on first call — fill every field the library actually reaches. In practice, start from an adapter and patch what you need: `KeepLib` is a whole database api, and no program should reimplement it just to change the clock.

---

## Injecting a Whole Library

Two dependencies are not behaviors but other libraries built with this same pattern: [`VerbLib`](/docs/verbdeps.Lib.md), the embedded Verb argv parser, and [`KeepLib`](/docs/keepdeps.Lib.md), the embedded Keep schema database. Because a contract is a struct of function fields and not an interface, the whole library fits in one plain field — no getter, no bridging type:

```go
type Deps struct {
	Now     func() time.Time
	VerbLib verbdeps.Lib // an embedded library, injected whole
	KeepLib keepdeps.Lib // and another one — the tracker's whole storage
}
```

The sandbox cannot import Verb or Keep — that would break [SandboxIsolation.md](/docs/SandboxIsolation.md) — so it declares a copy of each api in `sandbox/contracts/deps/verbdeps/` and `sandbox/contracts/deps/keepdeps/`. The adapter, outside the sandbox, initializes the real library and assigns its fields onto that copy:

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

The factory returns a **value** instead of a closure, because the field is a struct rather than a function — the one shape variation the [factory pattern](/docs/RULES.md#factory-pattern) allows. Building a `Deps` by hand means either filling `VerbLib` the same way or leaving it zero, in which case nothing may call through it. The same mechanic, seen from the other side, is what `bootstrap/` demonstrates: there, *this* library is the one being embedded.

### When the embedded library hands back structs

Verb's fields return plain values, so copying is one assignment per field. Keep's do not: `Lib.NewDatabase` returns a `KeepDatabase`, whose `GetSchema` returns a `SchemaInstance`, whose `NewItem` returns a `SchemaItem`. Assigning those functions straight across would leak the embedded library's types into the sandbox, so the factory wraps each one in a closure that converts what it returns:

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

The conversion runs in both directions — `toKeepProps` on the way in, `fromKeepDatabase` on the way out — and recurses as far as the api tree goes. That work is the price of the wall, and it is paid once per adapter, outside the sandbox. It stays mechanical for the same reason the copy itself does: every type on both sides is a struct of plain fields.

---

## Propagation

`lib.New` delegates to the constructor in `sandbox/internal/lib/`, which stores the `Deps` on the `api.Lib` struct itself and then runs the factories over it:

```go
// sandbox/new.go
func New(d deps.Deps) api.Lib {
	return internallib.New(d)
}

// sandbox/internal/lib/lib.go
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.AddCategory = AddCategoryFactory(&l)
	l.AddSpend = AddSpendFactory(&l)
	l.Balance = BalanceFactory(&l)
	// ... one assignment per field
	return l
}
```

The carrier is the **closure**. Each factory returns a closure that reads `l.Deps` when the field is called, and `New` assigns it:

```go
// sandbox/internal/lib/lib.go
func GetCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		record, ok := store.FindCategory(l.Deps, name)
		if !ok {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}
```

Every object the lib creates receives the same `Deps`, passed into the object package's `New` constructor, which stores it on the object's own api struct before running that object's factories:

```go
// sandbox/internal/lib/lib.go — inside GetCategoryFactory's closure
c := category.New(l.Deps, record)
```

So a dependency injected once is reachable from anywhere in the object graph — that is why a `Category` can read the injected database in `ListTransactions`, and why each `Transaction` it hands back can delete itself through that same database.

```
standard.New() ──▶ deps.Deps ──▶ lib.New(deps) ──▶ api.Lib
                                                     │
                                        GetCategory() (propagates Deps)
                                                     ▼
                                               api.Category
                                                     │
                                   ListTransactions() (propagates Deps)
                                                     ▼
                                             api.Transaction
```
