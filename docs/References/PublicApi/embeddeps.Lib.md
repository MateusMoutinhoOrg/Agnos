# `embeddeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	ReadFile             func(path string) ([]byte, error)
	ListFiles            func(path string) ([]string, error)
	ListFilesRecursively func(path string) ([]string, error)
	RenderTemplate       func(path string, vars interface{}) ([]byte, error)
}
```

## Description

The embedded-asset library injected whole as `Deps.Embeddeps`. It is read-only by design: assets ship with the program, and nothing in the sandbox ever writes one back. The sandbox asks for an asset by path the way it would ask a filesystem — `"all/sandbox/new.go"`, `"deplist/iodeps"` — and where the bytes come from is the adapter's decision: the `embeddeps` adapter lib serves them out of the `assets/` tree compiled into the binary by `//go:embed all:*`, so an installed `agnos` carries every template wherever it runs. Every path is slash-separated and relative to the asset root; the root itself is `"."`.

Generated files come from `utils.RenderGroup`, which lists a group with `ListFilesRecursively`, reads each asset with `ReadFile`, and renders it as a Go `text/template` with one `vars` map and the `render` / `copy` native functions — see [BuildPipeline.md](/docs/References/BuildPipeline.md#asset-groups-in-order). `RenderTemplate` bundles the read-and-execute step for a single asset with no native functions; `add-command` uses it for its scaffolds. Installed by the `embeddeps` dep, which also carries `assets/asset.go`.

## Fields

| Field | Description |
| :--- | :--- |
| `ReadFile(path)` | The whole content of one asset. An error is a packaging mistake, not a user mistake. |
| `ListFiles(path)` | The files directly inside `path`, lexical order, relative to it. Directories are not descended into or reported. |
| `ListFilesRecursively(path)` | Every file at or below `path`, lexical order, as slash paths relative to it. |
| `RenderTemplate(path, vars)` | Parses the asset as a Go `text/template` and executes it over `vars`. |

## Examples

```go
// List the deps dep-install can render, exactly as `agnos dep-list` does.
files, _ := deps.Embeddeps.ListFilesRecursively("deplist")
// files: ["argvdeps/adapters/libs/verb/Verb.go", "argvdeps/sandbox/deps/argvdeps/argvdeps.go", …]

// Render one template the way build does.
content, _ := deps.Embeddeps.RenderTemplate("all/sandbox/new.go", map[string]any{
	"Module": "github.com/you/my-tool", "HasDeps": true, "Binds": []string{"CliBind"},
})
```
