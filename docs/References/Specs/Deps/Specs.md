# Deps Specification

## Description
Defines the required shape of the dependency contract in `sandbox/contracts/deps/deps.go`. This spec describes **how the contract must be declared**, not which concrete dependencies a real library needs.

### Rules
- `deps.go` must declare a single `Deps` **struct** — the one contract every adapter fills. Contracts are structs, never interfaces; see [StructContracts.md](/docs/References/StructContracts.md).
- Every dependency must be a **function field** on that struct, so adapters inject behavior rather than data.
- Field names must be descriptive and exported; an unexported field cannot be filled by an adapter in another package.
- A dependency that is itself a library built with this pattern is declared as a **copy of that library's api struct** and injected as a plain field — never as a getter function. See `bootstrap/sandbox/contracts/deps/` for a worked case.
- `deps.go` must not import anything from `adapters/`, `examples/libraryExamples/`, `sandbox/internal/`, or `sandbox` (the entry point) — the contract stays free of implementations.

## Structure
1. **Package clause**: `package deps`.
2. **`Deps` struct**: a set of exported fields, one per requirement.
3. Each field is either a `Name func(...) ...` describing a single injectable behavior, or a nested contract struct for an injected library.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
