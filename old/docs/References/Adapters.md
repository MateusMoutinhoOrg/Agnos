# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory that runs one `<Field>Factory` per field of the contract and returns the filled contract struct, ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md) — the same [factory pattern](/docs/References/StructContracts.md#factories-fill-the-fields) the sandbox uses. Any single field can be replaced before injection — see [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#overwrinting-a-adapter-function). To build a new adapter, follow [HandleDependencies.md → Creating a adapter in repo](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-repo).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Real wall clock; `Printf` to standard output; embedded Verb parser over `os.Args[1:]`; embedded Keep database on the filesystem, one file per key under a caller-chosen base path; the project's assets compiled into the binary and served whole; the filesystem over `os`/`path/filepath`; HTTP over `net/http` with a bounded round trip | You want the default, with categories and transactions surviving across runs |

An adapter filling `Printf` with a buffer and `VerbLib` with a fixed argument vector is what makes the command-line interface itself — `api.Lib.Sandboxmain` — runnable without a terminal.

---

## Embedded Libraries

`Deps` carries four fields that are not behaviors but whole libraries: [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), the embedded Verb argv parser, [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md), the embedded Keep schema database, [`EmbedDeps`](/docs/References/PublicApi/embeddeps.Lib.md), the assets compiled into the binary, and [`IoLib`](/docs/References/PublicApi/iodeps.Lib.md), the filesystem. Every adapter must fill them, because the sandbox cannot import Verb itself — it holds only a copy of Verb's api in `sandbox/contracts/deps/verbdeps/`. An adapter's `VerbLibFactory` initializes the real library and assigns its fields onto that copy, which is why it returns a **value** rather than a closure. The `standard` adapter reads the process's command line.

`KeepLibFactory` works the same way with one addition: Keep's fields hand back further api structs (`KeepDatabase`, `SchemaInstance`, `SchemaItem`), so instead of assigning them straight across, the factory wraps each in a closure that converts the returned struct into the sandbox's copy — nothing of the embedded library ever crosses the wall. The `standard` adapter wires Keep's filesystem adapter, so the tracker's categories and transactions survive across runs.

`EmbedDepsFactory` fills the third one, in a file of its own — `adapters/standard/embed.go`. The library it stands for is not a module but Go's own `embed` machinery, which the sandbox may not use either: the `standard` adapter compiles the [`/assets/`](/docs/References/Structure.md#assets) tree into the binary and serves it from that tree's root. Another adapter could serve the same three functions from a directory on disk or from a network store, and the library would read different files without changing a line.

---

## Standing Capabilities

Three fields are declared and filled like any other dependency, but the financial tracker never calls them: [`EmbedDeps`](/docs/References/PublicApi/embeddeps.Lib.md), [`IoLib`](/docs/References/PublicApi/iodeps.Lib.md), and [`NewRequest`](/docs/References/PublicApi/requestdeps.Request.md). They ship because this repository is a template — a derived library reads embedded assets, touches the filesystem, or speaks HTTP without designing a contract for it first.

An adapter must fill them regardless. An unfilled field is a nil function the compiler does not catch, and it panics on first use, so "the current library does not call it" is not a reason to skip one.

| Field | Filled by | `standard` |
|-------|-----------|------------|
| `EmbedDeps` | `EmbedDepsFactory` in `adapters/standard/embed.go` | The [`/assets/`](/docs/References/Structure.md#assets) tree compiled into the binary, served from its root |
| `IoLib` | `IoLibFactory` in `adapters/standard/io.go` | `os` and `path/filepath`, resolving whatever path the caller passes |
| `NewRequest` | `NewRequestFactory` in `adapters/standard/request.go` | `net/http`, with every round trip bounded by a timeout the sandbox cannot set for itself |
