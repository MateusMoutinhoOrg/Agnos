package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The set-dep example: re-copy a remote dep after its contract moved
//
// It calls the same action `agnos set-dep` calls, and writes only inside
// TestDir. The setup is the one in add-remote-dep: TestDir/remote is the repo
// being installed and TestDir/app the consumer, wired with a `replace` so the
// pair can be developed side by side. Here the remote contract then gains a
// field, and SetDep is what carries it across — the copy and the generated shim
// together.
func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox

	remote_module := "example/remote"
	app_module := "example/app"

	if err := lib.Actions.Start(api.StartProps{
		Path:        "TestDir/remote",
		ProjectName: "Remote",
		Module:      &remote_module,
	}); err != nil {
		panic(err)
	}

	write("TestDir/remote/sandbox/api/greeter.go", greeterBefore)

	if err := lib.Actions.Build(api.BuildProps{Path: "TestDir/remote", Runtime: api.RuntimeGo}); err != nil {
		panic(err)
	}

	if err := lib.Actions.Start(api.StartProps{
		Path:        "TestDir/app",
		ProjectName: "App",
		Module:      &app_module,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.DepsInit("TestDir/app"); err != nil {
		panic(err)
	}

	appendTo("TestDir/app/go.mod", "\nrequire example/remote v0.0.1\n\nreplace example/remote => ../remote\n")

	if err := lib.Actions.AddDep(api.AddDepProps{
		Path: "TestDir/app",
		Dep:  "example/remote",
		As:   "remote",
	}); err != nil {
		panic(err)
	}

	// The remote contract moves: one more field on the props it takes.
	write("TestDir/remote/sandbox/api/greeter.go", greeterAfter)

	if err := lib.Actions.Build(api.BuildProps{Path: "TestDir/remote", Runtime: api.RuntimeGo}); err != nil {
		panic(err)
	}

	if err := lib.Actions.SetDep(api.SetDepProps{
		Path:    "TestDir/app",
		Dep:     "remote",
		Version: "v0.0.1",
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the paths this example asserts, copied out of
	// TestDir. The cli side copies the same set.
	if err := os.CopyFS("AssertDir/sandbox/deps", os.DirFS("TestDir/app/sandbox/deps")); err != nil {
		panic(err)
	}
	if err := os.CopyFS("AssertDir/adapters", os.DirFS("TestDir/app/adapters")); err != nil {
		panic(err)
	}
}

// write puts one hand-written file of the remote repo in place, creating the
// directory it lives in.
func write(path string, content string) {
	if err := os.MkdirAll(dir(path), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}
}

// appendTo adds the require and replace the consumer needs to reach the repo
// beside it.
func appendTo(path string, content string) {
	existing, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, append(existing, []byte(content)...), 0o644); err != nil {
		panic(err)
	}
}

// dir is the directory part of a slash-separated path.
func dir(path string) string {
	for index := len(path) - 1; index >= 0; index-- {
		if path[index] == '/' {
			return path[:index]
		}
	}
	return "."
}

const greeterBefore = `package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) string
}

// GreetProps names who is greeted.
type GreetProps struct {
	Name string
}
`

const greeterAfter = `package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) string
}

// GreetProps names who is greeted and how loudly.
type GreetProps struct {
	Name string
	Loud bool
}
`
