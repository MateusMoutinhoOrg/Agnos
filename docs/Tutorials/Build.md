# Build the CLI for Every Architecture

## Description
Cross-compile `agnos` into a binary for each supported operating system and architecture. Every target is one `go build` with `GOOS`/`GOARCH` set: the Go runtime cross-compiles on its own, so no container, no cross-compiler, and no SDK of the target platform is needed. Regenerating this checkout with itself before a release is [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md); installing a release is [InstallCli.md](/docs/Tutorials/InstallCli.md).

### Rules
- The **Go runtime is the only requirement** — never reach for a container runtime, a packager, or a C toolchain.
- Every artifact goes to `release/`, which is git-ignored.
- `CGO_ENABLED=0` on every target: the binary must not link against the building machine's libc.
- Compile `./cmd/main` and nothing wider. `./...` would try to compile `assets/`, which holds Go **templates**, not Go.

---

## Workflow

### Build a single target

1. Pick a target from the table below:

   | `GOOS`/`GOARCH` | Output | Platform |
   |-----------------|--------|----------|
   | `linux`/`amd64` | `release/linux86.out` | Linux, 64-bit Intel/AMD |
   | `linux`/`arm64` | `release/linuxarm64.out` | Linux, 64-bit ARM |
   | `linux`/`386` | `release/linuxi32.out` | Linux, 32-bit Intel |
   | `windows`/`amd64` | `release/windows86.exe` | Windows, 64-bit Intel/AMD |
   | `windows`/`386` | `release/windowsi32.exe` | Windows, 32-bit Intel |
   | `darwin`/`amd64` | `release/mac86.bin` | macOS, Intel |
   | `darwin`/`arm64` | `release/macarm64.bin` | macOS, Apple Silicon |

2. Build it from the repository root:

   ```bash
   mkdir -p release
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/linux86.out ./cmd/main
   ```

### Build every target at once

3. Loop over the table:

   ```bash
   mkdir -p release
   build() { CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -o "release/$3" ./cmd/main && echo "built release/$3"; }
   build linux   amd64 linux86.out
   build linux   arm64 linuxarm64.out
   build linux   386   linuxi32.out
   build windows amd64 windows86.exe
   build windows 386   windowsi32.exe
   build darwin  amd64 mac86.bin
   build darwin  arm64 macarm64.bin
   ```

4. Collect the artifacts from `release/`:

   ```bash
   ls release/
   ```

### Build with `agnos compile`

An installed `agnos` does every step above in one command — it runs `build` (schema
check, re-render, `go build` over the tree), creates `release/`, then cross-compiles each
target with `CGO_ENABLED=0`:

```bash
agnos compile --target all                        # every target in the table
agnos compile --target linux86 --target macarm64  # a subset; --target repeats
agnos compile --target windows86 --path ./my-tool # another project
```

The target names are `linux86`, `linuxarm64`, `linuxi32`, `mac86`, `macarm64`,
`windows86`, `windowsi32` — the same rows as the table above. An unknown name is a usage
error. Adding a row to that table means teaching `agnos compile` the name too: the map is
`targets` in [sandbox/internal/actions/compile/compile.go](/sandbox/internal/actions/compile/compile.go).

### Add a new target

5. Run `go tool dist list` to see every `GOOS`/`GOARCH` pair the installed Go runtime supports, and pick one.

6. Build it with the same three variables changed — the pair and the output name:

   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o release/linuxarmv7.out ./cmd/main
   ```

7. Add the pair to the target table above.
