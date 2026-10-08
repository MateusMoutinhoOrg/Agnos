package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The list-extensions example: what agnos generates for this project
//
// It calls the same action `agnos list-extensions` calls, and writes only inside test-dir.
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

	extensions, err := lib.Actions.ListExtensions(api.ListExtensionsProps{Path: "test-dir"})
	if err != nil {
		panic(err)
	}
	for _, extension := range extensions {
		fmt.Println(extension)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// test-dir. The cli side copies the same set.
	if err := copyFile("test-dir/AgnosConfig/extensions.yaml", "assert-dir/AgnosConfig/extensions.yaml"); err != nil {
		panic(err)
	}
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
