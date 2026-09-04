# Manage the Dependencies of a Project

## Description
Covers the dependency-injection layer of a scaffolded project: turning it on and off with `deps-init` / `deps-purge`, and adding or removing capabilities with `dep-install` / `dep-remove`. What each installable dep brings is listed in [DepList](/docs/DepList/doc.md); why the sandbox receives everything through `Deps` is explained in [SandboxIsolation](/docs/SandboxIsolation/doc.md); writing a capability of your own instead of installing one is [AddAdapterLib](/docs/AddAdapterLib/doc.md).

### Rules
- A dep is named after the **contract** it installs under `sandbox/deps/<dep>/`, never after the adapter implementing it: the dep installing `sandbox/deps/argvdeps` is `argvdeps`, even though its adapter package is `verb`.
- Installing renders the dep's files, syncs `go.mod` when the dep pulls an external module, persists, and then runs `build`; removing deletes exactly those files plus any directory the removal emptied, strips the `require`, and runs `build` with the `none` runtime.
- After `dep-install` of a dep with an external module (`argvdeps`, `dbdeps`), `go mod tidy` — which `build` runs for you — downloads it. A network-less machine sees that step fail.
- `cli-init` installs `std` and `argvdeps` on its own; `cli-purge` deliberately leaves them, since other code may use them.

---

## Workflow
1. Turn the layer on. `deps-init` creates `sandbox/deps/` and `adapters/`, and the follow-up build renders `sandbox/deps/deps.go` and `adapters/availables/standard/new.go` — both empty of fields until a dep is installed — and switches `sandbox/new.go` to the `New(deps *deps.Deps)` signature:
   ```bash
   agnos deps-init
   ```
2. List what can be installed. The result goes to standard output, one name per line, so it can be piped:
   ```bash
   agnos dep-list
   agnos dep-list 2>/dev/null | grep deps
   ```
3. Install a capability. The contract lands in `sandbox/deps/<dep>/`, the implementation in `adapters/libs/<lib>/`, and the rebuild regenerates `Deps` (one `<Title> <dep>.Lib` field) and `standard.New` (one `<lib>.Bind(&deps)` call) from the directories that now exist:
   ```bash
   agnos dep-install iodeps
   agnos dep-install requestdeps
   ```
4. Check the field name the sandbox reaches the capability by: it is the title-cased directory name, and nothing else — `sandbox/deps/iodeps/` is `deps.Iodeps`. Open `sandbox/deps/deps.go` to see the struct as generated.
5. Install a dep that brings an external module. `dep-install` writes the pinned `require` into `go.mod` and the build's `go mod tidy` downloads it:
   ```bash
   agnos dep-install argvdeps
   grep Verb go.mod
   ```
6. Remove a capability the project stopped using. The follow-up build renders only, because sandbox code may still refer to the field — remove those references and rebuild yourself:
   ```bash
   agnos dep-remove requestdeps
   agnos build
   ```
7. Turn the layer off entirely. `deps-purge` removes `sandbox/deps/` and `adapters/` whole and renders `sandbox/new.go` back to `New()`. A project with a CLI layer cannot live without deps — run `cli-purge` first, or reinstall with `deps-init`:
   ```bash
   agnos deps-purge
   ```
