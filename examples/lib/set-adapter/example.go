package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The set-adapter example: point one available at another adapter
//
// It calls the same action `agnos set-adapter` calls, and writes only inside TestDir.
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

	if err := lib.Actions.AddAdapter(api.AddAdapterProps{Path: "TestDir", Adapter: "reflectsort"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddAvailable("TestDir", "lambda"); err != nil {
		panic(err)
	}

	if err := lib.Actions.SetAdapter(api.SetAdapterProps{
		Path:      "TestDir",
		Dep:       "sortdeps",
		Adapter:   "reflectsort",
		Available: "lambda",
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/adapters", os.DirFS("TestDir/adapters")); err != nil {
		panic(err)
	}
}
