# `deps.Goimportsdeps`

`sandbox/deps/goimportsdeps`

## `Sandbox`

Sandbox is the Go-source parser injected whole as the Deps.Goimportsdeps field. Every field errors when the given content is not parsable Go.

| Field | Type | Description |
| --- | --- | --- |
| `Parse` | `func(content string) (*File, error)` | Parse parses a whole Go source file into a File describing its package clause, imports and top-level declarations (functions, methods, types, constants and variables). |
| `GetPackageName` | `func(content string) (string, error)` | GetPackageName returns the name in the file's package clause. |
| `GetImports` | `func(content string) ([]string, error)` | GetImports returns every import path the file declares, in source order. |
| `Format` | `func(content string) (string, error)` | Format rewrites the given source in the canonical form of the Go toolchain — the same bytes `gofmt` would write: standard indentation, aligned struct fields and comment blocks, no trailing whitespace. Everything agnos generates goes through it, so a regenerated tree diffs to zero against one a formatting editor has saved. |

## `File`

File is the parsed view of one Go source file.

| Field | Type | Description |
| --- | --- | --- |
| `Package` | `string` | Package is the name from the package clause. |
| `Doc` | `string` | Doc is the file's leading doc comment, trimmed, "" when absent. |
| `Imports` | `[]Import` | Imports are the import specs, in source order. |
| `Functions` | `[]Function` | Functions are the top-level function and method declarations, in source order. |
| `Types` | `[]Type` | Types are the top-level type declarations, in source order. |
| `Constants` | `[]Value` | Constants are the names declared in top-level `const` blocks. |
| `Variables` | `[]Value` | Variables are the names declared in top-level `var` blocks. |

## `Import`

Import is a single import spec.

| Field | Type | Description |
| --- | --- | --- |
| `Alias` | `string` | Alias is the explicit local name ("_", "." or an identifier), "" when the import uses its default name. |
| `Path` | `string` | Path is the unquoted import path. |

## `Function`

Function is a top-level function or method declaration, or one method of an interface type.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the function or method name. |
| `Doc` | `string` | Doc is the declaration's doc comment, trimmed, "" when absent. |
| `Receiver` | `string` | Receiver is the receiver type name without a leading "*", "" for plain functions and interface methods. |
| `Pointer` | `bool` | Pointer reports whether the receiver is a pointer. |
| `Params` | `[]Param` | Params are the parameters, one entry per name (unnamed params get a single entry with an empty Name). |
| `Results` | `[]Param` | Results are the results, same expansion rule as Params. |
| `Exported` | `bool` | Exported reports whether Name is exported. |

## `Param`

Param is one parameter or result of a function or method.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the identifier, "" when the param is unnamed. |
| `Type` | `string` | Type is the type expression, rendered back to source. |

## `Type`

Type is a top-level type declaration.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the type name. |
| `Doc` | `string` | Doc is the declaration's doc comment, trimmed, "" when absent. |
| `Kind` | `string` | Kind is one of "struct", "interface", "alias" (`type A = B`) or "other" (any other defined type, e.g. `type ID int`). |
| `Fields` | `[]Field` | Fields are the struct fields, populated only when Kind is "struct". |
| `Methods` | `[]Function` | Methods are the interface methods, populated only when Kind is "interface"; embedded interfaces are skipped. |
| `Underlying` | `string` | Underlying is the underlying type expression, populated only when Kind is "alias" or "other". |
| `Exported` | `bool` | Exported reports whether Name is exported. |

## `Field`

Field is one field of a struct type.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the field name, "" for an embedded field. |
| `Type` | `string` | Type is the field type expression, rendered back to source. |
| `Tag` | `string` | Tag is the unquoted struct tag, "" when absent. |
| `Doc` | `string` | Doc is the field's doc comment, trimmed, "" when absent. |
| `Exported` | `bool` | Exported reports whether the field is exported (for an embedded field, whether its type name is exported). |

## `Value`

Value is one name declared in a top-level `const` or `var` block.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the identifier. |
| `Doc` | `string` | Doc is the spec's doc comment, trimmed, "" when absent. |
| `Type` | `string` | Type is the declared type expression, "" when the type is inferred. |
| `Value` | `string` | Value is the assigned expression rendered back to source, "" when the spec assigns nothing (a const repeating the previous expression, or a var declared by type alone). |
| `Exported` | `bool` | Exported reports whether Name is exported. |

[every contract](doc.md)
