# `deps.Embeddeps`

`sandbox/deps/embeddeps`

## `Sandbox`

Sandbox is the embedded-asset library injected whole as the Deps.EmbedDeps field. It is read-only by design: assets ship with the program, and nothing in the library ever writes one back. Every path is slash-separated and relative to the root of the asset tree the adapter serves — "report.tmpl", "templates/invoice.tmpl" — never an absolute path and never a path reaching outside that root, so the same call means the same asset whatever the adapter is backed by.

| Field | Type | Description |
| --- | --- | --- |
| `ReadFile` | `func(path string) ([]byte, error)` | ReadFile returns the whole content of one asset. The error reports an asset that does not exist or could not be read; callers inside the sandbox report it rather than assuming the bytes are there, because a missing asset is a packaging mistake, not a user mistake. |
| `ListFiles` | `func(path string) ([]string, error)` | ListFiles returns the names of the assets directly inside the given directory, in lexical order, relative to that directory. Nested directories are not descended into and are not reported. The root itself is addressed as ".". |
| `ListFilesRecursively` | `func(path string) ([]string, error)` | ListFilesRecursively returns every asset at or below the given directory, in lexical order, as slash-separated paths relative to that directory — "templates/invoice.tmpl" and not just "invoice.tmpl". Directories are never reported, only the files inside them. |
| `RenderTemplate` | `func(path string, vars interface{}) ([]byte, error)` | RenderTemplate reads the template from path, renders it using the given variables, and returns the resulting byte slice. |

[every contract](doc.md)
