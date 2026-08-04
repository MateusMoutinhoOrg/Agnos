# Public API

## Description
Index of every public-facing entry of the library. Callers hold **structs of function fields** declared in `sandbox/contracts/api` and `sandbox/contracts/deps`; the **factories** that fill those fields live in `sandbox/internal/` and are unreachable from outside `sandbox/`. See [StructContracts.md](/docs/Explanations/StructContracts.md).

---

## Structs

### [api.Lib](/docs/References/PublicApi/api.Lib.md)
The library entry point — a key/value cache with per-key TTL. Returned by `lib.New`; exposes all library functions as fields.

### [api.Entry](/docs/References/PublicApi/api.Entry.md)
A single cached record handed back by the library, with its dependencies already wired into `IsExpired`.

### [deps.Deps](/docs/References/PublicApi/deps.Deps.md)
The dependency contract every adapter must fill: the clock, the storage backend, and the embedded Verb argv parser.

### [verbdeps.Lib](/docs/References/PublicApi/verbdeps.Lib.md)
The sandbox's copy of the embedded Verb argv-parser library's api, injected whole as the `deps.Deps.VerbLib` field.

---

## Functions

### [lib.New](/docs/References/PublicApi/lib.New.md)
Injects a `deps.Deps` into the library and returns an `api.Lib`.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Creates a `deps.Deps` using the standard library adapter (JSON-file store + real clock).

### [memory.New](/docs/References/PublicApi/memory.New.md)
Creates a `deps.Deps` using the memory adapter (in-memory store + real clock).

---

## Fields

### [api.Lib.Deps / api.Entry.Deps](/docs/References/PublicApi/api.Deps.md)
The injected dependency set the struct was built with; read-only after construction.

### [api.Lib.Set](/docs/References/PublicApi/api.Set.md)
Stores a value under a key, expiring after the given number of seconds.

### [api.Lib.Get](/docs/References/PublicApi/api.Get.md)
Returns the live `api.Entry` stored under a key, or `false` on a miss or expiry.

### [api.Entry.Value](/docs/References/PublicApi/api.Value.md)
The cached value.

### [api.Entry.ExpiresAt](/docs/References/PublicApi/api.ExpiresAt.md)
The moment the entry stops being valid.

### [api.Entry.IsExpired](/docs/References/PublicApi/api.IsExpired.md)
Reports whether the injected clock has passed the entry's expiry.
