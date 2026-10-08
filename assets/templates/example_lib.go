package main

import (
	"fmt"
	"os"
{{if .HasDeps}}
	"{{ .Module }}/adapters/bindings/standard"{{end}}
	"{{ .Module }}/sandbox"
)

// The {{ .ExampleName }} example, run by `run-examples`.
//
// Write it the way a reader would type it: an example is documentation first
// and a check second. It runs with this directory as the working directory and
// test-dir is the only place it may write — so a failure is a panic here, the
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

	if err := os.MkdirAll("test-dir", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("test-dir/example.txt", []byte("the {{ .ExampleName }} example\n"), 0o644); err != nil {
		panic(err)
	}

	// What result.yaml records is assert-dir, not test-dir: copy into it the
	// directories this example is about, keeping the place each one holds in
	// the tree. test-dir stays as it is, for reading. Copy nothing and the run
	// fails — an example that asserts nothing passes for the wrong reason. The
	// cli side copies the same set, or the cli-vs-lib check breaks over the
	// copy itself.
	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}

	fmt.Println("wrote test-dir/example.txt")
}
