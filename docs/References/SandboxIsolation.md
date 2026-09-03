# Sandbox Isolation

## Description
Explains what "closed sandbox" means in practice, for Agnos-Cli and for every project it generates: the code in `sandbox/` reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package — so everything it can do is exactly what the injected `Deps` allows. `verify` enforces the wall mechanically. [StructContracts.md](/docs/References/StructContracts.md) explains how the contracts crossing it are shaped.

---

## The Three Trees

Three top-level directories, and the arrows only point one way:

```
adapters/  ──▶  sandbox/  ◀──  cmd/
(reaches the OS)  (closed)     (wires the two together)
```

- `sandbox/` is the library and, with it, the whole command-line interface. Closed: it imports only itself and OS-independent stdlib packages (`strings`, `sort`, `strconv`, …).
- `adapters/` is the only place `os`, `net/http`, `os/exec`, `embed`, a parser library or any third-party module may appear.
- `cmd/` is the only place an adapter and the sandbox are named in the same file.

`assets/` sits outside the wall too: files compiled into the binary, reached only through `deps.Embeddeps`.

Nothing inside `sandbox/` can be affected by which OS, filesystem, or network the program runs on, because it has no way to reach any of them.

---

## What the Wall Forbids

A file under `sandbox/` may not import:

| Forbidden | Why |
|-----------|-----|
| `adapters/…`, `cmd/…`, `assets/…` | Binding to one concrete implementation makes injection pointless. `verify`: *imports … which is outside sandbox/*. |
| Any third-party module | A dependency the caller cannot replace is one the caller cannot test around. |
| OS-bound stdlib (`os`, `net`, `os/exec`, `syscall`, `fmt` for printing, …) | The effect belongs in an adapter, reached through a `Deps` field. |
| `embed`, and the `//go:embed` directive | Compiling a file into the binary is a build-time, filesystem-bound effect; the bytes arrive through `deps.Embeddeps` like any other. |

Two packages are held to a stricter rule still: `sandbox/api/*` imports **nothing** but other `sandbox/api` packages — not even the standard library, not even `sandbox/deps` — because it is pure contract; and `sandbox/deps/*` imports only the standard library and other `sandbox/deps` packages, because each of them is a restated copy of an outside library's api.

Everything the sandbox needs from the outside world is a field on `Deps` — the only door in the wall — and every field is a whole sub-contract struct, never a bare function:

```go
// sandbox/deps/deps.go — generated from the sandbox/deps/<x>/ directories
type Deps struct {
	Argvdeps      argvdeps.Lib      // instead of importing the Verb argv parser
	Embeddeps     embeddeps.Lib     // instead of //go:embed and os.ReadFile
	Goimportsdeps goimportsdeps.Lib // instead of go/parser
	Iodeps        iodeps.Lib        // instead of os and path/filepath
	Requestdeps   requestdeps.Lib   // instead of net/http
	Rundeps       rundeps.Lib       // instead of os/exec
	Serializables serializables.Lib // instead of a YAML/JSON module
	Std           std.Lib           // instead of time.Now, fmt.Printf, os.Stderr
}
```

This is what lets the whole command-line interface live inside the wall: the generated `climain.go` reads the command line through `deps.Argvdeps` and prints through `deps.Std`, never touching `os.Args` or `os.Stdout`. The binary hands it an argument vector; a test hands it a fixed vector and a buffer instead.

---

## Why Every Door Is a Copy

The sandbox cannot import the Verb parser, `go/parser` or `net/http`, so for each it declares a **copy** of the api it needs under `sandbox/deps/<x>/` — a struct of function fields restating the shape, field for field — and the adapter lib under `adapters/libs/<x>/`, which lives outside the wall, initializes the real library and assigns its fields onto that copy:

```go
// adapters/libs/verb/Verb.go — outside the wall
func Bind(deps *deps.Deps) {
	deps.Argvdeps.New = func(args []string) argvdeps.Parser {
		inner := verblib.New(args)            // the real, third-party library
		return argvdeps.Parser{Args: inner.Args, IsPresent: inner.IsPresent /* … */}
	}
}
```

Copying is cheap precisely because the copied libraries expose structs of function fields rather than interfaces. The same mechanic ships as the installable deps of [DepList.md](/docs/References/DepList.md), each one a contract-plus-adapter pair.

---

## What the Wall Forbids in the Other Direction

`sandbox/internal/` holds every implementation — the actions, the generated CLI, the parsers, SmartIO — and Go's `internal/` rule makes it unreachable from outside `sandbox/`: an adapter or consumer that tries gets a compile error. The outside world sees exactly three packages:

| Package | Who imports it | For what |
|---------|----------------|----------|
| `sandbox` | `cmd/`, consumers | `sandbox.New(&deps) *api.Sandbox` — the single wiring point |
| `sandbox/deps` and `sandbox/deps/<x>` | adapters, consumers | the contract struct to fill |
| `sandbox/api` | consumers | the structs handed back — `Sandbox`, `Actions`, `Cli`, the props |

`sandbox/binds/` is the seam between the two: each file mirrors an `api` file and holds the one function that assigns the internal implementation onto the contract. `verify` requires the mirror and forbids anything but functions there.

---

## Why the Entry Point Lives Inside

`sandbox/new.go` — itself generated — stays inside the wall because it obeys the same rule: it names no adapter, accepting a contract struct and returning a contract struct:

```go
// sandbox/new.go — generated by `agnos build` from the files of sandbox/binds/
func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	binds.ActionsBind(deps, &self)
	binds.CliBind(deps, &self)
	return &self
}
```

The caller decides which implementation fills the fields flowing in:

```go
// cmd/main/main.go — the only place an adapter and the sandbox meet
deps := agnosadapter.New()
lib := agnoslib.New(&deps)
```

For why the contracts are structs rather than interfaces, continue to [StructContracts.md](/docs/References/StructContracts.md).
