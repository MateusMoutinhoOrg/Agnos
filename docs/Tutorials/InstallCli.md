# Install the CLI

## Description
Covers installing the `agnos-cli` command-line interface globally, so it runs from any directory and survives a terminal restart. To use it once installed, follow [UseCli.md](/docs/Tutorials/UseCli.md); to consume the same behavior as a Go library instead, follow [LibInitialization.md](/docs/Tutorials/LibInitialization.md).

### Rules
- You do **not** need Go installed to run `agnos-cli`. You can download the pre-compiled binary for your operating system directly.
- The installation commands use `curl` to fetch the latest release from the repository.

---

## Workflow

### macOS

1. Copy and paste the command below corresponding to your Mac's processor into your terminal to download and install the binary:

   **Apple Silicon (M1/M2/M3 / arm64)**
   ```bash
   curl -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/macarm64.bin -o agnos-cli && \
     chmod +x agnos-cli && \
     sudo mv agnos-cli /usr/local/bin/ && \
     agnos-cli version
   ```

   **Intel (x86_64)**
   ```bash
   curl -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/mac86.bin -o agnos-cli && \
     chmod +x agnos-cli && \
     sudo mv agnos-cli /usr/local/bin/ && \
     agnos-cli version
   ```

2. Open a **new terminal** to make sure the binary is available everywhere.

### Linux

1. Copy and paste the command below corresponding to your architecture into your terminal to download and install the binary:

   **x86_64 (64-bit)**
   ```bash
   curl -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/linux86.out -o agnos-cli && \
     chmod +x agnos-cli && \
     sudo mv agnos-cli /usr/local/bin/ && \
     agnos-cli version
   ```

   **ARM64**
   ```bash
   curl -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/linuxarm64.out -o agnos-cli && \
     chmod +x agnos-cli && \
     sudo mv agnos-cli /usr/local/bin/ && \
     agnos-cli version
   ```

   **x86 (32-bit)**
   ```bash
   curl -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/linuxi32.out -o agnos-cli && \
     chmod +x agnos-cli && \
     sudo mv agnos-cli /usr/local/bin/ && \
     agnos-cli version
   ```

2. Open a **new terminal** to make sure the binary is available everywhere.

### Windows (PowerShell)

1. Open PowerShell and copy and paste the command below corresponding to your architecture. It will download the binary, install it to your user profile, and update your `PATH`:

   **x86_64 (64-bit)**
   ```powershell
   $dir = "$HOME\.local\bin"; `
     New-Item -ItemType Directory -Force -Path $dir | Out-Null; `
     curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/windows86.exe -o "$dir\agnos-cli.exe"; `
     if ($env:PATH -notlike "*$dir*") { `
       [Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH', 'User') + ";$dir", 'User'); `
       $env:PATH += ";$dir"; `
       Write-Host "Added $dir to your PATH (restart the terminal for full effect)"; `
     }; `
     agnos-cli.exe version
   ```

   **x86 (32-bit)**
   ```powershell
   $dir = "$HOME\.local\bin"; `
     New-Item -ItemType Directory -Force -Path $dir | Out-Null; `
     curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest/download/windowsi32.exe -o "$dir\agnos-cli.exe"; `
     if ($env:PATH -notlike "*$dir*") { `
       [Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH', 'User') + ";$dir", 'User'); `
       $env:PATH += ";$dir"; `
       Write-Host "Added $dir to your PATH (restart the terminal for full effect)"; `
     }; `
     agnos-cli.exe version
   ```

2. Open a **new terminal** to make sure the binary is available everywhere.

### Verify after reboot

3. After restarting the machine (or opening a fresh terminal), confirm the binary is still found globally:
   ```bash
   # On macOS / Linux
   agnos-cli version
   ```
   ```powershell
   # On Windows
   agnos-cli.exe version
   ```

### Troubleshooting

4. If `curl` fails with a network error, check your internet connection and ensure GitHub is accessible.
5. If `agnos-cli version` says `command not found` after installing, the directory you moved the binary to (e.g., `/usr/local/bin` ou `$HOME\.local\bin`) is not in your `PATH`. Add it to your shell profile manually.

### Install from a Clone (Requires Go)

1. Use this instead of the steps above when you are working on the project itself and want the binary built from your checkout:
   ```bash
   go build -o $(go env GOPATH)/bin/agnos-cli ./cmd/main
   ```
2. Or skip installing entirely and run it straight from source:
   ```bash
   go run ./cmd/main category list
   ```
