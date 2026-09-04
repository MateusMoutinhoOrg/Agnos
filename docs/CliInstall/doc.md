# CliInstall

`agnos` is a single static binary: no runtime, no dependencies. Pick your platform,
paste the block, done. Go 1.25+ is needed only to build it from source.

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

**From a checkout** — needs Go 1.25+:

```bash
go build -o agnos ./cmd/main && sudo mv agnos /usr/local/bin/
```

The released binaries are the ones `agnos compile --target all` builds and `agnos publish`
uploads. `agnos version` prints the `version` of `AgnosConfig/project.yaml`,
`agnos help` every command — each one is listed in [Commands](../Commands/doc.md).
