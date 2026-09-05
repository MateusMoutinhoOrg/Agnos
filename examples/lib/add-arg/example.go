package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-arg example: declare one positional argument on a command
//
// It calls the same action `agnos add-arg` calls, and writes only inside TestDir.
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

	if err := lib.Actions.CliInit("TestDir"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddCommand("TestDir", "greet", "Greet someone", "Core"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddArg(api.FieldProps{
		Path:        "TestDir",
		Command:     "greet",
		Name:        "who",
		Type:        "string",
		Required:    true,
		Description: "the person to greet",
		Position:    -1,
	}); err != nil {
		panic(err)
	}
}
