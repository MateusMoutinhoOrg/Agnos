# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification pairs a `Specs.md` (the description) with a `sample` (a concrete file that satisfies it).

This index is the **only** place a specification is located from. Never browse `docs/References/Specs/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification.
- Every specification lives in its own directory under `docs/References/Specs/`, containing a `Specs.md` and a `sample` file.
- Creating, renaming, or deleting a specification requires updating this index in the same commit.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/References/Specs/GeneralDoc/Specs.md) · [sample](/docs/References/Specs/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/References/Specs/Readme/Specs.md) · [sample](/docs/References/Specs/Readme/sample.md) |
| Structure | `docs/References/Structure.md` | [Specs](/docs/References/Specs/Structure/Specs.md) · [sample](/docs/References/Specs/Structure/sample.md) |
| AdaptersDoc | `docs/References/Adapters.md` | [Specs](/docs/References/Specs/AdaptersDoc/Specs.md) · [sample](/docs/References/Specs/AdaptersDoc/sample.md) |
| Index | The index page of each theme — `docs/Index/<Theme>.md` | [Specs](/docs/References/Specs/Index/Specs.md) · [sample](/docs/References/Specs/Index/sample.md) |
| TutorialDocs | Any page under `docs/Tutorials/` — a single-goal workflow guide, e.g. `HandleLibrarySamples.md`, `ForkTemplate.md` | [Specs](/docs/References/Specs/TutorialDocs/Specs.md) · [sample](/docs/References/Specs/TutorialDocs/sample.md) |
| ReferenceDocs | Any other **reference** page under `docs/References/` — listable content: indexes, command lists, and the API detail pages under `docs/References/PublicApi/` — except this index and `docs/References/Specs/` | [Specs](/docs/References/Specs/ReferenceDocs/Specs.md) · [sample](/docs/References/Specs/ReferenceDocs/sample.md) |
| ExplanationDocs | Any **explanation** page under `docs/References/` — background on one mechanic, e.g. `SandboxIsolation.md` | [Specs](/docs/References/Specs/ExplanationDocs/Specs.md) · [sample](/docs/References/Specs/ExplanationDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial follows **both** GeneralDoc and TutorialDocs. AdaptersDoc likewise builds on ReferenceDocs — `Adapters.md` follows all three. `docs/Tutorials/` holds only tutorials; `docs/References/` holds reference and explanation pages; `docs/Index/` holds only theme indexes.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Factories | **Every** file declaring `<Field>Factory` functions — `sandbox/` and `adapters/` alike | [Specs](/docs/References/Specs/Factories/Specs.md) · [sample](/docs/References/Specs/Factories/sample.go) |
| Deps | `sandbox/contracts/deps/deps.go` | [Specs](/docs/References/Specs/Deps/Specs.md) · [sample](/docs/References/Specs/Deps/sample.go) |
| Outputs | `sandbox/contracts/api/api.go` | [Specs](/docs/References/Specs/Outputs/Specs.md) · [sample](/docs/References/Specs/Outputs/sample.go) |
| Adapters | `adapters/<name>/<name>.go` | [Specs](/docs/References/Specs/Adapters/Specs.md) · [sample](/docs/References/Specs/Adapters/sample.go) |
| LibFunctions | Factories filling `api.Lib` fields, in `sandbox/lib/` | [Specs](/docs/References/Specs/LibFunctions/Specs.md) · [sample](/docs/References/Specs/LibFunctions/sample.go) |
| LibObjects | Factories and constructors for objects the lib creates, in `sandbox/<object>/` | [Specs](/docs/References/Specs/LibObjects/Specs.md) · [sample](/docs/References/Specs/LibObjects/sample.go) |
| CliMain | `cmd/main/main.go` | [Specs](/docs/References/Specs/CliMain/Specs.md) · [sample](/docs/References/Specs/CliMain/sample.go) |
| LibraryExamples | `examples/libraryExamples/<example>/<example>.go` | [Specs](/docs/References/Specs/LibraryExamples/Specs.md) · [sample](/docs/References/Specs/LibraryExamples/sample.go) |
| CliExamples | `examples/cliExamples/<Name>.sh` | [Specs](/docs/References/Specs/CliExamples/Specs.md) · [sample](/docs/References/Specs/CliExamples/sample.sh) |

Factories applies on top of the others, as GeneralDoc does for documentation: an adapter follows **both** Factories and Adapters, and a lib function follows **both** Factories and LibFunctions.

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure.md](/docs/References/Structure.md) and, for `.md` files, [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).
3. If a row matches, read its `Specs.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates described in the relevant tutorial.
