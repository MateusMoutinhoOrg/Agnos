# Dependency Mechanics

## Description
Explains how the library receives its dependencies. `Deps` is a **struct of function fields**: the library declares the behaviors it needs as fields, and the caller supplies a value with every field filled — a pre-built adapter, an adapter with one field replaced, or a struct literal written by hand.

---

## The Contract

`sandbox/contracts/deps` declares what the library needs; `sandbox/contracts/api` declares what it hands back. Nothing else crosses the boundary.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now   func() time.Time
	Load  func(key string) (value string, expiresAtUnix int64, ok bool)
	Store func(key string, value string, expiresAtUnix int64)
}

// sandbox/contracts/api/api.go — what the library returns
type Lib struct {
	Deps deps.Deps
	Set  func(key string, value string, ttlSeconds int)
	Get  func(key string) (Entry, bool)
}
```

`lib.New(deps.Deps) api.Lib` is the single wiring point. Because both sides are plain structs, `sandbox` never imports an adapter, and callers never import `sandbox/internal/`. For why the contracts are structs rather than interfaces, see [StructContracts.md](/docs/Explanations/StructContracts.md).

`Deps` is the *only* door in the sandbox wall: since nothing under `sandbox/` may import an adapter, a third-party module, or an OS-bound standard-library package, every effect the library performs has to be a field on this struct. That constraint is what this page's mechanic exists to serve — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

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
	myDeps := agnosadapter.New("data.json")
	l := agnoslib.New(myDeps)

	l.Set("greeting", "hello world", 60)
	if entry, ok := l.Get("greeting"); ok {
		println(entry.Value)
	}
}
```

---

## Overwriting a Single Behavior

To keep an adapter but change one behavior, take the `deps.Deps` it returns and assign the field you want to replace. Every other field keeps the adapter's implementation — no embedding, no wrapper type.

This is exactly how you test expiry without waiting: keep the standard store, but replace `Now` so the clock is under your control.

```go
package main

import (
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Get the default implementation from an adapter
	myDeps := agnosadapter.New("data.json")

	// 2. Replace only the clock — Load and Store stay as the adapter built them
	now := time.Unix(0, 0)
	myDeps.Now = func() time.Time { return now }

	// 3. Inject — the lib sees a normal deps.Deps
	l := agnoslib.New(myDeps)
	l.Set("k", "v", 60) // expires 60s after time 0

	// 4. Jump the clock past expiry — no real waiting needed
	now = time.Unix(120, 0)
	_, ok := l.Get("k")
	println(ok) // false — expired
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
)

func main() {
	// 1. Build your own implementation, keeping records in a plain map
	store := map[string]string{}

	myDeps := agnosdeps.Deps{
		Now: time.Now,
		Load: func(key string) (string, int64, bool) {
			v, ok := store[key]
			return v, time.Now().Add(time.Hour).Unix(), ok
		},
		Store: func(key, value string, expiresAtUnix int64) {
			store[key] = value
		},
		// VerbLib is left zero here: nothing in this program parses
		// arguments. See "Injecting a Whole Library" below.
	}

	// 2. Inject it into the library
	l := agnoslib.New(myDeps)

	// 3. Use the library normally
	l.Set("k", "v", 60)
	if entry, ok := l.Get("k"); ok {
		println(entry.Value)
	}
}
```

> **Careful:** the compiler cannot tell you a field is missing. A `Deps` built by hand with an unfilled field holds a nil function that panics on first call — fill every field.

---

## Injecting a Whole Library

One dependency is not a behavior but another library built with this same pattern: [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), the embedded Verb argv parser. Because a contract is a struct of function fields and not an interface, the whole library fits in one plain field — no getter, no bridging type:

```go
type Deps struct {
	Now     func() time.Time
	Load    func(key string) (value string, expiresAtUnix int64, ok bool)
	Store   func(key string, value string, expiresAtUnix int64)
	VerbLib verbdeps.Lib // the embedded library, injected whole
}
```

The sandbox cannot import Verb — that would break [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) — so it declares a copy of Verb's api in `sandbox/contracts/deps/verbdeps/`. The adapter, outside the sandbox, initializes the real library and assigns its fields onto that copy:

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

The factory returns a **value** instead of a closure, because the field is a struct rather than a function — the one shape variation the [factory pattern](/docs/References/RULES.md#factory-pattern) allows. Building a `Deps` by hand means either filling `VerbLib` the same way or leaving it zero, in which case nothing may call through it. The same mechanic, seen from the other side, is what `bootstrap/` demonstrates: there, *this* library is the one being embedded.

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
	l.Set = SetFactory(&l)
	l.Get = GetFactory(&l)
	return l
}
```

The carrier is the **closure**. Each factory returns a closure that reads `l.Deps` when the field is called, and `New` assigns it:

```go
// sandbox/internal/lib/lib.go
func SetFactory(l *api.Lib) func(key string, value string, ttlSeconds int) {
	return func(key string, value string, ttlSeconds int) {
		expiresAt := l.Deps.Now().Add(time.Duration(ttlSeconds) * time.Second)
		l.Deps.Store(key, value, expiresAt.Unix())
	}
}
```

Every object the lib creates receives the same `Deps`, passed into the object package's `New` constructor, which stores it on the object's own api struct before running that object's factories:

```go
// sandbox/internal/lib/lib.go — inside GetFactory's closure
e := entry.New(l.Deps, value, expiresAtUnix)
```

So a dependency injected once is reachable from anywhere in the object graph — that is why an `Entry` can consult the injected clock in `IsExpired`.

```
standard.New() ──▶ deps.Deps ──▶ lib.New(deps) ──▶ api.Lib
                                                     │
                                             Get() (propagates Deps)
                                                     ▼
                                                api.Entry
```

To add a new behavior to the contract, follow [AddDependency.md](/docs/Tutorials/AddDependency.md).
