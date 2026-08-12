# Handle Documents

## Description
Covers creating, renaming, moving, and deleting `.md` files in [docs/](/docs/), and registering them across the project — the companion step nearly every other tutorial ends in.

### Rules
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs.md](/docs/References/Specs.md).
- A file governed by a specification must reproduce the shape that specification requires.
- Every page lives in `docs/Tutorials/` if it is a workflow, in `docs/References/` if it is a lookup or an explanation. There are no theme directories: the two directories are flat, and a page's file name must be unique across the one it lands in.
- Every page belongs to exactly one **theme** — `CliUsage`, `LibUsage`, `Development`, `Templating` — and the theme is expressed by the [theme index](/docs/References/Specs/Index/Specs.md) that lists it, in `docs/Index/<Theme>.md`, never by the page's location.
- Adding, renaming, or deleting a `.md` file requires updating its theme index in `docs/Index/` and [Structure.md](/docs/References/Structure.md) in the same commit.
- The [README.md](/README.md) links to theme indexes only: it changes when a **theme** is added, renamed, or removed, never for a single page.
- Moving a document to another theme is a row moving from one index to another; the file itself only moves if it changes between a workflow and a reference.
- Content still needed elsewhere must be moved before deletion, not lost.

---

## Add a Document
1. Identify the theme the document belongs to, and whether it is a **Tutorial** (a workflow) or a **Reference** (a lookup or an explanation).
2. Check [Specs.md](/docs/References/Specs.md) for the specifications matching the file, and read them before writing.
3. Create the `.md` file in `docs/Tutorials/` or `docs/References/` with a descriptive, topic-based name unique in that directory (e.g., `HandleSamples.md`, `PublicApi.md`). When two themes need the same topic, prefix the name with the theme's subject — `CliQuickStart.md` and `LibQuickStart.md`.
4. Write the content following those specifications, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
5. Add cross-references using **relative paths**, and add the reverse link in every document that should point back to this one.
6. Add a row to the theme's index in `docs/Index/`, in the `Tutorials` or the `References` table.
7. Register the file in [Structure.md](/docs/References/Structure.md).

---

## Rename or Move a Document
1. Rename or move the `.md` file, keeping a descriptive, topic-based name and landing in `docs/Tutorials/` or `docs/References/`.
2. Find every reference to the old path:
   ```bash
   grep -rn "OldName.md" --include="*.md" .
   ```
3. Update each reference to the new **relative path**, following the cross-reference rules of the GeneralDoc specification.
4. Update the document's row in the theme's index in `docs/Index/` — link text and link target. When the document changed theme, remove the row from the old index and add it to the new one.
5. Update the file's entry in [Structure.md](/docs/References/Structure.md).

---

## Delete a Document
1. Find every reference to the document:
   ```bash
   grep -rn "DocName.md" --include="*.md" .
   ```
2. For each reference, remove it or repoint it to the document that now covers the topic.
3. Delete the `.md` file.
4. Remove the document's row from its theme's index in `docs/Index/`.
5. Remove the file's entry from [Structure.md](/docs/References/Structure.md).
