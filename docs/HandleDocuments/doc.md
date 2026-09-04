# Handle Documents

## Description
Covers creating, renaming, moving, and deleting a doc in [docs/](/docs/), and registering it across the project — the companion step nearly every other tutorial ends in. Every index is generated, so registering a doc means writing its `props.yaml` and running a build.

### Rules
- A **doc** is a directory: `docs/<DocName>/` holding `doc.md` (what you write) and `props.yaml` (its declaration). A doc directory may hold sub-doc directories of the same shape, to any depth; any other file in it — a sample, an image — is an asset and is ignored by the indexes.
- `Index` is reserved at the first level: `docs/Index/` holds the generated theme indexes and is never a doc.
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs](/docs/Specs/doc.md). `props.yaml` follows [DocProps](/docs/Specs/DocProps/doc.md).
- A doc belongs to one or more **themes** of `<ProjectName>Config/themes.yaml`, named by id in its `props.yaml`. A first-level doc must name at least one; a sub-doc must name none — it is listed by its parent's `Index.md`.
- **Indexes are generated, never hand-written.** `agnos build` writes `docs/Index/<theme-id>.md` for every theme and `docs/<doc>/Index.md` for every doc that has sub-docs, from the `props.yaml` files. `docs/Index/` is rewritten whole on every build.
- `verify` fails on a missing or unparsable `props.yaml`, an unknown theme id, a first-level doc with no themes, a sub-doc with themes, and a theme no doc names.
- The [README.md](/README.md) links to theme indexes only: it changes when a **theme** is added, renamed, or removed, never for a single doc.
- A doc that gains sub-docs links to its own `Index.md` once from its `doc.md` — the generator writes the index, not the link to it.
- Content still needed elsewhere must be moved before deletion, not lost.
- `CLAUDE.md` is the harness specification for the coding agent, not documentation: a pattern established or changed in `docs/` is mirrored there in the same commit, and the reverse.

---

## Add a Document
1. Decide whether the doc is a **workflow** (numbered steps toward one goal), a **lookup**, or an **explanation** — the kind picks its specification, not its location. Every doc is a directory of `docs/`.
2. Decide where it belongs: a first-level `docs/<DocName>/` for a doc a theme index lists, or `docs/<Parent>/<DocName>/` for one that only makes sense inside another doc (a symbol page, a specification).
3. Scaffold it with `add-doc`, which writes both files and runs the build:
   ```bash
   # a first-level doc: one --theme per theme id of themes.yaml, repeatable
   agnos add-doc HandleReports --theme development \
     --description "How a report is written and regenerated"

   # a sub-doc: nested name, and no theme — its parent's Index.md lists it
   agnos add-doc HandleReports/Layout \
     --description "The sections every generated report carries"
   ```
   The `name` key is derived from the last segment (`HandleReports` → `Handle Reports`); a segment with its own punctuation (`api.AddDoc`) is kept verbatim. Edit `props.yaml` afterwards to change it or to add an `order`.
4. Check [Specs](/docs/Specs/doc.md) for the specifications matching the files you are about to write, and read them first.
5. Replace the stub `doc.md`, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
6. Add cross-references using **repository-rooted paths** to the other doc's `doc.md`, and add the reverse link in every doc that should point back to this one.
7. Run a build so the indexes catch up:
   ```bash
   agnos build
   ```
8. Register the doc in [Structure](/docs/Structure/doc.md) if it is a new structural component.

---

## Rename or Move a Document
1. Rename or move the doc **directory** — `doc.md`, `props.yaml` and every asset travel with it.
2. Update `name` (and `order`, if the position changes) in its `props.yaml`, and the H1 of its `doc.md` to match.
3. Find every reference to the old path:
   ```bash
   grep -rn "OldName/doc.md" --include="*.md" . --exclude-dir=old
   ```
4. Update each reference to the new repository-rooted path, following the cross-reference rules of the GeneralDoc specification.
5. When the doc changes theme, edit the `themes` list in its `props.yaml`; moving it under another doc means deleting that list instead.
6. Run `agnos build` and update the doc's entry in [Structure](/docs/Structure/doc.md) if it is explicitly listed.

---

## Delete a Document
1. Find every reference to the doc:
   ```bash
   grep -rn "DocName/doc.md" --include="*.md" . --exclude-dir=old
   ```
2. For each reference, remove it or repoint it to the doc that now covers the topic.
3. Delete the doc directory, sub-docs and assets included — `remove-doc` does it and rebuilds:
   ```bash
   agnos remove-doc HandleReports
   ```
4. The theme indexes are rewritten whole by that build, and the parent's `Index.md` loses the row.
5. Remove the doc's entry from [Structure](/docs/Structure/doc.md) if it is explicitly listed, and drop its theme from `themes.yaml` if it was the last doc of that theme — `verify` rejects a theme with no docs.
