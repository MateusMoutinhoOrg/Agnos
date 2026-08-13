# Struct Contracts

## Description
Explains why every contract in this project — `deps.Deps` and everything in `sandbox/contracts/api` — is a **struct of function fields** instead of an interface, what that buys, and what it costs. Second stage of the Development learning path, after [SandboxIsolation.md](/docs/References/SandboxIsolation.md); [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md) shows the mechanic in motion.

---

## The Shape

A contract is a struct whose fields are functions. The library declares the shape; whoever fills it decides the behavior.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Now       func() time.Time
	VerbLib   verbdeps.Lib
	KeepLib   keepdeps.Lib
	EmbedDeps embeddeps.Lib
}

// sandbox/contracts/api/api.go — what the library hands back
type Lib struct {
	Deps        deps.Deps
	AddCategory func(name string) (Category, bool)
	Balance     func() int64
	// ... one field per library function
}
```

Callers use both exactly as they would an interface — `l.AddCategory("groceries")` reads the same either way. The api struct carries its own `Deps`, which removes the need for an internal mirror type: the struct handed to the caller is the same struct the library reads its dependencies from.

---

## Factories Fill the Fields

`sandbox/` holds **factories**: functions taking a pointer to an api struct and returning a closure for one of its fields. The package's `New` constructor assigns them all:

```go
// sandbox/lib/lib.go
func GetCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		record, ok := store.FindCategory(l.Deps, name) // l.Deps read at call time
		if !ok {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}

// New is the constructor and the factory aggregate in one —
// the one place that must stay complete.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.AddCategory = AddCategoryFactory(&l)
	l.GetCategory = GetCategoryFactory(&l)
	l.Balance = BalanceFactory(&l)
	// ... one assignment per field
	return l
}
```

Two properties are load-bearing:

- **One field, one factory.** A factory does nothing but return the closure; `New` is the only place that has to assign them all.
- **Deps are read through the pointer, never copied into the closure.** Capturing `d` directly would freeze the dependency at construction; reading `l.Deps` keeps the struct authoritative.

---

## Adapters Fill Their Contract the Same Way

Only the **carrier** changes: the adapter struct holds the configuration its closures read and declares the contract they fill.

```go
// adapters/standard/standard.go
type StandardAdapter struct {
	Deps         deps.Deps // the contract the factories assign into
	args         []string  // the state the closures read
	keepBasePath string
}

func NowFactory(s *StandardAdapter) func() time.Time {
	return func() time.Time { return time.Now() }
}

// New returns the contract struct, never the concrete adapter type.
func New(basePath string) deps.Deps {
	adapter := &StandardAdapter{args: os.Args[1:], keepBasePath: basePath}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.VerbLib = VerbLibFactory(adapter)
	adapter.Deps.KeepLib = KeepLibFactory(adapter)
	adapter.Deps.EmbedDeps = EmbedDepsFactory(adapter)
	return adapter.Deps
}
```

Binding a method into a field would work in Go, but the project forbids it: one shape for filling contracts means one place to check completeness — the `New` at the bottom of the file — on both sides of the wall. The rule is binding; see [RULES.md](/docs/References/RULES.md#factory-pattern) and the [Factories](/docs/References/Specs/Factories/Specs.md) specification. A field that is not a function has its factory return a value rather than a closure.

---

## Replacing One Behavior

With an interface, overriding one method means declaring a wrapper type. With a struct, it is an assignment — the everyday testing path:

```go
myDeps := standard.New("trackerdata")
myDeps.Now = func() time.Time { return time.Unix(0, 0) } // control the clock
l := lib.New(myDeps) // everything else keeps the adapter's implementation
```

---


## What It Costs

**The compiler no longer checks completeness.** An unfilled field compiles fine and panics on the first call with a nil dereference. That moves one guarantee from the compiler to the author:

- An adapter's `New` must call a factory for **every** field of `deps.Deps`.
- A package's `New` must call **every** field factory of its api struct.
- Adding a contract field means visiting every adapter — see [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).

**`Deps` is read-only once the struct is returned.** The closures captured the struct the factories ran over, so patch the `deps.Deps` value *before* calling `lib.New`:

```go
l.Deps.Now = func() time.Time { return time.Unix(0, 0) } // does nothing

myDeps.Now = func() time.Time { return time.Unix(0, 0) } // patch here instead
l = agnoslib.New(myDeps)
```

In exchange, there is no partial-implementation ambiguity at the call site: a filled contract is a value that can be copied, patched field by field, and passed on.
