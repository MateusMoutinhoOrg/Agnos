# `api.DocProps`

**Type:** Struct

## Definition

```go
type DocProps struct {
	Path        string
	Name        string
	Description string
	Themes      []string
}
```

## Description

Describes one doc to create under `docs/`, for [`api.Actions.AddDoc`](/docs/PublicApi/api.Actions/doc.md#adddoc--removedoc). `Name` is the doc's directory, `/`-nested to place it under its parent (`PublicApi/api.AddDoc`), and the doc's `name` key is derived from its last segment. `Themes` are theme **ids** of `<ProjectName>Config/themes.yaml`: required on a first-level doc, refused on a sub-doc, which is listed only by its parent's `Index.md`. The keys these become are governed by [DocProps](/docs/Specs/DocProps/doc.md), the specification of the `props.yaml` the action writes.

## Fields

| Field | Description |
| :--- | :--- |
| `Path string` | The project directory. |
| `Name string` | The doc's directory under `docs/`, `/`-nested for a sub-doc. `Index` is refused at the first level. |
| `Description string` | The one-line summary every index lists the doc with. Required. |
| `Themes []string` | Theme ids from `themes.yaml`. Required on a first-level doc, refused on a sub-doc. |

## Examples

```go
err := lib.Actions.AddDoc(api.DocProps{
	Path: "./my-tool", Name: "HandleReports",
	Description: "How a report is written and regenerated",
	Themes:      []string{"development"},
})

// A sub-doc: no themes, and its parent must already be a doc.
err = lib.Actions.AddDoc(api.DocProps{
	Path: "./my-tool", Name: "HandleReports/Layout",
	Description: "The sections every generated report carries",
})
```
