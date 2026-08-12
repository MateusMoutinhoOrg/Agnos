# Agnos-Cli

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go CLI template demonstrating **Dependency Injection** — the whole command-line interface lives inside a closed library.

---

## Overview

Agnos-Cli is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. It uses a **Dependency Injection** pattern in which:

- **`/sandbox/contracts/`** defines the `Deps` contract every adapter must fill and the `api` structs the library hands back.
- **`/adapters/`** contains opinionated, concrete implementations of the `Deps` contract.
- **`/sandbox/internal/`** contains the pure library logic as factories filling the `api` contract structs — it never imports concrete implementations.
- **`/sandbox/`** is the entry point: it takes a `Deps` and returns an `api.Lib`, whose `Sandboxmain` field is the command-line interface itself.
- **`/cmd/main/`** is the installed binary: it wires an adapter into the lib, calls `Sandboxmain`, and exits with its return.

This design ensures the interface and the library behind it remain portable, testable, and easy to extend without modifying their core.

---

## Doc Index

Documentation is split into four themes, one index page each under `docs/Index/`, listing that theme's **Tutorials** — step-by-step workflows — and its **References** — explanations and lookups. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [CLI Usage](/docs/Index/CliUsage.md) | For end users: installing the binary, running it, and every command it takes. |
| [Library Usage](/docs/Index/LibUsage.md) | For library consumers: installing the module, creating deps, and calling the Go API. |
| [Development](/docs/Index/Development.md) | For contributors: the rules, the mechanics, the workflows, and the specifications. |
| [Templating](/docs/Index/Templating.md) | For template users: forking, renaming, and adapting this structure into a new CLI. |

New here? [CLI Usage → InstallCli.md](/docs/Tutorials/InstallCli.md) installs the binary; [Library Usage → LibInitialization.md](/docs/Tutorials/LibInitialization.md) initializes the library.

---

## License

This project is licensed under the [Unlicense](./LICENSE).
