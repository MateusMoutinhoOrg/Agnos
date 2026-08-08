# Handle Documents

## Description
Covers creating, renaming, moving, and deleting `.md` files in [docs/](../), and registering them across the project.

### Rules
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs.md](/docs/Specs.md).
- A file governed by a specification must reproduce the shape that specification requires.
- Adding, renaming, or deleting a `.md` file requires updating the [README.md](/README.md) Doc Index and [Structure.md](/docs/Structure.md) in the same commit.
- Moving a document to another category must respect the category boundaries.
- Content still needed elsewhere must be moved before deletion, not lost.

---

## Add a Document
1. Identify the category the document belongs to.
2. Check [Specs.md](/docs/Specs.md) for the specifications matching the file, and read them before writing.
3. Create the `.md` file in the chosen directory with a descriptive, topic-based name (e.g., `HandleSamples.md`, `PublicApi.md`).
4. Write the content following those specifications, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
5. Add cross-references using **relative paths**, and add the reverse link in every document that should point back to this one.
6. Add an entry to the [README.md](/README.md) Doc Index, in every **theme** table where the document fits.
7. Register the file in [Structure.md](/docs/Structure.md).

---

## Rename or Move a Document
1. Rename or move the `.md` file, keeping a descriptive, topic-based name.
2. Find every reference to the old path:
   ```bash
   grep -rn "OldName.md" --include="*.md" .
   ```
3. Update each reference to the new **relative path**, following the cross-reference rules of the GeneralDoc specification.
4. Update the document's entry in the [README.md](/README.md) Doc Index — link text and link target, in **every** theme table that lists it.
5. Update the file's entry in [Structure.md](/docs/Structure.md).

---

## Delete a Document
1. Find every reference to the document:
   ```bash
   grep -rn "DocName.md" --include="*.md" .
   ```
2. For each reference, remove it or repoint it to the document that now covers the topic.
3. Delete the `.md` file.
4. Remove the document's entry from the [README.md](/README.md) Doc Index — from **every** theme table that lists it.
5. Remove the file's entry from [Structure.md](/docs/Structure.md).
