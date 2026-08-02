# Adapters

## Description
Lists every adapter shipped with the library — the opinionated `deps.Deps` implementations under `adapters/` — and when to use each one. To build a new adapter, follow [AddAdapter.md](/docs/Tutorials/AddAdapter.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Text-file store under the OS temp dir; real wall clock | You want the zero-config default, with values surviving across runs |
| `frozen` | [frozen.New](/docs/References/PublicApi/frozen.New.md) | In-memory store; clock frozen at a chosen time | You need deterministic expiry in tests, without real waiting |
