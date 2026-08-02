# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

An OS-independent Go library template demonstrating **Dependency Injection** with a clean separation between pure library logic and adapter implementations.

---

## Overview

Agnos is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. The library itself lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package. Everything it can do arrives through an injected `Deps`.

```
adapters/  ──▶  sandbox/  ◀──  examples/
(reaches the OS)  (closed)     (wires the two together)
```

- **`/sandbox/`** is the closed library and its single entry point: it takes a `Deps` and returns an `api.Lib`.
  - **`/sandbox/contracts/`** holds the public types everything is wired through — the `Deps` contract every adapter must fill, and the `api` structs the library hands back. Contracts are **structs of function fields**, never interfaces. This is the only part of the sandbox the outside world imports.
  - **`/sandbox/internal/`** holds the pure library logic as **factories** — functions that take a pointer to an `api` struct and fill its function fields with closures reading that struct's `Deps`. It declares no types and is unreachable from outside `sandbox/`.
- **`/adapters/`** sits outside the sandbox and holds opinionated, concrete implementations of the `Deps` contract, filled by the **same factories** the sandbox uses — the carrier is the adapter struct rather than an `api` struct. This is the only place OS-bound and third-party code is allowed.
- **`/examples/`** sits outside the sandbox too, and is the only place an adapter and the library are wired together.

This design ensures the library remains portable, testable, and easy to extend without modifying its core. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) for the full mechanic and [StructContracts.md](/docs/Explanations/StructContracts.md) for why the contracts are structs and how factories fill them.

---

## Quick Start

**1. Install the library:**
```bash
go get github.com/MateusMoutinhoOrg/Agnos@v0.0.7
```

**2. Create a `main.go` file:**
```go
package main

import (
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
    // 1. Create deps via an adapter (the "opinionated" part:
    //    JSON-file store + real clock)
    deps := agnosadapter.New("kvcache.json")

    // 2. Inject deps into the pure library — a key/value cache with TTL
    l := agnoslib.New(deps)

    // 3. Use the library — it never knows which adapter is behind the scenes
    l.Set("greeting", "hello world", 60)
    if entry, ok := l.Get("greeting"); ok {
        println(entry.Value)
    }
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

For consuming the lib as a user: install it, use the cache, and understand what the API offers.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/LibInitialization.md](/docs/Tutorials/LibInitialization.md) | Install the lib, create deps via an adapter, and run a first program | Tutorial |
| [/docs/Tutorials/CacheValue.md](/docs/Tutorials/CacheValue.md) | Store a value with Set, read it back with Get, and handle a cache miss | Tutorial |
| [/docs/Tutorials/InspectEntryExpiration.md](/docs/Tutorials/InspectEntryExpiration.md) | Read an entry's expiry with ExpiresAt and check it live with IsExpired | Tutorial |
| [/docs/Tutorials/RunSample.md](/docs/Tutorials/RunSample.md) | Browse and run the executable samples in the examples/ directory | Tutorial |
| [/docs/References/PublicApi.md](/docs/References/PublicApi.md) | Index of all public structs, functions, and fields with detail links | Reference |
| [/docs/References/Adapters.md](/docs/References/Adapters.md) | Lists every shipped adapter and when to use each one | Reference |
| [/docs/Explanations/DepsMechanic.md](/docs/Explanations/DepsMechanic.md) | How the dependency-injection mechanism works, including custom setups | Explanation |
| [/docs/Explanations/SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md) | Why the library lives in a closed sandbox and what it may not import | Explanation |
| [/docs/Explanations/StructContracts.md](/docs/Explanations/StructContracts.md) | Why every contract is a struct of function fields, and how factories fill them | Explanation |

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
| [KvCacheSample](/examples/KvCacheSample/KvCacheSample.go) | Store and read a value from the TTL key/value cache |
| [AuthTokenSample](/examples/AuthTokenSample/AuthTokenSample.go) | Cache a short-lived auth token, re-authenticating only after it expires |
| [WebPageSample](/examples/WebPageSample/WebPageSample.go) | Cache a mutable web page for 60s, hitting the network only on a miss |

---

## Extending the Library

Adding new lib functionality and exposing it in the public API.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/Tutorials/AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) | Declare a function field on api.Lib and write the factory that fills it | Tutorial |
| [/docs/Tutorials/AddLibObject.md](/docs/Tutorials/AddLibObject.md) | Add an object created by the lib, with its deps propagated by its New constructor | Tutorial |
| [/docs/Tutorials/ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md) | Publish a lib function, object, or field in the public API index | Tutorial |
| [/docs/References/PublicApi.md](/docs/References/PublicApi.md) | Index of all public structs, functions, and fields with detail links | Reference |

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

Turning the template into a real library of your own.

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

This project is licensed under the [MIT License](./LICENSE).
