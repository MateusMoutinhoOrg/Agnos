# CLI Usage

## Description
Index of the documentation for people who drive `agnos` from a terminal: installing the binary, scaffolding a project with it, declaring that project's commands, and looking up what each `agnos` command does. Working inside the project it wrote is indexed by [GeneratedProject.md](/docs/Index/GeneratedProject.md); running the same operations from Go code is indexed by [LibUsage.md](/docs/Index/LibUsage.md).

Every command below ends by regenerating and compiling the project it touched, so a step that finishes is a step that left the project buildable.

---

## Tutorials

- [InstallCli.md](/docs/Tutorials/InstallCli.md)
  - **description:** Install the CLI globally, or build and run it from a checkout
  - [macOS](/docs/Tutorials/InstallCli.md#macos)
  - [Linux](/docs/Tutorials/InstallCli.md#linux)
  - [Windows (PowerShell)](/docs/Tutorials/InstallCli.md#windows-powershell)
  - [Verify after reboot](/docs/Tutorials/InstallCli.md#verify-after-reboot)
  - [Troubleshooting](/docs/Tutorials/InstallCli.md#troubleshooting)
  - [Install from a Clone (Requires Go)](/docs/Tutorials/InstallCli.md#install-from-a-clone-requires-go)
- [ScaffoldProject.md](/docs/Tutorials/ScaffoldProject.md)
  - **description:** Take an empty directory to a compiling CLI with one command of your own
- [ShapeCommands.md](/docs/Tutorials/ShapeCommands.md)
  - **description:** Declare a command's flags, args, aliases and help from the command line
  - [Add a flag](/docs/Tutorials/ShapeCommands.md#add-a-flag)
  - [Add a positional argument](/docs/Tutorials/ShapeCommands.md#add-a-positional-argument)
  - [Rewrite the command-level keys](/docs/Tutorials/ShapeCommands.md#rewrite-the-command-level-keys)
  - [Remove a field or a command](/docs/Tutorials/ShapeCommands.md#remove-a-field-or-a-command)
- [ManageDeps.md](/docs/Tutorials/ManageDeps.md)
  - **description:** Turn the dependency layer on and off, and install or remove capabilities
- [RegenerateProject.md](/docs/Tutorials/RegenerateProject.md)
  - **description:** Check a project against the schema and re-render its generated files

---

## References

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
- [DepList.md](/docs/References/DepList.md)
  - **description:** Every installable dep, the contract and adapter it brings, and its module
  - [Deps](/docs/References/DepList.md#deps)
- [GeneratedFiles.md](/docs/References/GeneratedFiles.md)
  - **description:** Every file `agnos` writes into a project, and whether builds overwrite it
  - [Written by `start`](/docs/References/GeneratedFiles.md#written-by-start)
  - [Written by every `build`](/docs/References/GeneratedFiles.md#written-by-every-build)
  - [Written when `sandbox/deps/` exists](/docs/References/GeneratedFiles.md#written-when-sandboxdeps-exists)
  - [Written by `dep-install <dep>`](/docs/References/GeneratedFiles.md#written-by-dep-install-dep)
  - [Written when `sandbox/internal/cli/` exists](/docs/References/GeneratedFiles.md#written-when-sandboxinternalcli-exists)
  - [Written by `add-command <name>`](/docs/References/GeneratedFiles.md#written-by-add-command-name)
