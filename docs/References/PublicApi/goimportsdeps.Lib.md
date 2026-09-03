# `goimportsdeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	Parse          func(content string) (*File, error)
	GetPackageName func(content string) (string, error)
	GetImports     func(content string) ([]string, error)
}

type File struct {
	Package   string
	Doc       string
	Imports   []Import
	Functions []Function
	Types     []Type
	Constants []Value
	Variables []Value
}
```

## Description

A Go source reader injected whole as `Deps.Goimportsdeps` — the sandbox's copy of what `go/parser`, `go/ast` and `go/printer` provide, kept outside the wall because those packages are OS-adjacent. `Parse` returns the whole top-level view of one file: package clause, imports (`Import{Alias, Path}`), functions and methods (`Function{Name, Doc, Receiver, Pointer, Params, Results, Exported}`), types (`Type{Name, Doc, Kind, Fields, Methods, Underlying, Exported}` with `Kind` one of `struct`, `interface`, `alias`, `other`), and the names of top-level `const` and `var` blocks (`Value{Name, Doc, Type, Exported}`). Every field errors when the content is not parsable Go.

It is wired by the `goimportsdeps` dep and the standard adapter, but **no action uses it yet**; `verify` parses Go with its own helpers. It ships as a standing capability for a generated project — or a future Agnos — that inspects Go source.

## Fields

| Field | Description |
| :--- | :--- |
| `Parse(content)` | The full `*File` view. |
| `GetPackageName(content)` | The name in the package clause. |
| `GetImports(content)` | Every import path, in source order, unquoted. Parses the import section only. |

## Examples

```go
imports, err := deps.Goimportsdeps.GetImports(source)
if err != nil {
	return err // not Go
}
for _, path := range imports {
	if strings.HasPrefix(path, module+"/adapters") {
		violations = append(violations, file+" imports "+path)
	}
}
```
