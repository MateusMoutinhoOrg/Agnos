# Sandbox Isolation

## Description
Explains what "closed sandbox" means in practice: the library in `sandbox/` reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package — so everything it can do is exactly what the injected `Deps` allows. First stage of the Development learning path; [StructContracts.md](/docs/References/StructContracts.md) explains how the contracts crossing this wall are shaped.

---

## The Three Trees

Three top-level directories, and the arrows only point one way:

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

- `sandbox/` is the library. Closed: it imports only itself and OS-independent stdlib packages (`time`, `strings`, `errors`, …).
- `adapters/` is the only place `os`, `net`, a database driver, or any third-party module may appear.
- `cmd/` and `examples/libraryExamples/` are the only places an adapter and the sandbox are named in the same file.

Nothing inside `sandbox/` can be affected by which OS, filesystem, or network the program runs on, because it has no way to reach any of them.

---

## What the Wall Forbids

A file under `sandbox/` may not import:

| Forbidden | Why |
|-----------|-----|
| `adapters/…` | Binding to one concrete implementation makes injection pointless. |
| `cmd/…`, `examples/libraryExamples/…` | Consumers of the library are never part of it. |
| Any third-party module | A dependency the caller cannot replace is one the caller cannot test around. |
| OS-bound stdlib (`os`, `net`, `os/exec`, `syscall`, …) | The effect belongs in an adapter, reached through a `Deps` field. |
| `embed`, and the `//go:embed` directive | Compiling a file into the binary is a build-time, filesystem-bound effect; the bytes arrive through a `Deps` field like any other. |

Everything the library needs from the outside world is a function field on `Deps` — the only door in the wall:

```go
// sandbox/contracts/deps/deps.go
type Deps struct {
	Now       func() time.Time // instead of time.Now()
	Printf    func(format string, a ...any) (int, error) // instead of fmt.Printf
	VerbLib   verbdeps.Lib     // instead of importing the Verb argv parser
	KeepLib   keepdeps.Lib     // instead of importing the Keep database
	EmbedDeps embeddeps.Lib    // instead of //go:embed and os.ReadFile
}
```

This is what lets the whole command-line interface live inside the wall: `api.Lib.Sandboxmain` reads the command line through `Deps.VerbLib` and prints through `Deps.Printf`, never touching `os.Args` or `os.Stdout`. The binary hands it an argument vector; a test can hand it a fixed vector and a buffer instead.

The wall reaches the interface's *words* too, not only its streams: the usage screen, the version, and every message come from `Deps.EmbedDeps`, so `sandbox/` holds no display text at all — see [EmbeddedAssets.md](/docs/References/EmbeddedAssets.md).

`VerbLib`, `KeepLib` and `EmbedDeps` are the same door in a different shape: each is foreign to this sandbox, so instead of importing it the sandbox declares a copy of its api in `sandbox/contracts/deps/verbdeps/`, `keepdeps/` and `embeddeps/` and lets the adapter fill it — see [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#injecting-a-whole-library).

To add a new door, follow [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).

---

## What the Wall Forbids in the Other Direction

`sandbox/` holds the factories that fill the contract structs, and Go's `internal/` rule makes it unreachable from outside `sandbox/` — an adapter or consumer that tries gets a compile error. The outside world sees exactly three packages:

| Package | Who imports it | For what |
|---------|----------------|----------|
| `sandbox` (package `lib`) | consumers, examples | `lib.New(deps) api.Lib` — the single wiring point |
| `sandbox/contracts/deps` | adapters, consumers | the contract struct to fill |
| `sandbox/contracts/api` | consumers, examples | the structs handed back |

The api structs the caller holds are the same ones the library works on — a consumer can read `l.Deps` — but the factories that turned it into behavior stay unreachable, so they can be renamed or restructured without breaking a consumer.

---

## Why the Entry Point Lives Inside

`sandbox/new.go` stays inside the wall because it obeys the same rule — it names no adapter, accepting a contract struct and returning a contract struct:

```go
// sandbox/new.go
func New(d deps.Deps) api.Lib {
	return internallib.New(d)
}
```

The caller decides which implementation fills the fields flowing in:

```go
// This line lives outside the wall — the only place an adapter and the sandbox meet.
l := agnoslib.New(agnosadapter.New("data.json"))
```

For why the contracts are structs rather than interfaces, continue to [StructContracts.md](/docs/References/StructContracts.md).
