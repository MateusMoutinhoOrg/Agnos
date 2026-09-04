# Index Specification

## Description
Defines the shape of the **generated indexes** of the documentation tree: one theme index per theme under `docs/Index/<theme-id>.md`, and one `Index.md` inside every doc directory that has sub-docs. Both are written by `agnos build` from the `props.yaml` of the docs they list — they are never edited by hand. This specification exists so a reader can tell a correct index from a stale one, and so the templates that emit them (`assets/templates/theme_index.md`, `assets/templates/doc_index.md`) keep one shape.

A theme is a reader goal, not a directory: every doc is a first-level directory of `docs/`, and the index is what groups them — see [DocProps](/docs/Specs/DocProps/doc.md).

### Rules
- Every page must comply with [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).
- **Generated, never hand-written.** `docs/Index/` belongs to the generator end to end: it is removed and rewritten on every build, so an index of a deleted theme never survives. A doc's `Index.md` is rewritten on every build too. Editing either is pointless.
- One theme index per entry of `<ProjectName>Config/themes.yaml`, named after that theme's `id` — `cli-usage.md`, not `CliUsage.md`.
- A theme index lists **every** first-level doc whose `props.yaml` names that theme id, and nothing else. Sub-docs are listed only by their parent's `Index.md`.
- A doc directory gets an `Index.md` only when it has sub-docs; it lists its **direct** sub-docs, at any depth of the tree.
- Rows are ordered by `order`, then by name. A doc with no `order` is listed after every ordered one.
- Each row is `| [<Doc Name>](<link>) | <description> |`, with the name and description taken verbatim from the doc's `props.yaml` and the link repository-rooted at that doc's own `doc.md`.
- The reader reaches a doc only through an index: the [README](/README.md) links to the theme indexes, a theme index links to its docs, and a doc with sub-docs links to its `Index.md` once from its `doc.md`.

## Structure
1. **Title** (H1): the theme's `name` — or, in a doc's `Index.md`, the doc's `name` followed by `Index`.
2. **Description line**: the theme's or doc's `description`, verbatim, on the line under the title.
3. **Table**: `Doc | Description`, one row per listed doc, in index order.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/Index/sample.md).
