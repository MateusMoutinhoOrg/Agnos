# LibUsage
Documentation for developers consuming Agnos as a Go module - wiring an adapter into the sandbox, calling the same actions the CLI exposes from code, and the public API

| Doc | Description |
| --- | --- |
| [Library Initialization](/docs/LibInitialization/doc.md) | Install the module, create deps via an adapter, and scaffold a project from Go |
| [Shape a Project's Commands from Go](/docs/ShapeProjectFromGo/doc.md) | Declare a project's commands, flags and args through the action props structs |
| [Compose a Custom Deps](/docs/ComposeDeps/doc.md) | Replace one field for a test, or assemble a different mix of adapter libs |
| [Public API](/docs/PublicApi/doc.md) | Every public-facing entry of the library, grouped by role, with a page per symbol |
| [Adapters](/docs/Adapters/doc.md) | Every adapter lib and assembly shipped, what backs it, and when to use it |
| [Struct Contracts](/docs/StructContracts/doc.md) | Why every contract is a struct of function fields, and how binders fill them |
