# Index Specification

## Description
Defines the required shape of a **theme index** — one page per theme under `docs/Index/` (`CliUsage.md`, `LibUsage.md`, `Development.md`, `Templating.md`). A theme is a reader goal, not a directory: pages live flat in `docs/Tutorials/` and `docs/References/`, and the index is what groups them. A theme index is the single entry point of its theme: the [README.md](/README.md) links only to indexes, and an index links to every page of its theme.

An index entry is a **nested list**, not a table row: the page, its one-line description, and a link to each topic the page covers. The topic links are the point — they let a reader see what is inside a page without opening it.

### Rules
- Every page must comply with [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).
- One index per theme, named `<Theme>.md`, living directly inside `docs/Index/`.
- The index lists **every** `.md` file of its theme: one entry per page of the theme in `docs/Tutorials/`, one entry per page of the theme in `docs/References/`. No orphans. Pages already indexed by another page — the API details under `PublicApi/`, the specifications under `Specs/` — are covered through that page's entry.
- A page is listed by the index of **its own** themes, and can be listed in multiple indexes if it spans multiple themes.
- Each entry is a top-level list item shaped `[FileName.md](/docs/Tutorials/FileName.md)` — or `/docs/References/` for a reference page — holding one nested `- **description:** …` line of 50–100 characters saying what the reader gets, followed by one nested link per topic section of that page.
- **Topic sections only.** The skeleton headings every page shares are never listed: `Description`, `Rules`, `Workflow`, `Full Code`. An entry lists the headings that name a subject *inside* the page, and nothing else.
- Topics are the page's `##` headings. When a page's only body section is `## Workflow` and that workflow is split into named `###` subsections, those subsections are the topics instead.
- A page with no topic section — a single linear workflow — carries its `**description:**` line alone.
- A section link's text is the heading as written; its target is the page path plus the heading's anchor (lowercased, spaces to hyphens, other punctuation dropped) — `/docs/References/Commands.md#exit-codes`.
- Link text, link target, and section anchors must match the real file and its real headings.
- Entries are ordered by reader need: what a newcomer opens first comes first. Section links keep the order they have in the page.
- Creating, renaming, moving, or deleting a `.md` file requires updating its theme index in the same commit — see [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md). Renaming a section of a listed page updates that page's entry in the same commit.

## Structure
1. **Title** (H1): the theme name.
2. **`## Description`**: one short paragraph on what the theme covers and who it is for, linking to the neighboring theme indexes.
3. **`---`**: horizontal rule separating the header from the listings.
4. **`## Tutorials`**: one entry per workflow page of the theme.
5. **`---`**: horizontal rule.
6. **`## References`**: one entry per explanation and lookup page of the theme, in the same shape.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/Index/sample.md).
