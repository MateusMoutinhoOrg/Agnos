# Contribution Rules

Rules to follow when contributing to this project. Every file must also be shaped by the specification that governs it — locate it in [Specs.md](/docs/References/Specs.md).

---

## Tutorials Guide
Before making anything, open the theme index matching your goal — [CliUsage](/docs/Index/CliUsage.md), [LibUsage](/docs/Index/LibUsage.md), [Development](/docs/Index/Development.md), or [Templating](/docs/Index/Templating.md), all linked from the [README.md](/README.md) — and search its `Tutorials` section for a tutorial about what you want to do. Each entry lists the tutorial's own sections, so you can tell what it covers without opening it. If there is one, follow it; if there isn't, you need to create one following the spec defined in [TutorialDocs](/docs/References/Specs/TutorialDocs/).


## Specification Compliance

Before creating or editing any file, read [Specs.md](/docs/References/Specs.md) and check whether the file matches an **Applies To** entry. If it does, create or edit it following the specification that entry points to — reproduce the shape it requires, using its `sample` as reference.

---

## Sandbox Isolation

[sandbox/](/sandbox/) is a closed sandbox. No file inside it may import [adapters/](/adapters/), [examples/libraryExamples/](/examples/libraryExamples/), a third-party module, or an OS-bound standard-library package (`os`, `net`, `os/exec`, `syscall`, …). Every such effect must be declared as a function field on the `Deps` contract and reached through `l.Deps`, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md). The mechanic is explained in [SandboxIsolation.md](/docs/References/SandboxIsolation.md).

Contracts are **structs of function fields**, never interfaces — in [sandbox/contracts/deps](/sandbox/contracts/deps/) and [sandbox/contracts/api](/sandbox/contracts/api/) alike. Every type in the project is declared in `sandbox/contracts/`; [sandbox/internal/](/sandbox/internal/) declares no types at all. See [StructContracts.md](/docs/References/StructContracts.md).

---

## Factory Pattern

Every object this project hands out — an `api` struct built inside the sandbox, a `deps.Deps` filled by an adapter outside it — is built by its package's `New(...)` constructor running **factories**, never by methods bound into fields, never by an internal mirror type, and never by filling a struct literal at the call site. When you write or edit any file holding `<Field>Factory` functions, follow the [Factories](/docs/References/Specs/Factories/Specs.md) specification on top of the one governing that file's tree.

A factory takes a pointer to the **carrier** — the struct holding the state the closure reads — and returns exactly one field's value; the caller assigns it:

```go
// sandbox/internal/lib/ — the carrier is the api struct being filled
func GetCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		record, ok := store.FindCategory(l.Deps, name)
		if !ok {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}

func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.GetCategory = GetCategoryFactory(&l)
	return l
}

// adapters/<name>/ — the carrier is the adapter, whose Deps field is the contract
func NowFactory(s *StandardAdapter) func() time.Time {
	return func() time.Time { return time.Now() }
}

func New() deps.Deps {
	s := &StandardAdapter{}
	s.Deps.Now = NowFactory(s)
	return s.Deps
}
```

Four rules follow, and none of them is checked by the compiler:

- Every api struct whose behavior needs dependencies declares a `Deps deps.Deps` field, and closures reach dependencies through it — `l.Deps.<Field>(...)`, read inside the closure, never captured at factory time. Every adapter struct declares the same field, as the contract its factories fill.
- Every field factory must be called and its return value assigned from its package's `New(...)` constructor, which is the factory aggregate — there is no separate `Factory` function. A field no factory fills stays nil and panics on first call.
- A `New` constructor returns the filled **contract struct** by value — `api.Lib`, `api.<Object>`, or `adapter.Deps` — never the carrier type of an adapter.
- The `Deps` field is **read-only once the struct is returned**: closures capture the struct the factories ran over, so a caller patching `Deps` on a copy changes nothing. Patch the `deps.Deps` value before calling `lib.New`.

Conversely, nothing outside the sandbox may reach into it beyond its three public packages: `sandbox` (package `lib`), `sandbox/contracts/deps`, and `sandbox/contracts/api`.

---

## Import Aliases

Any file that **consumes** the library from outside it — [examples/libraryExamples/](/examples/libraryExamples/) and third-party consumers — imports it under `agnos`-prefixed aliases, so each call site says which layer it belongs to:

| Import | Alias |
|--------|-------|
| `adapters/<name>` | `agnosadapter` |
| `sandbox` | `agnoslib` |
| `sandbox/contracts/api` | `agnostypes` |
| `sandbox/contracts/deps` | `agnosdeps` |

```go
import (
	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)
```

Files that belong to the library — everything under `sandbox/` and its own [adapters/](/adapters/) — keep the plain package names (`api`, `deps`): there the prefix would be noise, since the import is already local.

---

## File Changes

Before creating, deleting, or renaming any file or directory, read [Structure.md](/docs/References/Structure.md) and check whether the change affects the project structure. If it does, update [Structure.md](/docs/References/Structure.md) in the same commit.

---

## Specification Changes

When you create, delete, or rename a specification inside [Specs/](/docs/References/Specs/), you MUST adapt all the files that match the spec's Applies To rule, and update the index in [Specs.md](/docs/References/Specs.md).

---

## Documentation Changes

When you create, delete, or rename a `.md` file, update its theme index in `docs/Index/` — see [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md). The [README.md](/README.md) Doc Index changes only when a **theme** is added, renamed, or removed.

---

## Sample Changes

When you create, delete, or rename a sample, update the reference page that lists it: [SamplesList.md](/docs/References/SamplesList.md) for a script in [examples/cliExamples/](/examples/cliExamples/), [ApiSamplesList.md](/docs/References/ApiSamplesList.md) for a program in [examples/libraryExamples/](/examples/libraryExamples/).

---

## Interface Changes

The command-line interface lives **inside** the sandbox, as the `Sandboxmain` field of `api.Lib`, dispatching in [sandbox/internal/cli/](/sandbox/internal/cli/). The binary in [cmd/main/](/cmd/main/) wires an adapter into the library, calls that field, and exits with its return — it must never branch on a command, parse a flag, or print anything of its own, or the behavior stops being reachable from any other front end.

When you add or change a command or a flag, update [Commands.md](/docs/References/Commands.md) and the usage screen in [assets/usages.txt](/assets/usages.txt) in the same commit, following [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md).

---

## Display Text

No text a user reads is written inside [sandbox/](/sandbox/). The usage screen, the version, and every message the interface prints are files under [assets/](/assets/), reached through the injected `Deps.EmbedDeps` contract; the sandbox addresses them by path and holds only format skeletons and diagnostics. The mechanic is explained in [EmbeddedAssets.md](/docs/References/EmbeddedAssets.md).

When you make the library print a line it did not print before, add the file under [assets/](/assets/) and the constant naming it in `sandbox/internal/cli/cli.go`, following [HandleAssets.md](/docs/Tutorials/HandleAssets.md), and update the `/assets/` table in [Structure.md](/docs/References/Structure.md#assets) in the same commit. The `//go:embed all:*` directive in [assets/asset.go](/assets/asset.go) takes the whole directory, so a new asset never needs a new pattern — but a path the code asks for and no file answers still fails at runtime, never at build time.
