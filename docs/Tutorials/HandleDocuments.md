# Handle Documents

## Description
Covers creating, renaming, moving, and deleting `.md` files in [docs/](/docs/), and registering them across the project — the companion step nearly every other tutorial ends in.

### Rules
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs.md](/docs/References/Specs.md).
- A file governed by a specification must reproduce the shape that specification requires.
- Workflow pages live in `docs/Tutorials/`. Lookups and explanations live in `docs/References/`. While `Tutorials/` is strictly flat, `References/` contains two subdirectories for specific listings: `PublicApi/` for API symbols and `Specs/` for specifications. A page's file name must be unique inside the directory it lands in.
- Every page belongs to exactly one **theme** — `CliUsage`, `LibUsage`, `Development`, `Templating` — and the theme is expressed by the index that lists it, never by the page's location.
- Adding, renaming, or deleting a `.md` file requires updating its index in the same commit. For flat pages, update the theme index in `docs/Index/` and [Structure.md](/docs/References/Structure.md). For subdirectory pages, update their specific index instead ([PublicApi.md](/docs/References/PublicApi.md) or [Specs.md](/docs/References/Specs.md)).
- An index entry lists the page's own topic sections, so adding, renaming, or removing a section of an indexed page updates that entry in the same commit.
- The [README.md](/README.md) links to theme indexes only: it changes when a **theme** is added, renamed, or removed, never for a single page.
- Moving a document to another theme is an entry moving from one index to another; the file itself only moves if it changes between a workflow and a reference.
- Content still needed elsewhere must be moved before deletion, not lost.

---

## Add a Document
1. Identify the theme the document belongs to, and whether it is a **Tutorial** (a workflow) or a **Reference** (a lookup or an explanation).
2. Check [Specs.md](/docs/References/Specs.md) for the specifications matching the file, and read them before writing.
3. Create the `.md` file in `docs/Tutorials/`, `docs/References/`, or its specific subdirectory (`docs/References/PublicApi/` or `docs/References/Specs/`), with a unique, topic-based name. When two themes need the same topic in a flat directory, prefix the name with the theme's subject.
4. Write the content following those specifications, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
5. Add cross-references using **relative paths**, and add the reverse link in every document that should point back to this one.
6. Add an entry to the document's index:
   - For flat pages: add it to its theme index in `docs/Index/`, under `Tutorials` or `References` with a nested `**description:**` line and one nested link per topic section — see the [Index](/docs/References/Specs/Index/Specs.md) specification.
   - For subdirectory pages: add it to `docs/References/PublicApi.md` or `docs/References/Specs.md` respectively.
7. Register the file in [Structure.md](/docs/References/Structure.md) if it is a new structural component.

---

## Rename or Move a Document
1. Rename or move the `.md` file, keeping a descriptive, topic-based name and landing in its appropriate directory (`Tutorials/`, `References/`, `PublicApi/`, or `Specs/`).
2. Find every reference to the old path:
   ```bash
   grep -rn "OldName.md" --include="*.md" .
   ```
3. Update each reference to the new **relative path**, following the cross-reference rules of the GeneralDoc specification.
4. Update the document's entry in its index (theme index in `docs/Index/`, or `PublicApi.md`/`Specs.md`) — link text, link target, and every section link under it. When the document changes theme, remove the entry from the old index and add it to the new one.
5. Update the file's entry in [Structure.md](/docs/References/Structure.md) if it is explicitly listed.

---

## Delete a Document
1. Find every reference to the document:
   ```bash
   grep -rn "DocName.md" --include="*.md" .
   ```
2. For each reference, remove it or repoint it to the document that now covers the topic.
3. Delete the `.md` file.
4. Remove the document's entry from its index (theme index in `docs/Index/`, or `PublicApi.md`/`Specs.md`).
5. Remove the file's entry from [Structure.md](/docs/References/Structure.md) if it is explicitly listed.
