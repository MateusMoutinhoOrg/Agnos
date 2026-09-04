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

`agnos` is a single static binary. Pick your platform, paste the block, done. Go 1.25+ is
only needed to build the projects it scaffolds, not to run `agnos` itself.

**macOS (Apple Silicon)**

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/macarm64.bin -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**macOS (Intel)**

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/mac86.bin -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**Linux (amd64)**

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/linux86.out -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**Linux (arm64)**

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/linuxarm64.out -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**Linux (32-bit)**

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/linuxi32.out -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**Windows (64-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/windows86.exe -o "$dir\agnos.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**Windows (32-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/windowsi32.exe -o "$dir\agnos.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**From source** — needs Go 1.25+:

```bash
git clone https://github.com/MateusMoutinhoOrg/Agnos.git && cd Agnos
go run ./cmd/main local-install
```

More detail in [Install](docs/Install/doc.md).


## Documentation Index

| Name | Description |
| --- | --- |
| [CliUsage](docs/Index/cli-usage.md) | Driving agnos from a terminal - install, scaffold, declare commands, reference |
| [LibUsage](docs/Index/lib-usage.md) | Using Agnos as a Go module - deps injection, actions, public API |
| [Development](docs/Index/development.md) | Changing this repository - schema, build mechanics, recipes |

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
