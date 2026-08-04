# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory that runs one `<Field>Factory` per field of the contract and returns the filled contract struct, ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md) — the same [factory pattern](/docs/References/RULES.md#factory-pattern) the sandbox uses. Any single field can be replaced before injection — see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md). To build a new adapter, follow [AddAdapter.md](/docs/Tutorials/AddAdapter.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Single JSON file at a caller-chosen path; real wall clock; embedded Verb parser over `os.Args[1:]`; embedded Keep database on the filesystem, under a `keepdata` directory beside the JSON file | You want the default, with values surviving across runs |
| `memory` | [memory.New](/docs/References/PublicApi/memory.New.md) | In-memory map guarded by a mutex; real wall clock; embedded Verb parser over caller-supplied args; embedded Keep database in memory | You want the fastest store and don't need values after the process exits |

---

## Embedded Libraries

`Deps` carries two fields that are not behaviors but whole libraries: [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md), the embedded Verb argv parser, and [`KeepLib`](/docs/References/PublicApi/keepdeps.Lib.md), the embedded Keep schema database. Every adapter must fill it, because the sandbox cannot import Verb itself — it holds only a copy of Verb's api in `sandbox/contracts/deps/verbdeps/`. An adapter's `VerbLibFactory` initializes the real library and assigns its fields onto that copy, which is why it returns a **value** rather than a closure. The only choice an adapter makes is which argument vector to parse: `standard` reads the process's command line, `memory` takes it from the caller.

`KeepLibFactory` works the same way with one addition: Keep's fields hand back further api structs (`KeepDatabase`, `SchemaInstance`, `SchemaItem`), so instead of assigning them straight across, the factory wraps each in a closure that converts the returned struct into the sandbox's copy — nothing of the embedded library ever crosses the wall. The choice each adapter makes is which storage backend Keep itself runs on: `standard` wires Keep's filesystem adapter, `memory` wires Keep's native in-memory adapter.
