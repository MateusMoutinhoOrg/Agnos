# `deps.Templatedeps`

`sandbox/deps/templatedeps`

## `Sandbox`

Sandbox is the template engine injected whole as the Deps.Templatedeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Render` | `func(props RenderProps) (string, error)` | Render parses one template source and executes it over the given vars, returning the result. The error reports a source that does not parse or an execution that failed — a native function returning an error included. |

## `RenderProps`

RenderProps is the whole input of one render.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the template name, used in error messages only. |
| `Source` | `string` | Source is the template text to parse and execute. |
| `Vars` | `any` | Vars is the value the template renders over, reached as `.`. |
| `Funcs` | `map[string]any` | Funcs are the native functions the template may call, by the name each is registered under. A value must be a function the engine accepts: one return value, or one return value and an error. |

[every contract](doc.md)
