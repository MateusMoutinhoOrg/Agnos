# LibUsage

`agnos` is a Go module before it is anything else: every feature lives in `sandbox/`
and is reachable from any Go program that imports it.

```bash
go get github.com/MateusMoutinhoOrg/Agnos@latest
```

## Wiring

`sandbox/` performs no OS effects of its own — filesystem, clock, stdout, processes all
arrive through a `deps.Deps` struct. `adapters/bindings/standard` builds the ready-made
assembly, and `sandbox.New` turns it into the API object, which carries the deps on
`Sandbox.Deps` — so everything inside reaches them through the api it was handed.

```go
package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
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
| `lib.Config` | `api.Config` |

`lib.Cli.Commands` (`[]api.Command`) is the command surface itself: every command
the project declares, each carrying its flags, its args and the `Handler` that runs it.
`api.BindCommand(&command)` copies one into the command a single run binds to, so a caller
drives a command without a command line — bind the values into the copy's `Items` and call
`copy.Handler(copy)`.

[PublicApi](../PublicApi/doc.md) lists every one of them — signatures, props structs and
dependency contracts — generated from `sandbox/api/` itself on every build.

## Custom deps

Every sub-contract is a struct of function fields, so any of them can be swapped for a
test double, an in-memory implementation or an instrumented wrapper. Patch fields **before**
`sandbox.New(&deps)`: the constructors capture the pointer.

```go
deps := standard.New()

var out bytes.Buffer
deps.StdDeps.Printf = func(f string, a ...any) (int, error) {
	return fmt.Fprintf(&out, f, a...)
}

lib := sandbox.New(&deps)
```

The contracts available to patch:

| Field | Contract package |
| --- | --- |
| `deps.OpinionatedAgnosCli` | `sandbox/deps/OpinionatedAgnosCli` |
| `deps.ArgvDeps` | `sandbox/deps/argvdeps` |
| `deps.EmbedDeps` | `sandbox/deps/embeddeps` |
| `deps.GoimportsDeps` | `sandbox/deps/goimportsdeps` |
| `deps.HashDeps` | `sandbox/deps/hashdeps` |
| `deps.InterviewDeps` | `sandbox/deps/interviewdeps` |
| `deps.IoDeps` | `sandbox/deps/iodeps` |
| `deps.ReflectDeps` | `sandbox/deps/reflectdeps` |
| `deps.RunDeps` | `sandbox/deps/rundeps` |
| `deps.SerializableDeps` | `sandbox/deps/serializabledeps` |
| `deps.ServerDeps` | `sandbox/deps/serverdeps` |
| `deps.SortDeps` | `sandbox/deps/sortdeps` |
| `deps.StdDeps` | `sandbox/deps/stddeps` |
| `deps.StringsDeps` | `sandbox/deps/stringsdeps` |
| `deps.TemplateDeps` | `sandbox/deps/templatedeps` |

Each one is filled by a matching implementation under `adapters/impls/`, every package
exposing the same `Bind(deps *deps.Deps)` entry point:

| Adapter lib | Binder |
| --- | --- |
| `adapters/impls/OpinionatedAgnosCli` | `OpinionatedAgnosCli.Bind(&deps)` |
| `adapters/impls/goembed` | `goembed.Bind(&deps)` |
| `adapters/impls/nethttpserver` | `nethttpserver.Bind(&deps)` |
| `adapters/impls/osexecrun` | `osexecrun.Bind(&deps)` |
| `adapters/impls/osio` | `osio.Bind(&deps)` |
| `adapters/impls/osstd` | `osstd.Bind(&deps)` |
| `adapters/impls/sha256hash` | `sha256hash.Bind(&deps)` |
| `adapters/impls/stdargv` | `stdargv.Bind(&deps)` |
| `adapters/impls/stdgoimports` | `stdgoimports.Bind(&deps)` |
| `adapters/impls/stdreflect` | `stdreflect.Bind(&deps)` |
| `adapters/impls/stdserializable` | `stdserializable.Bind(&deps)` |
| `adapters/impls/stdsort` | `stdsort.Bind(&deps)` |
| `adapters/impls/stdstrings` | `stdstrings.Bind(&deps)` |
| `adapters/impls/texttemplate` | `texttemplate.Bind(&deps)` |
| `adapters/impls/ttyinterview` | `ttyinterview.Bind(&deps)` |

Starting from `standard.New()` is the safe default: an unfilled field is a nil func that
panics on first call. For a permanent mix, write your own
`adapters/bindings/<name>/generated.new.go` binding only the libs you want — `standard/generated.new.go` is
regenerated on every build, while other directories under `bindings/` are left alone.

`sandbox/api` is pure contract and `sandbox/` never touches the OS, so both are safe to import
anywhere; the rest of the rules a caller can count on are in [Rules](../Rules/doc.md#layers),
and [DepList](../DepList/doc.md) lists every contract that can be added.
