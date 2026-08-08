# Fork This Repository as a Template

## Description
Covers using this repository as a GitHub template to start a **new** dependency-injected library. To convert a library that already exists, follow [AdaptExistingLib.md](/docs/AdaptExistingLib.md) instead.

### Rules
- Read [RULES.md](/docs/RULES.md) and [Structure.md](/docs/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/Structure.md): public contract structs in `sandbox/contracts/`, internal types and logic in `sandbox/internal/`, concrete dependencies in `adapters/`, the entry point in `sandbox/`, and the installed binary in `cmd/main/`. The command-line interface belongs to the library, as the `Sandboxmain` field of `api.Lib`, never to the binary. Contracts are structs of function fields, never interfaces — see [StructContracts.md](/docs/StructContracts.md).
- Every file of the template has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/TemplateFileActions.md); the steps below follow that order.
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/Specs.md).
- The fork is not complete until the final checklist in the last workflow step passes.

---

## Workflow
1. On the GitHub repository page, click **"Use this template"** and create the new repository.
2. Rename the module to the new GitHub path, following [RenameModule.md](/docs/RenameModule.md).
3. Leave every **[Copy](/docs/TemplateFileActions.md#copy)** file untouched — they describe the structure, not the library.
4. Rewrite [sandbox/contracts/deps/deps.go](../sandbox/contracts/deps/deps.go) with the dependencies the new library requires, following [HandleDependencies.md](/docs/HandleDependencies.md).
5. Rewrite [sandbox/contracts/api/api.go](../sandbox/contracts/api/api.go) with the `Lib` struct and one struct per object the new library hands back, following [HandleLibElements.md](/docs/HandleLibElements.md) and [HandleLibElements.md](/docs/HandleLibElements.md).
6. Rewrite [adapters/standard/standard.go](../adapters/standard/standard.go) so the default adapter fills every field of the new contract, following [HandleAdapters.md](/docs/HandleAdapters.md).
7. Create the new library logic in [sandbox/internal/](../sandbox/internal/) — the lib's factories plus one package per object — following [HandleLibElements.md](/docs/HandleLibElements.md) and [HandleLibElements.md](/docs/HandleLibElements.md).
8. Create any additional adapter in [adapters/](../adapters/), following [HandleAdapters.md](/docs/HandleAdapters.md).
9. Create the new samples: the Go programs in [examples/libraryExamples/](../examples/libraryExamples/), following [HandleSamples.md](/docs/HandleSamples.md), and the shell scripts in [examples/cliExamples/](../examples/cliExamples/), following [HandleCliExamples.md](/docs/HandleCliExamples.md).
10. Create the new detail pages in [docs/References/PublicApi/](../References/PublicApi/) and rewrite [PublicApi.md](/docs/PublicApi.md), following [ExposePublicApi.md](/docs/ExposePublicApi.md).
11. Delete every remaining **[Delete](/docs/TemplateFileActions.md#delete)** file — the example internal logic, adapters, samples, and API pages the new library replaced. For `.md` files, follow [HandleDocuments.md](/docs/HandleDocuments.md).
12. Rewrite [docs/References/Structure.md](/docs/Structure.md) to describe the resulting layout.
13. Create the tutorials specific to the new library — one page per workflow its maintainers will repeat — following [HandleDocuments.md](/docs/HandleDocuments.md) and the [TutorialDocs specification](/docs/Meta/TutorialDocs/Specs.md). The template tutorials cover the structure only; they do not document the library's own use cases.
14. Create any reference page the library needs beyond the public API, following [HandleDocuments.md](/docs/HandleDocuments.md) and the [ReferenceDocs specification](/docs/Meta/ReferenceDocs/Specs.md).
15. Rewrite the [README.md](/README.md): overview, both quick starts, badges, Doc Index, and the two Examples sections.
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
- Every created or rewritten file matches its specification from [Specs.md](/docs/Specs.md).
- The `README.md` Doc Index lists every `.md` file, and the CLI Examples and Library Examples sections list every sample.
- `cmd/main/main.go` wires, calls `Sandboxmain`, and exits — it branches on no command, parses no flag, and prints nothing.
