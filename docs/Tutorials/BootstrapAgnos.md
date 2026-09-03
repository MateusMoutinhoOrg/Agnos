# Bootstrap Agnos with Itself

## Description
Covers regenerating this repository with the `agnos` compiled **from this checkout**, then rebuilding and installing the result. Agnos is one of the projects `agnos build` renders — its `sandbox/deps/deps.go`, `adapters/availables/standard/new.go`, `sandbox/internal/cli/climain.go` and every `entries.go` are generated — so after any change to a template, a collector, a `sandbox/deps/<x>/`, an `adapters/libs/<x>/` or an `entries.yaml`, the tree must be re-rendered. Cross-compiling for other platforms is [Build.md](/docs/Tutorials/Build.md); the pipeline that runs is [BuildPipeline.md](/docs/References/BuildPipeline.md).

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
5. Rebuild for your machine from the regenerated tree and install it. `CGO_ENABLED=0` keeps the binary from linking against the building machine's libc. Pick the row for your platform:

   | Platform | Build command | Install |
   |----------|---------------|---------|
   | macOS, Apple Silicon | `CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o release/macarm64.bin ./cmd/main` | `sudo mv release/macarm64.bin /usr/local/bin/agnos` |
   | macOS, Intel | `CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release/mac86.bin ./cmd/main` | `sudo mv release/mac86.bin /usr/local/bin/agnos` |
   | Linux, 64-bit Intel/AMD | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/linux86.out ./cmd/main` | `sudo mv release/linux86.out /usr/local/bin/agnos` |
   | Linux, 64-bit ARM | `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o release/linuxarm64.out ./cmd/main` | `sudo mv release/linuxarm64.out /usr/local/bin/agnos` |
   | Linux, 32-bit Intel | `CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o release/linuxi32.out ./cmd/main` | `sudo mv release/linuxi32.out /usr/local/bin/agnos` |
   | Windows, 64-bit Intel/AMD | `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o release/windows86.exe ./cmd/main` | move `release/windows86.exe` onto a directory in your `PATH` as `agnos.exe` |
   | Windows, 32-bit Intel | `CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o release/windowsi32.exe ./cmd/main` | move `release/windowsi32.exe` onto a directory in your `PATH` as `agnos.exe` |

   When building for the machine you are on, the `GOOS`/`GOARCH` prefix is optional — a bare `CGO_ENABLED=0 go build -o release/agnos ./cmd/main` targets the host. Then drop the bootstrap binary and confirm the install:
   ```bash
   rm -f release/bootstrap.bin
   agnos version
   ```
6. Bump the release version in `AgnosConfig/project.yaml` when shipping. `build` regenerates `sandbox/internal/config/config.go` from it, so the `Version` constant and `agnos version` follow.
