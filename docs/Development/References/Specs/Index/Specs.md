# Index Specification

## Description
Defines the required shape of a **theme index** — the `Index.md` at the root of each theme directory under `docs/` (`CliUsage/`, `LibUsage/`, `Development/`, `Templating/`). A theme index is the single entry point of its theme: the [README.md](/README.md) links only to indexes, and an index links to every page of its theme.

### Rules
- Every page must comply with [GeneralDoc](/docs/Development/References/Specs/GeneralDoc/Specs.md).
- One index per theme, named `Index.md`, living directly inside the theme directory.
- The index lists **every** `.md` file of its theme: one row per page in `Tutorials/`, one row per page in `References/`. No orphans. Pages already indexed by another page — the API details under `PublicApi/`, the specifications under `Specs/` — are covered through that page's row.
- A page is listed by the index of **its own** theme. A page from another theme is linked in prose, never given a row.
- Each row is `[FileName.md](/docs/<Theme>/<Tutorials|References>/<FileName>.md)` plus a one-line description, 50–100 characters, saying what the reader gets.
- Link text and link target must match the real file location.
- Rows are ordered by reader need: what a newcomer opens first comes first.
- Creating, renaming, moving, or deleting a `.md` file requires updating its theme index in the same commit — see [HandleDocuments.md](/docs/Development/Tutorials/HandleDocuments.md).

## Structure
1. **Title** (H1): the theme name.
2. **`## Description`**: one short paragraph on what the theme covers and who it is for, linking to the neighboring theme indexes.
3. **`---`**: horizontal rule separating the header from the tables.
4. **`## Tutorials`**: a table of every workflow page of the theme, with `Doc` and `Description` columns.
5. **`---`**: horizontal rule.
6. **`## References`**: a table of every explanation and lookup page of the theme, with the same columns.

> **Note**: For a concrete example, refer to [sample.md](/docs/Development/References/Specs/Index/sample.md).
