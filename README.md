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
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

The CLI is `api.Lib.Sandboxmain` — one field of the library like any other. The installed binary in **`/cmd/main/`** holds no command, no flag, and no output of its own: it just wires an adapter into the library and calls that one field.

- **`/sandbox/`**: The closed library taking a `Deps` and returning an `api.Lib`.
- **`/adapters/`**: Concrete implementations of the `Deps` contract.
- **`/cmd/`** & **`/examples/libraryExamples/`**: Places where an adapter and the library are wired together.

See [SandboxIsolation.md](/docs/SandboxIsolation.md) and [StructContracts.md](/docs/StructContracts.md) for the full mechanic.

---

## Navigation

| Section | Description |
|---|---|
| [CLI Usage](#cli-usage) | For end users: installing the binary, driving it from a terminal, and usage documentation. |
| [API Usage](#api-usage) | For library consumers: installing the module, creating deps, and using the Go API. |
| [Development](#development) | For contributors: required reading, adding functions, commands, managing dependencies and documentation. |
| [Template](#template) | For template users: forking, renaming, and adapting this structure to a new CLI. |

---

## CLI Usage

This section is for users who want to run the CLI from their terminal.

### Quick Start CLI

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

### Documentation

| Doc | Description |
| --- | --- |
| [/docs/InstallCli.md](/docs/InstallCli.md) | Install the CLI globally, or build and run it from a checkout |
| [/docs/UseCli.md](/docs/UseCli.md) | Create categories, record transactions, and read balances from the terminal |
| [/docs/Cli.md](/docs/Cli.md) | Every command, flag, amount format, and exit code of the interface |
| [/docs/RunCliSample.md](/docs/RunCliSample.md) | How to run CLI examples from the source |
| [/docs/SamplesList.md](/docs/SamplesList.md) | A list of all examples that can be executed in CLI mode |

---

## API Usage

This section is for developers consuming the `Agnos-Cli` behavior as a Go library.

### Quick Start Library

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

### Documentation

| Doc | Description |
| --- | --- |
| [/docs/LibInitialization.md](/docs/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program |
| [/docs/Adapters.md](/docs/Adapters.md) | Every shipped adapter you can inject, and when to use each one |
| [/docs/HandleDependencies.md](/docs/HandleDependencies.md) | The `Deps` contract you inject: what each field is for, and how to swap or patch one |
| [/docs/ManageCategories.md](/docs/ManageCategories.md) | Create the categories transactions are tracked under, list them, and remove one |
| [/docs/TrackTransactions.md](/docs/TrackTransactions.md) | Record spend and received transactions, list them, and read a balance |
| [/docs/PublicApi.md](/docs/PublicApi.md) | Index of all public structs, functions, and fields with detail links |
| [/docs/RunApiSample.md](/docs/RunApiSample.md) | How to run API examples from the source |
| [/docs/ApiSamplesList.md](/docs/ApiSamplesList.md) | A list of all examples that can be executed in API mode |

---

## Development

This section is for contributors adding functionality, commands, or fixing bugs in the project.

> [!IMPORTANT]
> **Must Read before contributing.** The following documents are **required reading** for every developer. Do not open a pull request or make changes without first reading them:
>
> | Document | Why it's required |
> |----------|-------------------|
> | [Rules](/docs/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

### Learning Path

The documentation below is a reading path. The stages are ordered from what **every** change touches to what almost none does: read stage 1 once to absorb the mechanics, then jump to the stage matching your change.

#### 1. Mechanics — read once, before any change

Every contribution runs into these three mechanics; every tutorial below assumes them.

| Doc | Description |
| --- | --- |
| [/docs/SandboxIsolation.md](/docs/SandboxIsolation.md) | The sandbox wall: what `sandbox/` may not import, and why every effect is a `Deps` field |
| [/docs/StructContracts.md](/docs/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them |
| [/docs/HandleDependencies.md](/docs/HandleDependencies.md) | How injected deps travel the object graph, and how to add a field to the `Deps` contract |

#### 2. Everyday changes — most contributions live here

The typical feature is a library function plus the CLI command that calls it, published in the API index.

| Doc | Description |
| --- | --- |
| [/docs/HandleLibElements.md](/docs/HandleLibElements.md) | Add a function or an object to the library: declare, write the factory, register it, publish it |
| [/docs/HandleCliCommands.md](/docs/HandleCliCommands.md) | Add a command or a flag to the interface behind `api.Lib.Sandboxmain` |

#### 3. Companion updates — after most changes

New behavior is demonstrated by a sample and reflected in the docs, in the same commit.

| Doc | Description |
| --- | --- |
| [/docs/HandleCliExamples.md](/docs/HandleCliExamples.md) | Create and run shell scripts in `examples/cliExamples/` driving the built CLI |
| [/docs/HandleSamples.md](/docs/HandleSamples.md) | Create and run executable Go samples in `examples/libraryExamples/` |
| [/docs/HandleDocuments.md](/docs/HandleDocuments.md) | Create, rename, move, or delete a `.md` file without leaving broken references |

#### 4. Infrastructure — rare changes

The `Deps` implementations change far less often than the library they feed.

| Doc | Description |
| --- | --- |
| [/docs/Adapters.md](/docs/Adapters.md) | Every shipped adapter and when to use each one |
| [/docs/HandleAdapters.md](/docs/HandleAdapters.md) | Create a new opinionated implementation of the `Deps` contract |

---

## Template

This section is for those using the `Agnos-Cli` project as a template to bootstrap their own CLI. Pick **one** of the two workflows and follow it end to end — both are phased step lists, and both take each file's fate from the same per-file action table.

| Doc | Description |
| --- | --- |
| [/docs/ForkTemplate.md](/docs/ForkTemplate.md) | **Start here for a new library**: use this repo as a GitHub template |
| [/docs/AdaptExistingLib.md](/docs/AdaptExistingLib.md) | **Start here for an existing library**: convert it to this DI structure |
| [/docs/TemplateFileActions.md](/docs/TemplateFileActions.md) | The per-file action table both workflows follow: copy, create, rewrite, or delete |
| [/docs/RenameModule.md](/docs/RenameModule.md) | Rename the Go module path and update all internal imports — the first step of both workflows |

---

## License

This project is licensed under the [Unlicense](./LICENSE).
