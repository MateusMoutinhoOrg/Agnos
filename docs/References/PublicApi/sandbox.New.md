# `sandbox.New`

**Type:** Function

## Signature

```go
func New(deps *deps.Deps) *api.Sandbox
```

## Description

Injects a filled dependency contract into the sandbox and returns the [`api.Sandbox`](/docs/References/PublicApi/api.Sandbox.md) entry point. It creates an empty `api.Sandbox`, then runs every binder of `sandbox/binds/` over it — `ActionsBind`, `CliBind` — each assigning the internal implementations onto one contract struct, closing over the `deps` pointer so the sandbox reads dependencies at call time. This is the only wiring point: `sandbox` never imports an adapter, so the caller chooses which implementation to pass. The file is **generated** by `agnos build` from the listing of `sandbox/binds/`; in a project without `sandbox/deps/` the signature is `New() *api.Sandbox`.

Importers alias the package: `agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"`, matching the `agnosadapter` / `agnoslib` convention `cmd/main` uses.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `deps` | [`*deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A dependency contract with every field filled, usually built by [`standard.New`](/docs/References/PublicApi/standard.New.md). The pointer is kept: a field reassigned on it later is seen by the sandbox. |

## Returns

| Type | Description |
| :--- | :--- |
| [`*api.Sandbox`](/docs/References/PublicApi/api.Sandbox.md) | The ready-to-use library: `Actions` for every operation, `Cli` for the command-line interface. |

## Examples

```go
package main

import (
	"os"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	// 1. Build the deps through an adapter
	deps := agnosadapter.New()

	// 2. Inject them — the sandbox never knows which adapter is behind them
	lib := agnoslib.New(&deps)

	// 3. Run the whole command-line interface over the process's argv
	os.Exit(lib.Cli.CliMain(os.Args[1:]))
}
```
