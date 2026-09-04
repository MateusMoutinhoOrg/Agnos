# Requirements

Two tools, nothing else. Every recipe of [Workflow](../Workflow/doc.md) assumes both.

| Tool | Version | Needed for |
| --- | --- | --- |
| Go | 1.25+ | compiling this project; `agnos build` ends in a `go mod tidy` and a compile |
| agnos | latest | every generated file — the tree cannot be maintained by hand |

```bash
go version      # go1.25.0 or newer
agnos version
```

## Go 1.25+

A distro package is usually older than 1.25; prefer the official tarball or installer.

| Platform | Install |
| --- | --- |
| macOS | `brew install go`, or the `.pkg` from https://go.dev/dl/ |
| Linux | the tarball below |
| Windows | `winget install GoLang.Go`, or the `.msi` from https://go.dev/dl/ |

**Linux** — `<arch>` is `amd64`, `arm64` or `386`:

```bash
curl -sL https://go.dev/dl/go1.25.0.linux-<arch>.tar.gz -o go.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && . ~/.profile
go version
```

An existing Go is replaced, never upgraded in place: delete `/usr/local/go` (or the old
`.pkg` / `.msi` install) before unpacking a new one.

## agnos

A single static binary — no runtime, no dependencies. Pick the platform's asset:

| Platform | Binary |
| --- | --- |
| macOS arm64 | `macarm64.bin` |
| macOS amd64 | `mac86.bin` |
| Linux amd64 | `linux86.out` |
| Linux arm64 | `linuxarm64.out` |
| Linux 386 | `linuxi32.out` |
| Windows amd64 | `windows86.exe` |
| Windows 386 | `windowsi32.exe` |

**macOS / Linux** — replace `<binary>`:

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/<binary> -o agnos
chmod +x agnos && sudo mv agnos /usr/local/bin/
agnos version
```

**Windows** — PowerShell, replace `<binary>`:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/<binary> -o "$dir\agnos.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**From source** — needs Go 1.25+ first:

```bash
git clone https://github.com/MateusMoutinhoOrg/Agnos.git && cd Agnos
go run ./cmd/main local-install
```

## This repository

This project carries its own `assets/` template tree, so an installed `agnos` would rewrite it
to that older binary's shape. Never run one here — build the binary from this tree instead:

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build
```
