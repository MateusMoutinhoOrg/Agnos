# Get Started with the CLI

## Description
Installs the `agnos` binary and tracks a first budget with it, in two steps. Troubleshooting the installation is covered by [InstallCli.md](/docs/CliUsage/Tutorials/InstallCli.md); the full command surface is listed in [Commands.md](/docs/CliUsage/References/Commands.md); driving every operation from the terminal is covered by [UseCli.md](/docs/CliUsage/Tutorials/UseCli.md).

### Rules
- Go ≥ 1.22 must be installed.
- The binary is built from `cmd/main`, so `go install` names it `main` — rename it to `agnos` before use.
- Amounts are typed in major units (`84.50`), and direction comes from the command, never from a sign.

---

## Workflow

1. Install the CLI globally — pick your OS, copy the block, paste it in a terminal.

   **macOS / Linux (bash, zsh, etc.)**

   ```bash
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3 \
     && mv "$(go env GOPATH)/bin/main" "$(go env GOPATH)/bin/agnos" \
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
     && agnos version
   ```

   **Windows (PowerShell)**

   ```powershell
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3; `
     if ($?) { `
       $gobin = (go env GOPATH) + '\bin'; `
       Move-Item "$gobin\main.exe" "$gobin\agnos.exe" -Force; `
       if ($env:PATH -notlike "*$gobin*") { `
         [Environment]::SetEnvironmentVariable('PATH', `
           [Environment]::GetEnvironmentVariable('PATH','User') + ";$gobin", 'User'); `
         $env:PATH += ";$gobin"; `
       }; `
       agnos version `
     }
   ```

   The script also adds Go's binary directory to your `PATH` persistently if it isn't already there.

2. Track your first budget:

   ```bash
   agnos category add groceries
   agnos category add salary

   agnos received salary "august paycheck" 2500.00
   agnos spend groceries "weekly shopping" 84.50

   agnos transactions
   agnos balance            # 2415.50
   ```

3. Read [UseCli.md](/docs/CliUsage/Tutorials/UseCli.md) for the rest of the operations, and [Commands.md](/docs/CliUsage/References/Commands.md) for every command, flag, and exit code.
