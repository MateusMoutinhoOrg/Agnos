# Development

## Description
Index of the documentation for contributors changing this repository: the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy. Using the project is indexed by [CliUsage.md](/docs/Index/CliUsage.md) and [LibUsage.md](/docs/Index/LibUsage.md); working inside a project it generated is indexed by [GeneratedProject.md](/docs/Index/GeneratedProject.md).

> [!IMPORTANT]
> **Read before contributing.** [Structure.md](/docs/References/Structure.md) and [Specs.md](/docs/References/Specs.md) are required reading: they say **where** a change belongs and **how** the file you touch must be shaped. Agnos regenerates itself — every pattern is a template tomorrow, so consistency matters more than any individual feature.

---

## Tutorials

- [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md)
  - **description:** Regenerate this checkout with the binary compiled from it, then rebuild and install
- [Build.md](/docs/Tutorials/Build.md)
  - **description:** Cross-compile the CLI into a binary for each supported OS and architecture
  - [Build a single target](/docs/Tutorials/Build.md#build-a-single-target)
  - [Build every target at once](/docs/Tutorials/Build.md#build-every-target-at-once)
  - [Build with agnos compile](/docs/Tutorials/Build.md#build-with-agnos-compile)
  - [Add a new target](/docs/Tutorials/Build.md#add-a-new-target)
- [AddAction.md](/docs/Tutorials/AddAction.md)
  - **description:** Add a reusable operation: the two-file action, its contract field and its binder
- [AddAgnosCommand.md](/docs/Tutorials/AddAgnosCommand.md)
  - **description:** Expose an action on `agnos`'s own command line with the bootstrap binary
- [HandleDeplist.md](/docs/Tutorials/HandleDeplist.md)
  - **description:** Add a capability to the set `dep-install` can render, and pin its module
- [HandleAssetGroups.md](/docs/Tutorials/HandleAssetGroups.md)
  - **description:** Add a template to an asset group, and feed it a variable through a collector
  - [Add a template to a group](/docs/Tutorials/HandleAssetGroups.md#add-a-template-to-a-group)
  - [Add a variable through a collector](/docs/Tutorials/HandleAssetGroups.md#add-a-variable-through-a-collector)
  - [Add a single-file scaffold](/docs/Tutorials/HandleAssetGroups.md#add-a-single-file-scaffold)
- [HandleParsables.md](/docs/Tutorials/HandleParsables.md)
  - **description:** Add a five-file parser for one configuration file the actions read and write
- [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md)
  - **description:** Create, rename, move, or delete a `.md` file without leaving broken references
  - [Add a Document](/docs/Tutorials/HandleDocuments.md#add-a-document)
  - [Rename or Move a Document](/docs/Tutorials/HandleDocuments.md#rename-or-move-a-document)
  - [Delete a Document](/docs/Tutorials/HandleDocuments.md#delete-a-document)

---

## References

- [Structure.md](/docs/References/Structure.md)
  - **description:** The project's schema: which kind of file lives where, which are generated, and its spec
  - [Root](/docs/References/Structure.md#root)
  - [`/AgnosConfig/`](/docs/References/Structure.md#agnosconfig)
  - [`/sandbox/`](/docs/References/Structure.md#sandbox)
  - [`/adapters/`](/docs/References/Structure.md#adapters)
  - [`/assets/`](/docs/References/Structure.md#assets)
  - [`/cmd/`](/docs/References/Structure.md#cmd)
  - [`/release/`](/docs/References/Structure.md#release)
  - [`/old/`](/docs/References/Structure.md#old)
  - [`/docs/`](/docs/References/Structure.md#docs)
- [Specs.md](/docs/References/Specs.md)
  - **description:** Index of every specification and the files each one governs
  - [Documentation Specifications](/docs/References/Specs.md#documentation-specifications)
  - [Code Specifications](/docs/References/Specs.md#code-specifications)
- [Specs/](/docs/References/Specs/)
  - **description:** The specifications themselves — always reached through `Specs.md`, never browsed
- [BuildPipeline.md](/docs/References/BuildPipeline.md)
  - **description:** The verify gate, the collectors, the asset groups in order, persist, and the runtime
  - [The Verify Gate](/docs/References/BuildPipeline.md#the-verify-gate)
  - [Collectors](/docs/References/BuildPipeline.md#collectors)
  - [Template Variables](/docs/References/BuildPipeline.md#template-variables)
  - [Asset Groups, in Order](/docs/References/BuildPipeline.md#asset-groups-in-order)
  - [Persist, then Runtime](/docs/References/BuildPipeline.md#persist-then-runtime)
  - [Deps](/docs/References/BuildPipeline.md#deps)
  - [Self-Hosting](/docs/References/BuildPipeline.md#self-hosting)
- [SmartIO.md](/docs/References/SmartIO.md)
  - **description:** The transactional filesystem rooted at `--path` every action writes through
  - [One Root, Project-Relative Paths](/docs/References/SmartIO.md#one-root-project-relative-paths)
  - [The Transaction](/docs/References/SmartIO.md#the-transaction)
  - [Transaction-Aware Reads](/docs/References/SmartIO.md#transaction-aware-reads)
  - [Why Actions Compose](/docs/References/SmartIO.md#why-actions-compose)
- [SandboxIsolation.md](/docs/References/SandboxIsolation.md)
  - **description:** The sandbox wall: what `sandbox/` may not import, and why every effect is a dep
  - [The Three Trees](/docs/References/SandboxIsolation.md#the-three-trees)
  - [What the Wall Forbids](/docs/References/SandboxIsolation.md#what-the-wall-forbids)
  - [Why Every Door Is a Copy](/docs/References/SandboxIsolation.md#why-every-door-is-a-copy)
  - [What the Wall Forbids in the Other Direction](/docs/References/SandboxIsolation.md#what-the-wall-forbids-in-the-other-direction)
  - [Why the Entry Point Lives Inside](/docs/References/SandboxIsolation.md#why-the-entry-point-lives-inside)
- [StructContracts.md](/docs/References/StructContracts.md)
  - **description:** Why every contract is a struct of function fields, and how binders fill them
  - [The Shape](/docs/References/StructContracts.md#the-shape)
  - [Binders Fill the Fields](/docs/References/StructContracts.md#binders-fill-the-fields)
  - [Adapters Fill Their Contract the Same Way](/docs/References/StructContracts.md#adapters-fill-their-contract-the-same-way)
  - [Replacing One Behavior](/docs/References/StructContracts.md#replacing-one-behavior)
  - [What It Costs](/docs/References/StructContracts.md#what-it-costs)
- [CommandDispatch.md](/docs/References/CommandDispatch.md)
  - **description:** How `climain.go` is generated, parses argv into `Entries`, and reaches a handler
  - [Three Files per Command](/docs/References/CommandDispatch.md#three-files-per-command)
  - [CliMain](/docs/References/CommandDispatch.md#climain)
  - [The Dispatch Function](/docs/References/CommandDispatch.md#the-dispatch-function)
  - [The Help Command](/docs/References/CommandDispatch.md#the-help-command)
  - [What the Sandbox Never Touches](/docs/References/CommandDispatch.md#what-the-sandbox-never-touches)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every adapter lib and assembly shipped, what backs it, and when to use it
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
  - [Adapter Libs](/docs/References/Adapters.md#adapter-libs)
  - [Embedded Libraries](/docs/References/Adapters.md#embedded-libraries)
  - [Standing Capabilities](/docs/References/Adapters.md#standing-capabilities)
- [Commands.md](/docs/References/Commands.md)
  - **description:** Every command, flag, output channel and exit code of `agnos`
  - [Common Flags](/docs/References/Commands.md#common-flags)
  - [Core Commands](/docs/References/Commands.md#core-commands)
  - [Dependency System](/docs/References/Commands.md#dependency-system)
  - [Dependencies](/docs/References/Commands.md#dependencies)
  - [Cli System](/docs/References/Commands.md#cli-system)
  - [Info](/docs/References/Commands.md#info)
  - [Output Channels](/docs/References/Commands.md#output-channels)
  - [Exit Codes](/docs/References/Commands.md#exit-codes)
