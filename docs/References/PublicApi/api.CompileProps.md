# `api.CompileProps`

**Type:** Struct

## Definition

```go
type CompileProps struct {
	Path    string
	Targets []string
}
```

## Description

Describes one cross-compile run for [`api.Actions.Compile`](/docs/References/PublicApi/api.Actions.md#compile): the directory holding the project and the targets to build. `Compile` runs [`Build`](/docs/References/PublicApi/api.Actions.md#build) with `RuntimeGo` first, creates `release/`, then runs `go build -o release/<file> ./cmd/main` once per target with `CGO_ENABLED=0` and the target's `GOOS`/`GOARCH`.

Each entry of `Targets` is one of `linux86`, `linuxarm64`, `linuxi32`, `mac86`, `macarm64`, `windows86`, `windowsi32`, or `all` (which expands to every target, in that order). An unknown name or an empty slice is an error.

| Target | `GOOS`/`GOARCH` | Output file |
| :--- | :--- | :--- |
| `linux86` | `linux/amd64` | `release/linux86.out` |
| `linuxarm64` | `linux/arm64` | `release/linuxarm64.out` |
| `linuxi32` | `linux/386` | `release/linuxi32.out` |
| `mac86` | `darwin/amd64` | `release/mac86.bin` |
| `macarm64` | `darwin/arm64` | `release/macarm64.bin` |
| `windows86` | `windows/amd64` | `release/windows86.exe` |
| `windowsi32` | `windows/386` | `release/windowsi32.exe` |

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The project directory. `""` and `"."` both mean the current directory. |
| `Targets []string` | One or more target names, or `"all"`. Duplicates are ignored; written order is kept. |

## Examples

```go
err := lib.Actions.Compile(api.CompileProps{Path: "./my-tool", Targets: []string{"linux86", "windows86"}})

err = lib.Actions.Compile(api.CompileProps{Path: ".", Targets: []string{"all"}})
```
