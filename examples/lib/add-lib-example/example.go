package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-lib-example example: create an example under examples/lib/
//
// It calls the same action `agnos add-lib-example` calls, and writes only inside test-dir.
func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
	module := "Test"

	if err := lib.Actions.Start(api.StartProps{
		Path:        "test-dir",
		ProjectName: "Test",
		Module:      &module,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddLibExample(api.AddLibExampleProps{Path: "test-dir", Name: "greet"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/examples", os.DirFS("test-dir/examples")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir/docs/LibExamples", os.DirFS("test-dir/docs/LibExamples")); err != nil {
		panic(err)
	}
}
