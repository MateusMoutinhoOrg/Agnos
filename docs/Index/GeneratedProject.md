# Generated Project

## Description
Index of the documentation for people working **inside** a project `agnos` scaffolded: what its files are, which ones are theirs to write, how a command's handler receives its input, and how to give the sandbox a capability of their own. Driving `agnos` itself is indexed by [CliUsage.md](/docs/Index/CliUsage.md); contributing to Agnos-Cli — which is a project of this same shape — is indexed by [Development.md](/docs/Index/Development.md).

A generated project is a closed sandbox behind a generated command-line interface. The reader hand-writes exactly two kinds of file — `handler.go` for a command, and a contract-plus-lib pair for a capability — and `agnos build` writes the rest.

---

## Tutorials

- [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md)
  - **description:** Fill in a command's `handler.go`: typed entries, three output channels, exit codes
- [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md)
  - **description:** Give the sandbox an effect no installable dep provides, as a contract and a lib
- [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md)
  - **description:** Replace one field for a test, or assemble a different mix of adapter libs
- [ShapeCommands.md](/docs/Tutorials/ShapeCommands.md)
  - **description:** Declare a command's flags, args, aliases and help from the command line
  - [Add a flag](/docs/Tutorials/ShapeCommands.md#add-a-flag)
  - [Add a positional argument](/docs/Tutorials/ShapeCommands.md#add-a-positional-argument)
  - [Rewrite the command-level keys](/docs/Tutorials/ShapeCommands.md#rewrite-the-command-level-keys)
  - [Remove a field or a command](/docs/Tutorials/ShapeCommands.md#remove-a-field-or-a-command)
- [ManageDeps.md](/docs/Tutorials/ManageDeps.md)
  - **description:** Turn the dependency layer on and off, and install or remove capabilities

---

## References

- [EntriesYaml.md](/docs/References/EntriesYaml.md)
  - **description:** Every key of a command declaration and what the generated dispatch does with it
  - [Command Keys](/docs/References/EntriesYaml.md#command-keys)
  - [Field Keys](/docs/References/EntriesYaml.md#field-keys)
  - [Generated Types](/docs/References/EntriesYaml.md#generated-types)
  - [The quiet Flag](/docs/References/EntriesYaml.md#the-quiet-flag)
- [CommandDispatch.md](/docs/References/CommandDispatch.md)
  - **description:** How `climain.go` is generated, parses argv into `Entries`, and reaches a handler
  - [Three Files per Command](/docs/References/CommandDispatch.md#three-files-per-command)
  - [CliMain](/docs/References/CommandDispatch.md#climain)
  - [The Dispatch Function](/docs/References/CommandDispatch.md#the-dispatch-function)
  - [The Help Command](/docs/References/CommandDispatch.md#the-help-command)
  - [What the Sandbox Never Touches](/docs/References/CommandDispatch.md#what-the-sandbox-never-touches)
- [GeneratedFiles.md](/docs/References/GeneratedFiles.md)
  - **description:** Every file `agnos` writes into a project, and whether builds overwrite it
  - [Written by `start`](/docs/References/GeneratedFiles.md#written-by-start)
  - [Written by every `build`](/docs/References/GeneratedFiles.md#written-by-every-build)
  - [Written when `sandbox/deps/` exists](/docs/References/GeneratedFiles.md#written-when-sandboxdeps-exists)
  - [Written by `dep-install <dep>`](/docs/References/GeneratedFiles.md#written-by-dep-install-dep)
  - [Written when `sandbox/internal/cli/` exists](/docs/References/GeneratedFiles.md#written-when-sandboxinternalcli-exists)
  - [Written by `add-command <name>`](/docs/References/GeneratedFiles.md#written-by-add-command-name)
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
- [DepList.md](/docs/References/DepList.md)
  - **description:** Every installable dep, the contract and adapter it brings, and its module
  - [Deps](/docs/References/DepList.md#deps)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every adapter lib and assembly shipped, what backs it, and when to use it
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
  - [Adapter Libs](/docs/References/Adapters.md#adapter-libs)
  - [Embedded Libraries](/docs/References/Adapters.md#embedded-libraries)
  - [Standing Capabilities](/docs/References/Adapters.md#standing-capabilities)
