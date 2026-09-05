package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-doc example: create a doc directory under docs/
//
// It calls the same action `agnos add-doc` calls, and writes only inside TestDir.
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
}
