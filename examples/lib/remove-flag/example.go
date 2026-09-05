package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-flag example: delete one declared flag
//
// It calls the same action `agnos remove-flag` calls, and writes only inside TestDir.
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

	if err := lib.Actions.AddFlag(api.FieldProps{
		Path:        "TestDir",
		Command:     "greet",
		Name:        "name",
		Type:        "string",
		Default:     "world",
		Description: "who to greet",
		Position:    -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveFlag("TestDir", "greet", "name"); err != nil {
		panic(err)
	}
}
