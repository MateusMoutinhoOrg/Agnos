# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

A Go CLI that **scaffolds and regenerates other Go CLIs** built around a closed, dependency-injected sandbox.

---

## Overview

Agnos is a factory: `agnos start` writes a project skeleton, `agnos build` regenerates its generated files from embedded templates, and a handful of commands declare the project's whole command surface without editing a file by hand.

- **`/sandbox/`** is a closed sandbox: every effect arrives through an injected `Deps`.
- **`/adapters/`** is the only place OS-bound and third-party code lives.
- **`/assets/`** holds the Go templates every generated file is rendered from.
- **`/cmd/main/`** wires an adapter into the sandbox and exits with what the CLI returns.

---

## Doc Index

Documentation is split into four themes, one index page each under `docs/Index/`. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [CLI Usage](/docs/Index/CliUsage.md) | For end users: installing `agnos`, scaffolding a project, and every command it takes. |
| [Generated Project](/docs/Index/GeneratedProject.md) | For people working inside a project `agnos` wrote: its schema, contracts and handlers. |
| [Library Usage](/docs/Index/LibUsage.md) | For Go callers: running the same actions from code, and the public API. |
| [Development](/docs/Index/Development.md) | For contributors: the mechanics, the workflows, and the specifications. |

New here? [CLI Usage → InstallCli.md](/docs/Tutorials/InstallCli.md) installs the binary; [ScaffoldProject.md](/docs/Tutorials/ScaffoldProject.md) builds a first project.

---

## License

This project is licensed under the [Unlicense](./LICENSE).
