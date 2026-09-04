# Install

Go 1.25+ is required to build projects; `agnos` itself is a static binary.

| Platform | Binary |
|---|---|
| macOS arm64 | `macarm64.bin` |
| macOS amd64 | `mac86.bin` |
| Linux amd64 | `linux86.out` |
| Linux arm64 | `linuxarm64.out` |
| Linux 386 | `linuxi32.out` |
| Windows amd64 | `windows86.exe` |
| Windows 386 | `windowsi32.exe` |

macOS / Linux (replace `<binary>`):

```bash
curl -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/<binary> -o agnos && chmod +x agnos && sudo mv agnos /usr/local/bin/ && agnos version
```

Windows (PowerShell, replace `<binary>`):

```powershell
$dir="$HOME\.local\bin"; New-Item -ItemType Directory -Force -Path $dir | Out-Null
curl.exe -sL https://github.com/MateusMoutinhoOrg/Agnos/releases/latest/download/<binary> -o "$dir\agnos.exe"
[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ";$dir", 'User')
```

From a checkout (never regenerate this repo with an older installed binary, see [Contributing](/docs/Contributing/doc.md#bootstrap)):

```bash
git clone https://github.com/MateusMoutinhoOrg/Agnos.git && cd Agnos
go run ./cmd/main local-install      # builds and installs to /usr/local/bin (or ~/.local/bin on Windows)
```
