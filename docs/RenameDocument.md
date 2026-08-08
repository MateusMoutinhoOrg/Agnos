# Rename a Document

## Description
Covers renaming or moving an existing `.md` file in [docs/](../) without leaving broken references behind.

### Rules
- Renaming a `.md` file requires updating the [README.md](/README.md) Doc Index and [Structure.md](/docs/Structure.md) in the same commit.
- Moving a document to another category must respect the category boundaries described in [AddDocument.md](/docs/AddDocument.md).

---

## Workflow
1. Rename or move the `.md` file, keeping a descriptive, topic-based name.
2. Find every reference to the old path:
   ```bash
   grep -rn "OldName.md" --include="*.md" .
   ```
3. Update each reference to the new **relative path**, following the cross-reference rules of the GeneralDoc specification — locate it in [Specs.md](/docs/Specs.md).
4. Update the document's entry in the [README.md](/README.md) Doc Index — link text, link target, and `Type` column, in **every** theme table that lists it.
5. Update the file's entry in [Structure.md](/docs/Structure.md).
