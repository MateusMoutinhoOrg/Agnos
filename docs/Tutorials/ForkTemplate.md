# Fork This Repository as a Template

## Description
Covers using this repository as a GitHub template to start a **new** dependency-injected library. To convert a library that already exists, follow [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) instead. The steps are grouped into phases so progress is easy to track; every phase takes each file's action — **Copy**, **Create**, **Rewrite**, or **Delete** — from [TemplateFileActions.md](/docs/References/TemplateFileActions.md).

### Rules
- Read [Structure.md](/docs/References/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/References/Structure.md): public contract structs in `sandbox/contracts/`, internal factories in `sandbox/`, concrete dependencies in `adapters/`, the entry point in `sandbox/`, and the installed binary in `cmd/main/`. The command-line interface belongs to the library, as the `Sandboxmain` field of `api.Lib`, never to the binary. Contracts are structs of function fields, never interfaces — see [StructContracts.md](/docs/References/StructContracts.md).
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/References/Specs.md).
- The fork is not complete until the final checklist in the last workflow step passes.

---

## Workflow

### Phase 1 — Create the repository
1. On the GitHub repository page, click **"Use this template"** and create the new repository.
2. Rename the module to the new GitHub path, following [RenameModule.md](/docs/Tutorials/RenameModule.md).
3. Leave every **[Copy](/docs/References/TemplateFileActions.md#copy)** file untouched — they describe the structure, not the library.

### Phase 2 — Rewrite the contracts
4. Rewrite [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go) with the dependencies the new library requires, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).
5. Rewrite [sandbox/contracts/api/api.go](/sandbox/contracts/api/api.go) with the `Lib` struct and one struct per object the new library hands back, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).
6. Rewrite [adapters/standard/standard.go](/adapters/standard/standard.go) so the default adapter fills every field of the new contract, following [HandleAdapters.md](/docs/Tutorials/HandleAdapters.md).

### Phase 3 — Create the implementation
7. Create the new library logic in [sandbox/](/sandbox/) — the lib's factories plus one package per object — following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).
8. Create the command dispatch behind `Sandboxmain` in `sandbox/cli/`, following [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md).
9. Create any additional adapter in [adapters/](/adapters/), following [HandleAdapters.md](/docs/Tutorials/HandleAdapters.md).
10. Create the new samples: the Go programs in [examples/libraryExamples/](/examples/libraryExamples/), following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md), and the shell scripts in [examples/cliExamples/](/examples/cliExamples/), following [HandleCliExamples.md](/docs/Tutorials/HandleCliExamples.md).

### Phase 4 — Rewrite the documentation
11. Create the new API detail pages (`docs/References/PublicApi/<pkg>.<Symbol>.md`) and rewrite [PublicApi.md](/docs/References/PublicApi.md), following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).
12. Rewrite the remaining **[Rewrite](/docs/References/TemplateFileActions.md#rewrite)** docs with the new library's content: [Structure.md](/docs/References/Structure.md), [Commands.md](/docs/References/Commands.md), [Adapters.md](/docs/References/Adapters.md), and the usage guides ([InstallCli.md](/docs/Tutorials/InstallCli.md), [UseCli.md](/docs/Tutorials/UseCli.md), [LibInitialization.md](/docs/Tutorials/LibInitialization.md), [RunCliSample.md](/docs/Tutorials/RunCliSample.md), [RunApiSample.md](/docs/Tutorials/RunApiSample.md), [SamplesList.md](/docs/References/SamplesList.md), [ApiSamplesList.md](/docs/References/ApiSamplesList.md)).
13. Create the tutorials specific to the new library — one page per workflow its maintainers will repeat — following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md) and the [TutorialDocs specification](/docs/References/Specs/TutorialDocs/Specs.md). The generic guides carried over by Copy cover the structure only; they do not document the library's own use cases.
14. Create any reference page the library needs beyond the public API, following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md) and the [ReferenceDocs specification](/docs/References/Specs/ReferenceDocs/Specs.md).
15. Delete every remaining **[Delete](/docs/References/TemplateFileActions.md#delete)** file — the example internal logic, samples, and tracker docs the new library replaced. For `.md` files, follow [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md).
16. Rewrite the [README.md](/README.md) — overview, badges, and the Doc Index pointing at each theme index — then rewrite the four theme indexes (`docs/Index/<Theme>.md`) so each lists the new library's pages.

### Phase 5 — Verify
17. Build:
```bash
go build ./...
```
Then confirm every item below — the fork is only done when all pass:
- All library logic lives in `sandbox/`; no file there imports `os`, `net`, or a third-party implementation directly — every such call goes through `l.Deps`.
- `sandbox/contracts/deps/deps.go` declares one function field per injected call, and **every** adapter in `adapters/` fills all of them — the compiler does not check this.
- `sandbox/contracts/api/api.go` declares every public object as a struct with a `Deps` field, and every one of its function fields is filled by a factory registered in that package's `New` constructor.
- `sandbox/new.go` is the only wiring point, and it imports no adapter.
- Tutorials and reference pages specific to this library exist under `docs/Tutorials/` and `docs/References/`.
- Every created or rewritten file matches its specification from [Specs.md](/docs/References/Specs.md).
- Every theme index in `docs/Index/` lists every page of its theme, the `README.md` Doc Index lists every theme, and the samples lists cover every sample.
- `cmd/main/main.go` wires, calls `Sandboxmain`, and exits — it branches on no command, parses no flag, and prints nothing.
