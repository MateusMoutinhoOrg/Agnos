# Adapters Specification

## Description
Defines the required shape of an adapter in `adapters/<name>/<name>.go`. This spec describes **how an adapter must be built**, not which concrete dependencies it fills in — those come from the [Deps contract](/docs/References/Meta/Deps/Specs.md). An adapter fills its contract with **factories**, the same way `sandbox/internal/` fills an api struct: the [Factories](/docs/References/Meta/Factories/Specs.md) specification applies on top of this one.

### Rules
- Each adapter lives in its own directory under `adapters/` and uses a package named after that directory.
- The adapter declares a **struct** holding its configuration and state, leading with a `Deps deps.Deps` field — the contract its factories assign into.
- One factory per `Deps` field, named `<Field>Factory` and taking a single pointer to the adapter struct: `func NowFactory(s *StandardAdapter) func() time.Time`. Its body returns a single closure for `s.Deps.<Field>`; `New` assigns the return value.
- The closure reads the adapter's configuration through the pointer (`s.filePath`), so state is resolved at call time — that is what carries the adapter's state into the library.
- Fields are **never** filled by binding a method of the adapter. Methods may exist only as unexported helpers a closure calls (e.g. `loadAll`).
- Each adapter exposes a single `New(...) deps.Deps` constructor as its entry point: it builds the adapter instance, calls **every** field factory and assigns its return value into the matching `Deps` field, and returns `adapter.Deps` — the populated **contract struct**, never the concrete adapter type.
- `New` must fill **every** field of `deps.Deps`. A field whose factory's return value is never assigned is a nil function that panics on the first call — the compiler will not catch it, so filling all of them is on the author.
- An adapter may import `sandbox/contracts/deps` but must never import `sandbox/internal/`. It may import `sandbox` (the entry point) only when it needs to initialize an embedded library to inject — see [StructContracts.md](/docs/Explanations/StructContracts.md).
- The adapter is the **opinionated** layer, and it lives outside the sandbox: all concrete choices (OS-bound stdlib, third-party libs, config) live here and nowhere else. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: at least `sandbox/contracts/deps`, plus whatever the implementation needs.
3. **Adapter struct**: a `Deps deps.Deps` field followed by the adapter-specific configuration and state.
4. **Field factories**: `func <Field>Factory(a *<Name>Adapter) <FieldType>`, one per `Deps` requirement, each returning one closure for `a.Deps.<Field>`.
5. **`New(...) deps.Deps` constructor**: accepts adapter-specific configuration, builds the adapter instance, calls every field factory exactly once and assigns its return value into the matching field, and returns `adapter.Deps`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
