package apishape

import (
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
)

// Source is one file of a sandbox/api/ package: the name it is filed under and
// its content, kept together so a violation can name the file it is in.
type Source struct {
	Name    string
	Content string
}

// Api is one parsed sandbox/api/ package — the contract half of a repo, the
// half a consumer copies into sandbox/deps/<name>/. Every declaration is
// indexed by name because the shape rule and the shim generator ask the same
// question of every type expression: is this identifier a type of this package,
// a predeclared one, or something that cannot cross the boundary.
type Api struct {
	// Package is the name from the package clause of the first file.
	Package string

	// Files are the parsed sources, in the order they were given.
	Files []ParsedFile

	// Types are every top-level type declaration of the package, in file
	// order.
	Types []goimportsdeps.Type

	// ByName indexes Types by declared name.
	ByName map[string]goimportsdeps.Type
}

// ParsedFile is one source and what the Go parser made of it.
type ParsedFile struct {
	Name    string
	Content string
	Parsed  *goimportsdeps.File
}

// Declares reports whether name is a type this package declares — the one test
// that tells a type of the contract from a predeclared one.
func Declares(api *Api, name string) bool {
	_, ok := api.ByName[name]
	return ok
}

// IsStruct reports whether name is a struct type of this package. A struct is
// the only kind that needs a generated converter: its underlying type names
// the package's own types, so the two copies are never identical.
func IsStruct(api *Api, name string) bool {
	entry, ok := api.ByName[name]
	return ok && entry.Kind == "struct"
}
