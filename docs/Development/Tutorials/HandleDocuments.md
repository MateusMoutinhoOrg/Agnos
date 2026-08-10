# Handle Documents

## Description
Covers creating, renaming, moving, and deleting `.md` files in [docs/](/docs/), and registering them across the project — the companion step nearly every other tutorial ends in.

### Rules
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs.md](/docs/Development/References/Specs.md).
- A file governed by a specification must reproduce the shape that specification requires.
- Every page lives in a **theme** directory — `CliUsage/`, `LibUsage/`, `Development/`, `Templating/` — under `Tutorials/` if it is a workflow, under `References/` if it is a lookup or an explanation.
- Adding, renaming, or deleting a `.md` file requires updating its theme [Index.md](/docs/Development/References/Specs/Index/Specs.md) and [Structure.md](/docs/Development/References/Structure.md) in the same commit.
- The [README.md](/README.md) links to theme indexes only: it changes when a **theme** is added, renamed, or removed, never for a single page.
- Moving a document to another theme must respect the theme boundaries.
- Content still needed elsewhere must be moved before deletion, not lost.

---

## Add a Document
1. Identify the theme the document belongs to, and whether it is a **Tutorial** (a workflow) or a **Reference** (a lookup or an explanation).
2. Check [Specs.md](/docs/Development/References/Specs.md) for the specifications matching the file, and read them before writing.
3. Create the `.md` file in `docs/<Theme>/<Tutorials|References>/` with a descriptive, topic-based name (e.g., `HandleSamples.md`, `PublicApi.md`).
4. Write the content following those specifications, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
5. Add cross-references using **relative paths**, and add the reverse link in every document that should point back to this one.
6. Add a row to the theme's `Index.md`, in the `Tutorials` or the `References` table.
7. Register the file in [Structure.md](/docs/Development/References/Structure.md).

---

## Rename or Move a Document
1. Rename or move the `.md` file, keeping a descriptive, topic-based name and landing in a `Tutorials/` or `References/` directory of a theme.
2. Find every reference to the old path:
   ```bash
   grep -rn "OldName.md" --include="*.md" .
   ```
3. Update each reference to the new **relative path**, following the cross-reference rules of the GeneralDoc specification.
4. Update the document's row in the theme's `Index.md` — link text and link target. When the document changed theme, remove the row from the old index and add it to the new one.
5. Update the file's entry in [Structure.md](/docs/Development/References/Structure.md).

---

## Delete a Document
1. Find every reference to the document:
   ```bash
   grep -rn "DocName.md" --include="*.md" .
   ```
2. For each reference, remove it or repoint it to the document that now covers the topic.
3. Delete the `.md` file.
4. Remove the document's row from its theme's `Index.md`.
5. Remove the file's entry from [Structure.md](/docs/Development/References/Structure.md).
