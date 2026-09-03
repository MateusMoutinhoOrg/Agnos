# Bootstrap Agnos with Itself

## Description
Covers regenerating this repository with the `agnos` compiled **from this checkout**, then rebuilding and installing the result. Agnos-Cli is one of the projects `agnos build` renders — its `sandbox/deps/deps.go`, `adapters/availables/standard/new.go`, `sandbox/internal/cli/climain.go` and every `entries.go` are generated — so after any change to a template, a collector, a `sandbox/deps/<x>/`, an `adapters/libs/<x>/` or an `entries.yaml`, the tree must be re-rendered. Cross-compiling for other platforms is [Build.md](/docs/Tutorials/Build.md); the pipeline that runs is [BuildPipeline.md](/docs/References/BuildPipeline.md).

### Rules
- **Never run an already-installed `agnos build` against this repo after changing templates or collectors.** An older binary rewrites the tree to its own, stale shape. Compile first, and run *that* binary.
- `agnos build` must stay idempotent over this tree: running it twice changes nothing, and the result compiles. A change that breaks either is not finished.
- `go vet` and `go build` take `./cmd/... ./sandbox/... ./adapters/...`, never `./...` — `assets/` holds Go templates.
- All sandbox and adapter code reaches `Deps` fields by their mechanical names (`deps.Iodeps`, `deps.Std`, …), because the struct is regenerated from the directory listing.

---

## Workflow
1. Compile the current source, before any generation, into a bootstrap binary:
   ```bash
   go build -o release/bootstrap.bin ./cmd/main
   ```
2. Let that binary regenerate the tree. It runs `verify`, re-renders the generated files, and compiles the schema directories, so a broken template or collector fails here:
   ```bash
   ./release/bootstrap.bin build
   ```
3. Check the regeneration was a fixed point — the only diffs should be the ones you meant to make:
   ```bash
   git status --short
   git diff --stat
   ```
4. Run it a second time to prove idempotence; the tree must not change:
   ```bash
   ./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
   ```
5. Rebuild for your machine from the regenerated tree and install it. On an Intel Mac:
   ```bash
   CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release/mac86.bin ./cmd/main
   rm -f release/bootstrap.bin
   sudo mv release/mac86.bin /usr/local/bin/agnos
   agnos version
   ```
   `local_install.sh`, when present in your checkout, runs steps 1, 2 and 5 in this order.
6. Bump the release version in `AgnosConfig/project.yaml` when shipping. `build` regenerates `sandbox/internal/config/config.go` from it, so the `Version` constant and `agnos version` follow.
