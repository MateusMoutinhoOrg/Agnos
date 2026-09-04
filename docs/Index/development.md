# Development
Documentation for contributors changing this repository - the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy

| Doc | Description |
| --- | --- |
| [Add a Capability to a Project](/docs/AddAdapterLib/doc.md) | Write the contract and adapter pair for an effect no installable dep provides |
| [Bootstrap Agnos with Itself](/docs/BootstrapAgnos/doc.md) | Regenerate this checkout with the binary compiled from it, then rebuild and install |
| [Build the CLI for Every Architecture](/docs/Build/doc.md) | Cross-compile the CLI into a binary for each supported OS and architecture |
| [Add an Action](/docs/AddAction/doc.md) | Add a reusable operation: the two-file action, its contract field and its binder |
| [Add a Command to Agnos](/docs/AddAgnosCommand/doc.md) | Expose an action on `agnos`'s own command line with the bootstrap binary |
| [Add an Installable Dep](/docs/HandleDeplist/doc.md) | Add a capability to the set `dep-install` can render, and pin its module |
| [Handle Asset Groups and Collectors](/docs/HandleAssetGroups/doc.md) | Add a template to an asset group, and feed it a variable through a collector |
| [Add a Parsable Config](/docs/HandleParsables/doc.md) | Add a five-file parser for one configuration file the actions read and write |
| [Handle Documents](/docs/HandleDocuments/doc.md) | Create, rename, move, or delete a doc without leaving broken indexes or links |
| [CLI Commands](/docs/Commands/doc.md) | Every command, flag, output channel and exit code of `agnos` |
| [Adapters](/docs/Adapters/doc.md) | Every adapter lib and assembly shipped, what backs it, and when to use it |
| [Struct Contracts](/docs/StructContracts/doc.md) | Why every contract is a struct of function fields, and how binders fill them |
| [Project Structure](/docs/Structure/doc.md) | The project's schema: which kind of file lives where, which are generated, and its spec |
| [Specifications](/docs/Specs/doc.md) | Every specification and the files each one governs, with a page per specification |
| [The Build Pipeline](/docs/BuildPipeline/doc.md) | The verify gate, the collectors, the asset groups in order, persist, and the runtime |
| [SmartIO](/docs/SmartIO/doc.md) | The transactional filesystem rooted at `--path` every action writes through |
| [Sandbox Isolation](/docs/SandboxIsolation/doc.md) | The sandbox wall: what `sandbox/` may not import, and why every effect is a dep |
| [Command Dispatch](/docs/CommandDispatch/doc.md) | How `climain.go` is generated, parses argv into `Entries`, and reaches a handler |
