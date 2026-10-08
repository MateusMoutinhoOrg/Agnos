package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-doc example: delete a doc directory
//
// It calls the same action `agnos remove-doc` calls, and writes only inside test-dir.
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

	if err := lib.Actions.AddDoc(api.AddDocProps{
		Path:        "test-dir",
		Name:        "Report",
		Description: "How a report is written",
		Themes:      []string{"reference"},
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveDoc(api.RemoveDocProps{Path: "test-dir", Name: "Report"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/docs", os.DirFS("test-dir/docs")); err != nil {
		panic(err)
	}
}
