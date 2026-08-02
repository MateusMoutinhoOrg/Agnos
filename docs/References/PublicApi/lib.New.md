# `lib.New`

**Type:** Function

## Signature

```go
func New(d deps.Deps) api.Lib
```

## Description

Injects a filled dependency contract into the library and returns the [`api.Lib`](/docs/References/PublicApi/api.Lib.md) entry point. It stores the deps on the struct's `Deps` field, then runs the factories in `sandbox/internal/lib/` over it, each filling one function field with a closure that reads those deps at call time. This is the only wiring point: `sandbox` never imports an adapter, so the caller chooses which implementation to pass. The package is named `lib` and lives at `sandbox/`, so importers alias it: `agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"` — matching the `agnosadapter` / `agnoslib` / `agnostypes` alias convention used by every consumer of this module.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `d` | [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A dependency contract with every field filled, usually built by an adapter. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Lib`](/docs/References/PublicApi/api.Lib.md) | A ready-to-use library instance carrying the injected deps. |

## Examples

### Basic Initialization

```go
package main

import (
	"log"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	// 1. Build the deps through an adapter
	deps := agnosadapter.New("cache.json")

	// 2. Inject them into the library
	l := agnoslib.New(deps)

	// The library instance 'l' is now ready for use.
	log.Println("Library successfully initialized:", l.Set != nil)
}
```
