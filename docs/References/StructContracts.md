# Struct Contracts

## Description
Explains why every contract in this project — `deps.Deps` and its sub-contracts, `api.Sandbox` and its `Actions` and `Cli` — is a **struct of function fields** instead of an interface, how binders fill them on both sides of the sandbox wall, and what that costs. It assumes [SandboxIsolation.md](/docs/References/SandboxIsolation.md); [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md) shows the mechanic in motion.

---

## The Shape

A contract is a struct whose fields are functions. The sandbox declares the shape; whoever fills it decides the behavior.

```go
// sandbox/deps/std/std.go — what the sandbox needs from the runtime
type Lib struct {
	Now    func() time.Time
	Printf func(format string, a ...any) (n int, err error)
	Log    func(format string, a ...any) (n int, err error)
	Error  func(format string, a ...any) (n int, err error)
	Errorf func(format string, a ...any) error
}

// sandbox/api/actions.go — what the sandbox hands back
type Actions struct {
	Build      func(props BuildProps) error
	Verify     func(path string) error
	Start      func(props StartProps) error
	DepInstall func(path string, dep string) error
	// … one field per action
}
```

Callers use both exactly as they would an interface — `lib.Actions.Verify("./my-tool")` reads the same either way. `deps.Deps` itself carries **only** sub-contract structs, never a bare function: the loose runtime functions are gathered into `std.Lib`, so every field of `Deps` is the same kind of thing and a template can emit the struct from a directory listing.

---

## Binders Fill the Fields

On the sandbox side, one package per `api` file holds one **binder** that assigns the internal implementation onto the contract. The generated constructor runs them all:

```go
// sandbox/binds/actions.go — the one place every action is assigned
func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Verify = func(path string) error {
		return verifyAction.Verify(deps, path)
	}
	sandbox.Actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(deps, props)
	}
	// … one assignment per field
}

// sandbox/new.go — generated: one binds.<X>Bind call per file of sandbox/binds/
func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	binds.ActionsBind(deps, &self)
	binds.CliBind(deps, &self)
	return &self
}
```

Two properties are load-bearing:

- **One contract file, one binder file, same name.** `api/actions.go` ↔ `binds/actions.go`, `api/cli.go` ↔ `binds/cli.go`. `verify` enforces the mirror and that a binder declares nothing but functions, so completeness is checked in one place per contract.
- **Deps are read through the pointer, never copied.** The closures capture `deps *deps.Deps`, so a field reassigned on that struct before `New` — or through the same pointer after — is what the sandbox sees.

---

## Adapters Fill Their Contract the Same Way

Outside the wall, the shape is identical: one package per implementation exporting a single binder named `Bind`, and a generated assembly running them all:

```go
// adapters/libs/rundeps/Run.go
func Bind(deps *deps.Deps) {
	deps.Rundeps = rundeps.Lib{Run: run}
}

// adapters/availables/standard/new.go — generated: one <lib>.Bind per dir of adapters/libs/
func New() deps.Deps {
	deps := deps.Deps{}
	embeddeps.Bind(&deps)
	iodeps.Bind(&deps)
	rundeps.Bind(&deps)
	std.Bind(&deps)
	verb.Bind(&deps)
	// …
	return deps
}
```

Binding a method into a field would work in Go, but the project forbids it: one uniform name (`Bind`) and one shape mean a template can generate the assembly from a directory listing, which is exactly what `CollectAdapterLibs` does. A sub-contract whose real object is created per call — an HTTP request, an argv parser — exposes that constructor as a field of its `Lib` (`requestdeps.Lib.NewRequest`, `argvdeps.Lib.New`) and is bound the same way.

---

## Replacing One Behavior

With an interface, overriding one method means declaring a wrapper type. With a struct, it is an assignment — the everyday testing path:

```go
deps := standard.New()
deps.Std.Printf = func(format string, a ...any) (int, error) {
	return fmt.Fprintf(&buffer, format, a...) // capture the CLI's output
}
lib := sandbox.New(&deps) // everything else keeps the adapter's implementation
```

Composing a whole assembly of your own is [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md).

---

## What It Costs

**The compiler no longer checks completeness.** An unfilled field compiles fine and panics on the first call with a nil dereference. That moves one guarantee from the compiler to the author — and, in this project, to the generator:

- `standard.New` binds **every** lib under `adapters/libs/`, because it is generated from the listing; a hand-written assembly must be compared against `deps.Deps` field by field.
- Every binder in `sandbox/binds/` must assign **every** field of its `api` struct; adding an action means visiting both files in the same commit — see [AddAction.md](/docs/Tutorials/AddAction.md).

**Field names are not a choice.** `Deps.Iodeps`, `Deps.Std`, `Deps.Argvdeps` are the title-cased directory names, regenerated on every build; code spelling them otherwise stops compiling at the next `agnos build`.

In exchange, there is no partial-implementation ambiguity at the call site: a filled contract is a value that can be copied, patched field by field, passed on — and, above all, **generated**.
