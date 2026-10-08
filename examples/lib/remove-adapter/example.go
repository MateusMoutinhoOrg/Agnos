package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-adapter example: uninstall one adapter, keeping the contract
//
// It calls the same action `agnos remove-adapter` calls, and writes only inside test-dir.
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

	if err := lib.Actions.AddDep(api.AddDepProps{Path: "test-dir", Dep: "sortdeps"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddAdapter(api.AddAdapterProps{Path: "test-dir", Adapter: "reflectsort"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveAdapter(api.RemoveAdapterProps{Path: "test-dir", Adapter: "reflectsort"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/adapters", os.DirFS("test-dir/adapters")); err != nil {
		panic(err)
	}
}
