package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The server-route example: add the server layer, declare a route and one
// entry of every place its declaration holds something.
//
// It calls the same actions `agnos server-init`, `agnos add-route`,
// `agnos add-path`, `agnos add-parameter`, `agnos set-body` and
// `agnos add-body-field` call, and writes only inside test-dir.
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

	if err := lib.Actions.ServerInit(api.ServerInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddRoute(api.AddRouteProps{
		Path:        "test-dir",
		Name:        "create-user",
		Methods:     []string{"POST"},
		Trigger:     "/users/",
		TriggerType: "prefix",
		Summary:     "Create a user under a tenant",
		Category:    "Users",
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddPath(api.AddPathProps{
		Path: "test-dir", Route: "create-user", Name: "tenant", Start: "1", End: "1", Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddParameter(api.AddParameterProps{
		Path: "test-dir", Route: "create-user", Name: "authorization", Sources: []string{"header"}, Required: true, Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddParameter(api.AddParameterProps{
		Path: "test-dir", Route: "create-user", Name: "page", Type: "number", Default: "1", Position: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.SetBody(api.SetBodyProps{
		Path: "test-dir", Route: "create-user", Type: "json", Required: true, MaxBytes: -1,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddBodyField(api.AddBodyFieldProps{
		Path: "test-dir", Route: "create-user", Name: "email", Type: "string", Required: true, Format: "email",
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the declaration this example wrote and the
	// api.Route and Input build generated from it — the same set the cli
	// side copies.
	assert_dir := "assert-dir/sandbox/internal/routes/create_user"
	if err := os.MkdirAll(assert_dir, 0o755); err != nil {
		panic(err)
	}
	for _, file := range []string{"route.yaml", "new.go", "input.go"} {
		content, err := os.ReadFile("test-dir/sandbox/internal/routes/create_user/" + file)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(assert_dir+"/"+file, content, 0o644); err != nil {
			panic(err)
		}
	}
}
