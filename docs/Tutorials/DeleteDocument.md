# Delete a Document

## Description
Covers removing a `.md` file from [docs/](../) and clearing every reference to it.

### Rules
- Deleting a `.md` file requires updating the [README.md](/README.md) Doc Index and [Structure.md](/docs/References/Structure.md) in the same commit.
- Content still needed elsewhere must be moved before deletion, not lost — if the goal is to keep the content under another name, follow [RenameDocument.md](/docs/Tutorials/RenameDocument.md) instead.

---

## Workflow
1. Find every reference to the document:
   ```bash
   grep -rn "DocName.md" --include="*.md" .
   ```
2. For each reference, remove it or repoint it to the document that now covers the topic.
3. Delete the `.md` file.
4. Remove the document's entry from the [README.md](/README.md) Doc Index — from **every** theme table that lists it.
5. Remove the file's entry from [Structure.md](/docs/References/Structure.md).
