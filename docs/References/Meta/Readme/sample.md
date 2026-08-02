# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

An OS-independent Go library template demonstrating **Dependency Injection** with a clean separation between pure library logic and adapter implementations.

---

## Overview

Agnos is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. It uses a **Dependency Injection** pattern in which:

- **`/sandbox/contracts/`** defines the `Deps` contract every adapter must fill and the `api` structs the library hands back.
- **`/adapters/`** contains opinionated, concrete implementations of the `Deps` contract.
- **`/sandbox/internal/`** contains the pure library logic as factories filling the `api` contract structs — it never imports concrete implementations.
- **`/sandbox/`** is the entry point: it takes a `Deps` and returns an `api.Lib`.

This design ensures the library remains portable, testable, and easy to extend without modifying its core.

---

## Quick Start

**1. Install the library:**
```bash
go get github.com/MateusMoutinhoOrg/Agnos
```

**2. Create a `main.go` file:**
```go
package main

import (
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
    // 1. Create deps via an adapter (the "opinionated" part)
    deps := agnosadapter.New(3)

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
> | [Rules](/docs/References/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/References/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/References/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

## Library Usage

For consuming the lib as a user: install it, run a first program, and understand its core mechanic.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/LibInitialization.md](/docs/Tutorials/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program | Tutorial |
| [/docs/Tutorials/RunSample.md](/docs/Tutorials/RunSample.md) | Browse and run the executable samples in the examples/ directory | Tutorial |
| [/docs/Explanations/DepsMechanic.md](/docs/Explanations/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups | Explanation |

---

## Samples

Creating and running the example programs under `examples/`.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/RunSample.md](/docs/Tutorials/RunSample.md) | Browse and run the executable samples in the examples/ directory | Tutorial |
| [/docs/Tutorials/AddSample.md](/docs/Tutorials/AddSample.md) | Create a runnable sample in examples/ and register it in the README | Tutorial |

### Available Samples

| Sample | Description |
|----------|-------------|
| [ExampleSample](/examples/ExampleSample/ExampleSample.go) | How to use the library |

---

## Dependency Management

Working with the `Deps` contract and the adapters that satisfy it.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/AddDependency.md](/docs/Tutorials/AddDependency.md) | Add a field to the Deps contract and implement it in every adapter | Tutorial |
| [/docs/Tutorials/AddAdapter.md](/docs/Tutorials/AddAdapter.md) | Create a new opinionated implementation of the Deps contract | Tutorial |
| [/docs/Explanations/DepsMechanic.md](/docs/Explanations/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups | Explanation |

---

## Project Rules & Structure

The binding conventions every change to this repo must follow.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/References/RULES.md](/docs/References/RULES.md) | The binding contribution rules and their required companion updates | Reference |
| [/docs/References/Structure.md](/docs/References/Structure.md) | The project's directory layout and the purpose of each component | Reference |
| [/docs/References/Specs.md](/docs/References/Specs.md) | Lists every specification and the files each one governs | Reference |

---

## License

This project is licensed under the [MIT License](./LICENSE).
