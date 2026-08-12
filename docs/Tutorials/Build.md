# Build the CLI for Every Architecture

## Description
Cross-compile the CLI into a binary for each supported operating system and architecture, using the scripts under [scripts/](/docs/References/Structure.md#scripts). Every script is a thin wrapper over `go build`: the Go runtime cross-compiles on its own, so no container, no cross-compiler, and no SDK of the target platform is needed. Installing the CLI on the machine you are on is [InstallCli.md](/docs/Tutorials/InstallCli.md); renaming the module before a release is [RenameModule.md](/docs/Tutorials/RenameModule.md).

### Rules
- The **Go runtime is the only requirement** — a script must never reach for a container runtime, a packager, or a C toolchain.
- Every script writes its artifact to `release/`, which is git-ignored, and creates the directory if it is missing.
- Scripts resolve the repository root from their own path, so they run from any working directory.
- `CGO_ENABLED=0` on every target: the binary must not link against the building machine's libc.

---

## Workflow

### Build a single target

1. Pick a target from the table below:

   | Script | `GOOS`/`GOARCH` | Output | Platform |
   |--------|-----------------|--------|----------|
   | `scripts/linux86.sh` | `linux`/`amd64` | `release/linux86.out` | Linux, 64-bit Intel/AMD |
   | `scripts/linuxarm64.sh` | `linux`/`arm64` | `release/linuxarm64.out` | Linux, 64-bit ARM |
   | `scripts/linuxi32.sh` | `linux`/`386` | `release/linuxi32.out` | Linux, 32-bit Intel |
   | `scripts/windows86.sh` | `windows`/`amd64` | `release/windows86.exe` | Windows, 64-bit Intel/AMD |
   | `scripts/windowsi32.sh` | `windows`/`386` | `release/windowsi32.exe` | Windows, 32-bit Intel |
   | `scripts/mac86.sh` | `darwin`/`amd64` | `release/mac86.bin` | macOS, Intel |
   | `scripts/macarm64.sh` | `darwin`/`arm64` | `release/macarm64.bin` | macOS, Apple Silicon |

2. Run that script:

   ```bash
   bash ./scripts/linux86.sh
   ```

   It prints the artifact it wrote:

   ```text
   built release/linux86.out
   ```

### Build every target at once

3. Run `all.sh`, which runs every script in the table in order:

   ```bash
   bash ./scripts/all.sh
   ```

4. Collect the artifacts from `release/`:

   ```bash
   ls release/
   ```

### Add a new target

5. Copy the closest existing script to `scripts/<target>.sh`, and change the three things that differ: the header comment, the `GOOS`/`GOARCH` pair, and the output name.

   ```bash
   # scripts/linuxarmv7.sh — build the CLI for Linux armv7 into release/linuxarmv7.out.
   CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 \
   	go build -C "$root" -o release/linuxarmv7.out ./cmd/main
   ```

   Run `go tool dist list` to see every `GOOS`/`GOARCH` pair the installed Go runtime supports.

6. Make it executable:

   ```bash
   chmod +x ./scripts/<target>.sh
   ```

7. Add the target to the `targets` array in `scripts/all.sh`, so building everything includes it.

8. Register the new script in the target table above, and in [Structure.md](/docs/References/Structure.md#scripts) — required by [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md).
