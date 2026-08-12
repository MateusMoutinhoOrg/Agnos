# Factories Specification

## Description
Defines the **factory pattern**, the single way any struct of function fields is filled in this project — on both sides of the sandbox wall. A factory takes a pointer to the struct that carries the state and returns exactly one field's value; the package's `New` constructor assigns that return value into the field. This spec describes the shape every factory shares; the per-tree specifications ([LibFunctions](/docs/References/Specs/LibFunctions/Specs.md), [LibObjects](/docs/References/Specs/LibObjects/Specs.md), [Adapters](/docs/References/Specs/Adapters/Specs.md)) build on it and add what is specific to their tree.

### Rules

#### Carrier and Target
- Every factory takes **one** parameter: a pointer to the **carrier** — the struct holding the state the closure needs. Inside `sandbox/internal/`, the carrier is the api struct being filled (`*api.Lib`, `*api.<Object>`). Inside `adapters/`, the carrier is the adapter struct, which declares a `Deps deps.Deps` field `New` assigns the factory's return value into.
- A factory returns the field's type. Its only job is to build and return that value; assignment happens at the call site in `New`.

#### One Field, One Factory
- One factory per field, named `<Field>Factory` after the field it fills — `GetCategoryFactory` fills `GetCategory`, `NowFactory` fills `Deps.Now`.
- The body returns a single value for that field, and touches no other field.
- A function field's factory returns a **closure**; a plain struct field's factory returns a **value** (e.g. `TrackerLibFactory` in `bootstrap/adapters/standard/`).
- The returned closure's signature must match the field's declaration in the contract exactly.

#### State Is Read Through the Pointer
- The closure reads state through the carrier pointer — `l.Deps.Now()`, `s.filePath` — so the value is resolved when the field is **called**, never captured when the factory ran.
- Copying a dependency into a local variable at factory time freezes it, and defeats the pattern.

#### `New` Is the Aggregate
- Every package exposing a filled struct declares a `New(...)` constructor that builds the carrier, calls **every** factory in the package exactly once and assigns its return value into the matching field, and returns the filled struct **by value** — `api.Lib`, `api.<Object>`, or `deps.Deps`. There is no separate `Factory` aggregate function.
- `New` never returns the carrier type itself when the carrier is an adapter: it returns `adapter.Deps`, the contract struct.
- Completeness is **not** checked by the compiler. A field no factory's return value is assigned into stays nil and panics on the first call, so keeping `New` complete is the author's job. See [StructContracts.md](/docs/References/StructContracts.md).

#### No Methods in Place of Factories
- A field is never filled by binding a method of the carrier. Methods may exist only as unexported helpers a closure calls (e.g. `fromKeepDatabase` in `adapters/standard/`).
- There is no internal mirror type and no `Api()` projection anywhere in the project.

## Structure
1. **Carrier struct**: the state holder — an api struct in [sandbox/contracts/api](/sandbox/contracts/api/), or the adapter struct with its `Deps` field in `adapters/<name>/<name>.go`.
2. **Field factories**: `func <Field>Factory(c *<Carrier>) <FieldType>`, one per field, each returning one closure or value.
3. **Doc comment**: one sentence per factory, naming the field it fills and what the returned value does.
4. **`New` constructor**: builds the carrier, calls every field factory exactly once and assigns its return value into the matching field, and returns the filled struct by value.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
