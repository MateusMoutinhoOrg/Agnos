# Public API

## Description
Index of all public-facing components of the library, grouped by kind, with links to their detail pages.

---

## Interfaces

### [api.Lib](/docs/References/PublicApi/api.Lib.md)
The library entry point. Returned by `lib.New`; exposes all library methods.

### `api.ExampleLibObject`
An object created by the library with its dependencies automatically wired in. Shown unlinked because this sample's library is hypothetical and no detail page exists — an entry links to its page only when that page is real.

---

## Functions

### [lib.New](/docs/References/PublicApi/lib.New.md)
Injects a `deps.Deps` implementation into the library and returns an `api.Lib`.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Creates a `deps.Deps` using the standard library adapter.
