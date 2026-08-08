# Template File Actions

## Description
Lists every file and directory of this template and the action it takes when the template becomes a new library — whether by forking ([ForkTemplate.md](/docs/ForkTemplate.md)) or by adapting an existing library ([AdaptExistingLib.md](/docs/AdaptExistingLib.md)). Each file falls into exactly one action:

| Action | Meaning |
|--------|---------|
| **[Copy](#copy)** | Taken as-is — it describes the structure, not the library |
| **[Rewrite](#rewrite)** | The path and shape stay; the content becomes the new library's |
| **[Create](#create)** | Written from scratch for the new library |
| **[Delete](#delete)** | The template's example content, removed once the new library's own files exist |

To locate any file: find its exact path below; if it is not listed by name, it falls under the `*` row of its directory. `docs/` is flat, so its `.md` files are listed **by name** — the generic guides are copied, the library-shaped pages are rewritten, and the example-specific pages are deleted.

---

## Copy

Taken as-is from the template. They describe the structure itself, not the financial tracker, so they carry over unchanged. Adapting them is allowed but never required.

Copying these files carries over the template's **generic** guides and specifications only. The new library must still **create** its own case-specific tutorials and reference pages — see [Create](#create).

| Path | Description |
|------|-------------|
| `docs/Meta/*` | The specifications every file of the new library must be shaped by |
| `docs/RULES.md` | The binding contribution rules |
| `docs/Specs.md` | The index locating each specification |
| `docs/ForkTemplate.md`, `docs/AdaptExistingLib.md`, `docs/RenameModule.md`, `docs/TemplateFileActions.md` | The template workflows and this page |
| `docs/SandboxIsolation.md`, `docs/StructContracts.md` | The explanations of the structure's mechanics |
| `docs/HandleDependencies.md`, `docs/HandleLibElements.md`, `docs/HandleCliCommands.md`, `docs/HandleAdapters.md`, `docs/HandleSamples.md`, `docs/HandleCliExamples.md`, `docs/HandleDocuments.md`, `docs/ExposePublicApi.md` | The generic workflow guides for extending any library built on this structure |
| `sandbox/new.go` | The `New` constructor storing `Deps` on `api.Lib` and running the internal factories over it |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes. Every rewritten file must be shaped by the specification in its row, located through [Specs.md](/docs/Specs.md).

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, both quick starts, badges, Doc Index, and Examples sections | Readme |
| `sandbox/contracts/deps/deps.go` | The `Deps` function fields the new library requires | Deps |
| `sandbox/contracts/api/api.go` | The `Lib` struct and one struct per object the new library hands back | Outputs |
| `adapters/standard/standard.go` | The default adapter, filling the new `Deps` contract | Adapters |
| `cmd/main/main.go` | The new library's entry point: wire, call `Sandboxmain`, exit | CliMain |
| `docs/Structure.md` | The layout of the new library | Structure |
| `docs/Cli.md` | The commands, flags, and exit codes of the new library's interface | ReferenceDocs |
| `docs/PublicApi.md` | The index of the new public API entries | ReferenceDocs |
| `docs/Adapters.md` | The adapters the new library ships | AdaptersDoc |
| `docs/InstallCli.md`, `docs/UseCli.md` | Installing and driving the new library's CLI | TutorialDocs |
| `docs/LibInitialization.md` | Installing the new library and wiring its standard adapter | TutorialDocs |
| `docs/RunCliSample.md`, `docs/RunApiSample.md` | Running the new library's samples | TutorialDocs |
| `docs/SamplesList.md`, `docs/ApiSamplesList.md` | The new library's own samples | ReferenceDocs |

---

## Create

Written from scratch for the library being built or adapted. Nothing of the template's content survives here — the example files occupying these paths are removed by **[Delete](#delete)**. Every created file must be shaped by the specification in its row, located through [Specs.md](/docs/Specs.md).

| Path | Description | Specification |
|------|-------------|---------------|
| `sandbox/internal/lib/*` | The lib's field factories and the `New` constructor running them all, reaching every dependency through `l.Deps` | LibFunctions |
| `sandbox/internal/<object>/*` | One package per object the library hands back: its field factories and the `New` constructor running them all | LibObjects |
| `sandbox/internal/cli/*` | The command dispatch behind `api.Lib.Sandboxmain`, its usage screen, and its operand parsing | |
| `adapters/<name>/<name>.go` | One adapter per additional opinionated implementation of the `Deps` contract | Adapters |
| `examples/libraryExamples/<example>/<example>.go` | One runnable Go sample per demonstrated use case | LibraryExamples |
| `examples/cliExamples/<Name>.sh` | One shell script per goal demonstrated against the built CLI | CliExamples |
| `docs/PublicApi/<pkg>.<Symbol>.md` | One detail page per public API entry | ReferenceDocs |
| `docs/<Goal>.md` | One tutorial per workflow specific to the new library — the generic guides carried over by **[Copy](#copy)** do **not** fulfill this | TutorialDocs |
| `docs/<Name>.md` | Any reference page the new library needs beyond the public API index | ReferenceDocs |

---

## Delete

The template's example content — the financial tracker. Removed once the new library's own files exist. For `.md` files, follow [HandleDocuments.md](/docs/HandleDocuments.md) so the README Doc Index stays in sync.

| Path | Description |
|------|-------------|
| `sandbox/internal/*` | The tracker's lib factories, object packages, CLI dispatch, and store helpers — replaced by **[Create](#create)** |
| `sandbox/contracts/deps/verbdeps/`, `sandbox/contracts/deps/keepdeps/` | The sandbox copies of the embedded Verb and Keep libraries — keep one only if the new library embeds the same library |
| `examples/libraryExamples/*` | The tracker's Go samples |
| `examples/cliExamples/*` | The tracker's CLI scripts |
| `docs/PublicApi/*` | The tracker's public API detail pages |
| `docs/ManageCategories.md`, `docs/TrackTransactions.md` | The tracker's domain tutorials |
| `bootstrap/*` | The embedded-library demonstration — keep it only as reference while the new library embeds another Agnos-style library |
