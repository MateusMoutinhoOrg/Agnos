# Expose in the Public API

## Description
Covers publishing a library function, object, or field in the public API index at [PublicApi.md](/docs/PublicApi.md).

### Rules
- Every public-facing entry must be listed in [PublicApi.md](/docs/PublicApi.md).
- Detail pages live in [docs/References/PublicApi/](../References/PublicApi/) and are named `<pkg>.<Symbol>.md`.
- Adding a detail page requires updating [Structure.md](/docs/Structure.md) and the [README.md](/README.md) Doc Index.

---

## Workflow
1. Open [PublicApi.md](/docs/PublicApi.md).
2. Add the struct, function, or field to the section matching its kind, with a one-line description. An object is public only through its `sandbox/contracts/api` struct — never document the `sandbox/internal/` type as the entry.
3. Create the detail page under [docs/References/PublicApi/](../References/PublicApi/), named `<pkg>.<Symbol>.md` after the package the symbol is declared in (e.g., `api.Get.md`), following [HandleDocuments.md](/docs/HandleDocuments.md).
4. Link the new detail page from its entry in [PublicApi.md](/docs/PublicApi.md).
5. Register the detail page in [Structure.md](/docs/Structure.md).
