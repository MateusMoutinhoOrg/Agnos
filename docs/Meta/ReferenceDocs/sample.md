# Public API

## Description
Index of all public-facing components of the library, grouped by kind, with links to their detail pages.

---

## Interfaces

### [api.Lib](/docs/PublicApi/api.Lib.md)
The library entry point. Returned by `lib.New`; exposes all library methods.

### [api.ExampleLibObject](/docs/PublicApi/api.ExampleLibObject.md)
An object created by the library with its dependencies automatically wired in.

---

## Functions

### [lib.New](/docs/PublicApi/lib.New.md)
Injects a `deps.Deps` implementation into the library and returns an `api.Lib`.

### [standard.New](/docs/PublicApi/standard.New.md)
Creates a `deps.Deps` using the standard library adapter.
