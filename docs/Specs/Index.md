# Specifications Index
Every specification and the files each one governs, with a page per specification

| Doc | Description |
| --- | --- |
| [GeneralDoc](/docs/Specs/GeneralDoc/doc.md) | Every `doc.md` and every other `.md` file in the project |
| [Readme](/docs/Specs/Readme/doc.md) | Root `README.md` |
| [DocProps](/docs/Specs/DocProps/doc.md) | Every `docs/**/props.yaml` |
| [Index](/docs/Specs/Index/doc.md) | The generated indexes — `docs/Index/<theme-id>.md` and every `docs/**/Index.md` |
| [TutorialDocs](/docs/Specs/TutorialDocs/doc.md) | The `doc.md` of a workflow doc — a single-goal guide, e.g. `ScaffoldProject`, `AddAction` |
| [ReferenceDocs](/docs/Specs/ReferenceDocs/doc.md) | The `doc.md` of a lookup doc — listable content, e.g. `Commands`, `DepList`, `PublicApi` and its sub-docs |
| [ExplanationDocs](/docs/Specs/ExplanationDocs/doc.md) | The `doc.md` of an explanation doc — background on one mechanic, e.g. `SandboxIsolation`, `BuildPipeline` |
| [Structure](/docs/Specs/Structure/doc.md) | `docs/Structure/doc.md` |
| [AdaptersDoc](/docs/Specs/AdaptersDoc/doc.md) | `docs/Adapters/doc.md` |
| [Contract](/docs/Specs/Contract/doc.md) | Hand-written files of `sandbox/api/` |
| [Binder](/docs/Specs/Binder/doc.md) | Hand-written files of `sandbox/binds/` |
| [DepsContract](/docs/Specs/DepsContract/doc.md) | `sandbox/deps/<x>/<x>.go`, and its copy under `assets/deplist/` |
| [AdapterLib](/docs/Specs/AdapterLib/doc.md) | `adapters/libs/<lib>/*.go`, and its copy under `assets/deplist/` |
| [Available](/docs/Specs/Available/doc.md) | Hand-written `adapters/availables/<name>/new.go` |
| [CommandEntries](/docs/Specs/CommandEntries/doc.md) | `sandbox/internal/commands/<name>/entries.yaml` |
| [CommandHandler](/docs/Specs/CommandHandler/doc.md) | `sandbox/internal/commands/<name>/handler.go` |
| [Action](/docs/Specs/Action/doc.md) | `sandbox/internal/actions/<name>/<name>.go` and `<name>_internal.go` |
| [Collector](/docs/Specs/Collector/doc.md) | `sandbox/internal/actions/build/collect_*.go` |
| [Parsable](/docs/Specs/Parsable/doc.md) | Every file of `sandbox/internal/parsables/<name>conf/` |
| [Dep](/docs/Specs/Dep/doc.md) | The directory `assets/deplist/<dep>/` and its entry in `assets/depsversion.yaml` |
| [AssetTemplate](/docs/Specs/AssetTemplate/doc.md) | Every file under `assets/<group>/**` and `assets/templates/*` |
| [CliMain](/docs/Specs/CliMain/doc.md) | `assets/cli/cmd/main/main.go`, rendered to `cmd/main/main.go` |
