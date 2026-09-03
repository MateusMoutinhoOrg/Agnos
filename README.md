# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Unlicense-green)](LICENSE)

A Go CLI that **scaffolds and regenerates other Go CLIs** — each one a closed, dependency-injected sandbox behind a command-line interface generated from declarations.

> [!WARNING]
> **This project is still under active development.** Its patterns, commands and generated
> output change frequently and without notice. Using it is **not recommended** unless you
> are an experienced developer or team comfortable reading the source, tracking breaking
> changes, and fixing generated code by hand.

---

## Overview

Agnos (`agnos`) is a **factory**. `agnos start` writes a project skeleton; `agnos build` regenerates every generated file of that project from templates embedded in the binary; and a handful of commands — `add-command`, `add-flag`, `add-arg`, `dep-install` — declare the project's whole command surface and capability set without a file being edited by hand. Only two kinds of file are ever hand-written in a generated project: a command's `handler.go`, and a contract-plus-adapter pair for a capability of its own.

The core of every project, this one included, lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself. Everything it can do arrives through an injected `Deps`, a struct of sub-contract structs reconstructed from a directory listing on every build.

```
adapters/  ──▶  sandbox/  ◀──  cmd/
(reaches the OS)  (closed)     (wires the two together)
```

- **`/sandbox/`**: the closed sandbox — the actions, the generated command dispatch, and the contracts everything is wired through.
- **`/adapters/`**: the only place OS-bound and third-party code lives — one isolated lib per contract, and the assembly wiring them together.
- **`/assets/`**: the Go templates every generated file is rendered from, and the installable deps.
- **`/cmd/main/`**: wires an adapter into the sandbox and exits with what the CLI returns. Holds no logic.

Agnos is one of the projects it builds: its own generated files are rendered in place by `agnos build`, and the result compiles. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md) and [BuildPipeline.md](/docs/References/BuildPipeline.md) for the full mechanic.

---

## Doc Index

Documentation is split into four themes, one index page each under `docs/Index/`, listing that theme's **Tutorials** — step-by-step workflows — and its **References** — explanations and lookups. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [CLI Usage](/docs/Index/CliUsage.md) | For end users: installing `agnos`, scaffolding a project, and every command it takes. |
| [Generated Project](/docs/Index/GeneratedProject.md) | For people working inside a project `agnos` wrote: its files, contracts and handlers. |
| [Library Usage](/docs/Index/LibUsage.md) | For Go callers: running the same actions from code, and the public API. |
| [Development](/docs/Index/Development.md) | For contributors: the rules, the mechanics, the workflows, and the specifications. |

New here? [CLI Usage → InstallCli.md](/docs/Tutorials/InstallCli.md) installs the binary; [ScaffoldProject.md](/docs/Tutorials/ScaffoldProject.md) takes an empty directory to a running CLI.

---

## License

This project is licensed under the [Unlicense](./LICENSE).
