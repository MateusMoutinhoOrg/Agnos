# CliInstall

{{ if .HasAssets }}`{{.ProjectName}}` is a single static binary, but every command that writes a project runs the Go
toolchain on it, so Go {{.GoFloor}}+ must be on `PATH`. Pick your platform, paste the block, done.{{ else }}`{{.ProjectName}}` is a single static binary: no runtime, no dependencies. Pick your platform,
paste the block, done. Go {{.GoFloor}}+ is needed only to build it from source.{{ end }}

**macOS (Apple Silicon)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/macarm64.bin -o {{.ProjectName}} && chmod +x {{.ProjectName}} && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

**macOS (Intel)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/mac86.bin -o {{.ProjectName}} && chmod +x {{.ProjectName}} && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

**Linux (amd64)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linux86.out -o {{.ProjectName}} && chmod +x {{.ProjectName}} && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

**Linux (arm64)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linuxarm64.out -o {{.ProjectName}} && chmod +x {{.ProjectName}} && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

**Linux (32-bit)**

```bash
curl -sL https://{{.Module}}/releases/latest/download/linuxi32.out -o {{.ProjectName}} && chmod +x {{.ProjectName}} && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

**Windows (64-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null; curl.exe -sL https://{{.Module}}/releases/latest/download/windows86.exe -o "$dir\{{.ProjectName}}.exe"; [Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**Windows (32-bit)** — PowerShell:

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null; curl.exe -sL https://{{.Module}}/releases/latest/download/windowsi32.exe -o "$dir\{{.ProjectName}}.exe"; [Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

**From a checkout** — needs Go {{.GoFloor}}+:

```bash
go build -o {{.ProjectName}} ./cmd/main && sudo mv {{.ProjectName}} /usr/local/bin/ && {{.ProjectName}} version
```

The released binaries are the ones `{{.GeneratorName}} compile --target all` builds and `{{.GeneratorName}} publish`
uploads. `{{.ProjectName}} version` prints the `version` of `{{.ConfigDir}}/project.yaml`,
`{{.ProjectName}} help` every command — each one is listed in [Commands](../Commands/doc.md).
