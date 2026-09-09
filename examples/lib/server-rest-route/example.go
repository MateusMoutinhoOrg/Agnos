package main

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The server-rest-route example: a route whose last segment takes the rest of
// the path, so one declaration serves /static/a, /static/a/b.png and anything
// deeper.
//
// It calls the same actions `agnos server-init`, `agnos add-route` and
// `agnos add-segment --array` call, and writes only inside TestDir.
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

	if err := lib.Actions.AddRoute("TestDir", "static", "GET", "/static", "Serve a file under /static", "Files"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddSegment(api.RouteFieldProps{
		Path: "TestDir", Route: "static", Name: "rest", Type: "string", Array: true, Position: -1,
	}); err != nil {
		panic(err)
	}

	// What result.yaml records: the declaration this example wrote, the
	// []string field build generated from it, and the dispatch, which is where
	// the route stops asking for an exact segment count — the same set the cli
	// side copies.
	copy_out := map[string][]string{
		"sandbox/internal/routes/static": {"route.yaml", "entries.go"},
		"sandbox/internal/server":        {"servermain.go"},
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
