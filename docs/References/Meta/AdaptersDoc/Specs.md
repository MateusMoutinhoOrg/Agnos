# AdaptersDoc Specification

## Description
Defines the required shape of `docs/References/Adapters.md` — the reference page listing every adapter shipped in `adapters/` and when to use each one. It builds on [ReferenceDocs](/docs/References/Meta/ReferenceDocs/Specs.md) and [GeneralDoc](/docs/References/Meta/GeneralDoc/Specs.md); both still apply. This spec describes **how the page must be shaped**, not how an adapter is built — that is the [Adapters](/docs/References/Meta/Adapters/Specs.md) code specification.

### Rules
- **Every** directory under `adapters/` must have exactly one row in the page's table. Creating, renaming, or deleting an adapter requires updating the page in the same commit.
- Each row must state the adapter's name, its `New` factory (linked to the detail page in `docs/References/PublicApi/` when one exists), how it fills each injected behavior, and when to use it.
- The page must not contain workflows — building an adapter is covered by [AddAdapter.md](/docs/Tutorials/AddAdapter.md), which the page links to instead.

## Structure
1. **Title** (H1): `Adapters`.
2. **`## Description`**: one short paragraph on what the page lists, linking to [AddAdapter.md](/docs/Tutorials/AddAdapter.md) for creating new adapters.
3. **`## Available Adapters`**: a Markdown table with the columns **Adapter**, **Factory**, **Behavior**, and **Use When** — one row per adapter directory.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Meta/AdaptersDoc/sample.md).
