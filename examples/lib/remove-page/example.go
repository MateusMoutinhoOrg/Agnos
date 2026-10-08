package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The remove-page example: drop a page, its route and its html both.
//
// It calls the same actions `agnos front-init`, `agnos add-page` and
// `agnos remove-page` call, and writes only inside test-dir.
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

	if err := lib.Actions.FrontInit(api.FrontInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	for _, page := range []api.AddPageProps{
		{Path: "test-dir", Name: "about", Title: "About"},
		{Path: "test-dir", Name: "blog/post"},
	} {
		if err := lib.Actions.AddPage(page); err != nil {
			panic(err)
		}
	}

	if err := lib.Actions.RemovePage(api.RemovePageProps{Path: "test-dir", Name: "about"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	copyTree("test-dir/assets/front", "assert-dir/assets/front")
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
