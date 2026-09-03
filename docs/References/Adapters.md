# Adapters

## Description
Lists every adapter shipped in `adapters/` — the isolated implementations under `adapters/libs/<lib>/`, one per `Deps` sub-contract, and the ready-made assemblies under `adapters/availables/<name>/` that wire them together. Every lib exports one binder, `Bind(deps *deps.Deps)`, that fills its one field; `standard.New()` is generated to call every one of them. Any single field can be replaced before injection — see [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md). To build a new lib, follow [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `availables/standard` | [standard.New](/docs/References/PublicApi/standard.New.md) | Binds every lib under `adapters/libs/`: the real clock and the process streams, the Verb parser over the given argv, the assets compiled into the binary, `os`/`path/filepath`, `net/http` with a bounded round trip, `os/exec`, `go/parser`, and `gopkg.in/yaml.v3` / `encoding/json` | You want the default — it is what `cmd/main/main.go` wires, and what every generated project's `standard` adapter is too |

`standard/new.go` is **generated** by `agnos build` from the directories of `adapters/libs/`; a different mix is a hand-written `availables/<name>/new.go` of your own, which `verify` allows and `build` leaves alone.

---

## Adapter Libs

One package per sub-contract. The lib's package name is the implementation's; the field it fills is the contract's.

| Lib | Fills | Backed by | Notes |
|-----|-------|-----------|-------|
| `libs/verb` | `Deps.Argvdeps` — [argvdeps.Lib](/docs/References/PublicApi/argvdeps.Lib.md) | `github.com/MateusMoutinhoOrg/Verb` | `New(args)` initializes Verb over the given vector and copies its fields onto the sandbox's `Parser` copy, one for one. |
| `libs/embeddeps` | `Deps.Embeddeps` — [embeddeps.Lib](/docs/References/PublicApi/embeddeps.Lib.md) | `embed`, `io/fs`, `text/template` over `assets.Files` | Serves the `assets/` tree compiled into the binary from its root; `RenderTemplate` parses an asset as a Go template. The one lib importing a package of this module outside the sandbox. |
| `libs/goimportsdeps` | `Deps.Goimportsdeps` — [goimportsdeps.Lib](/docs/References/PublicApi/goimportsdeps.Lib.md) | `go/parser`, `go/ast`, `go/printer` | Parses Go source into package, imports and top-level declarations. Wired, not yet used by any action. |
| `libs/iodeps` | `Deps.Iodeps` — [iodeps.Lib](/docs/References/PublicApi/iodeps.Lib.md) | `os`, `path/filepath` | `WriteFile` creates missing parents; predicates report `false` rather than an error; `RemoveDir` removes files too. |
| `libs/requestdeps` | `Deps.Requestdeps` — [requestdeps.Lib](/docs/References/PublicApi/requestdeps.Lib.md) | `net/http` | `NewRequest(url)` builds one request per call; every round trip is bounded by a 30-second timeout the sandbox cannot set for itself. |
| `libs/rundeps` | `Deps.Rundeps` — [rundeps.Lib](/docs/References/PublicApi/rundeps.Lib.md) | `os/exec` | Merges stdout and stderr in order; a non-zero exit comes back in `Result.ExitCode`, only a program that could not start is an `error`. |
| `libs/serializables` | `Deps.Serializables` — [serializables.Lib](/docs/References/PublicApi/serializables.Lib.md) | `gopkg.in/yaml.v3`, `encoding/json` | Generic values wrapped as `SerializibleObject`s. Every parsable config reads and writes through it. |
| `libs/std` | `Deps.Std` — [std.Lib](/docs/References/PublicApi/std.Lib.md) | `time`, `fmt`, `os.Stdout`, `os.Stderr` | `Printf` to stdout; `Log` and `Error` to stderr; `Errorf` wraps `fmt.Errorf`. |

Every lib is mirrored as an installable dep under `assets/deplist/`, named after the **contract** it fills — see [DepList.md](/docs/References/DepList.md). The `dbdeps` dep (lib `keep`) ships there without being installed in this repository.

---

## Embedded Libraries

Two fields are not behaviors but whole third-party libraries: `Argvdeps`, the Verb argv parser, and — in projects that install `dbdeps` — the Keep schema database. Every adapter must fill them, because the sandbox cannot import the library itself: it holds only a copy of the api in `sandbox/deps/<x>/`. The lib initializes the real library and assigns its fields onto that copy, field for field, so nothing of the embedded library ever crosses the wall.

`Embeddeps` is the same shape for a different reason: the library it stands for is Go's own `embed` machinery, which the sandbox may not use either. `libs/embeddeps` compiles the [`/assets/`](/docs/References/Structure.md#assets) tree into the binary and serves it from that tree's root; another lib could serve the same four functions from a directory on disk, and the sandbox would render different templates without changing a line.

---

## Standing Capabilities

`Goimportsdeps` and `Requestdeps` are declared and filled like every other dependency, but no action calls them yet. They ship because a generated project — and a future Agnos — reads Go source or speaks HTTP without designing a contract for it first.

An adapter must fill them regardless. An unfilled field is a nil function the compiler does not catch, and it panics on first use, so "the current code does not call it" is not a reason to skip one.

| Field | Filled by | `standard` |
|-------|-----------|------------|
| `Goimportsdeps` | `Bind` in `adapters/libs/goimportsdeps/GoImports.go` | `go/parser` over the given source text |
| `Requestdeps` | `Bind` in `adapters/libs/requestdeps/Request.go` | `net/http`, with every round trip bounded by a timeout |
