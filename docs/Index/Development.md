# Development

## Description
Index of the documentation for contributors changing this repository: the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy. Using the project is indexed by [CliUsage.md](/docs/Index/CliUsage.md) and [LibUsage.md](/docs/Index/LibUsage.md); turning the project into a new library is indexed by [Templating.md](/docs/Index/Templating.md).

> [!IMPORTANT]
> **Read before contributing.** [Structure.md](/docs/References/Structure.md) and [Specs.md](/docs/References/Specs.md) are required reading: they say **where** a change belongs and **how** the file you touch must be shaped.

---

## Tutorials

- [Build.md](/docs/Tutorials/Build.md)
  - **description:** Cross-compile the CLI into a binary for each supported OS and architecture
  - [Build a single target](/docs/Tutorials/Build.md#build-a-single-target)
  - [Build every target at once](/docs/Tutorials/Build.md#build-every-target-at-once)
  - [Add a new target](/docs/Tutorials/Build.md#add-a-new-target)
- [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md)
  - **description:** Add a function to the internal logic or public API: declare, write the factory, register, publish
  - [AddLibFunction](/docs/Tutorials/HandleLibElements.md#addlibfunction)
  - [AddPublicLibFunction](/docs/Tutorials/HandleLibElements.md#addpubliclibfunction)
- [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md)
  - **description:** How injected deps travel the object graph, and how to add a `Deps` field
  - [Find Dependencies Functions you can use](/docs/Tutorials/HandleDependencies.md#find-dependencies-functions-you-can-use)
  - [Add New Dependencie](/docs/Tutorials/HandleDependencies.md#add-new-dependencie)
  - [Overwrinting a adapter function](/docs/Tutorials/HandleDependencies.md#overwrinting-a-adapter-function)
  - [Creating a adapter in repo](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-repo)
  - [Creating a adapter in your project](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-your-project)
- [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md)
  - **description:** Add a command or a flag to the interface behind `api.Lib.Sandboxmain`
  - [Add CLI Command](/docs/Tutorials/HandleCliCommands.md#add-cli-command)
  - [Remove CLI Command](/docs/Tutorials/HandleCliCommands.md#remove-cli-command)
- [HandleAssets.md](/docs/Tutorials/HandleAssets.md)
  - **description:** Add or edit a file under `assets/` the library reads instead of a Go string
  - [Add Asset](/docs/Tutorials/HandleAssets.md#add-asset)
  - [ListAssets in Runtime](/docs/Tutorials/HandleAssets.md#listassets-in-runtime)
  - [Retrieve Asset in Runtime](/docs/Tutorials/HandleAssets.md#retrieve-asset-in-runtime)
- [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md)
  - **description:** Create and run executable Go samples
  - [Run a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#run-a-library-sample)
  - [Add a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample)
- [HandleCliExamples.md](/docs/Tutorials/HandleCliExamples.md)
  - **description:** Create and run shell scripts in `examples/cliExamples/` driving the built CLI
  - [Run a CLI Example](/docs/Tutorials/HandleCliExamples.md#run-a-cli-example)
  - [Add a CLI Example](/docs/Tutorials/HandleCliExamples.md#add-a-cli-example)
- [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md)
  - **description:** Create, rename, move, or delete a `.md` file without leaving broken references
  - [Add a Document](/docs/Tutorials/HandleDocuments.md#add-a-document)
  - [Rename or Move a Document](/docs/Tutorials/HandleDocuments.md#rename-or-move-a-document)
  - [Delete a Document](/docs/Tutorials/HandleDocuments.md#delete-a-document)

---

## References

- [Structure.md](/docs/References/Structure.md)
  - **description:** The project's schema: which kind of file lives where, and its spec
  - [Root](/docs/References/Structure.md#root)
  - [`/scripts/`](/docs/References/Structure.md#scripts)
  - [`/sandbox/`](/docs/References/Structure.md#sandbox)
  - [`/adapters/`](/docs/References/Structure.md#adapters)
  - [`/assets/`](/docs/References/Structure.md#assets)
  - [`/cmd/`](/docs/References/Structure.md#cmd)
  - [`/examples/cliExamples/`](/docs/References/Structure.md#examplescliexamples)
  - [`/examples/libraryExamples/`](/docs/References/Structure.md#exampleslibraryexamples)
  - [`/docs/`](/docs/References/Structure.md#docs)
- [Specs.md](/docs/References/Specs.md)
  - **description:** Index of every specification and the files each one governs
  - [Documentation Specifications](/docs/References/Specs.md#documentation-specifications)
  - [Code Specifications](/docs/References/Specs.md#code-specifications)
- [Specs/](/docs/References/Specs/)
  - **description:** The specifications themselves — always reached through `Specs.md`, never browsed
- [SandboxIsolation.md](/docs/References/SandboxIsolation.md)
  - **description:** The sandbox wall: what `sandbox/` may not import, and why effects are deps
  - [The Three Trees](/docs/References/SandboxIsolation.md#the-three-trees)
  - [What the Wall Forbids](/docs/References/SandboxIsolation.md#what-the-wall-forbids)
  - [What the Wall Forbids in the Other Direction](/docs/References/SandboxIsolation.md#what-the-wall-forbids-in-the-other-direction)
  - [Why the Entry Point Lives Inside](/docs/References/SandboxIsolation.md#why-the-entry-point-lives-inside)
- [StructContracts.md](/docs/References/StructContracts.md)
  - **description:** Why every contract is a struct of function fields, and how factories fill them
  - [The Shape](/docs/References/StructContracts.md#the-shape)
  - [Factories Fill the Fields](/docs/References/StructContracts.md#factories-fill-the-fields)
  - [Adapters Fill Their Contract the Same Way](/docs/References/StructContracts.md#adapters-fill-their-contract-the-same-way)
  - [Replacing One Behavior](/docs/References/StructContracts.md#replacing-one-behavior)
  - [Consuming a Library That Uses This Pattern](/docs/References/StructContracts.md#consuming-a-library-that-uses-this-pattern)
  - [What It Costs](/docs/References/StructContracts.md#what-it-costs)
