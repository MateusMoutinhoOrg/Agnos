# `deps.Sortdeps`

`sandbox/deps/sortdeps`

## `Sandbox`

Sandbox is the sorting library injected whole as the Deps.Sortdeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Strings` | `func(list []string)` | Strings sorts a slice of strings into increasing order. |
| `Slice` | `func(slice any, less func(i int, j int) bool)` | Slice sorts slice given the less function, which reports whether the element at index i must sort before the one at index j. The sort is not guaranteed to be stable; use SliceStable when equal elements must keep their original order. |
| `SliceStable` | `func(slice any, less func(i int, j int) bool)` | SliceStable sorts slice given the less function, keeping equal elements in their original order. |

[every contract](doc.md)
