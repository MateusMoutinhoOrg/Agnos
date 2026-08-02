# `api.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	Deps deps.Deps
	Set  func(key string, value string, ttlSeconds int)
	Get  func(key string) (Entry, bool)
}
```

## Description

The library entry point, returned by [`lib.New`](/docs/References/PublicApi/lib.New.md). It is a key/value cache with per-key expiry, exposed as a struct of function fields: `lib.New` stores the injected [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) in `Deps`, then runs the factories in `sandbox/internal/lib/`, each of which fills one function field with a closure reading `Deps` at call time. The same deps are propagated into every [`Entry`](/docs/References/PublicApi/api.Entry.md) the lib creates. Calling a field reads exactly like calling a method — `l.Set("k", "v", 60)`. See [StructContracts.md](/docs/Explanations/StructContracts.md).

`Deps` is exported because the library's own factories read it, but it is **read-only after construction**: the closures already captured the struct they were built over, so reassigning `Deps` here does not change behavior. Patch the `deps.Deps` value before calling `lib.New`.

## Fields

| Field | Description |
| :--- | :--- |
| [`Deps deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | The dependency set injected by `lib.New`; read-only after construction. |
| [`Set func(key string, value string, ttlSeconds int)`](/docs/References/PublicApi/api.Set.md) | Stores a value under a key, expiring after `ttlSeconds`. |
| [`Get func(key string) (Entry, bool)`](/docs/References/PublicApi/api.Get.md) | Returns the live `Entry` for a key, or `false` on a miss or expiry. |
