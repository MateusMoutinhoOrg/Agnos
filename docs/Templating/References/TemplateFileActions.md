# Template File Actions

## Description
Lists every file and directory of this template and the action it takes when the template becomes a new library — whether by forking ([ForkTemplate.md](/docs/Templating/Tutorials/ForkTemplate.md)) or by adapting an existing library ([AdaptExistingLib.md](/docs/Templating/Tutorials/AdaptExistingLib.md)). Each file falls into exactly one action:

| Action | Meaning |
|--------|---------|
| **[Copy](#copy)** | Taken as-is — it describes the structure, not the library |
| **[Rewrite](#rewrite)** | The path and shape stay; the content becomes the new library's |
| **[Create](#create)** | Written from scratch for the new library |
| **[Delete](#delete)** | The template's example content, removed once the new library's own files exist |

To locate any file: find its exact path below; if it is not listed by name, it falls under the `*` row of its directory. Under `docs/`, pages are listed **by name** inside their theme — the generic guides are copied, the library-shaped pages are rewritten, and the example-specific pages are deleted.

---

## Copy

Taken as-is from the template. They describe the structure itself, not the financial tracker, so they carry over unchanged. Adapting them is allowed but never required.

Copying these files carries over the template's **generic** guides and specifications only. The new library must still **create** its own case-specific tutorials and reference pages — see [Create](#create).

| Path | Description |
|------|-------------|
| `docs/Development/References/Specs/*` | The specifications every file of the new library must be shaped by |
| `docs/Development/References/RULES.md` | The binding contribution rules |
| `docs/Development/References/Specs.md` | The index locating each specification |
| `docs/Templating/Tutorials/ForkTemplate.md`, `docs/Templating/Tutorials/AdaptExistingLib.md`, `docs/Templating/Tutorials/RenameModule.md`, `docs/Templating/References/TemplateFileActions.md` | The template workflows and this page |
| `docs/Development/References/SandboxIsolation.md`, `docs/Development/References/StructContracts.md` | The explanations of the structure's mechanics |
| `docs/Development/Tutorials/HandleDependencies.md`, `docs/Development/Tutorials/HandleLibElements.md`, `docs/Development/Tutorials/HandleCliCommands.md`, `docs/Development/Tutorials/HandleAdapters.md`, `docs/Development/Tutorials/HandleAssets.md`, `docs/Development/Tutorials/HandleSamples.md`, `docs/Development/Tutorials/HandleCliExamples.md`, `docs/Development/Tutorials/HandleDocuments.md` | The generic workflow guides for extending any library built on this structure |
| `sandbox/new.go` | The `New` constructor storing `Deps` on `api.Lib` and running the internal factories over it |
| `sandbox/contracts/deps/embeddeps/embeddeps.go`, `adapters/standard/embed.go`, `assets/asset.go` | The asset mechanic: the read-only contract, the factory serving the compiled-in files, and the `//go:embed all:*` directive taking the whole asset tree — generic, whatever the new library displays |
| `docs/LibUsage/References/EmbeddedAssets.md` | The explanation of that mechanic |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes. Every rewritten file must be shaped by the specification in its row, located through [Specs.md](/docs/Development/References/Specs.md).

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, badges, and the Doc Index pointing at each theme index | Readme |
| `sandbox/contracts/deps/deps.go` | The `Deps` function fields the new library requires | Deps |
| `sandbox/contracts/api/api.go` | The `Lib` struct and one struct per object the new library hands back | Outputs |
| `adapters/standard/standard.go` | The default adapter, filling the new `Deps` contract | Adapters |
| `assets/version.txt` | The new library's version, as its interface reports it | |
| `cmd/main/main.go` | The new library's entry point: wire, call `Sandboxmain`, exit | CliMain |
| `docs/Development/References/Structure.md` | The layout of the new library | Structure |
| `docs/CliUsage/References/Commands.md` | The commands, flags, and exit codes of the new library's interface | ReferenceDocs |
| `docs/LibUsage/References/PublicApi.md` | The index of the new public API entries | ReferenceDocs |
| `docs/LibUsage/References/Adapters.md` | The adapters the new library ships | AdaptersDoc |
| `docs/CliUsage/Index.md`, `docs/LibUsage/Index.md`, `docs/Development/Index.md`, `docs/Templating/Index.md` | The new library's page list, one index per theme | Index |
| `docs/CliUsage/Tutorials/QuickStart.md`, `docs/LibUsage/Tutorials/QuickStart.md` | The new library's first-run walkthroughs, CLI and Go | TutorialDocs |
| `docs/CliUsage/Tutorials/InstallCli.md`, `docs/CliUsage/Tutorials/UseCli.md` | Installing and driving the new library's CLI | TutorialDocs |
| `docs/LibUsage/Tutorials/LibInitialization.md` | Installing the new library and wiring its standard adapter | TutorialDocs |
| `docs/CliUsage/Tutorials/RunCliSample.md`, `docs/LibUsage/Tutorials/RunApiSample.md` | Running the new library's samples | TutorialDocs |
| `docs/CliUsage/References/SamplesList.md`, `docs/LibUsage/References/ApiSamplesList.md` | The new library's own samples | ReferenceDocs |

---

## Create

Written from scratch for the library being built or adapted. Nothing of the template's content survives here — the example files occupying these paths are removed by **[Delete](#delete)**. Every created file must be shaped by the specification in its row, located through [Specs.md](/docs/Development/References/Specs.md).

| Path | Description | Specification |
|------|-------------|---------------|
| `sandbox/internal/lib/*` | The lib's field factories and the `New` constructor running them all, reaching every dependency through `l.Deps` | LibFunctions |
| `sandbox/internal/<object>/*` | One package per object the library hands back: its field factories and the `New` constructor running them all | LibObjects |
| `sandbox/internal/cli/*` | The command dispatch behind `api.Lib.Sandboxmain`, the paths of the text it prints, and its operand parsing | |
| `assets/usages.txt`, `assets/messages/*` | The new interface's usage screen and one file per line it prints | |
| `adapters/<name>/<name>.go` | One adapter per additional opinionated implementation of the `Deps` contract | Adapters |
| `examples/libraryExamples/<example>/<example>.go` | One runnable Go sample per demonstrated use case | LibraryExamples |
| `examples/cliExamples/<Name>.sh` | One shell script per goal demonstrated against the built CLI | CliExamples |
| `docs/LibUsage/References/PublicApi/<pkg>.<Symbol>.md` | One detail page per public API entry | ReferenceDocs |
| `docs/<Theme>/Tutorials/<Goal>.md` | One tutorial per workflow specific to the new library — the generic guides carried over by **[Copy](#copy)** do **not** fulfill this | TutorialDocs |
| `docs/<Theme>/References/<Name>.md` | Any reference page the new library needs beyond the public API index | ReferenceDocs |

---

## Delete

The template's example content — the financial tracker. Removed once the new library's own files exist. For `.md` files, follow [HandleDocuments.md](/docs/Development/Tutorials/HandleDocuments.md) so the theme indexes stay in sync.

| Path | Description |
|------|-------------|
| `sandbox/internal/*` | The tracker's lib factories, object packages, CLI dispatch, and store helpers — replaced by **[Create](#create)** |
| `sandbox/contracts/deps/verbdeps/`, `sandbox/contracts/deps/keepdeps/` | The sandbox copies of the embedded Verb and Keep libraries — keep one only if the new library embeds the same library |
| `examples/libraryExamples/*` | The tracker's Go samples |
| `examples/cliExamples/*` | The tracker's CLI scripts |
| `docs/LibUsage/References/PublicApi/*` | The tracker's public API detail pages |
| `assets/usages.txt`, `assets/messages/*` | The tracker's usage screen and messages — replaced by **[Create](#create)** |
| `docs/LibUsage/Tutorials/ManageCategories.md`, `docs/LibUsage/Tutorials/TrackTransactions.md` | The tracker's domain tutorials |
| `bootstrap/*` | The embedded-library demonstration — keep it only as reference while the new library embeds another Agnos-Cli-style library |
