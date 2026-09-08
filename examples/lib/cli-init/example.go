package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The cli-init example: add the cli layer to a project that has none
//
// It calls the same action `agnos cli-init` calls, and writes only inside TestDir.
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

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/cmd", os.DirFS("TestDir/cmd")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/sandbox/internal/cli", os.DirFS("TestDir/sandbox/internal/cli")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/sandbox/internal/commands", os.DirFS("TestDir/sandbox/internal/commands")); err != nil {
		panic(err)
	}
}
