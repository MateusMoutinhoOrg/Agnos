# DocProps Specification

## Description
Defines the required shape of `docs/<doc-name>/props.yaml` — the declaration that makes a directory a doc. It is the only place a doc's name, its one-line description and its themes are written: every index is generated from these files, so a doc that is not declared here is listed nowhere. Its parsed shape is `docpropsconf`, and the rules below are the ones [verify](/docs/RegenerateProject/doc.md) enforces.

### Rules
- Every doc directory — first level or nested — holds exactly one `props.yaml`, next to its `doc.md`. A missing or unparsable file fails `verify`.
- `name` is the doc as an index lists it: the H1 of its `doc.md`, minus any suffix every sibling repeats (`GeneralDoc Specification` is listed as `GeneralDoc`).
- `description` is one line, 50–100 characters, saying what the reader gets. It is the row description in every index that lists the doc, so it must read on its own.
- `themes` is a list of theme **ids** from `<ProjectName>Config/themes.yaml` (`cli-usage`, not `CliUsage`). It is **required** on a first-level doc — a doc reachable from no theme index fails `verify` — and **forbidden** on a sub-doc, which is listed only by its parent's `Index.md`.
- A theme id that is not declared in `themes.yaml` fails `verify`, and so does a theme that no doc names.
- `order` is optional: the position the doc takes in the indexes that list it. Docs with no `order` are listed after every ordered one, alphabetically by name. A doc in several themes carries one `order`, so pick a value that reads correctly in each of them.
- Adding, renaming, moving or deleting a doc is a `props.yaml` change plus a build — never a hand-edited index; see [HandleDocuments](/docs/HandleDocuments/doc.md).

## Structure
1. **`name`**: the doc's title.
2. **`description`**: the one-line summary.
3. **`themes`** *(first-level docs only)*: one theme id per list entry.
4. **`order`** *(optional)*: the integer position in its indexes.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/DocProps/sample.yaml).
