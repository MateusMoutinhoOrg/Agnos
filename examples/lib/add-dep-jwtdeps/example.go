package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-dep-jwtdeps example: install a catalog dep whose adapter needs a
// third-party module
//
// It calls the same action `agnos add-dep` calls, and writes only inside TestDir.
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

	if err := lib.Actions.AddDep(api.AddDepProps{Path: "TestDir", Dep: "jwtdeps"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/sandbox/deps", os.DirFS("TestDir/sandbox/deps")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/adapters", os.DirFS("TestDir/adapters")); err != nil {
		panic(err)
	}
	content, err := os.ReadFile("TestDir/go.mod")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("AssertDir/go.mod", content, 0o644); err != nil {
		panic(err)
	}
}
