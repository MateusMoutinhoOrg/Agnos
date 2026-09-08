package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-doc example: delete a doc directory
//
// It calls the same action `agnos remove-doc` calls, and writes only inside TestDir.
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

	if err := lib.Actions.AddDoc(api.DocProps{
		Path:        "TestDir",
		Name:        "Report",
		Description: "How a report is written",
		Themes:      []string{"reference"},
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveDoc("TestDir", "Report"); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/docs", os.DirFS("TestDir/docs")); err != nil {
		panic(err)
	}
}
