# AdapterLib Specification

## Description
Defines the required shape of an **adapter lib** — a package under `adapters/libs/<lib>/` filling exactly one field of `deps.Deps` with a real implementation — and of its installable copy under `assets/deplist/<x>/adapters/libs/<lib>/`. It is the outside half of the pair whose inside half is the [DepsContract](/docs/References/Specs/DepsContract/Specs.md) specification.

### Rules
- One directory, one package, named after the **implementation** (`verb`, `keep`) or, when there is no library name to borrow, after the contract (`iodeps`, `std`). File names are `<Topic>.go`, capitalized (`Io.go`, `Run.go`, `Verb.go`).
- Exactly one exported function, `Bind(deps *deps.Deps)`, that assigns exactly one field: `deps.<X> = <x>.Lib{…}` or, for a per-call constructor, `deps.<X>.New = …`. `agnos build` generates a call to `Bind` under that name; any other name is not called.
- The package imports `sandbox/deps` for the struct to fill, `sandbox/deps/<x>` for the type, and whatever it needs — this is the one place `os`, `net/http`, `os/exec`, `embed` or a third-party module may appear. It imports nothing else from `sandbox/`.
- Every field of the sub-contract is assigned. A field the implementation cannot provide is filled with a function that returns a clear error, never left nil.
- Conversion at the boundary — copying a real library's struct onto the sandbox's copy, wrapping a returned object — is done here, field for field, in unexported helpers below `Bind`.
- Behavior the contract cannot express but a caller must know (a timeout, a file mode, merged stdout/stderr) is a documented decision of the lib, stated in a comment.
- The copy under `assets/deplist/` is byte-identical apart from `{{.Module}}` in place of this module's path.

## Structure
1. **Package clause**: `package <lib>`.
2. **Imports**: `sandbox/deps`, `sandbox/deps/<x>`, and the real implementation's packages.
3. **`Bind`**: the one exported function, assigning the one field.
4. **Helpers**: unexported functions the assigned closures call, each documented with the contract field it fills.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/AdapterLib/sample.go).
