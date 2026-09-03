# Add an Installable Dep

## Description
Covers adding a capability to the set `agnos dep-install` can render into a project — a directory under `assets/deplist/<dep>/` mirroring the files the dep installs, plus its pinned module in `assets/depsversion.yaml` when it pulls one. The pair every dep is made of is described in [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md); how a dep is rendered and removed is [BuildPipeline.md](/docs/References/BuildPipeline.md#deps); the shape the directory must have is the Dep specification in [Specs.md](/docs/References/Specs.md).

### Rules
- The dep's name **is the contract it installs**: the directory under `sandbox/deps/` it carries. A dep whose adapter package has a different name (`verb` for `argvdeps`) is still named after the contract.
- The tree under `assets/deplist/<dep>/` mirrors the target project: `assets/deplist/<dep>/sandbox/deps/<dep>/<dep>.go` installs to `sandbox/deps/<dep>/<dep>.go`.
- Every file in it is a Go `text/template`. Import paths of the target module are written as `{{.Module}}/…`, never as this repository's module path.
- Every lib this repository ships is mirrored as a dep, and the two copies must stay identical apart from the module path. Changing `sandbox/deps/<x>/` or `adapters/libs/<x>/` here means changing the dep too.
- A dep that pulls an external module lists it in `assets/depsversion.yaml` as `<dep>: <module>@<version>`; a dep bundling only sandbox-copy code is absent from the file and leaves `go.mod` untouched.

---

## Workflow
1. Write the contract and the adapter lib in this repository first, following [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md), and bootstrap so Agnos itself uses them — see [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md).
2. Create the dep directory mirroring the two paths:
   ```bash
   mkdir -p assets/deplist/clockdeps/sandbox/deps/clockdeps assets/deplist/clockdeps/adapters/libs/clockdeps
   ```
3. Copy both files in and replace this module's path with the template variable:
   ```bash
   sed 's#github.com/MateusMoutinhoOrg/Agnos#{{.Module}}#g' sandbox/deps/clockdeps/clockdeps.go > assets/deplist/clockdeps/sandbox/deps/clockdeps/clockdeps.go
   sed 's#github.com/MateusMoutinhoOrg/Agnos#{{.Module}}#g' adapters/libs/clockdeps/Clock.go > assets/deplist/clockdeps/adapters/libs/clockdeps/Clock.go
   ```
4. Pin the external module, if the adapter lib imports one:
   ```yaml
   # assets/depsversion.yaml
   clockdeps: github.com/some/clock@v1.2.0
   ```
5. Include any extra file the capability needs at the path it must land on. The `embeddeps` dep, for instance, also carries `assets/asset.go`, because the `//go:embed` directive has to live next to the assets it embeds.
6. Rebuild the bootstrap binary — the `assets` package embeds `all:*`, so the new directory exists at runtime with no other change — and try it against a scratch project:
   ```bash
   go build -o release/bootstrap.bin ./cmd/main
   ./release/bootstrap.bin dep-list | grep clockdeps
   ./release/bootstrap.bin dep-install clockdeps --path ../scratch-project
   ./release/bootstrap.bin dep-remove clockdeps --path ../scratch-project
   ```
7. Add the dep's row to [DepList.md](/docs/References/DepList.md) and, when the contract is public API, a detail page under `docs/References/PublicApi/` listed by [PublicApi.md](/docs/References/PublicApi.md) — required by [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md).
