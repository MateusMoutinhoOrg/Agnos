package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-command example: delete a command and unwire its dispatch
//
// It calls the same action `agnos remove-command` calls, and writes only inside test-dir.
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

	if err := lib.Actions.CliInit(api.CliInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddCommand(api.AddCommandProps{Path: "test-dir", Name: "greet", Summary: "Greet someone", Category: "Core"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.RemoveCommand(api.RemoveCommandProps{Path: "test-dir", Name: "greet"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/sandbox/internal/commands", os.DirFS("test-dir/sandbox/internal/commands")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir/docs/Commands", os.DirFS("test-dir/docs/Commands")); err != nil {
		panic(err)
	}
}
