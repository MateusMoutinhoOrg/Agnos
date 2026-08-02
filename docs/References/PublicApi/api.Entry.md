# `api.Entry`

**Type:** Struct

## Definition

```go
type Entry struct {
	Deps      deps.Deps
	Value     string
	ExpiresAt time.Time
	IsExpired func() bool
}
```

## Description

A single cached record handed back by the library. `Deps` is the dependency set propagated from the parent [`api.Lib`](/docs/References/PublicApi/api.Lib.md); `Value` and `ExpiresAt` are plain data, fixed when the entry is built; `IsExpired` is a function field filled by a factory in `sandbox/internal/entry/`, whose closure reads `Deps.Now()` at call time. An `Entry` is always constructed via [`api.Lib.Get`](/docs/References/PublicApi/api.Get.md), which propagates the deps in. The factories behind it live inside the closed sandbox and are not reachable by callers.

On a miss, `Get` returns the **zero value** of this struct — there is no nil entry to check for; use the returned `bool`.

## Fields

| Field | Description |
| :--- | :--- |
| [`Deps deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | The dependency set propagated from the parent lib; read-only after construction. |
| [`Value string`](/docs/References/PublicApi/api.Value.md) | The cached value. |
| [`ExpiresAt time.Time`](/docs/References/PublicApi/api.ExpiresAt.md) | The moment the entry stops being valid. |
| [`IsExpired func() bool`](/docs/References/PublicApi/api.IsExpired.md) | Reports whether the injected clock has passed the entry's expiry. |
