# ReferenceDocs Specification

## Description
Defines the required shape of a **Reference** page — any page under a `docs/<Theme>/References/` directory — including the API detail pages under `docs/LibUsage/References/PublicApi/` — that lists enumerable items (structures, specs, commands, API entries) and is not one of the special documents ([RULES](/docs/Development/References/Specs/Rules/Specs.md), [Structure](/docs/Development/References/Specs/Structure/Specs.md), a theme [Index](/docs/Development/References/Specs/Index/Specs.md), the `Specs.md` index, or anything under `docs/Development/References/Specs/`). A reference page is meant to be **scanned**, not read linearly.

### Rules
- Every page must comply with [GeneralDoc](/docs/Development/References/Specs/GeneralDoc/Specs.md).
- The body must be **listable content**: Markdown tables or linked entry lists, one `##` section per group of items.
- Each entry must have a short description; when a detail page exists, the entry name must link to it.
- Reference pages must not contain workflows — link to the relevant tutorial instead.
- Every new page must be registered in the `Index.md` of its theme.

## Structure
1. **Title** (H1): the name of what is being listed.
2. **`## Description`**: one short paragraph on what the page lists.
3. **One `##` section per item group**, separated by `---`, each containing a table or a list of entries with a name (linked when a detail page exists) and a short description.

> **Note**: For a concrete example, refer to [sample.md](/docs/Development/References/Specs/ReferenceDocs/sample.md).
