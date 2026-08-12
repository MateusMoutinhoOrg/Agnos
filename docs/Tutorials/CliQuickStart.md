# Get Started with the CLI

## Description
Installs the `agnos-cli` binary and tracks a first budget with it, in two steps. Troubleshooting the installation is covered by [InstallCli.md](/docs/Tutorials/InstallCli.md); the full command surface is listed in [Commands.md](/docs/References/Commands.md); driving every operation from the terminal is covered by [UseCli.md](/docs/Tutorials/UseCli.md).

### Rules
- Go ≥ 1.22 must be installed.
- The binary is built from `cmd/main`, so `go install` names it `main` — rename it to `agnos-cli` before use.
- Amounts are typed in major units (`84.50`), and direction comes from the command, never from a sign.

---

## Workflow

1. Install the CLI globally — pick your OS, copy the block, paste it in a terminal.

   **macOS / Linux (bash, zsh, etc.)**

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

   **Windows (PowerShell)**

   ```powershell
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@v0.0.3; `
     if ($?) { `
       $gobin = (go env GOPATH) + '\bin'; `
       Move-Item "$gobin\main.exe" "$gobin\agnos-cli.exe" -Force; `
       if ($env:PATH -notlike "*$gobin*") { `
         [Environment]::SetEnvironmentVariable('PATH', `
           [Environment]::GetEnvironmentVariable('PATH','User') + ";$gobin", 'User'); `
         $env:PATH += ";$gobin"; `
       }; `
       agnos-cli version `
     }
   ```

   The script also adds Go's binary directory to your `PATH` persistently if it isn't already there.

2. Track your first budget:

   ```bash
   agnos-cli category add groceries
   agnos-cli category add salary

   agnos-cli received salary "august paycheck" 2500.00
   agnos-cli spend groceries "weekly shopping" 84.50

   agnos-cli transactions
   agnos-cli balance            # 2415.50
   ```

3. Read [UseCli.md](/docs/Tutorials/UseCli.md) for the rest of the operations, and [Commands.md](/docs/References/Commands.md) for every command, flag, and exit code.
