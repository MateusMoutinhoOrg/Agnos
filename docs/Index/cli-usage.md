# CliUsage
Documentation for people who drive agnos from a terminal - installing the binary, scaffolding a project, declaring its commands, and looking up what each command does

| Doc | Description |
| --- | --- |
| [Install the CLI](/docs/InstallCli/doc.md) | Install the CLI globally, or build and run it from a checkout |
| [Scaffold a Project](/docs/ScaffoldProject/doc.md) | Take an empty directory to a compiling CLI with one command of your own |
| [Shape a Command from the Command Line](/docs/ShapeCommands/doc.md) | Declare a command's flags, args, aliases and help from the command line |
| [Write a Command Handler](/docs/WriteCommandHandler/doc.md) | Fill in the one file of a command you write by hand, inside the sandbox |
| [Manage the Dependencies of a Project](/docs/ManageDeps/doc.md) | Turn the dependency layer on and off, and install or remove capabilities |
| [Add a Capability to a Project](/docs/AddAdapterLib/doc.md) | Write the contract and adapter pair for an effect no installable dep provides |
| [Regenerate and Check a Project](/docs/RegenerateProject/doc.md) | Check a project against the schema and re-render its generated files |
| [CLI Commands](/docs/Commands/doc.md) | Every command, flag, output channel and exit code of `agnos` |
| [The entries.yaml Declaration](/docs/EntriesYaml/doc.md) | Every key of a command declaration and what it makes the generated code do |
| [Installable Deps](/docs/DepList/doc.md) | Every installable dep, the contract and adapter it brings, and its module |
| [Generated Files](/docs/GeneratedFiles/doc.md) | Every file `agnos` writes into a project, and whether builds overwrite it |
