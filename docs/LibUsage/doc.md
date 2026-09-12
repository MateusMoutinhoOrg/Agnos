# LibUsage

`agnos` is a Go module before it is anything else: every feature lives in `sandbox/`
and is reachable from any Go program that imports it.

```bash
go get github.com/MateusMoutinhoOrg/Agnos@latest
```

## Wiring

`sandbox/` performs no OS effects of its own — filesystem, clock, stdout, processes all
arrive through a `deps.Deps` struct. `adapters/availables/standard` builds the ready-made
assembly, and `sandbox.New` turns it into the API object, which carries the deps on
`Sandbox.Deps` — so everything inside reaches them through the api it was handed.

```go
package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox

	_ = lib
}
```

## What the sandbox exposes

`*api.Sandbox` is a flat struct, one field per contract declared in `sandbox/api/`.
Everything callable from Go is behind one of them.

| Field | Type |
| --- | --- |
| `lib.Actions` | `api.Actions` |
| `lib.Cli` | `api.Cli` |
| `lib.Commands` | `api.Commands` |

[PublicApi](../PublicApi/doc.md) lists every one of them — signatures, props structs and
dependency contracts — generated from `sandbox/api/` itself on every build.

## Custom deps

Every sub-contract is a struct of function fields, so any of them can be swapped for a
test double, an in-memory implementation or an instrumented wrapper. Patch fields **before**
`sandbox.New(&deps)`: the binders capture the pointer.

```go
deps := standard.New()

var out bytes.Buffer
deps.Std.Printf = func(f string, a ...any) (int, error) {
	return fmt.Fprintf(&out, f, a...)
}

lib := sandbox.New(&deps)
```

The contracts available to patch:

| Field | Contract package |
| --- | --- |
| `deps.Argvdeps` | `sandbox/deps/argvdeps` |
| `deps.Embeddeps` | `sandbox/deps/embeddeps` |
| `deps.Goimportsdeps` | `sandbox/deps/goimportsdeps` |
| `deps.Hashdeps` | `sandbox/deps/hashdeps` |
| `deps.Iodeps` | `sandbox/deps/iodeps` |
| `deps.Rundeps` | `sandbox/deps/rundeps` |
| `deps.Serializables` | `sandbox/deps/serializables` |
| `deps.Serverdeps` | `sandbox/deps/serverdeps` |
| `deps.Sortdeps` | `sandbox/deps/sortdeps` |
| `deps.Std` | `sandbox/deps/std` |
| `deps.Stringsdeps` | `sandbox/deps/stringsdeps` |
| `deps.Templatedeps` | `sandbox/deps/templatedeps` |

Each one is filled by a matching implementation under `adapters/libs/`, every package
exposing the same `Bind(deps *deps.Deps)` entry point:

| Adapter lib | Binder |
| --- | --- |
| `adapters/libs/embeddeps` | `embeddeps.Bind(&deps)` |
| `adapters/libs/goimportsdeps` | `goimportsdeps.Bind(&deps)` |
| `adapters/libs/hashdeps` | `hashdeps.Bind(&deps)` |
| `adapters/libs/iodeps` | `iodeps.Bind(&deps)` |
| `adapters/libs/rundeps` | `rundeps.Bind(&deps)` |
| `adapters/libs/serializables` | `serializables.Bind(&deps)` |
| `adapters/libs/serverdeps` | `serverdeps.Bind(&deps)` |
| `adapters/libs/sortdeps` | `sortdeps.Bind(&deps)` |
| `adapters/libs/std` | `std.Bind(&deps)` |
| `adapters/libs/stringsdeps` | `stringsdeps.Bind(&deps)` |
| `adapters/libs/templatedeps` | `templatedeps.Bind(&deps)` |
| `adapters/libs/verb` | `verb.Bind(&deps)` |

Starting from `standard.New()` is the safe default: an unfilled field is a nil func that
panics on first call. For a permanent mix, write your own
`adapters/availables/<name>/new.go` binding only the libs you want — `standard/new.go` is
regenerated on every build, while other directories under `availables/` are left alone.

`sandbox/api` is pure contract and `sandbox/` never touches the OS, so both are safe to import
anywhere; the rest of the rules a caller can count on are in [Rules](../Rules/doc.md#layers),
and [DepList](../DepList/doc.md) lists every contract that can be added.
