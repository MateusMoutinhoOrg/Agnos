# agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go CLI that **scaffolds and regenerates other Go CLIs** — each one a closed, dependency-injected sandbox behind a command-line interface generated from declarations.

> [!WARNING]
> **This project is still under active development.** Its patterns, commands and generated
> output change frequently and without notice. Using it is **not recommended** unless you
> are an experienced developer or team comfortable reading the source, tracking breaking
> changes, and fixing generated code by hand.

---

## Overview

Agnos (`agnos`) is a **factory**. `agnos start` writes a project skeleton, `agnos build`
re-renders every generated file from templates embedded in the binary, and commands like
`add-command`, `add-flag` and `dep-install` declare the project's whole command surface
without a file being edited by hand. Only two things stay hand-written: a command's
`handler.go`, and any contract-plus-adapter pair of your own.

```
adapters/  ──▶  sandbox/  ◀──  cmd/
(reaches the OS)  (closed)     (wires the two together)
```

- **`/sandbox/`** — the closed core: actions, the generated dispatch, and the contracts
  everything is injected through. Reaches nothing outside itself.
- **`/adapters/`** — the only place OS-bound and third-party code lives.
- **`/assets/`** — the templates every generated file comes from, plus the installable deps.
- **`/cmd/main/`** — wires an adapter into the sandbox. Holds no logic.

Agnos builds itself: `agnos build` re-renders this repo in place and the result compiles.
Start with [Quickstart](docs/Quickstart/doc.md); the mechanics are in
[Structure](docs/Structure/doc.md) and [BuildPipeline](docs/BuildPipeline/doc.md).

## Installation

`agnos` is a single static binary; Go 1.25+ is only needed to build the projects it
scaffolds, not to run `agnos` itself. Every platform's block is in
[CliInstall](docs/CliInstall/doc.md).


## Documentation

### CliUsage

Driving agnos from a terminal - install, scaffold, declare commands, reference

| Doc | Description |
| --- | --- |
| [CliInstall](docs/CliInstall/doc.md) | Install agnos: a released binary per platform, or a build from source |
| [Install](docs/Install/doc.md) | Install the agnos binary, or run it from a checkout |
| [Quickstart](docs/Quickstart/doc.md) | Empty directory to a compiling CLI with one command, using only agnos commands |
| [Commands](docs/Commands/doc.md) | Every command of agnos, generated from the command declarations |
| [CliExamples](docs/CliExamples/doc.md) | Examples of the agnos cli, and the suite that runs them as tests |

### LibUsage

Using Agnos as a Go module - deps injection, actions, public API

| Doc | Description |
| --- | --- |
| [LibUsage](docs/LibUsage/doc.md) | Use agnos as a Go module: wire the deps, build the sandbox, call its API |
| [PublicApi](docs/PublicApi/doc.md) | Every exported symbol of agnos, generated from the contract sources and their doc comments |
| [LibExamples](docs/LibExamples/doc.md) | Examples of agnos as a Go module, and the suite that runs them as tests |

### Development

Changing this repository - schema, build mechanics, recipes

| Doc | Description |
| --- | --- |
| [Requirements](docs/Requirements/doc.md) | The two tools this project needs — Go and agnos — installed per platform |
| [Workflow](docs/Workflow/doc.md) | Every change this project takes and the agnos command that makes it |
| [Rules](docs/Rules/doc.md) | Every rule the generators, `verify` and the hand-written files must hold to |
| [Structure](docs/Structure/doc.md) | The project schema: what lives where, what is generated, what verify enforces |
| [BuildPipeline](docs/BuildPipeline/doc.md) | What build does: verify, collectors, template vars, asset groups, SmartIO persist, runtime, dispatch |
| [Contributing](docs/Contributing/doc.md) | Recipes specific to changing agnos itself: bootstrap, actions, installable deps, templates, collectors, parsables |

### Reference

Lookup tables - schemas, file formats, generated file listings

| Doc | Description |
| --- | --- |
| [EntriesYaml](docs/EntriesYaml/doc.md) | Every key of a command's entries.yaml and what the generated code does with it |
| [DepList](docs/DepList/doc.md) | Every dep `agnos dep-install` can add, the contract and adapter lib it brings, and what backs it |
| [GeneratedFiles](docs/GeneratedFiles/doc.md) | Every file agnos writes into this project and whether build overwrites it |
| [LibExamples](docs/LibExamples/doc.md) | Examples of agnos as a Go module, and the suite that runs them as tests |
| [CliExamples](docs/CliExamples/doc.md) | Examples of the agnos cli, and the suite that runs them as tests |

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

