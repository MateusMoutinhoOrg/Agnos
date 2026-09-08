package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The build example: regenerate every generated file of a project
//
// It calls the same action `agnos build` calls, and writes only inside TestDir.
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

	if err := lib.Actions.Build(api.BuildProps{Path: "TestDir", Runtime: api.RuntimeGo}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/docs", os.DirFS("TestDir/docs")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/AgnosConfig", os.DirFS("TestDir/AgnosConfig")); err != nil {
		panic(err)
	}
}
