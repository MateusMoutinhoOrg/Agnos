# Install the CLI

## Description
Covers installing the `agnos-cli` command-line interface globally, so it runs from any directory and survives a terminal restart. To use it once installed, follow [UseCli.md](/docs/CliUsage/Tutorials/UseCli.md); to consume the same behavior as a Go library instead, follow [LibInitialization.md](/docs/LibUsage/Tutorials/LibInitialization.md).

### Rules
- Go 1.22 or newer must be installed and available on your `PATH` (`go version` must print a version).
- The binary is built from [cmd/main](/cmd/main/), so `go install` names it `main` — the install commands below rename it to `agnos-cli`.
- The install commands detect whether `$(go env GOPATH)/bin` (or `%GOPATH%\bin` on Windows) is on your `PATH` and add it persistently if it is not.

---

## Workflow

### macOS / Linux

1. Copy and paste this entire block into a terminal (bash, zsh, or any POSIX shell):
   ```bash
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3 \
     && mv "$(go env GOPATH)/bin/main" "$(go env GOPATH)/bin/agnos-cli" \
     && { \
          case ":$PATH:" in \
            *":$(go env GOPATH)/bin:"*) ;; \
            *) \
              PROF="$HOME/.profile"; \
              [ -n "$ZSH_VERSION" ] && PROF="$HOME/.zshrc"; \
              [ -n "$BASH_VERSION" ] && [ -f "$HOME/.bashrc" ] && PROF="$HOME/.bashrc"; \
              echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> "$PROF"; \
              export PATH="$PATH:$(go env GOPATH)/bin"; \
              echo "Added GOPATH/bin to $PROF (open a new terminal or run: source $PROF)"; \
          esac; \
        } \
     && agnos-cli version
   ```
2. If `agnos-cli version` printed the version, the install is complete. Open a **new terminal** to make the PATH change available everywhere.

### Windows (PowerShell)

1. Copy and paste this entire block into PowerShell:
   ```powershell
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3; `
     if ($?) { `
       $gobin = (go env GOPATH) + '\bin'; `
       Move-Item "$gobin\main.exe" "$gobin\agnos-cli.exe" -Force; `
       if ($env:PATH -notlike "*$gobin*") { `
         [Environment]::SetEnvironmentVariable('PATH', `
           [Environment]::GetEnvironmentVariable('PATH','User') + ";$gobin", 'User'); `
         $env:PATH += ";$gobin"; `
         Write-Host "Added $gobin to your PATH (restart the terminal for full effect)"; `
       }; `
       agnos-cli version `
     }
   ```
2. If `agnos-cli version` printed the version, the install is complete. Open a **new terminal** to make the PATH change available everywhere.

### Verify after reboot

3. After restarting the machine (or opening a fresh terminal), confirm the binary is still found:
   ```bash
   agnos-cli version
   ```

### Troubleshooting

4. If `go install` fails with a network error, check your internet connection and that Go's proxy is reachable (`go env GOPROXY`).
5. If `agnos-cli version` says "command not found" after installing, `GOPATH/bin` is not on your PATH. Run `go env GOPATH` to find the directory, and add its `bin/` subdirectory to your shell profile manually:
   ```bash
   # macOS / Linux — replace .zshrc with .bashrc if you use bash
   echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
   source ~/.zshrc
   ```
   ```powershell
   # Windows PowerShell
   $gobin = (go env GOPATH) + '\bin'
   [Environment]::SetEnvironmentVariable('PATH',
     [Environment]::GetEnvironmentVariable('PATH','User') + ";$gobin", 'User')
   ```
6. If `mv` (or `Move-Item`) fails because `agnos-cli` already exists from a previous install, delete the old binary first and re-run the install command.

### Install from a Clone
Use this instead of the steps above when you are working on the project itself and want the binary built from your checkout:

1. Build it into Go's binary directory under the right name:
   ```bash
   go build -o "$(go env GOPATH)/bin/agnos-cli" ./cmd/main
   ```
2. Or skip installing entirely and run it straight from source:
   ```bash
   go run ./cmd/main category list
   ```
