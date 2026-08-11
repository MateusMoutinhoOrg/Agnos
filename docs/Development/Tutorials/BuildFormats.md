# Build for Different Formats

## Description
Cross-compile the CLI into platform-specific binaries and packages — Windows `.exe`, Linux binary, Debian `.deb`, and Fedora/RHEL `.rpm` — using the build script under [build/](/docs/Development/References/Structure.md#build). The script runs each build inside a container (Docker or Podman), so the host needs no cross-compiler toolchain. Renaming the module path across the whole repository is covered at the end.

### Rules
- The build script must be run from the **repository root**.
- A container runtime — Docker or Podman — must be installed and running.
- Build artifacts are written to `release/`, which is git-ignored.

---

## Workflow

### Build a single target

1. Pick a target from the table below:

   | Target | Output | Description |
   |--------|--------|-------------|
   | `i32` | `release/windows-i32.exe` | Windows 386 executable |
   | `linux86` | `release/linux86.out` | Linux amd64 binary |
   | `deb` | `release/x86-deb.deb` | Debian/Ubuntu package |
   | `rpm` | `release/x86-rpm.rpm` | Fedora/RHEL package |
   | `mac-arm64` | `release/mac-arm64.bin` | macOS Apple Silicon binary |
   | `mac-amd64` | `release/mac-amd64.bin` | macOS Intel binary |

2. Run the build script with that target. The script auto-detects Docker or Podman:

   ```sh
   go run ./build/main build linux86
   ```

### Build multiple targets at once

3. List every target you want in one command:

   ```sh
   go run ./build/main build rpm deb i32 linux86 mac-arm64 mac-amd64
   ```

### Choose a specific container provider

4. Pass `--provider` to skip auto-detection:

   ```sh
   go run ./build/main build i32 --provider docker
   ```

   ```sh
   go run ./build/main build rpm deb --provider podman
   ```

### Rename the module path

5. To rename the module across every `.go` file and `go.mod`:

   ```sh
   go run ./build/main rename github.com/MateusMoutinhoOrg/NewName
   ```

   The command walks the repository, replaces every occurrence of the current module path with the new one, and prints each updated file.
