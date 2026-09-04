# `api.StartProps`

**Type:** Struct

## Definition

```go
type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}
```

## Description

Describes one scaffold for [`api.Actions.Start`](/docs/PublicApi/api.Actions/doc.md#start). `ProjectName` becomes `name` in `AgnosConfig/project.yaml` and, title-cased, the `config.ProjectName` constant; lowercased, it is the binary name the generated `help` screens print. `Module` is a pointer so "not given" is distinguishable from an empty string: when set, `go.mod` is written with it and the installed toolchain's version (`go env GOVERSION`, falling back to `1.25.0`); the `start` command requires it when the directory has no `go.mod`.

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The directory to scaffold into. |
| `ProjectName string` | The project's name. Required. |
| `Module *string` | The Go module path written into `go.mod`. `nil` leaves `go.mod` alone. |
| `Force bool` | Overwrite an existing `go.mod` instead of refusing. |

## Examples

```go
module := "github.com/you/my-tool"
err := lib.Actions.Start(api.StartProps{
	Path:        "./my-tool",
	ProjectName: "my-tool",
	Module:      &module,
})
```
