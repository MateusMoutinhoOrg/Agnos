# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory that runs one `<Field>Factory` per field of the contract and returns the filled contract struct, ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md) — the same [factory pattern](/docs/References/StructContracts.md#factories-fill-the-fields) the sandbox uses. Any single field can be replaced before injection — see [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md). To build a new adapter, follow [HandleAdapters.md](/docs/Tutorials/HandleAdapters.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Real wall clock; `Printf` to standard output; embedded Verb parser over `os.Args[1:]`; embedded Keep database on the filesystem, one file per key under a caller-chosen base path; the project's assets compiled into the binary and served whole | You want the default, with categories and transactions surviving across runs |

An adapter filling `Printf` with a buffer and `VerbLib` with a fixed argument vector is what makes the command-line interface itself — `api.Lib.Sandboxmain` — runnable without a terminal.

---

## Embedded Libraries

`Deps` carries three fields that are not behaviors but whole libraries: [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), the embedded Verb argv parser, [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md), the embedded Keep schema database, and [`EmbedDeps`](/docs/References/PublicApi/embeddeps.Lib.md), the embedded assets the interface reads its text from. Every adapter must fill them, because the sandbox cannot import Verb itself — it holds only a copy of Verb's api in `sandbox/contracts/deps/verbdeps/`. An adapter's `VerbLibFactory` initializes the real library and assigns its fields onto that copy, which is why it returns a **value** rather than a closure. The `standard` adapter reads the process's command line.

`KeepLibFactory` works the same way with one addition: Keep's fields hand back further api structs (`KeepDatabase`, `SchemaInstance`, `SchemaItem`), so instead of assigning them straight across, the factory wraps each in a closure that converts the returned struct into the sandbox's copy — nothing of the embedded library ever crosses the wall. The `standard` adapter wires Keep's filesystem adapter, so the tracker's categories and transactions survive across runs.

`EmbedDepsFactory` fills the third one, in a file of its own — `adapters/standard/embed.go`. The library it stands for is not a module but Go's own `embed` machinery, which the sandbox may not use either: the `standard` adapter compiles the [`/assets/`](/docs/References/Structure.md#assets) tree into the binary and serves it from that tree's root. Another adapter could serve the same three functions from a directory on disk or from a translation store, and the library would print different words without changing a line.
