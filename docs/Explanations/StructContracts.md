# Struct Contracts

## Description
Explains why every contract in this project — `deps.Deps` and everything in `sandbox/contracts/api` — is a **struct of function fields** instead of an interface, what that buys, and what it costs.

---

## The Shape

A contract is a struct whose fields are functions. The library declares the shape; whoever fills it decides the behavior.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now   func() time.Time
	Load  func(key string) (value string, expiresAtUnix int64, ok bool)
	Store func(key string, value string, expiresAtUnix int64)
}

// sandbox/contracts/api/api.go — what the library hands back
type Lib struct {
	Deps deps.Deps
	Set  func(key string, value string, ttlSeconds int)
	Get  func(key string) (Entry, bool)
}
```

Callers use both exactly as they would an interface — `l.Set("k", "v", 60)` reads the same whether `Set` is a method or a field holding a function. The difference only shows up at the wiring points.

The api struct carries its own `Deps`. That is what removes the need for a second, internal mirror type: the struct handed to the caller is the same struct the library's own code reads its dependencies from.

---

## Factories Fill the Fields

There is no internal type and no method set. `sandbox/internal/` holds **factories**: functions that take a pointer to an api struct and return a closure for one of its function fields; the caller assigns the result.

```go
// sandbox/internal/lib/lib.go
func SetFactory(l *api.Lib) func(key string, value string, ttlSeconds int) {
	return func(key string, value string, ttlSeconds int) {
		expiresAt := l.Deps.Now().Add(time.Duration(ttlSeconds) * time.Second)
		l.Deps.Store(key, value, expiresAt.Unix())
	}
}

// New is the package's constructor and its factory aggregate in one — the
// one place that must stay complete.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.Set = SetFactory(&l)
	l.Get = GetFactory(&l)
	return l
}
```

```go
// sandbox/new.go — the whole wiring
func New(d deps.Deps) api.Lib {
	return internallib.New(d)
}
```

The closure captures `l`, so `l.Deps.Now()` is resolved when `Set` is *called*, not when the factory ran. That is what carries the injected deps into behavior the caller can hold, while `sandbox/internal/` stays unreachable from outside.

Two properties follow, and both are load-bearing:

- **One field, one factory.** A factory does nothing but return the closure; the package's `New` constructor is the only place that has to assign them all.
- **Deps are read through the pointer, never copied into the closure.** Capturing `d` directly would freeze the dependency at construction; reading `l.Deps` keeps the struct authoritative.

---

## Adapters Fill Their Contract the Same Way

The pattern does not stop at the sandbox wall. An adapter fills `deps.Deps` with factories too — only the **carrier** changes: instead of an api struct holding the deps it reads, the carrier is the adapter struct, holding the configuration its closures read and declaring the contract they fill.

```go
// adapters/standard/standard.go
type StandardAdapter struct {
	Deps     deps.Deps // the contract the factories assign into
	mu       sync.RWMutex
	filePath string // the state the closures read
}

func NowFactory(s *StandardAdapter) func() time.Time {
	return func() time.Time { return time.Now() }
}

// New is the adapter's factory aggregate: it returns the contract struct,
// never the concrete adapter type.
func New(filePath string) deps.Deps {
	adapter := &StandardAdapter{filePath: filePath}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Load = LoadFactory(adapter)
	adapter.Deps.Store = StoreFactory(adapter)
	return adapter.Deps
}
```

Binding a method into a field would work in Go, but the project does not do it: one shape for filling struct contracts means one place to look for completeness — the `New` at the bottom of the file — on both sides of the wall. The rule is binding; see [RULES.md](/docs/References/RULES.md#factory-pattern) and the [Factories](/docs/References/Meta/Factories/Specs.md) specification.

A field that is not a function has its factory return a value rather than a closure — that is how `bootstrap/adapters/standard/` hands the whole embedded library over in `CacheLibFactory`.

---

## Replacing One Behavior

With an interface, overriding a single method means declaring a type that embeds the original and shadows one method. With a struct, it is an assignment:

```go
myDeps := standard.New("data.json")

// Keep the adapter's Load and Store; control the clock
myDeps.Now = func() time.Time { return time.Unix(0, 0) }

l := lib.New(myDeps)
```

This is the everyday testing path: freeze time, stub one call, leave the rest of the adapter alone.

---

## Consuming a Library That Uses This Pattern

The payoff shows when one Agnos-style library depends on another. The consuming sandbox may not import the embedded library — that would be a third-party import — so it restates the shape it needs in its own contracts, and the adapter fills it.

```go
// bootstrap/sandbox/contracts/deps/agnosdeps/agnosdeps.go — inside the sandbox
type Lib struct {
	Set func(key string, value string, ttlSeconds int)
	Get func(key string) (Entry, bool)
}

// bootstrap/sandbox/contracts/deps/deps.go — the embedded lib is one plain field
type Deps struct {
	Println  func(a ...any)
	CacheLib agnosdeps.Lib
}
```

```go
// bootstrap/adapters/standard/standard.go — outside the sandbox
func CacheLibFactory(s *StandardAdapter) agnosdeps.Lib {
	inner := agnoslib.New(agnosadapter.New(s.cacheFilePath))

	return agnosdeps.Lib{
		Set: inner.Set, // identical signature: assigned straight across
		Get: func(key string) (agnosdeps.Entry, bool) {
			// only Get needs a wrapper, because Entry is a distinct named type on each side
			entry, found := inner.Get(key)
			if !found {
				return agnosdeps.Entry{}, false
			}
			return agnosdeps.Entry{Value: entry.Value, ExpiresAt: entry.ExpiresAt, IsExpired: entry.IsExpired}, true
		},
	}
}
// s.Deps.CacheLib = CacheLibFactory(s) — assigned by New, like every other factory
```

Every field whose signature already matches is copied with one line. With interfaces, the same wiring needs a bridging type declaring every method, even the ones that pass straight through. The `bootstrap/` tree is the worked example — see [Structure.md](/docs/References/Structure.md#bootstrap).

---

## What It Costs

The compiler no longer checks completeness. An interface implementation that misses a method fails to build; a struct contract with an unfilled field compiles fine and panics on the first call with a nil-pointer dereference.

That moves one guarantee from the compiler to the author:

- An adapter's `New` must call a factory for **every** field of `deps.Deps`.
- A package's `New` constructor must call **every** field factory of its api struct.
- Adding a field to a contract means visiting every adapter — see [AddDependency.md](/docs/Tutorials/AddDependency.md).

There is a second cost, specific to factories: **the `Deps` field is read-only once the struct is returned.** The closures capture the struct the factories ran over, so patching a copy has no effect on behavior.

```go
l := agnoslib.New(myDeps)
l.Deps.Now = func() time.Time { return time.Unix(0, 0) } // does nothing

myDeps.Now = func() time.Time { return time.Unix(0, 0) } // patch here instead
l = agnoslib.New(myDeps)
```

Replace behavior on the `deps.Deps` value **before** calling `lib.New` — that path is shown above and in [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

In exchange, there is no partial-implementation ambiguity at the call site: a filled contract is a value that can be copied, patched field by field, and passed on.
