# Adapt a Pre-Existing Library

## Description
Covers converting a library that already exists into this project's dependency-injected structure. To start a new library from scratch, follow [ForkTemplate.md](/docs/ForkTemplate.md) instead. The steps are grouped into phases so progress is easy to track; every phase takes each file's action — **Copy**, **Create**, **Rewrite**, or **Delete** — from [TemplateFileActions.md](/docs/TemplateFileActions.md).

### Rules
- Read [RULES.md](/docs/RULES.md) and [Structure.md](/docs/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/Structure.md): public contract structs in `sandbox/contracts/`, internal factories in `sandbox/internal/`, concrete dependencies in `adapters/`, the entry point in `sandbox/`, and the installed binary in `cmd/main/`. The command-line interface belongs to the library, as the `Sandboxmain` field of `api.Lib`, never to the binary. Contracts are structs of function fields, never interfaces — see [StructContracts.md](/docs/StructContracts.md).
- The pre-existing package layout does **not** survive: all library logic ends up in `sandbox/internal/`, calling every OS-bound and third-party dependency through `l.Deps`. Code left in its original packages, or still calling `os`/`net`/third-party APIs directly, is not adapted.
- Every public type the library returns becomes a contract struct in `sandbox/contracts/api`, whose function fields are filled by factories in `sandbox/internal/`. A type still declared in `sandbox/internal/` and handed back to callers is not adapted.
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/Specs.md). A file that ignores its specification is not adapted.
- The adaptation is not complete until the final checklist in the last workflow step passes.

---

## Workflow

### Phase 1 — Bring the structure in
1. Recreate this project's directory layout inside the library being converted, using [Structure.md](/docs/Structure.md) as reference.
2. Copy every **[Copy](/docs/TemplateFileActions.md#copy)** file into the library unchanged — the specifications, rules, generic guides, explanations, and `sandbox/new.go`.

### Phase 2 — Rewrite the contracts
3. Rewrite `sandbox/contracts/deps/deps.go` with the OS-bound and third-party calls the library must receive as dependencies, one function field each, following [HandleDependencies.md](/docs/HandleDependencies.md).
4. Rewrite `sandbox/contracts/api/api.go` with the `Lib` struct and one struct per type the library hands back, following [HandleLibElements.md](/docs/HandleLibElements.md).
5. Rewrite `adapters/standard/standard.go` so the default adapter fills every field of that contract with the library's current behavior, following [HandleAdapters.md](/docs/HandleAdapters.md).

### Phase 3 — Move the code into the sandbox
6. Rewrite the existing library code into `sandbox/internal/`: move each source file in, turn each public function into a `<Field>Factory(l *api.Lib)` that returns a closure for the matching api field, assign every factory's return value from the package's `New` constructor, and replace **every** OS-bound or third-party call with a call through `l.Deps.<Field>(...)`, following [HandleLibElements.md](/docs/HandleLibElements.md). Do not keep the code in its original packages, leave methods on internal types, or leave direct calls in place.
7. Create the command dispatch behind `Sandboxmain` in `sandbox/internal/cli/`, following [HandleCliCommands.md](/docs/HandleCliCommands.md).
8. Create any additional adapter in `adapters/`, following [HandleAdapters.md](/docs/HandleAdapters.md).
9. Create the samples demonstrating the converted entry points: the Go programs in `examples/libraryExamples/`, following [HandleSamples.md](/docs/HandleSamples.md), and the shell scripts in `examples/cliExamples/`, following [HandleCliExamples.md](/docs/HandleCliExamples.md).

### Phase 4 — Rewrite the documentation
10. Create the API detail pages (`docs/PublicApi/<pkg>.<Symbol>.md`) and rewrite `docs/PublicApi.md`, following [ExposePublicApi.md](/docs/ExposePublicApi.md).
11. Rewrite the remaining **[Rewrite](/docs/TemplateFileActions.md#rewrite)** docs with the converted library's content: `docs/Structure.md`, `docs/Cli.md`, `docs/Adapters.md`, and the usage guides ([InstallCli.md](/docs/InstallCli.md), [UseCli.md](/docs/UseCli.md), [LibInitialization.md](/docs/LibInitialization.md), [RunCliSample.md](/docs/RunCliSample.md), [RunApiSample.md](/docs/RunApiSample.md), [SamplesList.md](/docs/SamplesList.md), [ApiSamplesList.md](/docs/ApiSamplesList.md)).
12. Create the tutorials specific to the converted library — one page per workflow its maintainers will repeat (e.g. adding a domain object, extending a feature, releasing) — following [HandleDocuments.md](/docs/HandleDocuments.md) and the [TutorialDocs specification](/docs/Meta/TutorialDocs/Specs.md). The generic guides copied in Phase 1 cover the structure only; they do not document the library's own use cases.
13. Create any reference page the library needs beyond the public API, following [HandleDocuments.md](/docs/HandleDocuments.md) and the [ReferenceDocs specification](/docs/Meta/ReferenceDocs/Specs.md).
14. Delete every **[Delete](/docs/TemplateFileActions.md#delete)** file carried over from the template, plus the pre-existing code the converted lib replaced. For `.md` files, follow [HandleDocuments.md](/docs/HandleDocuments.md).
15. Rewrite the `README.md`: overview, both quick starts, badges, Doc Index, and the two Examples sections.

### Phase 5 — Verify
16. Build:
```bash
go build ./...
```
Then confirm every item below — the adaptation is only done when all pass:
- All library logic lives in `sandbox/internal/`; no file there imports `os`, `net`, or a third-party implementation directly — every such call goes through `l.Deps`.
- `sandbox/contracts/deps/deps.go` declares one function field per injected call, and **every** adapter in `adapters/` fills all of them — the compiler does not check this.
- `sandbox/contracts/api/api.go` declares every public object as a struct with a `Deps` field, and every one of its function fields is filled by a factory registered in that package's `New` constructor.
- `sandbox/new.go` is the only wiring point, and it imports no adapter.
- Tutorials and reference pages specific to this library exist under `docs/`.
- Every created or rewritten file matches its specification from [Specs.md](/docs/Specs.md).
- The `README.md` Doc Index lists every `.md` file, and the CLI Examples and Library Examples sections list every sample.
- `cmd/main/main.go` wires, calls `Sandboxmain`, and exits — it branches on no command, parses no flag, and prints nothing.
