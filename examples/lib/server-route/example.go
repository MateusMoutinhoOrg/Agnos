package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The server-route example: add the server layer, declare a route and one
// field of each origin.
//
// It calls the same actions `agnos server-init`, `agnos add-route` and
// `agnos add-field` call, and writes only inside TestDir.
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

	fields := []api.RouteFieldProps{
		{Path: "TestDir", Route: "create-user", Name: "tenant", In: "path", Type: "string", Required: true, Position: -1},
		{Path: "TestDir", Route: "create-user", Name: "authorization", In: "header", Type: "string", Required: true, Position: -1},
		{Path: "TestDir", Route: "create-user", Name: "email", In: "body", Type: "string", Required: true, Format: "email", Position: -1},
	}
	for _, field := range fields {
		if err := lib.Actions.AddField(field); err != nil {
			panic(err)
		}
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
