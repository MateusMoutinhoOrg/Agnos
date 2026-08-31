# AdaptersDoc Specification

## Description
Defines the required shape of `docs/References/Adapters.md` — the reference page listing every adapter shipped in `adapters/` and when to use each one. It builds on [ReferenceDocs](/docs/References/Specs/ReferenceDocs/Specs.md) and [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md); both still apply. This spec describes **how the page must be shaped**, not how an adapter is built — that is the [Adapters](/docs/References/Specs/Adapters/Specs.md) code specification.

### Rules
- **Every** directory under `adapters/` must have exactly one row in the page's table. Creating, renaming, or deleting an adapter requires updating the page in the same commit.
- Each row must state the adapter's name, its `New` factory (linked to the `docs/References/PublicApi/<pkg>.<Symbol>.md` detail page when one exists), how it fills each injected behavior, and when to use it.
- The page must not contain workflows — building an adapter is covered by [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-repo), which the page links to instead.
- A `Deps` field carrying a whole library, rather than a single behavior, must be explained under `## Embedded Libraries` — including why its factory returns a value rather than a closure.
- A `Deps` field no current library function calls must be listed under `## Standing Capabilities`, stating which factory fills it and what the shipped adapter backs it with. An adapter must fill every field regardless, so the page must say so.

## Structure
1. **Title** (H1): `Adapters`.
2. **`## Description`**: one short paragraph on what the page lists, linking to [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#creating-a-adapter-in-repo) for creating new adapters.
3. **`## Available Adapters`**: a Markdown table with the columns **Adapter**, **Factory**, **Behavior**, and **Use When** — one row per adapter directory.
4. **`## Embedded Libraries`**: how each whole-library field is filled and converted at the boundary.
5. **`## Standing Capabilities`** *(when the contract carries fields no library function calls)*: a Markdown table with the columns **Field**, **Filled by**, and one column per shipped adapter.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/AdaptersDoc/sample.md).
