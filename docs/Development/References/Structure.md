# Project Structure

This document maps the project **schema** — the kinds of files the project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/Development/References/Specs.md) to get its description and sample.

The project is a **CLI** whose interface lives inside the library. It is split into three top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the pure library, and the command-line interface with it. Nothing inside it may import `adapters/`, `cmd/`, `examples/libraryExamples/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation.md](/docs/Development/References/SandboxIsolation.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each adapter imports `sandbox/contracts/deps` and nothing else from the sandbox.
- **`/cmd/`** and **`/examples/libraryExamples/`** sit outside the sandbox too, and are the only places where an adapter and the sandbox meet — `cmd/` for the installable binary, `examples/libraryExamples/` for the runnable Go samples.
- **`/assets/`** sits outside the sandbox as well: the text the library displays, compiled into the binary and reached only through the injected `Deps.EmbedDeps` contract, so the sandbox holds no display text of its own.

Because the interface is `api.Lib.Sandboxmain` — one field of the library like any other — the binary in `cmd/main/` holds no command, no flag, and no output of its own: it wires, runs, and exits.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and the Doc Index pointing at each theme's `Index.md` | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition and dependencies | |
| `.gitignore` | Intentionally untracked files to ignore | |

---

## `/scripts/`
The cross-platform build scripts. One shell script per OS/architecture target, each a thin wrapper over `go build` with the target's `GOOS`/`GOARCH` set — the Go runtime cross-compiles on its own, so nothing here needs a container runtime or a cross-compiler. Every script resolves the repository root from its own path and writes its artifact to `release/`, which is git-ignored. Building is [Build.md](/docs/Development/Tutorials/Build.md).

| File | Description | Spec |
|------|-------------|------|
| `all.sh` | Runs every target script below, in order | |
| `linux86.sh` | Builds for Linux amd64, producing `release/linux86.out` | |
| `linuxarm64.sh` | Builds for Linux arm64, producing `release/linuxarm64.out` | |
| `linuxi32.sh` | Builds for Linux 386, producing `release/linuxi32.out` | |
| `windows86.sh` | Builds for Windows amd64, producing `release/windows86.exe` | |
| `windowsi32.sh` | Builds for Windows 386, producing `release/windowsi32.exe` | |
| `mac86.sh` | Builds for macOS Intel (amd64), producing `release/mac86.bin` | |
| `macarm64.sh` | Builds for macOS Apple Silicon (arm64), producing `release/macarm64.bin` | |

**Build one target:**
```sh
bash ./scripts/linux86.sh
```

**Build every target:**
```sh
bash ./scripts/all.sh
```

---

## `/sandbox/`
The closed sandbox — the pure library. It holds its own entry point, the contracts everything is wired through, and the internal implementation. It reaches nothing outside itself: every OS-bound or third-party effect arrives through the injected `Deps`. Its package is named `lib`, so consumers import it as `lib "…/sandbox"` and call `lib.New`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor storing `Deps` on `api.Lib` and running the internal factories over it | |

### `/sandbox/contracts/`
The structs the rest of the project is wired through — the only part of the sandbox anything outside it may import. Contracts hold the project's **public types** and are structs of function fields, never interfaces; see [StructContracts.md](/docs/Development/References/StructContracts.md). Contracts import nothing from `adapters/` or `sandbox/internal/`.

#### `/sandbox/contracts/deps/`
The contract every adapter must fill.

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | The `Deps` struct, one function field per injectable behavior, plus one plain field per embedded library | Deps |

##### `/sandbox/contracts/deps/verbdeps/`
The sandbox's copy of the embedded [Verb](https://github.com/MateusMoutinhoOrg/Verb) argv-parser library's public api. The sandbox may not import Verb — that would be a third-party import — so it restates the shape it needs, field for field; the adapter, outside the sandbox, is what fills it. Same mechanic as `bootstrap/sandbox/contracts/deps/agnosdeps/`.

| File | Description | Spec |
|------|-------------|------|
| `verbdeps.go` | Copy of the embedded Verb library's `api.Lib` struct, injected whole as the `Deps.VerbLib` field | |

##### `/sandbox/contracts/deps/keepdeps/`
The sandbox's copy of the embedded [Keep](https://github.com/MateusMoutinhoOrg/Keep) schema-database library's public api, for the same reason `verbdeps/` exists. Keep's api is a tree of structs — `Lib` hands back a `KeepDatabase`, which hands back a `SchemaInstance`, which hands back `SchemaItem`s — so the copy restates the whole tree, and the adapter converts each returned struct as it crosses.

| File | Description | Spec |
|------|-------------|------|
| `keepdeps.go` | Copy of the embedded Keep library's api structs and constants, injected whole as the `Deps.KeepLib` field | |

##### `/sandbox/contracts/deps/embeddeps/`
The sandbox's copy of an embedded-asset library's public api, for the same reason the two above exist: reading a file is an OS-bound effect, and compiling one into a binary needs the `//go:embed` directive, so the sandbox may declare neither. It is how the library reaches the text under [`/assets/`](#assets) — explained in [EmbeddedAssets.md](/docs/LibUsage/References/EmbeddedAssets.md).

| File | Description | Spec |
|------|-------------|------|
| `embeddeps.go` | Copy of the asset-reading api — `ReadFile`, `ListFiles`, `ListFilesRecursively` — injected whole as the `Deps.EmbedDeps` field | |

#### `/sandbox/contracts/api/`
The structs the library hands back to callers.

| File | Description | Spec |
|------|-------------|------|
| `api.go` | The `Lib` entry-point struct plus one struct per object the lib creates, each carrying a `Deps` field | Outputs |

### `/sandbox/internal/`
**Factories only** — no types. Each package here holds the functions that take a pointer to an [`api`](#sandboxcontractsapi) struct and return closures reading that struct's `Deps`, which the package's `New` constructor assigns into the matching function fields. Types never live here; they stay in `contracts/`. Go's `internal/` rule makes this tree unreachable from outside `sandbox/`, so neither consumers nor `adapters/` can reach in — the sandbox wall is enforced by the compiler, not by convention alone.

#### `/sandbox/internal/lib/`
The entry-point implementation. The `internal/` parent already marks it private, so the package carries no `internal_` prefix.

| File | Description | Spec |
|------|-------------|------|
| `lib.go` | One `<Field>Factory(l *api.Lib)` per lib function — `Sandboxmain` among them — each returning a closure, plus the `New(d deps.Deps) api.Lib` constructor that assigns every factory's return value and runs them all | LibFunctions |

#### `/sandbox/internal/cli/`
The command-line interface itself: the command dispatch `Sandboxmain` delegates to, the paths of the text it prints, and the reading of amounts off the command line. It reads the command line through `deps.Deps.VerbLib`, takes every word it displays from `deps.Deps.EmbedDeps`, and writes every line through `deps.Deps.Printf`, so the whole interface stays inside the closed sandbox and holds no display text of its own. Like `store/`, it is neither an object nor the entry point, so no specification governs it, and it declares **no types and no factories**.

| File | Description | Spec |
|------|-------------|------|
| `cli.go` | The `Run(l *api.Lib, args []string) int` dispatch, one helper per command group, the asset paths and message names the interface prints, and the amount parser | |

#### `/sandbox/internal/<object>/`
One package per object the library creates, named after the object itself — `category/` and `transaction/` for this library.

| File | Description | Spec |
|------|-------------|------|
| `<object>.go` | The object's `<Field>Factory` functions, each returning a closure, plus the `New(d deps.Deps, …) api.<Object>` constructor that propagates `Deps` and assigns every factory's return value | LibObjects |

#### `/sandbox/internal/store/`
Shared helpers over the injected database, used by `internal/lib/` and by every object package: the schema the tracker's records are persisted under, the lookups that reach it, and the encoding of a transaction's reference. It declares **no types and no factories** — it is the one internal package that is neither an object nor the entry point, so no specification governs it.

| File | Description | Spec |
|------|-------------|------|
| `store.go` | The database `Props`, the field-name constants, and the helpers that read a stored record through `deps.Deps.KeepLib` | |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract, each providing a distinct concrete behavior. This is where OS-bound and third-party code lives; an adapter imports `sandbox/contracts/deps` and nothing else from `sandbox/`. An adapter fills its contract with the same **factories** [`sandbox/internal/`](#sandboxinternal) uses — the carrier is the adapter struct, which declares the `Deps` field the factories' return values are assigned into.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory(a *<Name>Adapter)` per `Deps` field returning a closure, plus the `New(...) deps.Deps` constructor that assigns every factory's return value and runs them all | Adapters |
| `<field>.go` | One factory, split out of `<name>.go` when it carries conversion helpers of its own — `embed.go` in `standard/`, which wraps the compiled-in [`/assets/`](#assets) into the `Deps.EmbedDeps` contract. `New` still calls it | Factories |

---

## `/assets/`
Outside the sandbox. The files the library displays instead of holding them as Go strings: the interface's text today, images and item templates as the project grows. They are compiled into the binary, so an installed `agnos-cli` carries its own help screen with no files beside it, and they are reached only through the injected `Deps.EmbedDeps` contract — never imported by the sandbox. The mechanic is explained in [EmbeddedAssets.md](/docs/LibUsage/References/EmbeddedAssets.md); adding one is [HandleAssets.md](/docs/Development/Tutorials/HandleAssets.md).

This directory is a Go package for one reason: a `//go:embed` directive can only reach files inside its own package directory, so the directive has to sit next to the assets. That single directive is `//go:embed all:*`, which takes **every** file in the tree, so a new asset needs no change to it — put the file here and it exists at runtime.

| File | Description | Spec |
|------|-------------|------|
| `asset.go` | Package `assets`: the `//go:embed all:*` directive and the `Files` embedded filesystem the standard adapter serves | |
| `version.txt` | The interface version reported by `agnos-cli version` and `--version` | |
| `usages.txt` | The help screen, printed for `help`, for `--help`, and after any refused command line | |
| `messages/<name>.txt` | One file per line the interface can print, named after what it reports; a `Printf` format when it names a value | |

---

## `/cmd/`
Outside the sandbox. The executables the project ships — for this project, the single binary a user installs.

### `/cmd/main/`

| File | Description | Spec |
|------|-------------|------|
| `main.go` | Self-contained `package main` that wires an adapter into the lib, calls `api.Lib.Sandboxmain(os.Args[1:])`, and exits with its return | CliMain |

**Run the CLI from source:**
```sh
go run ./cmd/main <command> [arguments]
```

**Install it globally:**
```sh
go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3
```

---

## `/examples/cliExamples/`
Outside the sandbox. Shell scripts driving the built binary the way a user would from a terminal. Each one builds the CLI into a scratch directory and points it at a budget of its own, so nothing a script does touches the records in the user's home.

| File | Description | Spec |
|------|-------------|------|
| `<Name>.sh` | Self-contained shell script demonstrating one goal against the built CLI | CliExamples |

**Run a CLI example:**
```sh
bash ./examples/cliExamples/ManageCategories.sh
```

---

## `/examples/libraryExamples/`
Outside the sandbox. Runnable Go examples demonstrating how to use the library from code, when the CLI is not what the caller wants.

### `/examples/libraryExamples/<example>/`

| File | Description | Spec |
|------|-------------|------|
| `<example>.go` | Self-contained `package main` wiring an adapter into the lib | LibraryExamples |

**Run an example:**
```sh
go run ./examples/libraryExamples/<example>/<example>.go
```

---

## `/bootstrap/`
A second, self-contained Agnos-Cli library — same three trees (`sandbox/`, `adapters/`, `libraryExamples/`) and the same rules — demonstrating how one Agnos-Cli-compliant library **embeds** another. Its sandbox reaches nothing outside itself, so it never imports the root library: the embedded library arrives as one plain `Deps` field. Explained in [Bootstrap.md](/docs/LibUsage/References/Bootstrap.md).

| Path | Description |
|------|-------------|
| `sandbox/contracts/deps/deps.go` | The `Deps` struct, including `TrackerLib` — the embedded library, held as a locally declared contract struct |
| `sandbox/contracts/deps/agnosdeps/agnosdeps.go` | Copy of the embedded library's `api` structs, declared inside the sandbox so the sandbox never imports the embedded library |
| `adapters/<name>/<name>.go` | Its `TrackerLibFactory` initializes the embedded library with the embedded library's own adapter, and copies its `api` fields onto the local `agnosdeps` ones |
| `libraryExamples/<example>/<example>.go` | Self-contained `package main` wiring a bootstrap adapter into the bootstrap lib |

The copying lives in the adapter because only code outside the sandbox may import the embedded library. Because both sides are structs of function fields, the copy is field assignment: a wrapper is needed only where a named type differs between the two declarations. See [StructContracts.md](/docs/Development/References/StructContracts.md).

---

## `/docs/`
Documentation of the project, split by **theme** — what the reader wants to accomplish. One top-level directory per theme, each shaped the same way: an `Index.md` listing the theme's pages, a `Tutorials/` directory holding its workflows, and a `References/` directory holding its explanations and lookups. The [README](/README.md) links to the four indexes and to nothing else inside `docs/`.

| Directory | Description |
|-----------|-------------|
| `CliUsage/` | Installing the binary, driving it from a terminal, and the command surface |
| `LibUsage/` | Consuming the same behavior as a Go library, and the public API |
| `Development/` | Contributing: the rules, the mechanics, the workflows, the specifications |
| `Templating/` | Turning this repository into another library |

Every theme directory holds the same three slots:

| File | Description | Spec |
|------|-------------|------|
| `Index.md` | The theme's entry point: one table of its Tutorials, one of its References | Index |
| `Tutorials/<Goal>.md` | One page per workflow, its title phrased as the action it performs | TutorialDocs |
| `References/<Name>.md` | One page per lookup table or explained mechanic | ReferenceDocs / ExplanationDocs |

### `/docs/CliUsage/`
The command-line interface, from the reader's side.

| File | Description | Spec |
|------|-------------|------|
| `Tutorials/QuickStart.md` | Install the binary and track a first budget | TutorialDocs |
| `Tutorials/InstallCli.md` | Install the CLI globally, or build it from a checkout | TutorialDocs |
| `Tutorials/UseCli.md` | Drive every operation of the tracker from the terminal | TutorialDocs |
| `Tutorials/RunCliSample.md` | Run a shipped shell example from the source tree | TutorialDocs |
| `References/Commands.md` | Every command, flag, and exit code of the interface | ReferenceDocs |
| `References/SamplesList.md` | Every example under `examples/cliExamples/` | ReferenceDocs |

### `/docs/LibUsage/`
The same behavior consumed from Go code.

| File | Description | Spec |
|------|-------------|------|
| `Tutorials/QuickStart.md` | Install the module and run the smallest program using it | TutorialDocs |
| `Tutorials/LibInitialization.md` | Create deps via an adapter and inject them into the lib | TutorialDocs |
| `Tutorials/<Goal>.md` | One page per library workflow — categories, transactions, samples | TutorialDocs |
| `References/PublicApi.md` | Index of all public-facing components, linking to their detail pages | ReferenceDocs |
| `References/Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |
| `References/ApiSamplesList.md` | Every example under `examples/libraryExamples/` | ReferenceDocs |
| `References/Bootstrap.md` | How a library built this way embeds another one | ExplanationDocs |
| `References/EmbeddedAssets.md` | Where the text the library displays comes from, and how to serve your own | ExplanationDocs |

#### `/docs/LibUsage/References/PublicApi/`
One detail page per public-facing component, indexed by `PublicApi.md`. Reach a page through that index rather than by browsing the directory.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | One detail page per public struct, function, or field, named after the package the symbol is declared in | ReferenceDocs |

### `/docs/Development/`
Everything a contributor must read before changing the repository.

| File | Description | Spec |
|------|-------------|------|
| `Tutorials/Handle<Subject>.md` | One page per maintenance workflow — lib elements, commands, deps, adapters, assets, samples, documents | TutorialDocs |
| `Tutorials/Build.md` | Cross-compile the CLI for every target with the scripts in [`/scripts/`](#scripts) | TutorialDocs |
| `References/RULES.md` | Rules to follow when contributing to this project | Rules |
| `References/Structure.md` | The project's schema and the purpose of each component | Structure |
| `References/Specs.md` | Index of every specification and the files each one governs | |
| `References/SandboxIsolation.md` | What the sandbox may not import, and why every effect is a dep | ExplanationDocs |
| `References/StructContracts.md` | Why contracts are structs of function fields, and how factories fill them | ExplanationDocs |

#### `/docs/Development/References/Specs/`
The specifications describing how each kind of file in the project must be shaped. Never browse this directory — locate a specification by reading `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/Specs.md` | The required shape of the artifact the specification governs | |
| `<Spec>/sample.<ext>` | Concrete reference implementation of the specification | |

### `/docs/Templating/`
Turning this repository into another library.

| File | Description | Spec |
|------|-------------|------|
| `Tutorials/ForkTemplate.md` | Use this repo as a GitHub template for a new library | TutorialDocs |
| `Tutorials/AdaptExistingLib.md` | Convert an existing library to this DI structure | TutorialDocs |
| `Tutorials/RenameModule.md` | Rename the module path and update every internal import | TutorialDocs |
| `References/TemplateFileActions.md` | The action each template file takes when forking or adapting | ReferenceDocs |
