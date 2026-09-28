# `deps.Reflectdeps`

`sandbox/deps/reflectdeps`

## `Sandbox`

Sandbox is the reflection library injected whole as the Deps.Reflectdeps field.

| Field | Type | Description |
| --- | --- | --- |
| `NumIn` | `func(fn any) int` | NumIn returns how many parameters the function fn takes, -1 when fn is not a function. |
| `NewIn` | `func(fn any, index int) any` | NewIn builds a fresh value for the parameter at index of the function fn: a pointer to a new zero value when that parameter is a pointer, the zero value itself otherwise. It returns nil when fn is not a function or index is out of range. |
| `Call` | `func(fn any, args []any) []any` | Call calls the function fn with args, one per parameter, and returns what it returned. A nil arg is passed as the zero value of its parameter. It returns nil when fn is not a function or the argument count differs. |
| `NumField` | `func(target any) int` | NumField returns how many fields the struct target points at, -1 when target is not a pointer to a struct. |
| `FieldName` | `func(target any, index int) string` | FieldName returns the Go name of the field at index of the struct target points at, "" when there is none. |
| `FieldTag` | `func(target any, index int, key string) string` | FieldTag returns the value of key in the tag of the field at index of the struct target points at, "" when the tag or the field is absent. |
| `FieldKind` | `func(target any, index int) string` | FieldKind returns the kind of the field at index of the struct target points at, spelled as the standard library spells it ("string", "float64", "bool", "slice", "struct"), "" when there is none. |
| `SetField` | `func(target any, index int, value any) error` | SetField stores value in the field at index of the struct target points at, and reports an error when the field is absent, not settable, or of a type value cannot be assigned to. |

[every contract](doc.md)
