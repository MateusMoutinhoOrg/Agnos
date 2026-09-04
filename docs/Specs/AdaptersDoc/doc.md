# AdaptersDoc Specification

## Description
Defines the required shape of `docs/Adapters/doc.md` — the lookup doc listing every adapter lib under `adapters/libs/` and every assembly under `adapters/availables/`, and when to use each one. It builds on [ReferenceDocs](/docs/Specs/ReferenceDocs/doc.md) and [GeneralDoc](/docs/Specs/GeneralDoc/doc.md); both still apply. This spec describes **how the page must be shaped**, not how a lib is built — that is the [AdapterLib](/docs/Specs/AdapterLib/doc.md) code specification.

### Rules
- **Every** directory under `adapters/availables/` must have exactly one row in the `Available Adapters` table, and **every** directory under `adapters/libs/` exactly one row in the `Adapter Libs` table. Creating, renaming, or deleting either requires updating the page in the same commit.
- An assembly row states its name, its `New` factory (linked to the `docs/PublicApi/<pkg>.<Symbol>/doc.md` sub-doc when one exists), what it binds, and when to use it. A lib row states its package, the `Deps` field it fills (linked to the sub-contract's detail page), what backs it, and any behavior a caller must know.
- The page must not contain workflows — building a lib is covered by [AddAdapterLib](/docs/AddAdapterLib/doc.md), which the page links to instead.
- A `Deps` field carrying a whole outside library, rather than a set of behaviors written here, must be explained under `## Embedded Libraries`.
- A `Deps` field no current action calls must be listed under `## Standing Capabilities`, stating which binder fills it and what the shipped assembly backs it with. An assembly must fill every field regardless, so the page must say so.

## Structure
1. **Title** (H1): `Adapters`.
2. **`## Description`**: one short paragraph on what the page lists, linking to [AddAdapterLib](/docs/AddAdapterLib/doc.md) for creating new libs.
3. **`## Available Adapters`**: a Markdown table with the columns **Adapter**, **Factory**, **Behavior**, and **Use When** — one row per `availables/` directory.
4. **`## Adapter Libs`**: a Markdown table with the columns **Lib**, **Fills**, **Backed by**, and **Notes** — one row per `libs/` directory.
5. **`## Embedded Libraries`**: how each whole-library field is filled and converted at the boundary.
6. **`## Standing Capabilities`** *(when the contract carries fields no action calls)*: a Markdown table with the columns **Field**, **Filled by**, and one column per shipped assembly.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/AdaptersDoc/sample.md).
