# ReferenceDocs Specification

## Description
Defines the required shape of a **Reference** page — any `.md` file in `docs/` — including the API detail pages under `docs/PublicApi/` — that lists enumerable items (structures, specs, commands, API entries) and is not one of the special documents ([RULES](/docs/Meta/Rules/Specs.md), [Structure](/docs/Meta/Structure/Specs.md), the `Specs.md` index, or anything under `docs/Meta/`). A reference page is meant to be **scanned**, not read linearly.

### Rules
- Every page must comply with [GeneralDoc](/docs/Meta/GeneralDoc/Specs.md).
- The body must be **listable content**: Markdown tables or linked entry lists, one `##` section per group of items.
- Each entry must have a short description; when a detail page exists, the entry name must link to it.
- Reference pages must not contain workflows — link to the relevant tutorial instead.
- Every new page must be registered in the [README.md](/README.md) Doc Index.

## Structure
1. **Title** (H1): the name of what is being listed.
2. **`## Description`**: one short paragraph on what the page lists.
3. **One `##` section per item group**, separated by `---`, each containing a table or a list of entries with a name (linked when a detail page exists) and a short description.

> **Note**: For a concrete example, refer to [sample.md](/docs/Meta/ReferenceDocs/sample.md).
