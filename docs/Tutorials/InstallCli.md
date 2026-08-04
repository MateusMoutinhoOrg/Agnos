# Install the CLI

## Description
Covers installing the `agnos` command-line interface globally, so it runs from any directory. To use it once installed, follow [UseCli.md](/docs/Tutorials/UseCli.md); to consume the same behavior as a Go library instead, follow [LibInitialization.md](/docs/Tutorials/LibInitialization.md).

### Rules
- Go 1.22 or newer must be installed, and `$(go env GOPATH)/bin` must be on your `PATH`.
- The binary is built from [cmd/main](/cmd/main/), so `go install` names it `main` — the second half of the command below renames it to `agnos`.

---

## Workflow
1. Install the binary and give it its name:
   ```bash
   go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@latest && \
     mv "$(go env GOPATH)/bin/main" "$(go env GOPATH)/bin/agnos"
   ```
2. Put Go's binary directory on your `PATH` if it is not already there, adding the line to your shell profile so it survives a new terminal:
   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```
3. Check the installation:
   ```bash
   agnos version
   ```
4. Choose where the records are kept, if the default of `.agnos` in your home directory is not what you want:
   ```bash
   export AGNOS_DATA="$HOME/budgets/personal"
   ```
5. Track a first transaction, following [UseCli.md](/docs/Tutorials/UseCli.md).

### Install from a Clone
Use this instead of step 1 when you are working on the project itself and want the binary built from your checkout:

1. Build it into Go's binary directory under the right name:
   ```bash
   go build -o "$(go env GOPATH)/bin/agnos" ./cmd/main
   ```
2. Or skip installing entirely and run it straight from source:
   ```bash
   go run ./cmd/main category list
   ```
