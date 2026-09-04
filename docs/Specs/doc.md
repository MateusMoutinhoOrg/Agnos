# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification is a sub-doc of this one: a `doc.md` (the description) next to a `sample` file (a concrete file that satisfies it), listed in [Index.md](/docs/Specs/Index.md).

This index is the **only** place a specification is located from. Never browse `docs/Specs/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification.
- Every specification is one sub-doc directory of `docs/Specs/`, holding `doc.md`, `props.yaml` and a `sample` file. The sample is an asset: it sits next to `doc.md` and is not itself a doc.
- Creating, renaming, or deleting a specification requires updating this index in the same commit — see [HandleDocuments](/docs/HandleDocuments/doc.md).
- A slot marked *generated* in [Structure](/docs/Structure/doc.md) has no specification: its shape is the template's under `assets/`, which the AssetTemplate specification governs instead.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/Specs/GeneralDoc/doc.md) · [sample](/docs/Specs/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/Specs/Readme/doc.md) · [sample](/docs/Specs/Readme/sample.md) |
| DocProps | Every `docs/**/props.yaml` | [Specs](/docs/Specs/DocProps/doc.md) · [sample](/docs/Specs/DocProps/sample.yaml) |
| Index | The generated indexes — `docs/Index/<theme-id>.md` and every `docs/**/Index.md` | [Specs](/docs/Specs/Index/doc.md) · [sample](/docs/Specs/Index/sample.md) |
| Structure | `docs/Structure/doc.md` | [Specs](/docs/Specs/Structure/doc.md) · [sample](/docs/Specs/Structure/sample.md) |
| AdaptersDoc | `docs/Adapters/doc.md` | [Specs](/docs/Specs/AdaptersDoc/doc.md) · [sample](/docs/Specs/AdaptersDoc/sample.md) |
| TutorialDocs | The `doc.md` of a **workflow** doc — one goal, numbered steps, e.g. `ScaffoldProject`, `AddAction` | [Specs](/docs/Specs/TutorialDocs/doc.md) · [sample](/docs/Specs/TutorialDocs/sample.md) |
| ReferenceDocs | The `doc.md` of a **lookup** doc — listable content: `Commands`, `DepList`, `EntriesYaml`, `GeneratedFiles`, `PublicApi` and its sub-docs — except this index and its own sub-docs | [Specs](/docs/Specs/ReferenceDocs/doc.md) · [sample](/docs/Specs/ReferenceDocs/sample.md) |
| ExplanationDocs | The `doc.md` of an **explanation** doc — background on one mechanic, e.g. `SandboxIsolation`, `BuildPipeline` | [Specs](/docs/Specs/ExplanationDocs/doc.md) · [sample](/docs/Specs/ExplanationDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial doc follows **both** GeneralDoc and TutorialDocs. AdaptersDoc likewise builds on ReferenceDocs — `docs/Adapters/doc.md` follows all three. Which of the three kinds a doc is comes from the page itself, not from where it sits: every doc is a directory of `docs/`.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Contract | Hand-written files of `sandbox/api/` — `actions.go` | [Specs](/docs/Specs/Contract/doc.md) · [sample](/docs/Specs/Contract/sample.go) |
| Binder | Hand-written files of `sandbox/binds/` — `actions.go` | [Specs](/docs/Specs/Binder/doc.md) · [sample](/docs/Specs/Binder/sample.go) |
| DepsContract | `sandbox/deps/<x>/<x>.go`, and its copy under `assets/deplist/<x>/sandbox/deps/<x>/` | [Specs](/docs/Specs/DepsContract/doc.md) · [sample](/docs/Specs/DepsContract/sample.go) |
| AdapterLib | `adapters/libs/<lib>/*.go`, and its copy under `assets/deplist/<x>/adapters/libs/<lib>/` | [Specs](/docs/Specs/AdapterLib/doc.md) · [sample](/docs/Specs/AdapterLib/sample.go) |
| Available | Hand-written `adapters/availables/<name>/new.go` (not `standard`, which is generated) | [Specs](/docs/Specs/Available/doc.md) · [sample](/docs/Specs/Available/sample.go) |
| CommandEntries | `sandbox/internal/commands/<name>/entries.yaml` | [Specs](/docs/Specs/CommandEntries/doc.md) · [sample](/docs/Specs/CommandEntries/sample.yaml) |
| CommandHandler | `sandbox/internal/commands/<name>/handler.go` | [Specs](/docs/Specs/CommandHandler/doc.md) · [sample](/docs/Specs/CommandHandler/sample.go) |
| Action | `sandbox/internal/actions/<name>/<name>.go` and `<name>_internal.go` | [Specs](/docs/Specs/Action/doc.md) · [sample](/docs/Specs/Action/sample.go) |
| Collector | `sandbox/internal/actions/build/collect_*.go` | [Specs](/docs/Specs/Collector/doc.md) · [sample](/docs/Specs/Collector/sample.go) |
| Parsable | Every file of `sandbox/internal/parsables/<name>conf/` | [Specs](/docs/Specs/Parsable/doc.md) · [sample](/docs/Specs/Parsable/sample.go) |
| Dep | The directory `assets/deplist/<dep>/` and its entry in `assets/depsversion.yaml` | [Specs](/docs/Specs/Dep/doc.md) · [sample](/docs/Specs/Dep/sample.md) |
| AssetTemplate | Every file under `assets/<group>/**` and `assets/templates/*` | [Specs](/docs/Specs/AssetTemplate/doc.md) · [sample](/docs/Specs/AssetTemplate/sample.go) |
| CliMain | `assets/cli/cmd/main/main.go`, rendered to `cmd/main/main.go` | [Specs](/docs/Specs/CliMain/doc.md) · [sample](/docs/Specs/CliMain/sample.go) |

AssetTemplate applies on top of the others for every file under `assets/`: the `iodeps` copy under `assets/deplist/` follows **both** AdapterLib and AssetTemplate.

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure](/docs/Structure/doc.md) and, for `.md` files, [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).
3. If a row matches, read its `doc.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates described in the relevant tutorial.
