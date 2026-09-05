package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The set-command example: rewrite the command-level keys of a command
//
// It calls the same action `agnos set-command` calls, and writes only inside TestDir.
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

	if err := lib.Actions.SetCommand(api.CommandProps{
		Path:            "TestDir",
		Command:         "greet",
		LongDescription: "Greets one person by name.",
		Examples:        []string{"greet --name World"},
	}); err != nil {
		panic(err)
	}
}
