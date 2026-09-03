# DepsContract Specification

## Description
Defines the required shape of a **sub-contract** — `sandbox/deps/<x>/<x>.go`, the sandbox's copy of one outside library's api, injected whole as the `Deps.<X>` field — and of its installable copy under `assets/deplist/<x>/sandbox/deps/<x>/`. Why every effect is a copy is [SandboxIsolation.md](/docs/References/SandboxIsolation.md#why-every-door-is-a-copy).

### Rules
- One directory, one package, one file, all named `<x>` — lowercase, ending in `deps` by convention (`iodeps`, `rundeps`) unless the library has a name of its own (`serializables`, `std`). The `Deps` field is the title-cased directory name and is generated, never chosen.
- The file imports only the standard library and other `sandbox/deps` packages. `verify` rejects anything else — in particular the real library being copied.
- One exported struct named `Lib`, of function fields only, each documented with what it promises and how it reports failure. `Lib` is what the adapter lib assigns.
- A capability whose real object is created per call exposes a constructor field on `Lib` (`New func(args []string) Parser`, `NewRequest func(url string) Request`) and declares the returned object as another struct of function fields in the same file.
- Supporting data types (`RunProps`, `Result`, `File`) are plain structs of data fields, documented per field, declared after `Lib`.
- No method, no constructor, no state: the package is declaration only. The adapter lib fills it; nothing in the sandbox constructs it.
- The copy under `assets/deplist/` is byte-identical apart from `{{.Module}}` in place of this module's path.

## Structure
1. **Package clause**: `package <x>`.
2. **Package comment**: why this library is copied rather than imported, and what the standard adapter backs it with.
3. **`type Lib struct`**: the function fields, documented.
4. **Supporting types** *(optional)*: per-call objects as structs of function fields; data carriers as structs of data fields.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/DepsContract/sample.go).
