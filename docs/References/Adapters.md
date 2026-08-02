# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. Every adapter exposes a `New(...) deps.Deps` factory that runs one `<Field>Factory` per field of the contract and returns the filled contract struct, ready to be passed to [`lib.New`](/docs/References/PublicApi/lib.New.md) — the same [factory pattern](/docs/References/RULES.md#factory-pattern) the sandbox uses. Any single field can be replaced before injection — see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md). To build a new adapter, follow [AddAdapter.md](/docs/Tutorials/AddAdapter.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Single JSON file at a caller-chosen path; real wall clock | You want the default, with values surviving across runs |
| `memory` | [memory.New](/docs/References/PublicApi/memory.New.md) | In-memory map guarded by a mutex; real wall clock | You want the fastest store and don't need values after the process exits |
