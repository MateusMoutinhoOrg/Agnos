# Dep Specification

## Description
Defines the required shape of an installable **dep** — a directory `assets/deplist/<dep>/` that `agnos dep-install <dep>` renders into a project, plus its entry in `assets/depsversion.yaml` when it pulls an external module. Adding one is [HandleDeplist.md](/docs/Tutorials/HandleDeplist.md); the shipped deps are listed in [DepList.md](/docs/References/DepList.md).

### Rules
- The directory is named after the **contract** it installs — the `sandbox/deps/<dep>/` it carries — never after the adapter lib (`argvdeps`, not `verb`).
- The tree under it mirrors the target project: every file sits at the path it will be written to. A dep carries at least `sandbox/deps/<dep>/<dep>.go` (a [DepsContract](/docs/References/Specs/DepsContract/Specs.md)) and `adapters/libs/<lib>/<File>.go` (an [AdapterLib](/docs/References/Specs/AdapterLib/Specs.md)), and any other file the capability needs where it must land (`embeddeps` carries `assets/asset.go`).
- Every file is a Go `text/template` following the [AssetTemplate](/docs/References/Specs/AssetTemplate/Specs.md) specification; the target module path is `{{.Module}}`, and that is the only variable a dep uses.
- A dep whose adapter lib imports an external module has one line in `assets/depsversion.yaml`, `<dep>: <module>@<version>`; a dep bundling only sandbox-copy code has none.
- Every lib this repository ships is mirrored as a dep, and the copies are byte-identical to `sandbox/deps/<x>/` and `adapters/libs/<x>/` apart from the module path. A change to either side is a change to both.
- Installing must leave the project compiling after `dep-install`'s follow-up build, with nothing else required of the user but the `go mod tidy` that build runs.

## Structure
1. **`assets/deplist/<dep>/sandbox/deps/<dep>/<dep>.go`**: the contract.
2. **`assets/deplist/<dep>/adapters/libs/<lib>/<File>.go`**: the implementation.
3. **Other files** *(optional)*: at their target paths.
4. **`assets/depsversion.yaml` entry** *(when an external module is pulled)*.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/Dep/sample.md), which lists the tree of one shipped dep.
