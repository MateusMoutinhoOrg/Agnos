# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification pairs a `Specs.md` (the description) with a `sample` (a concrete file that satisfies it).

This index is the **only** place a specification is located from. Never browse `docs/Meta/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification — see [RULES.md](/docs/RULES.md#specification-compliance).
- Every specification lives in its own directory under `docs/Meta/`, containing a `Specs.md` and a `sample` file.
- Creating, renaming, or deleting a specification requires updating this index in the same commit.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/Meta/GeneralDoc/Specs.md) · [sample](/docs/Meta/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/Meta/Readme/Specs.md) · [sample](/docs/Meta/Readme/sample.md) |
| Rules | `docs/RULES.md` | [Specs](/docs/Meta/Rules/Specs.md) · [sample](/docs/Meta/Rules/sample.md) |
| Structure | `docs/Structure.md` | [Specs](/docs/Meta/Structure/Specs.md) · [sample](/docs/Meta/Structure/sample.md) |
| AdaptersDoc | `docs/Adapters.md` | [Specs](/docs/Meta/AdaptersDoc/Specs.md) · [sample](/docs/Meta/AdaptersDoc/sample.md) |
| ReferenceDocs | Any other **reference** page in `docs/` — listable content: indexes, API detail pages, command lists — except this index and `docs/Meta/` | [Specs](/docs/Meta/ReferenceDocs/Specs.md) · [sample](/docs/Meta/ReferenceDocs/sample.md) |
| ExplanationDocs | Any **explanation** page in `docs/` — background on one mechanic, e.g. `SandboxIsolation.md` | [Specs](/docs/Meta/ExplanationDocs/Specs.md) · [sample](/docs/Meta/ExplanationDocs/sample.md) |
| TutorialDocs | Any **tutorial** in `docs/` — a single-goal workflow guide, e.g. `HandleSamples.md`, `ForkTemplate.md` | [Specs](/docs/Meta/TutorialDocs/Specs.md) · [sample](/docs/Meta/TutorialDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial follows **both** GeneralDoc and TutorialDocs. AdaptersDoc likewise builds on ReferenceDocs — `Adapters.md` follows all three.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Factories | **Every** file declaring `<Field>Factory` functions — `sandbox/internal/` and `adapters/` alike | [Specs](/docs/Meta/Factories/Specs.md) · [sample](./Meta/Factories/sample.go) |
| Deps | `sandbox/contracts/deps/deps.go` | [Specs](/docs/Meta/Deps/Specs.md) · [sample](./Meta/Deps/sample.go) |
| Outputs | `sandbox/contracts/api/api.go` | [Specs](/docs/Meta/Outputs/Specs.md) · [sample](./Meta/Outputs/sample.go) |
| Adapters | `adapters/<name>/<name>.go` | [Specs](/docs/Meta/Adapters/Specs.md) · [sample](./Meta/Adapters/sample.go) |
| LibFunctions | Factories filling `api.Lib` fields, in `sandbox/internal/lib/` | [Specs](/docs/Meta/LibFunctions/Specs.md) · [sample](./Meta/LibFunctions/sample.go) |
| LibObjects | Factories and constructors for objects the lib creates, in `sandbox/internal/<object>/` | [Specs](/docs/Meta/LibObjects/Specs.md) · [sample](./Meta/LibObjects/sample.go) |
| CliMain | `cmd/main/main.go` | [Specs](/docs/Meta/CliMain/Specs.md) · [sample](./Meta/CliMain/sample.go) |
| LibraryExamples | `examples/libraryExamples/<example>/<example>.go` | [Specs](/docs/Meta/LibraryExamples/Specs.md) · [sample](./Meta/LibraryExamples/sample.go) |
| CliExamples | `examples/cliExamples/<Name>.sh` | [Specs](/docs/Meta/CliExamples/Specs.md) · [sample](./Meta/CliExamples/sample.sh) |

Factories applies on top of the others, as GeneralDoc does for documentation: an adapter follows **both** Factories and Adapters, and a lib function follows **both** Factories and LibFunctions.

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure.md](/docs/Structure.md) and, for `.md` files, [GeneralDoc](/docs/Meta/GeneralDoc/Specs.md).
3. If a row matches, read its `Specs.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates required by [RULES.md](/docs/RULES.md).
