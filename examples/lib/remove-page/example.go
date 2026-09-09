package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-page example: drop a page, its route and its html both.
//
// It calls the same actions `agnos front-init`, `agnos add-page` and
// `agnos remove-page` call, and writes only inside TestDir.
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

	for _, page := range []api.PageProps{
		{Path: "TestDir", Name: "home", Trigger: "/", Title: "Home"},
		{Path: "TestDir", Name: "about", Title: "About"},
	} {
		if err := lib.Actions.AddPage(page); err != nil {
			panic(err)
		}
	}

	if err := lib.Actions.RemovePage("TestDir", "about"); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/routes",
		"assets/frontend/pages",
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
}
