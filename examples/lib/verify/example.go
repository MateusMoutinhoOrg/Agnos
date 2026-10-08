package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The verify example: check a project against the schema, writing nothing
//
// It calls the same action `agnos verify` calls, and writes only inside test-dir.
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

	if err := lib.Actions.Verify(api.VerifyProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	// verify writes nothing, so what this asserts is the config it read,
	// untouched.
	if err := os.CopyFS("assert-dir/AgnosConfig", os.DirFS("test-dir/AgnosConfig")); err != nil {
		panic(err)
	}
}
