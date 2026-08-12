# Agnos-Cli

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos-Cli.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos-Cli)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos-Cli)](https://github.com/MateusMoutinhoOrg/Agnos-Cli/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

An OS-independent Go **CLI template** — a command-line financial tracker whose entire interface lives inside a closed, dependency-injected library.

---

## Overview

Agnos-Cli is a financial tracker you drive from the terminal. It is built as a structured Go template demonstrating how to build a **CLI** whose behavior is fully decoupled from the process hosting it. The program itself lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself. Everything it can do arrives through an injected `Deps`.

```
adapters/  ──▶  sandbox/  ◀──  cmd/, examples/libraryExamples/
(reaches the OS)  (closed)     (wire the two together)
```

The CLI is `api.Lib.Sandboxmain` — one field of the library like any other. The installed binary in **`/cmd/main/`** holds no command, no flag, and no output of its own: it just wires an adapter into the library and calls that one field.

- **`/sandbox/`**: The closed library taking a `Deps` and returning an `api.Lib`.
- **`/adapters/`**: Concrete implementations of the `Deps` contract.
- **`/cmd/`** & **`/examples/libraryExamples/`**: Places where an adapter and the library are wired together.
- **`/assets/`**: The text the interface displays, compiled into the binary and reached only through the injected `Deps` — so even the help screen is not written inside the sandbox.

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
