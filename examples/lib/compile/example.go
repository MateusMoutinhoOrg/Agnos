package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The compile example: cross-compile a project's cmd/main into release/
//
// It calls the same action `agnos compile` calls, and writes only inside test-dir.
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

	if err := lib.Actions.Compile(api.CompileProps{Path: "test-dir", Targets: []string{"linux86"}}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	// The binary in release/ is machine-specific and no golden, so what this
	// asserts is the source it was built from, untouched.
	if err := os.CopyFS("assert-dir/cmd", os.DirFS("test-dir/cmd")); err != nil {
		panic(err)
	}
}
