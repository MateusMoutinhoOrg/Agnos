# Development

## Description
Index of the documentation for contributors changing this repository: the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy. Using the project is indexed by [CliUsage.md](/docs/Index/CliUsage.md) and [LibUsage.md](/docs/Index/LibUsage.md); turning the project into a new library is indexed by [Templating.md](/docs/Index/Templating.md).

> [!IMPORTANT]
> **Read before contributing.** [RULES.md](/docs/References/RULES.md), [Structure.md](/docs/References/Structure.md), and [Specs.md](/docs/References/Specs.md) are required reading: they say what is allowed, **where** a change belongs, and **how** the file you touch must be shaped.

---

## Tutorials

| Doc | Description |
| --- | --- |
| [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md) | How injected deps travel the object graph, and how to add a `Deps` field |
| [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md) | Add a function or an object: declare, write the factory, register, publish |
| [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md) | Add a command or a flag to the interface behind `api.Lib.Sandboxmain` |
| [HandleAdapters.md](/docs/Tutorials/HandleAdapters.md) | Create a new opinionated implementation of the `Deps` contract |
| [HandleAssets.md](/docs/Tutorials/HandleAssets.md) | Add or edit a file under `assets/` the library reads instead of a Go string |
| [HandleSamples.md](/docs/Tutorials/HandleSamples.md) | Create and run executable Go samples in `examples/libraryExamples/` |
| [HandleCliExamples.md](/docs/Tutorials/HandleCliExamples.md) | Create and run shell scripts in `examples/cliExamples/` driving the built CLI |
| [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md) | Create, rename, move, or delete a `.md` file without leaving broken references |
| [Build.md](/docs/Tutorials/Build.md) | Cross-compile the CLI into a binary for each supported OS and architecture |

---

## References

| Doc | Description |
| --- | --- |
| [RULES.md](/docs/References/RULES.md) | The binding rules every change must follow to be accepted |
| [Structure.md](/docs/References/Structure.md) | The project's schema: which kind of file lives where, and its spec |
| [Specs.md](/docs/References/Specs.md) | Index of every specification and the files each one governs |
| [Specs/](/docs/References/Specs/) | The specifications themselves — always reached through `Specs.md`, never browsed |
| [SandboxIsolation.md](/docs/References/SandboxIsolation.md) | The sandbox wall: what `sandbox/` may not import, and why effects are deps |
| [StructContracts.md](/docs/References/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them |
