# CliInstall

`{{.Name}}` is a single static binary: no runtime, no dependencies. Pick your platform,
paste the block, done. Go 1.25+ is needed only to build it from source.

**macOS (Apple Silicon)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/macarm64.bin -o {{.Name}}
chmod +x {{.Name}} && sudo mv {{.Name}} /usr/local/bin/
{{.Name}} version
```

**macOS (Intel)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/mac86.bin -o {{.Name}}
chmod +x {{.Name}} && sudo mv {{.Name}} /usr/local/bin/
{{.Name}} version
```

**Linux (amd64)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linux86.out -o {{.Name}}
chmod +x {{.Name}} && sudo mv {{.Name}} /usr/local/bin/
{{.Name}} version
```

**Linux (arm64)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linuxarm64.out -o {{.Name}}
chmod +x {{.Name}} && sudo mv {{.Name}} /usr/local/bin/
{{.Name}} version
```

**Linux (32-bit)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linuxi32.out -o {{.Name}}
chmod +x {{.Name}} && sudo mv {{.Name}} /usr/local/bin/
{{.Name}} version
```

**Windows (64-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://{{.Module}}/releases/latest/download/windows86.exe -o "$dir\{{.Name}}.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**Windows (32-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://{{.Module}}/releases/latest/download/windowsi32.exe -o "$dir\{{.Name}}.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**From a checkout** — needs Go 1.25+:

```bash
go build -o {{.Name}} ./cmd/main && sudo mv {{.Name}} /usr/local/bin/
```

The released binaries are the ones `agnos compile --target all` builds and `agnos publish`
uploads. `{{.Name}} version` prints the `version` of `{{.ConfigDir}}/project.yaml`,
`{{.Name}} help` every command — each one is listed in [Commands](../Commands/doc.md).
