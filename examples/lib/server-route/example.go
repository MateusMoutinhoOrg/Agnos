package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The server-route example: add the server layer, declare a route and one
// entry of every place its declaration holds something.
//
// It calls the same actions `agnos server-init`, `agnos add-route`,
// `agnos add-segment`, `agnos add-header`, `agnos add-param`, `agnos set-body`
// and `agnos add-body-field` call, and writes only inside TestDir.
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

	if err := lib.Actions.ServerInit("TestDir"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddRoute("TestDir", "create-user", "POST", "/users", "Create a user under a tenant", "Users"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddSegment(api.RouteFieldProps{
		Path: "TestDir", Route: "create-user", Name: "tenant", Type: "string", Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddHeader(api.RouteFieldProps{
		Path: "TestDir", Route: "create-user", Name: "authorization", Type: "string", Required: true, Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddParam(api.RouteFieldProps{
		Path: "TestDir", Route: "create-user", Name: "page", Type: "int", Default: "1", Min: "1", Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.SetBody(api.RouteBodyProps{
		Path: "TestDir", Route: "create-user", Type: "json", Required: true, MaxBytes: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddBodyField(api.RouteBodyFieldProps{
		Path: "TestDir", Route: "create-user", Name: "email", Type: "string", Required: true, Format: "email",
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the declaration this example wrote and the
	// struct build generated from it — the same set the cli side copies.
	assert_dir := "AssertDir/sandbox/internal/routes/create_user"
	if err := os.MkdirAll(assert_dir, 0o755); err != nil {
		panic(err)
	}
	for _, file := range []string{"route.yaml", "entries.go"} {
		content, err := os.ReadFile("TestDir/sandbox/internal/routes/create_user/" + file)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(assert_dir+"/"+file, content, 0o644); err != nil {
			panic(err)
		}
	}
}
