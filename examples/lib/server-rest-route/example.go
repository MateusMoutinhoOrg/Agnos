package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The server-rest-route example: a route whose second path reads every segment
// after the mount, so one declaration serves /static/a, /static/a/b.png and
// anything deeper.
//
// It calls the same actions `agnos server-init`, `agnos add-route`,
// `agnos set-path` and `agnos add-path` call, and writes only inside TestDir.
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

	if err := lib.Actions.AddRoute(api.AddRouteProps{
		Path:     "TestDir",
		Name:     "static",
		Methods:  []string{"GET"},
		Trigger:  "/static",
		Help:     "Serve a file under /static",
		Category: "Files",
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.SetPath(api.RoutePathEditProps{
		Path: "TestDir", Route: "static", Id: "Route", End: "0",
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddPath(api.RoutePathProps{
		Path: "TestDir", Route: "static", Id: "rest", Start: "1", End: "-1", Position: -1,
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the declaration this example wrote and the
	// api.Route and Entries build generated from it: the first path fixes
	// segment 0, the second reads segment 1 to the last — the same set the
	// cli side copies.
	copy_out := map[string][]string{
		"sandbox/internal/routeslist/static": {"route.yaml", "new.go", "entries.go"},
	}
	for dir, files := range copy_out {
		if err := os.MkdirAll("AssertDir/"+dir, 0o755); err != nil {
			panic(err)
		}
		for _, file := range files {
			content, err := os.ReadFile("TestDir/" + dir + "/" + file)
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile("AssertDir/"+dir+"/"+file, content, 0o644); err != nil {
				panic(err)
			}
		}
	}
}
