# {{.ProjectName}}

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go CLI that **scaffolds and regenerates other Go CLIs** — each one a closed, dependency-injected sandbox behind a command-line interface generated from declarations.

<p align="center">
  <img src=".github/assets/logo.png" alt="Agnos Logo" width="200"/>
</p>

> [!WARNING]
> **Under Development (Status: Pre-Beta)**
>
> - **Expected Beta:** end of October
> - **Expected Stable:** end of 2026
>
> From pre-beta on, breaking changes are avoided: a change to a declaration's schema, a command's flags, a public contract or the generated output keeps existing projects building, and one that cannot is named in its release notes. Patterns may still change before stable, so pin a version and read the release notes before upgrading.

---

## Overview

Agnos (`agnos`) is a **factory**. `agnos start` writes a project skeleton, `agnos build`
re-renders every generated file from templates embedded in the binary, and commands like
`add-command`, `add-flag` and `add-dep` declare the project's whole command surface
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

`agnos` is a single static binary, but every command that writes a project runs the Go
toolchain on it (`go mod tidy`, `go build`), so Go 1.25+ must be on `PATH`. Every platform's
block is in
[CliInstall](docs/CliInstall/doc.md).
