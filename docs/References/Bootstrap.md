# Bootstrap: Embedding an Agnos-Cli-Style Library

## Description
Explains how a library built with this pattern consumes **another** Agnos-Cli-style library as a dependency. The [/bootstrap/](/bootstrap/) tree is the worked example: a second, self-contained Agnos-Cli library whose only real dependency is the root financial tracker, injected through its `Deps` like any other behavior.

---

## The Bootstrap Tree

`bootstrap/` is a complete Agnos-Cli library in miniature — the same trees and the same rules as the root, minus the CLI:

| Path | Description |
|------|-------------|
| `sandbox/new.go` | The wiring point: `New(deps.Deps) api.Lib`, delegating to `sandbox/internal/lib` |
| `sandbox/contracts/deps/deps.go` | The `Deps` contract: `Println` plus `TrackerLib`, the embedded library held as one plain struct field |
| `sandbox/contracts/deps/agnosdeps/agnosdeps.go` | The sandbox's **copy** of the embedded library's api structs (`Lib`, `Category`, `Transaction` and the kind constants) |
| `sandbox/contracts/api/api.go` | The output struct: `api.Lib` with `Deps` and one demonstration field, `TestFunc` |
| `sandbox/internal/lib/lib.go` | `TestFuncFactory` and the `New` constructor running it |
| `adapters/standard/standard.go` | Fills the contract; its `TrackerLibFactory` is where the embedded library is initialized and copied across |
| `libraryExamples/Test/test.go` | Self-contained `package main` wiring the bootstrap adapter into the bootstrap lib |

Run the demonstration from the repository root:

```bash
go run ./bootstrap/libraryExamples/Test/test.go   # writes bootstrap-trackerdata/, gitignored
```

---

## The Sandbox Wall Applies to Libraries Too

The [sandbox isolation rule](/docs/References/SandboxIsolation.md) makes no exception for a dependency that happens to be another Agnos-Cli-style library: importing it from `bootstrap/sandbox/` would be a third-party import. So the consuming sandbox **restates the shape it needs** as its own contract — the `agnosdeps` package is a field-for-field copy of the embedded library's api, declared inside the consuming sandbox and owned by it:

```go
// bootstrap/sandbox/contracts/deps/agnosdeps/agnosdeps.go — inside the sandbox
type Lib struct {
	AddCategory func(name string) (Category, bool) // returns the *local* Category copy
	Balance     func() int64
	// ... one field per embedded api field the sandbox needs
}
```

The copy deliberately omits the embedded library's `Deps` field: the sandbox has no business reading the dependencies the embedded library was wired with.

The root library uses the identical mechanic in the other direction — `sandbox/contracts/deps/verbdeps/` and `keepdeps/` are its copies of the embedded Verb and Keep libraries. See [Structure.md](/docs/References/Structure.md#bootstrap).

---

## Declaring the Dependency

Because contracts are [structs of function fields](/docs/References/StructContracts.md), a whole library is just one more field — a plain struct, never a getter function:

```go
// bootstrap/sandbox/contracts/deps/deps.go
type Deps struct {
	Println    func(a ...any)
	TrackerLib agnosdeps.Lib // the embedded library, already initialized by the adapter
}
```

This is a rule of the Deps specification — locate it through [Specs.md](/docs/References/Specs.md).

---

## Filling It from the Adapter

The adapter, outside the sandbox, is the only place allowed to import the embedded library. Its factory initializes the embedded library with the embedded library's **own** adapter, then copies the api onto the local `agnosdeps` structs. Because the field is a struct rather than a function, this factory returns a **value**, not a closure — the one case the [Factories specification](/docs/References/Specs/Factories/Specs.md) treats differently:

```go
// bootstrap/adapters/standard/standard.go — outside the sandbox
func TrackerLibFactory(s *StandardAdapter) agnosdeps.Lib {
	inner := agnoslib.New(agnosadapter.New(s.trackerBasePath))
	return agnosdeps.Lib{
		Balance: inner.Balance, // identical signature: assigned straight across
		AddCategory: func(name string) (agnosdeps.Category, bool) {
			category, ok := inner.AddCategory(name) // wrapper only where a
			return fromCategory(category), ok       // named type differs
		},
		// ...
	}
}
```

Two cases, and only two:

- **Signatures match** — plain assignment, one line per field (`inner.Balance`).
- **A named type differs** — the embedded field hands back an api struct of its own, so the adapter wraps it in a closure converting the returned struct too (`fromCategory`, `fromTransaction`), recursively for structs those structs return.

The conversion is the price of the wall, paid once per adapter — see [What It Costs](/docs/References/StructContracts.md#consuming-a-library-that-uses-this-pattern) and [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).

---

## Reaching It from the Sandbox

Inside the sandbox the embedded library is indistinguishable from any other dependency — a field on `l.Deps`, called without any import:

```go
// bootstrap/sandbox/internal/lib/lib.go
func TestFuncFactory(l *api.Lib) func() {
	return func() {
		tracker := l.Deps.TrackerLib // the whole embedded library

		category, created := tracker.AddCategory("groceries")
		if !created {
			l.Deps.Println("groceries: could not create the category")
			return
		}
		category.AddSpend("weekly shopping", 8450) // cents, per the money convention

		l.Deps.Println("balance:", tracker.Balance())
	}
}
```

Swapping the embedded library for a fake in a test is the usual one-line patch on the `deps.Deps` value before `lib.New` — see [Replacing One Behavior](/docs/References/StructContracts.md#replacing-one-behavior).

---

## Embedding This Library in Your Own

To depend on the Agnos-Cli tracker (or any Agnos-Cli-style library) from a library of your own, reproduce the three pieces the bootstrap tree demonstrates:

| Piece | Bootstrap reference |
|-------|---------------------|
| A local copy of the embedded api, inside your sandbox | `bootstrap/sandbox/contracts/deps/agnosdeps/` |
| A plain struct field for it on your `Deps` contract | `TrackerLib` in `bootstrap/sandbox/contracts/deps/deps.go` |
| A value-returning factory in every adapter, converting where types differ | `TrackerLibFactory` in `bootstrap/adapters/standard/` |

Consumers of *your* library then follow the [Import Aliases rule](/docs/References/RULES.md#import-aliases) with your own prefix (`bootstrapadapter`, `bootstraplib`), exactly as `bootstrap/libraryExamples/Test/test.go` does. For turning this repository into that library in the first place, follow [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) or [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) — [TemplateFileActions.md](/docs/References/TemplateFileActions.md) lists the fate of every `bootstrap/*` file there.
