# `deps.Deps`

**Type:** Struct

## Definition

```go
type Deps struct {
	Now   func() time.Time
	Load  func(key string) (value string, expiresAtUnix int64, ok bool)
	Store func(key string, value string, expiresAtUnix int64)
}
```

## Description

The dependency contract every adapter must fill. Each field is one injectable behavior the cache needs: `Now` is the clock (injecting it makes TTL expiry testable without real waiting) and `Load`/`Store` are the storage backend (an adapter can keep records in memory, on disk, or anywhere). A filled `Deps` is built by an adapter — see [`standard.New`](/docs/References/PublicApi/standard.New.md) — and passed to [`lib.New`](/docs/References/PublicApi/lib.New.md).

Because it is a struct and not an interface, a value returned by an adapter can be patched field by field before injection, and a custom contract needs no type declaration at all. The trade-off: the compiler cannot detect a field you forgot to fill — it stays nil and panics on first call. See [StructContracts.md](/docs/Explanations/StructContracts.md) and [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Now func() time.Time` | Returns the current time, used to stamp and check expiry. |
| `Load func(key string) (value string, expiresAtUnix int64, ok bool)` | Fetches a stored record; `ok` is false when the key is absent. |
| `Store func(key string, value string, expiresAtUnix int64)` | Persists a record and its unix expiry timestamp. |

## Examples

```go
package main

import (
	"fmt"
	"time"

	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
	agnosdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
)

func main() {
	// A deps.Deps built without any adapter. Its clock never moves, so
	// expiry is fully deterministic — ideal for tests.
	now := time.Unix(0, 0)
	store := map[string]string{} // value only; expiry ignored for brevity

	d := agnosdeps.Deps{
		Now: func() time.Time { return now },
		Load: func(key string) (string, int64, bool) {
			v, ok := store[key]
			return v, now.Add(time.Hour).Unix(), ok
		},
		Store: func(key, value string, expiresAtUnix int64) {
			store[key] = value
		},
	}

	l := agnoslib.New(d)

	l.Set("k", "v", 60)
	entry, _ := l.Get("k")
	fmt.Println(entry.Value) // v
}
```
