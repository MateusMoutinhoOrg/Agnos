package main

import (
	"fmt"
{{if .HasDeps}}
	"{{ .Module }}/adapters/availables/standard"{{end}}
	"{{ .Module }}/sandbox"
)

// The {{ .Name }} example, run by `exec-test`.
//
// Write it the way a reader would type it: an example is documentation first
// and a check second. It runs with this directory as the working directory and
// TestDir is the only place it may write. Everything it prints, the status it
// exits with and every file it leaves in TestDir become result.yaml — so a
// failure is a panic here, the lib's counterpart of the cli's error message.
func main() {
{{if .HasDeps}}
	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
{{else}}
	lib := sandbox.New() // *api.Sandbox
{{end}}
	// Call lib.Actions.<Action> here.
	_ = lib

	fmt.Println("the {{ .Name }} example")
}
