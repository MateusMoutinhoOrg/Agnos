//go:build ignore

// This file is an illustrative sample, not part of the build.
package deps

// Deps is the dependency contract every adapter must fill.
// Each field is one injectable behavior the library needs.
type Deps struct {
	ExampleDepFunctionA func() int
}
