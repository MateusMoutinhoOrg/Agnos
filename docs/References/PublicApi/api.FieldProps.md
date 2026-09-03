# `api.FieldProps`

**Type:** Struct

## Definition

```go
type FieldProps struct {
	Path        string
	Command     string
	Name        string
	Identifiers []string
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
	Array       bool
	Min         string
	Max         string
	Position    int
}
```

## Description

Describes one flag or positional argument to add to a command's `entries.yaml`, for [`api.Actions.AddFlag`](/docs/References/PublicApi/api.Actions.md#addflag--removeflag--addarg--removearg) and `AddArg`. `Default`, `Min` and `Max` are the **raw literals** as typed on the command line — `""` means unset — so the action can tell "not given" from a zero value and validate the literal against `Type`. `Identifiers` is ignored by `AddArg` (positionals bind by order) and defaults to `--<Name>` for `AddFlag`. The keys these become, and what the generated dispatch does with each, are in [EntriesYaml.md](/docs/References/EntriesYaml.md#field-keys).

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The project directory. |
| `Command string` | The command receiving the field — its identifier or its package name. |
| `Name string` | The field name; becomes the Go field (`out-file` → `OutFile`). |
| `Identifiers []string` | Flag spellings (`--out`, `-o`). Flags only. |
| `Description string` | Help text printed under the field. |
| `Examples []string` | Usage lines printed under the field. |
| `Type string` | `string` (default when empty), `boolean`, `int`, `float`. |
| `Default string` | Literal assigned when absent. Excludes `Required`. |
| `Required bool` | Absence is a usage error. Refused on booleans and with `Default`. |
| `Array bool` | Collect every occurrence into a slice. An array arg must stay last. |
| `Min string`, `Max string` | Bounds for `int` / `float`. |
| `Position int` | Zero-based index to insert at; `< 0` appends. |

## Examples

```go
err := lib.Actions.AddArg(api.FieldProps{
	Path: "./my-tool", Command: "greet", Name: "times",
	Type: "int", Description: "how many times",
	Default: "1", Min: "1", Position: -1,
})
```
