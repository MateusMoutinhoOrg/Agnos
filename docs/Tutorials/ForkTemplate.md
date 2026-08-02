# Fork This Repository as a Template

## Description
Covers using this repository as a GitHub template to start a **new** dependency-injected library. To convert a library that already exists, follow [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) instead.

### Rules
- Read [RULES.md](/docs/References/RULES.md) and [Structure.md](/docs/References/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/References/Structure.md): public contract structs in `sandbox/contracts/`, internal types and logic in `sandbox/internal/`, concrete dependencies in `adapters/`, and the entry point in `sandbox/`. Contracts are structs of function fields, never interfaces — see [StructContracts.md](/docs/Explanations/StructContracts.md).
- Every file of the template has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/References/TemplateFileActions.md); the steps below follow that order.
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/References/Specs.md).
- The fork is not complete until the final checklist in the last workflow step passes.

---

## Workflow
1. On the GitHub repository page, click **"Use this template"** and create the new repository.
2. Rename the module to the new GitHub path, following [RenameModule.md](/docs/Tutorials/RenameModule.md).
3. Leave every **[Copy](/docs/References/TemplateFileActions.md#copy)** file untouched — they describe the structure, not the library.
4. Rewrite [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go) with the dependencies the new library requires, following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. Rewrite [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go) with the `Lib` struct and one struct per object the new library hands back, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) and [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
6. Rewrite [adapters/standard/standard.go](../../adapters/standard/standard.go) so the default adapter fills every field of the new contract, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
7. Create the new library logic in [sandbox/internal/](../../sandbox/internal/) — the lib's factories plus one package per object — following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) and [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
8. Create any additional adapter in [adapters/](../../adapters/), following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
9. Create the new samples in [examples/](../../examples/), following [AddSample.md](/docs/Tutorials/AddSample.md).
10. Create the new detail pages in [docs/References/PublicApi/](../References/PublicApi/) and rewrite [PublicApi.md](/docs/References/PublicApi.md), following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
11. Delete every remaining **[Delete](/docs/References/TemplateFileActions.md#delete)** file — the example internal logic, adapters, samples, and API pages the new library replaced. For `.md` files, follow [DeleteDocument.md](/docs/Tutorials/DeleteDocument.md).
12. Rewrite [docs/References/Structure.md](/docs/References/Structure.md) to describe the resulting layout.
13. Create the tutorials specific to the new library — one page per workflow its maintainers will repeat — following [AddDocument.md](/docs/Tutorials/AddDocument.md) and the [TutorialDocs specification](/docs/References/Meta/TutorialDocs/Specs.md). The template tutorials cover the structure only; they do not document the library's own use cases.
14. Create any reference page the library needs beyond the public API, following [AddDocument.md](/docs/Tutorials/AddDocument.md) and the [ReferenceDocs specification](/docs/References/Meta/ReferenceDocs/Specs.md).
15. Rewrite the [README.md](/README.md): overview, quick start, badges, Doc Index, and Samples section.
16. Verify the result:
```bash
go build ./...
```
Then confirm every item below — the fork is only done when all pass:
- All library logic lives in `sandbox/internal/`; no file there imports `os`, `net`, or a third-party implementation directly — every such call goes through `l.Deps`.
- `sandbox/contracts/deps/deps.go` declares one function field per injected call, and **every** adapter in `adapters/` fills all of them — the compiler does not check this.
- `sandbox/contracts/api/api.go` declares every public object as a struct with a `Deps` field, and every one of its function fields is filled by a factory registered in that package's `New` constructor.
- `sandbox/new.go` is the only wiring point, and it imports no adapter.
- Tutorials and reference pages specific to this library exist under `docs/Tutorials/` and `docs/References/`.
- Every created or rewritten file matches its specification from [Specs.md](/docs/References/Specs.md).
- The `README.md` Doc Index lists every `.md` file and the Samples section lists every sample.
