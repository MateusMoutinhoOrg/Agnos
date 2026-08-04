# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go **CLI template** — a command-line financial tracker whose entire interface lives inside a closed, dependency-injected library.

---

## Overview

Agnos is a financial tracker you drive from the terminal: categories holding spend and received transactions, persisted through an injected schema database. It is built as a structured Go template showing how to build a **CLI** whose behavior is fully decoupled from the process hosting it. The program itself lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package. Everything it can do arrives through an injected `Deps`.

```
adapters/  ──▶  sandbox/  ◀──  cmd/, libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

The command-line interface is `api.Lib.Sandboxmain` — one field of the library like any other. It reads the command line through the injected argv parser and prints through the injected `Printf`, so the installed binary in **`/cmd/main/`** holds no command, no flag, and no output of its own: it wires an adapter into the library, calls that one field, and exits with what it returns.

- **`/sandbox/`** is the closed library and its single entry point: it takes a `Deps` and returns an `api.Lib`.
  - **`/sandbox/contracts/`** holds the public types everything is wired through — the `Deps` contract every adapter must fill, and the `api` structs the library hands back. Contracts are **structs of function fields**, never interfaces. This is the only part of the sandbox the outside world imports.
  - **`/sandbox/internal/`** holds the pure library logic as **factories** — functions that take a pointer to an `api` struct and fill its function fields with closures reading that struct's `Deps` — plus `cli/`, the command dispatch behind `Sandboxmain`. It declares no types and is unreachable from outside `sandbox/`.
- **`/adapters/`** sits outside the sandbox and holds opinionated, concrete implementations of the `Deps` contract, filled by the **same factories** the sandbox uses — the carrier is the adapter struct rather than an `api` struct. This is the only place OS-bound and third-party code is allowed.
- **`/cmd/`** and **`/libraryExamples/`** sit outside the sandbox too, and are the only places an adapter and the library are wired together.

Consuming Agnos as a Go library still works and is fully documented — it is simply the background feature. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) for the full mechanic and [StructContracts.md](/docs/Explanations/StructContracts.md) for why the contracts are structs and how factories fill them.

---

## Quick Start CLI

**1. Install the CLI globally** — copy, paste, run:

```bash
go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@latest && \
  mv "$(go env GOPATH)/bin/main" "$(go env GOPATH)/bin/agnos" && \
  agnos version
```

> Needs Go 1.22+ and `$(go env GOPATH)/bin` on your `PATH`. The binary is built from `cmd/main`, so `go install` names it `main` — the `mv` gives it the name you actually type.

**2. Track your first budget:**

```bash
agnos category add groceries
agnos category add salary

agnos received salary "august paycheck" 2500.00
agnos spend groceries "weekly shopping" 84.50

agnos transactions
agnos balance            # 2415.50
```

Every command, flag, and exit code is listed in [Cli.md](/docs/References/Cli.md). Records live in `.agnos` in your home directory, or wherever `AGNOS_DATA` points.

---

## Quick Start Library

**1. Install the module:**
```bash
go get github.com/MateusMoutinhoOrg/Agnos-Cli@latest
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
> | [Rules](/docs/References/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/References/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/References/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

## CLI Usage

Installing the `agnos` binary, driving it from a terminal, and adding commands to it.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/InstallCli.md](/docs/Tutorials/InstallCli.md) | Install the CLI globally, or build and run it from a checkout | Tutorial |
| [/docs/Tutorials/UseCli.md](/docs/Tutorials/UseCli.md) | Create categories, record transactions, and read balances from the terminal | Tutorial |
| [/docs/Tutorials/AddCliCommand.md](/docs/Tutorials/AddCliCommand.md) | Add a command or a flag to the interface behind api.Lib.Sandboxmain | Tutorial |
| [/docs/References/Cli.md](/docs/References/Cli.md) | Every command, flag, amount format, and exit code of the interface | Reference |
| [/docs/References/PublicApi/api.Sandboxmain.md](/docs/References/PublicApi/api.Sandboxmain.md) | The one library field the whole command-line interface lives behind | Reference |
| [/docs/Explanations/SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) | Why the interface lives in a closed sandbox and what it may not import | Explanation |

---

## CLI Examples

Shell scripts driving the built binary the way a user would — each builds the CLI itself and runs against a budget of its own.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/RunCliExample.md](/docs/Tutorials/RunCliExample.md) | Run the shell scripts in cliExamples/ and read their transcripts | Tutorial |
| [/docs/Tutorials/AddCliExample.md](/docs/Tutorials/AddCliExample.md) | Write a cliExamples/ script and register it in the README | Tutorial |

### Available CLI Examples

| Sample | Description |
|----------|-------------|
| [example1.sh](/cliExamples/example1.sh) | Set up a budget: create the categories, list them, and drop one |
| [example2.sh](/cliExamples/example2.sh) | Track money: record spend and received, list them, read balances |
| [example3.sh](/cliExamples/example3.sh) | Script the CLI: quiet output, exit codes, and piping listings into text tools |

---

## Library Usage

Consuming the same behavior from Go code: install the module, track your money, and understand what the API offers.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/LibInitialization.md](/docs/Tutorials/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program | Tutorial |
| [/docs/Tutorials/ManageCategories.md](/docs/Tutorials/ManageCategories.md) | Create the categories transactions are tracked under, list them, and remove one | Tutorial |
| [/docs/Tutorials/TrackTransactions.md](/docs/Tutorials/TrackTransactions.md) | Record spend and received transactions, list them, and read a balance | Tutorial |
| [/docs/References/PublicApi.md](/docs/References/PublicApi.md) | Index of all public structs, functions, and fields with detail links | Reference |
| [/docs/References/Adapters.md](/docs/References/Adapters.md) | Lists every shipped adapter and when to use each one | Reference |
| [/docs/Explanations/DepsMechanic.md](/docs/Explanations/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups | Explanation |
| [/docs/Explanations/StructContracts.md](/docs/Explanations/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them | Explanation |

---

## Library Examples

Runnable Go programs under `libraryExamples/`, wiring an adapter into the lib.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/RunSample.md](/docs/Tutorials/RunSample.md) | Browse and run the executable samples in the libraryExamples/ directory | Tutorial |
| [/docs/Tutorials/AddSample.md](/docs/Tutorials/AddSample.md) | Create a runnable sample in libraryExamples/ and register it in the README | Tutorial |

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

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) | Declare a function field on api.Lib and write the factory that fills it | Tutorial |
| [/docs/Tutorials/AddLibObject.md](/docs/Tutorials/AddLibObject.md) | Add an object created by the lib, with its deps propagated by its New constructor | Tutorial |
| [/docs/Tutorials/AddCliCommand.md](/docs/Tutorials/AddCliCommand.md) | Add a command or a flag to the interface behind api.Lib.Sandboxmain | Tutorial |
| [/docs/Tutorials/ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md) | Publish a lib function, object, or field in the public API index | Tutorial |
| [/docs/References/PublicApi.md](/docs/References/PublicApi.md) | Index of all public structs, functions, and fields with detail links | Reference |
| [/docs/Explanations/SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) | Why the library lives in a closed sandbox and what it may not import | Explanation |

---

## Dependency Management

Working with the `Deps` contract and the adapters that satisfy it.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/AddDependency.md](/docs/Tutorials/AddDependency.md) | Add a field to the Deps contract and fill it in every adapter | Tutorial |
| [/docs/Tutorials/AddAdapter.md](/docs/Tutorials/AddAdapter.md) | Create a new opinionated implementation of the Deps contract | Tutorial |
| [/docs/References/Adapters.md](/docs/References/Adapters.md) | Lists every shipped adapter and when to use each one | Reference |
| [/docs/Explanations/DepsMechanic.md](/docs/Explanations/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups | Explanation |
| [/docs/Explanations/StructContracts.md](/docs/Explanations/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them | Explanation |

---

## Documentation Management

Maintaining the docs themselves: creating, renaming, and deleting `.md` files.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/AddDocument.md](/docs/Tutorials/AddDocument.md) | Create or update a .md file and register it in README and Structure | Tutorial |
| [/docs/Tutorials/RenameDocument.md](/docs/Tutorials/RenameDocument.md) | Rename or move a .md file without leaving broken references behind | Tutorial |
| [/docs/Tutorials/DeleteDocument.md](/docs/Tutorials/DeleteDocument.md) | Remove a .md file and clear every reference pointing to it | Tutorial |
| [/docs/References/Specs.md](/docs/References/Specs.md) | Lists every specification and the files each one governs | Reference |

---

## Template Adaptation

Turning the template into a CLI of your own.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) | Use this repo as a GitHub template to start a new DI library | Tutorial |
| [/docs/Tutorials/AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) | Convert a pre-existing library to this DI structure | Tutorial |
| [/docs/Tutorials/RenameModule.md](/docs/Tutorials/RenameModule.md) | Rename the Go module path and update all internal imports | Tutorial |
| [/docs/References/TemplateFileActions.md](/docs/References/TemplateFileActions.md) | The action each template file takes when adapting: copy, create, rewrite, or delete | Reference |

---

## Project Rules & Structure

The binding conventions every change to this repo must follow.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/References/RULES.md](/docs/References/RULES.md) | The binding contribution rules and their required companion updates | Reference |
| [/docs/References/Structure.md](/docs/References/Structure.md) | The project's directory layout and the purpose of each component | Reference |
| [/docs/References/Specs.md](/docs/References/Specs.md) | Lists every specification and the files each one governs | Reference |
| [/docs/Explanations/SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) | Why the library lives in a closed sandbox and what it may not import | Explanation |

---

## License

This project is licensed under the [Unlicense](./LICENSE).
