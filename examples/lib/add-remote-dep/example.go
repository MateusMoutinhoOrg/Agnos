package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The add-remote-dep example: install another agnos repo as a dep
//
// It calls the same action `agnos add-dep <module>` calls, and writes only
// inside TestDir. TestDir/remote is the repo being installed — every agnos repo
// is installable by construction — and TestDir/app is the consumer. A published
// repo is pinned with "<module>@<version>"; here the two live side by side,
// wired with a `replace`, which is how a pair of repos is developed together.
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

	write("TestDir/remote/sandbox/api/greeter.go", greeterApi)
	write("TestDir/remote/sandbox/binds/greeter.go", greeterBinds)
	write("TestDir/remote/sandbox/internal/greeter/greeter.go", greeterInternal)

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

const greeterApi = `package api

// Greeter is what this repo publishes.
type Greeter struct {
	// Greet returns the greeting for one name.
	Greet func(props GreetProps) Greeting
}

// GreetProps names who is greeted and how loudly.
type GreetProps struct {
	Name string
	Loud bool
}

// Greeting is what came back.
type Greeting struct {
	Text  string
	Words int
}
`

const greeterBinds = `package binds

import (
	api "example/remote/sandbox/api"
	greeter "example/remote/sandbox/internal/greeter"
)

func GreeterBind(sandbox *api.Sandbox) {
	sandbox.Greeter.Greet = func(props api.GreetProps) api.Greeting {
		return greeter.Greet(props)
	}
}
`

const greeterInternal = `package greeter

import (
	api "example/remote/sandbox/api"
)

// Greet builds the greeting for one name.
func Greet(props api.GreetProps) api.Greeting {
	text := "hello " + props.Name
	if props.Loud {
		text = text + "!"
	}
	return api.Greeting{Text: text, Words: 2}
}
`
