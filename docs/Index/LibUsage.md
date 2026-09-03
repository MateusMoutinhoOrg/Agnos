# Library Usage

## Description
Index of the documentation for developers consuming Agnos as a Go module: wiring an adapter into the sandbox, calling the same actions the CLI exposes from code, and looking up the public API. Driving the same behavior from a terminal is indexed by [CliUsage.md](/docs/Index/CliUsage.md); changing the library is indexed by [Development.md](/docs/Index/Development.md).

The library is always built the same way: `standard.New()` produces a `deps.Deps`, `sandbox.New(&deps)` injects it into the closed sandbox, and the returned `*api.Sandbox` carries every action and the whole command-line interface.

---

## Tutorials

- [LibInitialization.md](/docs/Tutorials/LibInitialization.md)
  - **description:** Install the module, create deps via an adapter, and scaffold a project from Go
- [ShapeProjectFromGo.md](/docs/Tutorials/ShapeProjectFromGo.md)
  - **description:** Declare a project's commands, flags and args through the action props structs
- [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md)
  - **description:** Replace one field for a test, or assemble a different mix of adapter libs

---

## References

- [PublicApi.md](/docs/References/PublicApi.md)
  - **description:** Index of every public-facing entry of the library, grouped by role, with detail pages
  - [Entry Points](/docs/References/PublicApi.md#entry-points)
  - [Core Interface](/docs/References/PublicApi.md#core-interface)
  - [Actions](/docs/References/PublicApi.md#actions)
  - [Props](/docs/References/PublicApi.md#props)
  - [Dependency Contracts](/docs/References/PublicApi.md#dependency-contracts)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every adapter lib and assembly shipped, what backs it, and when to use it
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
  - [Adapter Libs](/docs/References/Adapters.md#adapter-libs)
  - [Embedded Libraries](/docs/References/Adapters.md#embedded-libraries)
  - [Standing Capabilities](/docs/References/Adapters.md#standing-capabilities)
- [StructContracts.md](/docs/References/StructContracts.md)
  - **description:** Why every contract is a struct of function fields, and how binders fill them
  - [The Shape](/docs/References/StructContracts.md#the-shape)
  - [Binders Fill the Fields](/docs/References/StructContracts.md#binders-fill-the-fields)
  - [Adapters Fill Their Contract the Same Way](/docs/References/StructContracts.md#adapters-fill-their-contract-the-same-way)
  - [Replacing One Behavior](/docs/References/StructContracts.md#replacing-one-behavior)
  - [What It Costs](/docs/References/StructContracts.md#what-it-costs)
