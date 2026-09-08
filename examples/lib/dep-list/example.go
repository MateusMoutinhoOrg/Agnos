package main

import (
	"os"

	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The dep-list example: list the deps available to a project
//
// It calls the same action `agnos dep-list` calls, and writes only inside TestDir.
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

	if err := lib.Actions.DepInstall("TestDir", "iodeps"); err != nil {
		panic(err)
	}

	installed, err := lib.Actions.DepList("TestDir")
	if err != nil {
		panic(err)
	}
	for _, dep := range installed {
		fmt.Println(dep)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/sandbox/deps", os.DirFS("TestDir/sandbox/deps")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/adapters", os.DirFS("TestDir/adapters")); err != nil {
		panic(err)
	}
}
