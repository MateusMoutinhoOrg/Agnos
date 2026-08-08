# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go CLI template demonstrating **Dependency Injection** — the whole command-line interface lives inside a closed library.

---

## Overview

Agnos is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. It uses a **Dependency Injection** pattern in which:

- **`/sandbox/contracts/`** defines the `Deps` contract every adapter must fill and the `api` structs the library hands back.
- **`/adapters/`** contains opinionated, concrete implementations of the `Deps` contract.
- **`/sandbox/internal/`** contains the pure library logic as factories filling the `api` contract structs — it never imports concrete implementations.
- **`/sandbox/`** is the entry point: it takes a `Deps` and returns an `api.Lib`, whose `Sandboxmain` field is the command-line interface itself.
- **`/cmd/main/`** is the installed binary: it wires an adapter into the lib, calls `Sandboxmain`, and exits with its return.

This design ensures the interface and the library behind it remain portable, testable, and easy to extend without modifying their core.

---

## Quick Start CLI

**1. Install the CLI globally** — copy, paste, run:

```bash
go install github.com/MateusMoutinhoOrg/Agnos-Cli/cmd/main@latest && \
  mv "$(go env GOPATH)/bin/main" "$(go env GOPATH)/bin/examplecli" && \
  examplecli version
```

**2. Run a first command:**

```bash
examplecli do-the-thing --with-option
```

---

## Quick Start Library

**1. Install the library:**
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
    // 1. Create deps via an adapter (the "opinionated" part)
    deps := agnosadapter.New("data.json")

    // 2. Inject deps into the pure library
    l := agnoslib.New(deps)

    // 3. Use the library — it never knows which adapter is behind the scenes
    obj := l.NewExampleObject(1, "2")
    println(obj.ExampleObjectMethod())
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

Installing the binary, driving it from a terminal, and adding commands to it.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/InstallCli.md](/docs/InstallCli.md) | Install the CLI globally, or build and run it from a checkout | Tutorial |
| [/docs/UseCli.md](/docs/UseCli.md) | Run the interface's commands from a terminal, start to finish | Tutorial |
| [/docs/Cli.md](/docs/Cli.md) | Every command, flag, and exit code of the interface | Reference |

---

## CLI Examples

Shell scripts driving the built binary the way a user would.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/HandleCliExamples.md](/docs/HandleCliExamples.md) | Run the shell scripts in examples/cliExamples/ and read their transcripts | Tutorial |
| [/docs/HandleCliExamples.md](/docs/HandleCliExamples.md) | Write a examples/cliExamples/ script and register it in the README | Tutorial |

### Available CLI Examples

| Sample | Description |
|----------|-------------|
| [ManageCategories.sh](/examples/cliExamples/ManageCategories.sh) | What this script demonstrates against the built CLI |

---

## Library Usage

For consuming the lib as a user: install it, run a first program, and understand its core mechanic.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/LibInitialization.md](/docs/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program | Tutorial |
| [/docs/HandleDependencies.md](/docs/HandleDependencies.md) | How the dependency-injection mechanism works, including custom setups | Explanation |

---

## Library Examples

Runnable Go programs under `examples/libraryExamples/`, wiring an adapter into the lib.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/HandleSamples.md](/docs/HandleSamples.md) | Browse and run the executable samples in the examples/libraryExamples/ directory | Tutorial |
| [/docs/HandleSamples.md](/docs/HandleSamples.md) | Create a runnable sample in examples/libraryExamples/ and register it in the README | Tutorial |

### Available Library Examples

| Sample | Description |
|----------|-------------|
| [ExampleSample](/examples/libraryExamples/ExampleSample/ExampleSample.go) | How to use the library |

---

## Sandbox Management

Adding to the library inside the closed sandbox, and exposing what you add.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/HandleLibElements.md](/docs/HandleLibElements.md) | Declare a function field on api.Lib and write the factory that fills it | Tutorial |
| [/docs/HandleCliCommands.md](/docs/HandleCliCommands.md) | Add a command or a flag to the interface behind api.Lib.Sandboxmain | Tutorial |

---

## Dependency Management

Working with the `Deps` contract and the adapters that satisfy it.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/HandleDependencies.md](/docs/HandleDependencies.md) | Add a field to the Deps contract and implement it in every adapter | Tutorial |
| [/docs/HandleAdapters.md](/docs/HandleAdapters.md) | Create a new opinionated implementation of the Deps contract | Tutorial |
| [/docs/HandleDependencies.md](/docs/HandleDependencies.md) | How the dependency-injection mechanism works, including custom setups | Explanation |

---

## Project Rules & Structure

The binding conventions every change to this repo must follow.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/RULES.md](/docs/RULES.md) | The binding contribution rules and their required companion updates | Reference |
| [/docs/Structure.md](/docs/Structure.md) | The project's directory layout and the purpose of each component | Reference |
| [/docs/Specs.md](/docs/Specs.md) | Lists every specification and the files each one governs | Reference |

---

## License

This project is licensed under the [Unlicense](./LICENSE).
