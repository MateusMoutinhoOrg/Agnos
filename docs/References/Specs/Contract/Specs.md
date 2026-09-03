# Contract Specification

## Description
Defines the required shape of a hand-written file of `sandbox/api/` — a **contract**: a struct of function fields the sandbox hands back to callers, plus the props structs and constants those fields need. `actions.go` is the one such file today; `sandbox.go` and `cli.go` are generated. Why contracts are structs is [StructContracts.md](/docs/References/StructContracts.md).

### Rules
- `package api`. The file imports **nothing** — not the standard library, not `sandbox/deps`, not an external module. `verify` rejects any import outside `sandbox/api`. Signatures therefore use only builtin types and other `api` structs.
- One exported struct named after the file (`actions.go` → `Actions`) whose fields are **all** functions, one per operation, alphabetically or by subsystem; no methods, no state.
- A function taking more than three parameters takes one props struct instead, declared in the same file above the contract and named `<Operation>Props`. Optional-versus-empty values use a pointer (`Module *string`) or a raw string literal (`Default string`, `""` unset), never a sentinel.
- Constants a caller must spell (`RuntimeGo`, `ExitUsage`) are declared here, documented with what each one means.
- Every field is mirrored by exactly one assignment in `sandbox/binds/<same file>.go` — the Binder specification — and every new field or props struct gets a detail page under `docs/References/PublicApi/`.

## Structure
1. **Package clause**: `package api`.
2. **Constants** *(optional)*: documented `const` blocks the contract's fields take as values.
3. **Props structs** *(optional)*: one per multi-parameter operation, fields documented.
4. **The contract struct**: `type <Name> struct` of function fields.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/Contract/sample.go).
