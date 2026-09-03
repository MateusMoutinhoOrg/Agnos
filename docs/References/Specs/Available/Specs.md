# Available Specification

## Description
Defines the required shape of a hand-written **assembly** — `adapters/availables/<name>/new.go`, a ready-made `deps.Deps` wiring a chosen mix of adapter libs together. `availables/standard/new.go` has this shape but is **generated** by `agnos build` from the listing of `adapters/libs/`; this specification governs the assemblies a project adds beside it. Composing one is [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md).

### Rules
- One directory, one package, one file `new.go`, the package named after the directory.
- Exactly one exported function, `New(...) deps.Deps`, returning the contract struct by value — never the adapter's own type, of which there is none.
- The body creates an empty `deps.Deps{}`, calls one lib's `Bind(&deps)` per field it fills, patches any field it overrides **after** the binds, and returns. Configuration the assembly needs arrives as parameters of `New`.
- It imports `sandbox/deps` and adapter libs; it never imports `sandbox/` itself or anything under `sandbox/internal`.
- Every field of `deps.Deps` is filled, by a `Bind` or by an assignment. The compiler does not check this; the author compares against `sandbox/deps/deps.go` field by field.
- An assembly binding two libs that fill the same field is a mistake: the later call wins silently.

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: the adapter libs bound, and `sandbox/deps`.
3. **`New`**: empty struct, one `Bind` per lib, overrides, return.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/Available/sample.go).
