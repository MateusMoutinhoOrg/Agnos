# ReferenceDocs Specification

## Description
Defines the required shape of the `doc.md` of a **lookup** doc — a page listing enumerable items (commands, deps, keys, generated files, API entries), including the symbol sub-docs of `docs/PublicApi/`, and which is not one of the special documents ([Structure](/docs/Specs/Structure/doc.md), a generated [Index](/docs/Specs/Index/doc.md), the specifications index, or one of its sub-docs). A lookup doc is meant to be **scanned**, not read linearly.

### Rules
- Every page must comply with [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).
- The body must be **listable content**: Markdown tables or linked entry lists, one `##` section per group of items.
- Each entry must have a short description; when a detail page exists, the entry name must link to it.
- Reference pages must not contain workflows — link to the relevant tutorial instead.
- A symbol sub-doc of `PublicApi/` is titled with the qualified symbol in backticks, states its **Type**, and carries a `## Definition` (or `## Signature`) block, a `## Description`, a `## Fields` (or `## Parameters` / `## Returns`) table, and `## Examples`.
- Every new doc declares itself in its `props.yaml` — see [DocProps](/docs/Specs/DocProps/doc.md); a sub-doc carries no themes and is listed by its parent's generated `Index.md` and, for `PublicApi/`, by the grouped listing in its parent's `doc.md`.

## Structure
1. **Title** (H1): the name of what is being listed.
2. **`## Description`**: one short paragraph on what the page lists.
3. **One `##` section per item group**, separated by `---`, each containing a table or a list of entries with a name (linked when a detail page exists) and a short description.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/ReferenceDocs/sample.md).
