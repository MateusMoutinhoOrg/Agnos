package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The list-adapters example: list the adapters of the catalog and of the project
//
// It calls the same action `agnos list-adapters` calls, and writes only inside TestDir.
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

	if err := lib.Actions.DepsInit("TestDir"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddDep(api.AddDepProps{Path: "TestDir", Dep: "sortdeps"}); err != nil {
		panic(err)
	}

	adapters, err := lib.Actions.ListAdapters("TestDir")
	if err != nil {
		panic(err)
	}
	for _, adapter := range adapters {
		fmt.Println(adapter.Name, adapter.Dep, adapter.Installed, adapter.Availables)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/adapters", os.DirFS("TestDir/adapters")); err != nil {
		panic(err)
	}
}
