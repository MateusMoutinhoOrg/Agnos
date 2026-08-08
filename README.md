# Agnos-cli

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go **CLI template** — a command-line financial tracker whose entire interface lives inside a closed, dependency-injected library.

---

## Overview

Agnos-cli is a financial tracker you drive from the terminal. It is built as a structured Go template demonstrating how to build a **CLI** whose behavior is fully decoupled from the process hosting it. The program itself lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself. Everything it can do arrives through an injected `Deps`.

```
adapters/  ──▶  sandbox/  ◀──  cmd/, libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

The CLI is `api.Lib.Sandboxmain` — one field of the library like any other. The installed binary in **`/cmd/main/`** holds no command, no flag, and no output of its own: it just wires an adapter into the library and calls that one field.

- **`/sandbox/`**: The closed library taking a `Deps` and returning an `api.Lib`.
- **`/adapters/`**: Concrete implementations of the `Deps` contract.
- **`/cmd/`** & **`/libraryExamples/`**: Places where an adapter and the library are wired together.

See [SandboxIsolation.md](/docs/SandboxIsolation.md) and [StructContracts.md](/docs/StructContracts.md) for the full mechanic.

---

## Quick Start CLI

**1. Install the CLI globally** — pick your OS, copy the block, paste it in a terminal:

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
      Write-Host "Added $gobin to your PATH (restart the terminal for full effect)"; `
    }; `
    agnos version `
  }
```

> Needs Go ≥ 1.22 installed. The binary is built from `cmd/main`, so `go install` names it `main` — the rename gives it the name you actually type. The script also adds Go's binary directory to your PATH persistently if it isn't already there. See [InstallCli.md](/docs/InstallCli.md) for troubleshooting.

**2. Track your first budget:**

```bash
agnos category add groceries
agnos category add salary

agnos received salary "august paycheck" 2500.00
agnos spend groceries "weekly shopping" 84.50

agnos transactions
agnos balance            # 2415.50
```

Every command, flag, and exit code is listed in [Cli.md](/docs/Cli.md). Records live in `.agnos` in your home directory, or wherever `AGNOS_DATA` points.

---

## Quick Start Library

**1. Install the module:**
```bash
go get github.com/MateusMoutinhoOrg/Agnos-Cli@v0.0.3
```

**2. Create a `main.go` file:**
```go
package main

import (
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
    // 1. Create deps via an adapter (the "opinionated" part:
    //    real clock + standard output + a schema database on disk)
    deps := agnosadapter.New("trackerdata")

    // 2. Inject deps into the pure library — a financial tracker
    l := agnoslib.New(deps)

    // 3. Use the library — it never knows which adapter is behind the scenes.
    //    Amounts are in the smallest currency unit, so 8450 is 84.50.
    l.AddCategory("groceries")
    l.AddSpend("groceries", "weekly shopping", 8450)

    l.AddCategory("salary")
    l.AddReceived("salary", "august paycheck", 250000)

    println(l.Balance()) // 241550
}
```

**3. Run:**
```bash
go run main.go
```

---

> [!IMPORTANT]
> **Must Read before contributing.** The following documents are **required reading** for every developer. Do not open a pull request or make changes without first reading them:
>
> | Document | Why it's required |
> |----------|-------------------|
> | [Rules](/docs/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

## CLI Usage

Installing the `agnos` binary, driving it from a terminal, and adding commands to it.

| Doc | Description |
| --- | --- |
| [/docs/InstallCli.md](/docs/InstallCli.md) | Install the CLI globally, or build and run it from a checkout |
| [/docs/UseCli.md](/docs/UseCli.md) | Create categories, record transactions, and read balances from the terminal |
| [/docs/AddCliCommand.md](/docs/AddCliCommand.md) | Add a command or a flag to the interface behind api.Lib.Sandboxmain |
| [/docs/Cli.md](/docs/Cli.md) | Every command, flag, amount format, and exit code of the interface |
| [/docs/api.Sandboxmain.md](/docs/api.Sandboxmain.md) | The one library field the whole command-line interface lives behind |
| [/docs/SandboxIsolation.md](/docs/SandboxIsolation.md) | Why the interface lives in a closed sandbox and what it may not import |

---

## CLI Examples

Shell scripts driving the built binary the way a user would — each builds the CLI itself and runs against a budget of its own.

| Doc | Description |
| --- | --- |
| [/docs/RunCliExample.md](/docs/RunCliExample.md) | Run the shell scripts in cliExamples/ and read their transcripts |
| [/docs/AddCliExample.md](/docs/AddCliExample.md) | Write a cliExamples/ script and register it in the README |

### Available CLI Examples

| Sample | Description |
|----------|-------------|
| [ManageCategories.sh](/cliExamples/ManageCategories.sh) | Set up a budget: create the categories, list them, and drop one |
| [TrackTransactions.sh](/cliExamples/TrackTransactions.sh) | Track money: record spend and received, list them, read balances |
| [ScriptTheCli.sh](/cliExamples/ScriptTheCli.sh) | Script the CLI: quiet output, exit codes, and piping listings into text tools |

---

## Library Usage

Consuming the same behavior from Go code: install the module, track your money, and understand what the API offers.

| Doc | Description |
| --- | --- |
| [/docs/LibInitialization.md](/docs/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program |
| [/docs/ManageCategories.md](/docs/ManageCategories.md) | Create the categories transactions are tracked under, list them, and remove one |
| [/docs/TrackTransactions.md](/docs/TrackTransactions.md) | Record spend and received transactions, list them, and read a balance |
| [/docs/PublicApi.md](/docs/PublicApi.md) | Index of all public structs, functions, and fields with detail links |
| [/docs/Adapters.md](/docs/Adapters.md) | Lists every shipped adapter and when to use each one |
| [/docs/DepsMechanic.md](/docs/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups |
| [/docs/StructContracts.md](/docs/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them |

---

## Library Examples

Runnable Go programs under `libraryExamples/`, wiring an adapter into the lib.

| Doc | Description |
| --- | --- |
| [/docs/RunSample.md](/docs/RunSample.md) | Browse and run the executable samples in the libraryExamples/ directory |
| [/docs/AddSample.md](/docs/AddSample.md) | Create a runnable sample in libraryExamples/ and register it in the README |

### Available Library Examples

| Sample | Description |
|----------|-------------|
| [AddCategorySample](/libraryExamples/AddCategorySample/AddCategorySample.go) | Create the tracker's categories on disk and list them back |
| [TrackSpendSample](/libraryExamples/TrackSpendSample/TrackSpendSample.go) | Record spend and received transactions and read each category's balance |
| [ListTransactionsSample](/libraryExamples/ListTransactionsSample/ListTransactionsSample.go) | List every transaction across categories, remove one, and total the rest |
| [MainCallSample](/libraryExamples/MainCallSample/MainCallSample.go) | Run the whole CLI from a Go program by calling Sandboxmain and exiting with it |

---

## Sandbox Management

Adding to the library inside the closed sandbox, and exposing what you add.

| Doc | Description |
| --- | --- |
| [/docs/AddLibFunction.md](/docs/AddLibFunction.md) | Declare a function field on api.Lib and write the factory that fills it |
| [/docs/AddLibObject.md](/docs/AddLibObject.md) | Add an object created by the lib, with its deps propagated by its New constructor |
| [/docs/AddCliCommand.md](/docs/AddCliCommand.md) | Add a command or a flag to the interface behind api.Lib.Sandboxmain |
| [/docs/ExposePublicApi.md](/docs/ExposePublicApi.md) | Publish a lib function, object, or field in the public API index |
| [/docs/PublicApi.md](/docs/PublicApi.md) | Index of all public structs, functions, and fields with detail links |
| [/docs/SandboxIsolation.md](/docs/SandboxIsolation.md) | Why the library lives in a closed sandbox and what it may not import |

---

## Dependency Management

Working with the `Deps` contract and the adapters that satisfy it.

| Doc | Description |
| --- | --- |
| [/docs/AddDependency.md](/docs/AddDependency.md) | Add a field to the Deps contract and fill it in every adapter |
| [/docs/AddAdapter.md](/docs/AddAdapter.md) | Create a new opinionated implementation of the Deps contract |
| [/docs/Adapters.md](/docs/Adapters.md) | Lists every shipped adapter and when to use each one |
| [/docs/DepsMechanic.md](/docs/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups |
| [/docs/StructContracts.md](/docs/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them |

---

## Documentation Management

Maintaining the docs themselves: creating, renaming, and deleting `.md` files.

| Doc | Description |
| --- | --- |
| [/docs/AddDocument.md](/docs/AddDocument.md) | Create or update a .md file and register it in README and Structure |
| [/docs/RenameDocument.md](/docs/RenameDocument.md) | Rename or move a .md file without leaving broken references behind |
| [/docs/DeleteDocument.md](/docs/DeleteDocument.md) | Remove a .md file and clear every reference pointing to it |
| [/docs/Specs.md](/docs/Specs.md) | Lists every specification and the files each one governs |

---

## Template Adaptation

Turning the template into a CLI of your own.

| Doc | Description |
| --- | --- |
| [/docs/ForkTemplate.md](/docs/ForkTemplate.md) | Use this repo as a GitHub template to start a new DI library |
| [/docs/AdaptExistingLib.md](/docs/AdaptExistingLib.md) | Convert a pre-existing library to this DI structure |
| [/docs/RenameModule.md](/docs/RenameModule.md) | Rename the Go module path and update all internal imports |
| [/docs/TemplateFileActions.md](/docs/TemplateFileActions.md) | The action each template file takes when adapting: copy, create, rewrite, or delete |

---

## Project Rules & Structure

The binding conventions every change to this repo must follow.

| Doc | Description |
| --- | --- |
| [/docs/RULES.md](/docs/RULES.md) | The binding contribution rules and their required companion updates |
| [/docs/Structure.md](/docs/Structure.md) | The project's directory layout and the purpose of each component |
| [/docs/Specs.md](/docs/Specs.md) | Lists every specification and the files each one governs |
| [/docs/SandboxIsolation.md](/docs/SandboxIsolation.md) | Why the library lives in a closed sandbox and what it may not import |

---

## License

This project is licensed under the [Unlicense](./LICENSE).
