package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The build example: regenerate every generated file of a project
//
// It calls the same action `agnos build` calls, and writes only inside test-dir.
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

	if err := lib.Actions.Build(api.BuildProps{Path: "test-dir", Runtime: api.RuntimeGo}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/docs", os.DirFS("test-dir/docs")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir/AgnosConfig", os.DirFS("test-dir/AgnosConfig")); err != nil {
		panic(err)
	}
}
