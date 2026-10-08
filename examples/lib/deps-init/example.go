package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The deps-init example: add the dependency layer to a project that has none
//
// It calls the same action `agnos deps-init` calls, and writes only inside test-dir.
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

	if err := lib.Actions.DepsInit(api.DepsInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := os.CopyFS("assert-dir/sandbox/deps", os.DirFS("test-dir/sandbox/deps")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir/adapters", os.DirFS("test-dir/adapters")); err != nil {
		panic(err)
	}

	// The declaration the pair wrote: this is the whole of what tells the
	// build the mechanic is on or off from here.
	if err := copyExtensions(); err != nil {
		panic(err)
	}
}

// copyExtensions puts the project's extensions.yaml into assert-dir at the
// place it holds in the tree.
func copyExtensions() error {
	content, err := os.ReadFile("test-dir/AgnosConfig/extensions.yaml")
	if err != nil {
		return err
	}
	if err := os.MkdirAll("assert-dir/AgnosConfig", 0o755); err != nil {
		return err
	}
	return os.WriteFile("assert-dir/AgnosConfig/extensions.yaml", content, 0o644)
}
