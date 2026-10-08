package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The disable-extension example: stop generating one mechanic, keep its files
//
// It calls the same action `agnos disable-extension` calls, and writes only inside test-dir.
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

	if err := lib.Actions.DisableExtension(api.DisableExtensionProps{Path: "test-dir", Name: "readme"}); err != nil {
		panic(err)
	}

	// Nothing was removed: turning a mechanic off only stops the generation,
	// so README.md is still there and is the project's from here on.
	_, stat_err := os.Stat("test-dir/README.md")
	fmt.Println("README.md still here:", answer(stat_err == nil))

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := copyFile("test-dir/AgnosConfig/extensions.yaml", "assert-dir/AgnosConfig/extensions.yaml"); err != nil {
		panic(err)
	}
	if err := copyFile("test-dir/README.md", "assert-dir/README.md"); err != nil {
		panic(err)
	}
}

// answer words a yes/no the way the cli side echoes it.
func answer(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// copyFile writes one file of test-dir into assert-dir at the place it holds in
// the tree.
func copyFile(source string, dest string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dirOf(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0o644)
}

// dirOf is the directory part of a slash-separated path.
func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return "."
}
