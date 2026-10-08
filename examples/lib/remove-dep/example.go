package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-dep example: uninstall one installed dep
//
// It calls the same action `agnos remove-dep` calls, and writes only inside test-dir.
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

	if err := lib.Actions.DepsInit(api.DepsInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddDep(api.AddDepProps{Path: "test-dir", Dep: "iodeps"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveDep(api.RemoveDepProps{
		Path:         "test-dir",
		Dep:          "iodeps",
		WithAdapters: true,
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/sandbox/deps", os.DirFS("test-dir/sandbox/deps")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir/adapters", os.DirFS("test-dir/adapters")); err != nil {
		panic(err)
	}
}
