# `api.BuildProps`

**Type:** Struct

## Definition

```go
const (
	RuntimeGo   = "go"
	RuntimeNone = "none"
)

type BuildProps struct {
	Path    string
	Runtime string
}
```

## Description

Describes one (re)render of a project for [`api.Actions.Build`](/docs/References/PublicApi/api.Actions.md#build): the directory holding it and the runtime that then checks the result. `RuntimeGo` runs `go mod tidy` and `go build` over the schema directories that exist, so a build that returns `nil` is a build the Go toolchain accepted; `RuntimeNone` renders only. Any other string is an error. Every follow-up build inside another action names its runtime through this struct — adding actions pass `RuntimeGo`, removing actions pass `RuntimeNone`.

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The project directory. `""` and `"."` both mean the current directory. |
| `Runtime string` | `RuntimeGo`, `RuntimeNone`, or `""` (treated as `none`). |

## Examples

```go
err := lib.Actions.Build(api.BuildProps{Path: "./my-tool", Runtime: api.RuntimeGo})
```
