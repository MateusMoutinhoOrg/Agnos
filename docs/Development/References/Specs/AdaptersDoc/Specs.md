# AdaptersDoc Specification

## Description
Defines the required shape of `docs/LibUsage/References/Adapters.md` — the reference page listing every adapter shipped in `adapters/` and when to use each one. It builds on [ReferenceDocs](/docs/Development/References/Specs/ReferenceDocs/Specs.md) and [GeneralDoc](/docs/Development/References/Specs/GeneralDoc/Specs.md); both still apply. This spec describes **how the page must be shaped**, not how an adapter is built — that is the [Adapters](/docs/Development/References/Specs/Adapters/Specs.md) code specification.

### Rules
- **Every** directory under `adapters/` must have exactly one row in the page's table. Creating, renaming, or deleting an adapter requires updating the page in the same commit.
- Each row must state the adapter's name, its `New` factory (linked to the `docs/LibUsage/References/PublicApi/<pkg>.<Symbol>.md` detail page when one exists), how it fills each injected behavior, and when to use it.
- The page must not contain workflows — building an adapter is covered by [HandleAdapters.md](/docs/Development/Protocols/HandleAdapters.md), which the page links to instead.

## Structure
1. **Title** (H1): `Adapters`.
2. **`## Description`**: one short paragraph on what the page lists, linking to [HandleAdapters.md](/docs/Development/Protocols/HandleAdapters.md) for creating new adapters.
3. **`## Available Adapters`**: a Markdown table with the columns **Adapter**, **Factory**, **Behavior**, and **Use When** — one row per adapter directory.

> **Note**: For a concrete example, refer to [sample.md](/docs/Development/References/Specs/AdaptersDoc/sample.md).
