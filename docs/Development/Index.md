# Development

## Description
Index of the documentation for contributors changing this repository: the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy. Using the project is indexed by [CliUsage/Index.md](/docs/CliUsage/Index.md) and [LibUsage/Index.md](/docs/LibUsage/Index.md); turning the project into a new library is indexed by [Templating/Index.md](/docs/Templating/Index.md).

> [!IMPORTANT]
> **Read before contributing.** [RULES.md](/docs/Development/References/RULES.md), [Structure.md](/docs/Development/References/Structure.md), and [Specs.md](/docs/Development/References/Specs.md) are required reading: they say what is allowed, **where** a change belongs, and **how** the file you touch must be shaped.

---

## Protocols

| Doc | Description |
| --- | --- |
| [HandleDependencies.md](/docs/Development/Protocols/HandleDependencies.md) | How injected deps travel the object graph, and how to add a `Deps` field |
| [HandleLibElements.md](/docs/Development/Protocols/HandleLibElements.md) | Add a function or an object: declare, write the factory, register, publish |
| [HandleCliCommands.md](/docs/Development/Protocols/HandleCliCommands.md) | Add a command or a flag to the interface behind `api.Lib.Sandboxmain` |
| [HandleAdapters.md](/docs/Development/Protocols/HandleAdapters.md) | Create a new opinionated implementation of the `Deps` contract |
| [HandleSamples.md](/docs/Development/Protocols/HandleSamples.md) | Create and run executable Go samples in `examples/libraryExamples/` |
| [HandleCliExamples.md](/docs/Development/Protocols/HandleCliExamples.md) | Create and run shell scripts in `examples/cliExamples/` driving the built CLI |
| [HandleDocuments.md](/docs/Development/Protocols/HandleDocuments.md) | Create, rename, move, or delete a `.md` file without leaving broken references |

---

## References

| Doc | Description |
| --- | --- |
| [RULES.md](/docs/Development/References/RULES.md) | The binding rules every change must follow to be accepted |
| [Structure.md](/docs/Development/References/Structure.md) | The project's schema: which kind of file lives where, and its spec |
| [Specs.md](/docs/Development/References/Specs.md) | Index of every specification and the files each one governs |
| [Specs/](/docs/Development/References/Specs/) | The specifications themselves — always reached through `Specs.md`, never browsed |
| [SandboxIsolation.md](/docs/Development/References/SandboxIsolation.md) | The sandbox wall: what `sandbox/` may not import, and why effects are deps |
| [StructContracts.md](/docs/Development/References/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them |
