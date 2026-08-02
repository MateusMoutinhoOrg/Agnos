# LibFunctions Specification

## Description
Defines the required shape of a library function — a **factory** in `sandbox/internal/lib/` that returns the closure for one function field of the `Lib` struct declared in `sandbox/contracts/api/api.go`. A factory takes a pointer to the api struct and returns a closure; the caller (the package's `New` constructor) assigns it into the field, and the closure reaches dependencies only through that struct's `Deps`.

### Rules
- One factory per function field, named `<Field>Factory` and taking a single `*api.Lib` parameter: `func SetFactory(l *api.Lib) func(...)`. Its only job is to build and return the closure.
- The factory's body returns a closure for the field: `return func(...) ... { … }`. The closure's signature must match the field's declaration in `sandbox/contracts/api/api.go` exactly.
- Every factory must be called from the package's `New(d deps.Deps) api.Lib` constructor, which assigns its return value into the matching field and doubles as the factory aggregate — there is no separate `Factory` function. A field whose factory's return value is never assigned stays nil and panics on first call; the compiler does not check this.
- Dependencies are called **only** through `l.Deps.<Field>(...)` inside the closure — never construct or import a concrete implementation. Reading `l.Deps` inside the closure rather than capturing it at factory time is what lets the injected value stay authoritative.
- `sandbox/` is a closed sandbox: a factory must never import `adapters/`, `examples/`, a third-party module, or an OS-bound standard-library package (`os`, `net`, `syscall`, …). See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- A closure returning a library object returns that object's **api struct**, built by the object package's `New` constructor — see the [LibObjects](/docs/References/Meta/LibObjects/Specs.md) specification.
- A closure returning an optional object returns the api struct's **zero value** on the miss path; there is no nil struct to return.
- Factories and the fields they fill must have doc comments, and the fields must be listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package lib`.
2. **Field factory**: `func <Field>Factory(l *api.Lib) <FieldType>` returning a closure for `l.<Field>`, calling dependencies via `l.Deps.<Field>(...)`.
3. **Doc comment**: one sentence naming the field the factory fills and what the closure does.
4. **`New` constructor**: `func New(d deps.Deps) api.Lib` building `api.Lib{Deps: d}`, calling every field factory in the package exactly once and assigning its return value into the matching field, and returning the struct. `sandbox/new.go` does nothing but delegate to it.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
