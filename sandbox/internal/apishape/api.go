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

	// Mechanic is every type of the package that is a mechanic's surface
	// rather than the repo's own contract (see IsMechanic).
	Mechanic map[string]bool
}

// ParsedFile is one source and what the Go parser made of it.
type ParsedFile struct {
	Name    string
	Content string
	Parsed  *goimportsdeps.File
}

// Declares reports whether name is a type this package declares — the one test
// that tells a type of the contract from a predeclared one.
func Declares(shape *Api, name string) bool {
	_, ok := shape.ByName[name]
	return ok
}

// IsStruct reports whether name is a struct type of this package. A struct is
// the only kind that needs a generated converter: its underlying type names
// the package's own types, so the two copies are never identical.
func IsStruct(shape *Api, name string) bool {
	entry, ok := shape.ByName[name]
	return ok && entry.Kind == "struct"
}

// SandboxType is the root type of an api package: the one a consumer converts,
// and the one the converter plan is walked from.
const SandboxType = "Sandbox"

// DepsField is the field of Sandbox holding the repo's own dependency set. It
// is the one declaration of an api package that does not cross into a
// consumer: it names sandbox/deps, which is how this repo reaches the outside
// world and not part of what it offers. So the copy drops it, the shape rule
// skips it and no converter is written for it — a consumer installs the api of
// a repo, never its wiring.
const DepsField = "Deps"

// IsDepsWiring reports whether a field of a type is Sandbox.Deps, the one
// field that stays behind when the api is copied into a consumer.
func IsDepsWiring(typeName string, fieldName string) bool {
	return typeName == SandboxType && fieldName == DepsField
}

// OpinatedPrefix starts the last segment of an opinated lib's contract import
// path: sandbox/deps/OpinatedAgnos<X>, the dep that carries a mechanic.
const OpinatedPrefix = "OpinatedAgnos"

// IsMechanic reports whether a type of the package is a mechanic's surface: an
// alias of a type an opinated lib's contract declares — api.Command is
// opinatedagnoscli.Command — or a struct every field of which is one, the
// <x>sandbox.go part a mechanic adds. The lib owns those types and a consumer
// installs a repo's api, never a mechanic of it: like Sandbox.Deps, the copy
// drops them, the shape rule skips them and no converter is written for them.
func IsMechanic(shape *Api, name string) bool {
	return shape.Mechanic[name]
}

// IsMechanicField reports whether one field of a struct is typed with a
// mechanic type — CliSandbox embedded in Sandbox — and so stays behind with it.
func IsMechanicField(shape *Api, field goimportsdeps.Field) bool {
	return IsMechanic(shape, field.Type)
}

// IsMechanicFile reports whether one file declares nothing but mechanic
// surface: every type it declares is one and it declares no function. Its
// constants alias the lib's too, so the whole file stays behind.
func IsMechanicFile(shape *Api, file ParsedFile) bool {
	if len(file.Parsed.Types) == 0 || len(file.Parsed.Functions) > 0 {
		return false
	}
	for _, entry := range file.Parsed.Types {
		if !IsMechanic(shape, entry.Name) {
			return false
		}
	}
	return true
}

// WithoutMechanic is the package a consumer copies: every mechanic file left
// out, with the types it declared. The fields that name them stay on the
// struct that declares them, for the copy to strip and the converter to skip.
func WithoutMechanic(shape *Api) *Api {
	kept := &Api{Package: shape.Package, ByName: map[string]goimportsdeps.Type{}, Mechanic: shape.Mechanic}
	for _, file := range shape.Files {
		if IsMechanicFile(shape, file) {
			continue
		}
		kept.Files = append(kept.Files, file)
		for _, entry := range file.Parsed.Types {
			kept.Types = append(kept.Types, entry)
			kept.ByName[entry.Name] = entry
		}
	}
	return kept
}

// FieldName is the name a struct field is read and written by: its own, or
// for an embedded field — which Violations admits only for a struct of the
// package — the name of the type it embeds.
func FieldName(field goimportsdeps.Field) string {
	if field.Name == "" {
		return field.Type
	}
	return field.Name
}
