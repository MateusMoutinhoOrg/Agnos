# Agnos-Cli

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go **CLI template** — a command-line financial tracker whose entire interface lives inside a closed, dependency-injected library.

---

## Overview

Agnos-Cli is a **full CLI template** designed to be completely independent of the underlying operating system. It provides a complete harness and architectural foundation for building command-line applications whose behavior is fully decoupled from the hosting environment. Furthermore, the repository is designed to be self-teaching, providing **comprehensive tutorials for every kind of usecase** directly within its own documentation.

The core of the application lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself. Everything it can do arrives through an injected `Deps`.

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

The CLI is `api.Lib.Sandboxmain` — one field of the library like any other. The installed binary in **`/cmd/main/`** holds no command, no flag, and no output of its own: it just wires an adapter into the library and calls that one field.

- **`/sandbox/`**: The closed library taking a `Deps` and returning an `api.Lib`.
- **`/adapters/`**: Concrete implementations of the `Deps` contract.
- **`/cmd/`** & **`/examples/libraryExamples/`**: Places where an adapter and the library are wired together.
- **`/assets/`**: Files compiled into the binary and reached only through the injected `Deps` — templates, long-form text, images. Empty here by design, and wired end to end for a derived library to fill.

The `Deps` contract is wider than this tracker uses. Reading embedded assets, touching the filesystem, and speaking HTTP are declared and filled by the standard adapter even though the demonstration never calls them: they are capabilities a derived library gets for free. See [PublicApi.md](/docs/References/PublicApi.md) for every field.

See [SandboxIsolation.md](/docs/References/SandboxIsolation.md) and [StructContracts.md](/docs/References/StructContracts.md) for the full mechanic.

---

## Doc Index

Documentation is split into four themes, one index page each under `docs/Index/`, listing that theme's **Tutorials** — step-by-step workflows — and its **References** — explanations and lookups. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [CLI Usage](/docs/Index/CliUsage.md) | For end users: installing the binary, driving it from a terminal, and every command it takes. |
| [Library Usage](/docs/Index/LibUsage.md) | For library consumers: installing the module, creating deps, and calling the Go API. |
| [Development](/docs/Index/Development.md) | For contributors: the rules, the mechanics, the per-goal workflows, and the specifications. |
| [Templating](/docs/Index/Templating.md) | For template users: forking, renaming, and adapting this structure into a new CLI. |

New here? [CLI Usage → InstallCli.md](/docs/Tutorials/InstallCli.md) installs the binary; [Library Usage → LibInitialization.md](/docs/Tutorials/LibInitialization.md) initializes the library.

---

## License

This project is licensed under the [Unlicense](./LICENSE).
