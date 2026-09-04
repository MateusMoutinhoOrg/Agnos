# agnos

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

Agnos is one of the projects it builds: its own generated files are rendered in place by `agnos build`, and the result compiles. See [SandboxIsolation](docs/SandboxIsolation/doc.md) and [BuildPipeline](docs/BuildPipeline/doc.md) for the full mechanic.





## Documentation Index

| Name | Description |
| --- | --- |
| [CliUsage](docs/Index/cli-usage.md) | Documentation for people who drive agnos from a terminal - installing the binary, scaffolding a project, declaring its commands, and looking up what each command does |
| [LibUsage](docs/Index/lib-usage.md) | Documentation for developers consuming Agnos as a Go module - wiring an adapter into the sandbox, calling the same actions the CLI exposes from code, and the public API |
| [Development](docs/Index/development.md) | Documentation for contributors changing this repository - the binding rules, the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy |

## License


MIT License

Copyright (c) 2026 Mateus Moutinho

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
