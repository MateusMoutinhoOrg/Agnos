package goimportsdeps

// This package is the sandbox's *copy* of the api a Go-source parser exposes —
// the same mechanic as argvdeps, dbdeps, embeddeps, iodeps, requestdeps and
// std, for the same reason: parsing Go source pulls in `go/parser`, `go/ast`,
// `go/printer` and `go/token`, which are OS-adjacent enough to stay outside the
// sandbox. The contract is restated here, and the adapter — which lives outside
// the sandbox — is what fills it.

// Lib is the Go-source parser injected whole as the Deps.Goimportsdeps field.
// Every field errors when the given content is not parsable Go.
type Lib struct {
	// Parse parses a whole Go source file into a File describing its package
	// clause, imports and top-level declarations (functions, methods, types,
	// constants and variables).
	Parse func(content string) (*File, error)

	// GetPackageName returns the name in the file's package clause.
	GetPackageName func(content string) (string, error)

	// GetImports returns every import path the file declares, in source order.
	GetImports func(content string) ([]string, error)
}

// File is the parsed view of one Go source file.
type File struct {
	// Package is the name from the package clause.
	Package string
	// Doc is the file's leading doc comment, trimmed, "" when absent.
	Doc string
	// Imports are the import specs, in source order.
	Imports []Import
	// Functions are the top-level function and method declarations, in source
	// order.
	Functions []Function
	// Types are the top-level type declarations, in source order.
	Types []Type
	// Constants are the names declared in top-level `const` blocks.
	Constants []Value
	// Variables are the names declared in top-level `var` blocks.
	Variables []Value
}

// Import is a single import spec.
type Import struct {
	// Alias is the explicit local name ("_", "." or an identifier), "" when
	// the import uses its default name.
	Alias string
	// Path is the unquoted import path.
	Path string
}

// Function is a top-level function or method declaration, or one method of an
// interface type.
type Function struct {
	// Name is the function or method name.
	Name string
	// Doc is the declaration's doc comment, trimmed, "" when absent.
	Doc string
	// Receiver is the receiver type name without a leading "*", "" for plain
	// functions and interface methods.
	Receiver string
	// Pointer reports whether the receiver is a pointer.
	Pointer bool
	// Params are the parameters, one entry per name (unnamed params get a
	// single entry with an empty Name).
	Params []Param
	// Results are the results, same expansion rule as Params.
	Results []Param
	// Exported reports whether Name is exported.
	Exported bool
}

// Param is one parameter or result of a function or method.
type Param struct {
	// Name is the identifier, "" when the param is unnamed.
	Name string
	// Type is the type expression, rendered back to source.
	Type string
}

// Type is a top-level type declaration.
type Type struct {
	// Name is the type name.
	Name string
	// Doc is the declaration's doc comment, trimmed, "" when absent.
	Doc string
	// Kind is one of "struct", "interface", "alias" (`type A = B`) or "other"
	// (any other defined type, e.g. `type ID int`).
	Kind string
	// Fields are the struct fields, populated only when Kind is "struct".
	Fields []Field
	// Methods are the interface methods, populated only when Kind is
	// "interface"; embedded interfaces are skipped.
	Methods []Function
	// Underlying is the underlying type expression, populated only when Kind
	// is "alias" or "other".
	Underlying string
	// Exported reports whether Name is exported.
	Exported bool
}

// Field is one field of a struct type.
type Field struct {
	// Name is the field name, "" for an embedded field.
	Name string
	// Type is the field type expression, rendered back to source.
	Type string
	// Tag is the unquoted struct tag, "" when absent.
	Tag string
	// Doc is the field's doc comment, trimmed, "" when absent.
	Doc string
	// Exported reports whether the field is exported (for an embedded field,
	// whether its type name is exported).
	Exported bool
}

// Value is one name declared in a top-level `const` or `var` block.
type Value struct {
	// Name is the identifier.
	Name string
	// Doc is the spec's doc comment, trimmed, "" when absent.
	Doc string
	// Type is the declared type expression, "" when the type is inferred.
	Type string
	// Value is the assigned expression rendered back to source, "" when the
	// spec assigns nothing (a const repeating the previous expression, or a
	// var declared by type alone).
	Value string
	// Exported reports whether Name is exported.
	Exported bool
}
