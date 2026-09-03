# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification pairs a `Specs.md` (the description) with a `sample` (a concrete file that satisfies it).

This index is the **only** place a specification is located from. Never browse `docs/References/Specs/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification.
- Every specification lives in its own directory under `docs/References/Specs/`, containing a `Specs.md` and a `sample` file.
- Creating, renaming, or deleting a specification requires updating this index in the same commit.
- A slot marked *generated* in [Structure.md](/docs/References/Structure.md) has no specification: its shape is the template's under `assets/`, which the AssetTemplate specification governs instead.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/References/Specs/GeneralDoc/Specs.md) · [sample](/docs/References/Specs/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/References/Specs/Readme/Specs.md) · [sample](/docs/References/Specs/Readme/sample.md) |
| Structure | `docs/References/Structure.md` | [Specs](/docs/References/Specs/Structure/Specs.md) · [sample](/docs/References/Specs/Structure/sample.md) |
| AdaptersDoc | `docs/References/Adapters.md` | [Specs](/docs/References/Specs/AdaptersDoc/Specs.md) · [sample](/docs/References/Specs/AdaptersDoc/sample.md) |
| Index | The index page of each theme — `docs/Index/<Theme>.md` | [Specs](/docs/References/Specs/Index/Specs.md) · [sample](/docs/References/Specs/Index/sample.md) |
| TutorialDocs | Any page under `docs/Tutorials/` — a single-goal workflow guide, e.g. `ScaffoldProject.md`, `AddAction.md` | [Specs](/docs/References/Specs/TutorialDocs/Specs.md) · [sample](/docs/References/Specs/TutorialDocs/sample.md) |
| ReferenceDocs | Any other **reference** page under `docs/References/` — listable content: `Commands.md`, `DepList.md`, `EntriesYaml.md`, `GeneratedFiles.md`, `PublicApi.md` and the detail pages under `docs/References/PublicApi/` — except this index and `docs/References/Specs/` | [Specs](/docs/References/Specs/ReferenceDocs/Specs.md) · [sample](/docs/References/Specs/ReferenceDocs/sample.md) |
| ExplanationDocs | Any **explanation** page under `docs/References/` — background on one mechanic, e.g. `SandboxIsolation.md`, `BuildPipeline.md` | [Specs](/docs/References/Specs/ExplanationDocs/Specs.md) · [sample](/docs/References/Specs/ExplanationDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial follows **both** GeneralDoc and TutorialDocs. AdaptersDoc likewise builds on ReferenceDocs — `Adapters.md` follows all three. `docs/Tutorials/` holds only tutorials; `docs/References/` holds reference and explanation pages; `docs/Index/` holds only theme indexes.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Contract | Hand-written files of `sandbox/api/` — `actions.go` | [Specs](/docs/References/Specs/Contract/Specs.md) · [sample](/docs/References/Specs/Contract/sample.go) |
| Binder | Hand-written files of `sandbox/binds/` — `actions.go` | [Specs](/docs/References/Specs/Binder/Specs.md) · [sample](/docs/References/Specs/Binder/sample.go) |
| DepsContract | `sandbox/deps/<x>/<x>.go`, and its copy under `assets/deplist/<x>/sandbox/deps/<x>/` | [Specs](/docs/References/Specs/DepsContract/Specs.md) · [sample](/docs/References/Specs/DepsContract/sample.go) |
| AdapterLib | `adapters/libs/<lib>/*.go`, and its copy under `assets/deplist/<x>/adapters/libs/<lib>/` | [Specs](/docs/References/Specs/AdapterLib/Specs.md) · [sample](/docs/References/Specs/AdapterLib/sample.go) |
| Available | Hand-written `adapters/availables/<name>/new.go` (not `standard`, which is generated) | [Specs](/docs/References/Specs/Available/Specs.md) · [sample](/docs/References/Specs/Available/sample.go) |
| CommandEntries | `sandbox/internal/commands/<name>/entries.yaml` | [Specs](/docs/References/Specs/CommandEntries/Specs.md) · [sample](/docs/References/Specs/CommandEntries/sample.yaml) |
| CommandHandler | `sandbox/internal/commands/<name>/handler.go` | [Specs](/docs/References/Specs/CommandHandler/Specs.md) · [sample](/docs/References/Specs/CommandHandler/sample.go) |
| Action | `sandbox/internal/actions/<name>/<name>.go` and `<name>_internal.go` | [Specs](/docs/References/Specs/Action/Specs.md) · [sample](/docs/References/Specs/Action/sample.go) |
| Collector | `sandbox/internal/actions/build/collect_*.go` | [Specs](/docs/References/Specs/Collector/Specs.md) · [sample](/docs/References/Specs/Collector/sample.go) |
| Parsable | Every file of `sandbox/internal/parsables/<name>conf/` | [Specs](/docs/References/Specs/Parsable/Specs.md) · [sample](/docs/References/Specs/Parsable/sample.go) |
| Dep | The directory `assets/deplist/<dep>/` and its entry in `assets/depsversion.yaml` | [Specs](/docs/References/Specs/Dep/Specs.md) · [sample](/docs/References/Specs/Dep/sample.md) |
| AssetTemplate | Every file under `assets/<group>/**` and `assets/templates/*` | [Specs](/docs/References/Specs/AssetTemplate/Specs.md) · [sample](/docs/References/Specs/AssetTemplate/sample.go) |
| CliMain | `assets/cli/cmd/main/main.go`, rendered to `cmd/main/main.go` | [Specs](/docs/References/Specs/CliMain/Specs.md) · [sample](/docs/References/Specs/CliMain/sample.go) |

AssetTemplate applies on top of the others for every file under `assets/`: the `iodeps` copy under `assets/deplist/` follows **both** AdapterLib and AssetTemplate.

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure.md](/docs/References/Structure.md) and, for `.md` files, [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).
3. If a row matches, read its `Specs.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates described in the relevant tutorial.
