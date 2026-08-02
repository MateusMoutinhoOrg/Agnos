# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification pairs a `Specs.md` (the description) with a `sample` (a concrete file that satisfies it).

This index is the **only** place a specification is located from. Never browse `docs/References/Meta/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification — see [RULES.md](/docs/References/RULES.md#specification-compliance).
- Every specification lives in its own directory under `docs/References/Meta/`, containing a `Specs.md` and a `sample` file.
- Creating, renaming, or deleting a specification requires updating this index in the same commit.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/References/Meta/GeneralDoc/Specs.md) · [sample](/docs/References/Meta/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/References/Meta/Readme/Specs.md) · [sample](/docs/References/Meta/Readme/sample.md) |
| Rules | `docs/References/RULES.md` | [Specs](/docs/References/Meta/Rules/Specs.md) · [sample](/docs/References/Meta/Rules/sample.md) |
| Structure | `docs/References/Structure.md` | [Specs](/docs/References/Meta/Structure/Specs.md) · [sample](/docs/References/Meta/Structure/sample.md) |
| AdaptersDoc | `docs/References/Adapters.md` | [Specs](/docs/References/Meta/AdaptersDoc/Specs.md) · [sample](/docs/References/Meta/AdaptersDoc/sample.md) |
| ReferenceDocs | Any other `.md` under `docs/References/`, except this index and `docs/References/Meta/` | [Specs](/docs/References/Meta/ReferenceDocs/Specs.md) · [sample](/docs/References/Meta/ReferenceDocs/sample.md) |
| ExplanationDocs | Any `.md` under `docs/Explanations/` | [Specs](/docs/References/Meta/ExplanationDocs/Specs.md) · [sample](/docs/References/Meta/ExplanationDocs/sample.md) |
| TutorialDocs | Any `.md` under `docs/Tutorials/` | [Specs](/docs/References/Meta/TutorialDocs/Specs.md) · [sample](/docs/References/Meta/TutorialDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial follows **both** GeneralDoc and TutorialDocs. AdaptersDoc likewise builds on ReferenceDocs — `Adapters.md` follows all three.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Factories | **Every** file declaring `<Field>Factory` functions — `sandbox/internal/` and `adapters/` alike | [Specs](/docs/References/Meta/Factories/Specs.md) · [sample](./Meta/Factories/sample.go) |
| Deps | `sandbox/contracts/deps/deps.go` | [Specs](/docs/References/Meta/Deps/Specs.md) · [sample](./Meta/Deps/sample.go) |
| Outputs | `sandbox/contracts/api/api.go` | [Specs](/docs/References/Meta/Outputs/Specs.md) · [sample](./Meta/Outputs/sample.go) |
| Adapters | `adapters/<name>/<name>.go` | [Specs](/docs/References/Meta/Adapters/Specs.md) · [sample](./Meta/Adapters/sample.go) |
| LibFunctions | Factories filling `api.Lib` fields, in `sandbox/internal/lib/` | [Specs](/docs/References/Meta/LibFunctions/Specs.md) · [sample](./Meta/LibFunctions/sample.go) |
| LibObjects | Factories and constructors for objects the lib creates, in `sandbox/internal/<object>/` | [Specs](/docs/References/Meta/LibObjects/Specs.md) · [sample](./Meta/LibObjects/sample.go) |
| Examples | `examples/<example>/<example>.go` | [Specs](/docs/References/Meta/Examples/Specs.md) · [sample](./Meta/Examples/sample.go) |

Factories applies on top of the others, as GeneralDoc does for documentation: an adapter follows **both** Factories and Adapters, and a lib function follows **both** Factories and LibFunctions.

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure.md](/docs/References/Structure.md) and, for `.md` files, [GeneralDoc](/docs/References/Meta/GeneralDoc/Specs.md).
3. If a row matches, read its `Specs.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates required by [RULES.md](/docs/References/RULES.md).
