# `api.CommandProps`

**Type:** Struct

## Definition

```go
type CommandProps struct {
	Path            string
	Command         string
	Help            string
	Category        string
	LongDescription string
	Hidden          bool
	Visible         bool
	Identifiers     []string
	Examples        []string
}
```

## Description

Carries the command-level keys of `entries.yaml` that [`api.Actions.SetCommand`](/docs/References/PublicApi/api.Actions.md#addcommand--removecommand--setcommand) may rewrite. Empty strings leave the current value alone; `Identifiers` and `Examples` are appended and deduplicated; `Hidden` and `Visible` are the two directions of one switch, and setting neither leaves visibility untouched. The keys are described in [EntriesYaml.md](/docs/References/EntriesYaml.md#command-keys).

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The project directory. |
| `Command string` | The command to update — identifier or package name. |
| `Help string` | New one-line help. |
| `Category string` | New category for the `help` listing. |
| `LongDescription string` | New paragraph for `help <command>`. |
| `Hidden bool`, `Visible bool` | Hide from, or restore to, the general listing. |
| `Identifiers []string` | Extra verbs the command answers to. |
| `Examples []string` | Extra usage examples. |

## Examples

```go
err := lib.Actions.SetCommand(api.CommandProps{
	Path: "./my-tool", Command: "greet",
	Category:    "Core",
	Identifiers: []string{"hello"},
	Hidden:      true,
})
```
