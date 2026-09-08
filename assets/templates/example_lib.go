package main

import (
	"fmt"
	"os"
{{if .HasDeps}}
	"{{ .Module }}/adapters/availables/standard"{{end}}
	"{{ .Module }}/sandbox"
)

// The {{ .Name }} example, run by `exec-test`.
//
// Write it the way a reader would type it: an example is documentation first
// and a check second. It runs with this directory as the working directory and
// TestDir is the only place it may write — so a failure is a panic here, the
// lib's counterpart of the cli's error message.
func main() {
{{if .HasDeps}}
	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
{{else}}
	lib := sandbox.New() // *api.Sandbox
{{end}}
	// Call lib.Actions.<Action> here.
	_ = lib

	if err := os.MkdirAll("TestDir", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("TestDir/example.txt", []byte("the {{ .Name }} example\n"), 0o644); err != nil {
		panic(err)
	}

	// What result.yaml records is AssertDir, not TestDir: copy into it the
	// directories this example is about, keeping the place each one holds in
	// the tree. TestDir stays as it is, for reading. Copy nothing and the run
	// fails — an example that asserts nothing passes for the wrong reason. The
	// cli side copies the same set, or the cli-vs-lib check breaks over the
	// copy itself.
	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}

	fmt.Println("wrote TestDir/example.txt")
}
