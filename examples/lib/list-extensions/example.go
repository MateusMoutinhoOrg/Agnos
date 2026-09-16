package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The list-extensions example: what agnos generates for this project
//
// It calls the same action `agnos list-extensions` calls, and writes only inside TestDir.
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

	extensions, err := lib.Actions.ListExtensions("TestDir")
	if err != nil {
		panic(err)
	}
	for _, extension := range extensions {
		fmt.Println(extension)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := copyFile("TestDir/AgnosConfig/extensions.yaml", "AssertDir/AgnosConfig/extensions.yaml"); err != nil {
		panic(err)
	}
}

// copyFile writes one file of TestDir into AssertDir at the place it holds in
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
