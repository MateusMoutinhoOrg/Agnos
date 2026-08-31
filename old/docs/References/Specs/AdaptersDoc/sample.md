# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. To build a new adapter, follow [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-repo).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | `standard.New` | Text-file store under the OS temp dir; real wall clock | You want the zero-config default, with values surviving across runs |
| `frozen` | `frozen.New` | In-memory store; clock frozen at a chosen time | You need deterministic expiry in tests, without real waiting |

This sample's library is hypothetical, so its factories are shown unlinked: an entry links to its `docs/References/PublicApi/<pkg>.<Symbol>.md` page only when that page is real.

---

## Embedded Libraries

`Deps` carries one field that is not a behavior but a whole library: `StoreLib`, the embedded key-value store. Every adapter must fill it, because the sandbox cannot import the store itself — it holds only a copy of that library's api in `sandbox/contracts/deps/storedeps/`. An adapter's `StoreLibFactory` initializes the real library and assigns its fields onto that copy, which is why it returns a **value** rather than a closure.

---

## Standing Capabilities

`Clock` is declared and filled like any other dependency, but no library function calls it yet. It ships because the contract is meant to outlive this one library.

An adapter must fill it regardless. An unfilled field is a nil function the compiler does not catch, and it panics on first use.

| Field | Filled by | `standard` | `frozen` |
|-------|-----------|------------|----------|
| `Clock` | `ClockFactory` in `<adapter>/clock.go` | The real wall clock | A time fixed at construction |
