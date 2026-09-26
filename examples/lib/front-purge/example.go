package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The front-purge example: remove the front layer, keeping assets/frontend.
//
// It calls the same actions `agnos front-init`, `agnos add-page` and
// `agnos front-purge` call, and writes only inside TestDir.
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

	if err := lib.Actions.FrontInit("TestDir"); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddPage(api.PageProps{
		Path: "TestDir", Name: "about", Title: "About",
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.FrontPurge("TestDir"); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/routeslist",
		"assets/frontend",
	} {
		copyTree("TestDir/"+dir, "AssertDir/"+dir)
	}
}

// copyTree copies every file under source into dest, keeping the place each
// one holds in the tree.
func copyTree(source string, dest string) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dest, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0o644)
	})
	if err != nil {
		panic(err)
	}

	// The declaration the pair wrote: this is the whole of what tells the
	// build the mechanic is on or off from here.
	if err := copyExtensions(); err != nil {
		panic(err)
	}
}

// copyExtensions puts the project's extensions.yaml into AssertDir at the
// place it holds in the tree.
func copyExtensions() error {
	content, err := os.ReadFile("TestDir/AgnosConfig/extensions.yaml")
	if err != nil {
		return err
	}
	if err := os.MkdirAll("AssertDir/AgnosConfig", 0o755); err != nil {
		return err
	}
	return os.WriteFile("AssertDir/AgnosConfig/extensions.yaml", content, 0o644)
}
