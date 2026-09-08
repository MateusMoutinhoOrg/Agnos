package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-lib-example example: create an example under examples/lib/
//
// It calls the same action `agnos add-lib-example` calls, and writes only inside TestDir.
func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
	module := "Test"

	if err := lib.Actions.Start(api.StartProps{
		Path:        "TestDir",
		ProjectName: "Test",
		Module:      &module,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddLibExample("TestDir", "greet"); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/examples", os.DirFS("TestDir/examples")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/docs/LibExamples", os.DirFS("TestDir/docs/LibExamples")); err != nil {
		panic(err)
	}
}
