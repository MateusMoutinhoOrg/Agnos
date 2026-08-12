# Index Specification

## Description
Defines the required shape of a **theme index** — one page per theme under `docs/Index/` (`CliUsage.md`, `LibUsage.md`, `Development.md`, `Templating.md`). A theme is a reader goal, not a directory: pages live flat in `docs/Tutorials/` and `docs/References/`, and the index is what groups them. A theme index is the single entry point of its theme: the [README.md](/README.md) links only to indexes, and an index links to every page of its theme.

### Rules
- Every page must comply with [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).
- One index per theme, named `<Theme>.md`, living directly inside `docs/Index/`.
- The index lists **every** `.md` file of its theme: one row per page of the theme in `docs/Tutorials/`, one row per page of the theme in `docs/References/`. No orphans. Pages already indexed by another page — the API details under `PublicApi/`, the specifications under `Specs/` — are covered through that page's row.
- A page is listed by the index of **its own** theme, and by exactly one index. A page belonging to another theme is linked in prose, never given a row.
- Each row is `[FileName.md](/docs/<Tutorials|References>/<FileName>.md)` plus a one-line description, 50–100 characters, saying what the reader gets.
- Link text and link target must match the real file location.
- Rows are ordered by reader need: what a newcomer opens first comes first.
- Creating, renaming, moving, or deleting a `.md` file requires updating its theme index in the same commit — see [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md).

## Structure
1. **Title** (H1): the theme name.
2. **`## Description`**: one short paragraph on what the theme covers and who it is for, linking to the neighboring theme indexes.
3. **`---`**: horizontal rule separating the header from the tables.
4. **`## Tutorials`**: a table of every workflow page of the theme, with `Doc` and `Description` columns.
5. **`---`**: horizontal rule.
6. **`## References`**: a table of every explanation and lookup page of the theme, with the same columns.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/Index/sample.md).
