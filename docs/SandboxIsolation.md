# Sandbox Isolation

## Description
Explains why the library lives in `sandbox/` and what "closed sandbox" means in practice: the library reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package — so everything it can do is exactly what the injected `Deps` allows.

---

## The Three Trees

The project is split into three top-level directories, and the arrows only point one way:

```
adapters/  ──▶  sandbox/  ◀──  libraryExamples/
(reaches the OS)  (closed)     (wires the two together)
```

- `sandbox/` is the library. It is closed: it imports only itself and OS-independent standard-library packages (`time`, `strings`, `errors`, …).
- `adapters/` is outside the wall. It is the only place `os`, `net`, a database driver, or any third-party module may appear.
- `libraryExamples/` is outside the wall too, and is the only place an adapter and the sandbox are named in the same file.

The split is what makes the library OS-independent: nothing inside `sandbox/` can be affected by which operating system, filesystem, or network the program runs on, because it has no way to reach any of them.

---

## What the Wall Forbids

A file under `sandbox/` may not import:

| Forbidden | Why |
|-----------|-----|
| `adapters/…` | The library would bind itself to one concrete implementation, and injection would be pointless. |
| `cmd/…`, `libraryExamples/…` | The binary and the samples are consumers of the library, never part of it. |
| Any third-party module | A dependency the caller cannot replace is a dependency the caller cannot test around. |
| OS-bound stdlib (`os`, `net`, `os/exec`, `syscall`, …) | The effect belongs in an adapter, reached through a `Deps` field. |

Everything the library needs from the outside world is declared as a function field on `Deps`:

```go
// sandbox/contracts/deps/deps.go — the only door in the wall
type Deps struct {
	Now     func() time.Time // instead of time.Now()
	Printf  func(format string, a ...any) (int, error) // instead of fmt.Printf
	VerbLib verbdeps.Lib     // instead of importing the Verb library
	KeepLib keepdeps.Lib     // instead of importing the Keep library, and
	                         // with it every os.ReadFile the tracker needs
}
```

This is what lets the *whole command-line interface* live inside the wall. `api.Lib.Sandboxmain` reads the command line through `Deps.VerbLib` and writes every line through `Deps.Printf`, so it never touches `os.Args` or `os.Stdout` itself — the binary in `cmd/main/` hands it an argument vector and exits with its return, and a test can hand it a fixed vector and a buffer instead.

`VerbLib` and `KeepLib` are the same door in a different shape: the Verb argv parser and the Keep schema database are third-party modules to this sandbox, so instead of importing them the sandbox declares a copy of each api in `sandbox/contracts/deps/verbdeps/` and `sandbox/contracts/deps/keepdeps/` and lets the adapter fill it. Restating the shape is cheap because it is a struct of function fields — see [StructContracts.md](/docs/StructContracts.md) and [DepsMechanic.md](/docs/DepsMechanic.md#injecting-a-whole-library).

Inside the sandbox, the same behaviors are reached only through `l.Deps`:

```go
// sandbox/internal/lib/lib.go — no os, no net, no third party
func AddCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		categories, ok := store.Categories(l.Deps) // reaches l.Deps.KeepLib
		if !ok {
			return api.Category{}, false
		}
		record, err := categories.NewItem(map[string]any{
			store.NameField:      name,
			store.CreatedAtField: l.Deps.Now().Unix(),
		})
		if err != nil {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}
```

To add a new door, follow [AddDependency.md](/docs/AddDependency.md).

---

## What the Wall Forbids in the Other Direction

The wall is not only about what the sandbox imports — it also limits what the outside may reach into. `sandbox/internal/` holds the **factories** that fill the contract structs, and it is protected by Go's `internal/` rule: only packages rooted at `sandbox/` can import it. An adapter or a consumer that tries gets a compile error, not a convention warning.

Note what this does *not* hide. Because the factories fill the api structs directly, the structs the caller holds are the same ones the library works on — including their `Deps` field. What stays unreachable is the logic: a consumer can read `l.Deps`, but cannot call, replace, or re-run the factories that turned it into behavior.

So the outside world sees exactly three things:

| Package | Who imports it | For what |
|---------|----------------|----------|
| `sandbox` (package `lib`) | consumers, examples | `lib.New(deps) api.Lib` — the single wiring point |
| `sandbox/contracts/deps` | adapters, consumers | the contract struct to fill |
| `sandbox/contracts/api` | consumers, examples | the structs handed back |

Everything else in `sandbox/` is unreachable, which is why the factories can be renamed, split, or restructured without breaking a single consumer.

---

## Why the Entry Point Lives Inside

`sandbox/new.go` is the one file in the sandbox that consumers import directly, and it stays inside the wall because it obeys the same rule — it names no adapter:

```go
// sandbox/new.go
func New(d deps.Deps) api.Lib {
	return internallib.New(d)
}
```

It accepts a contract struct and returns a contract struct. The caller decides which implementation fills the fields flowing in, so the sandbox never learns what is behind the contract:

```go
import (
	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

// This line is in libraryExamples/, outside the wall — the only place
// an adapter and the sandbox meet.
l := agnoslib.New(agnosadapter.New("data.json"))
```

For how the injected value then travels through the object graph, see [DepsMechanic.md](/docs/DepsMechanic.md). For why the contracts are structs rather than interfaces, see [StructContracts.md](/docs/StructContracts.md).
