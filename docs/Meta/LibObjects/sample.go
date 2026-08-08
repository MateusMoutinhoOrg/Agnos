//go:build ignore

// This file is an illustrative sample, not part of the build.
package example_object

import (
	"strconv"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// ExampleObjectMethodFactory returns the closure that fills
// api.ExampleLibObject.ExampleObjectMethod, closed over o, so it reads the
// object's own properties — and its propagated Deps — at call time.
func ExampleObjectMethodFactory(o *api.ExampleLibObject) func() string {
	return func() string {
		propsConcatenated := strconv.Itoa(o.FirstProp) + o.SecondProp
		return "Props Concatenated: " + propsConcatenated
	}
}

// New builds an api.ExampleLibObject from its properties, propagating the
// library's Deps into it, and runs every object factory over it to fill its
// function fields. Adding a function field to api.ExampleLibObject means
// adding its factory call here.
func New(d deps.Deps, firstProp int, secondProp string) api.ExampleLibObject {
	o := api.ExampleLibObject{
		Deps:       d,
		FirstProp:  firstProp,
		SecondProp: secondProp,
	}
	o.ExampleObjectMethod = ExampleObjectMethodFactory(&o)
	return o
}

// ---------------------------------------------------------------
// The constructor factory below lives in sandbox/internal/lib/, where
// the object is created with the parent lib's Deps propagated into it.
// ---------------------------------------------------------------

// NewExampleObjectFactory returns the closure that fills
// api.Lib.NewExampleObject, creating an ExampleLibObject with the library's
// injected dependencies automatically wired in.
func NewExampleObjectFactory(l *api.Lib) func(firstProp int, secondProp string) api.ExampleLibObject {
	return func(firstProp int, secondProp string) api.ExampleLibObject {
		return example_object.New(l.Deps, firstProp, secondProp)
	}
}
