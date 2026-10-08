package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The front-purge example: remove the front layer, keeping assets/front.
//
// It calls the same actions `agnos front-init`, `agnos add-page` and
// `agnos front-purge` call, and writes only inside test-dir.
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

	if err := lib.Actions.AddPage(api.AddPageProps{
		Path: "test-dir", Name: "about", Title: "About",
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.FrontPurge(api.FrontPurgeProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/routes",
		"assets/front",
	} {
		copyTree("test-dir/"+dir, "assert-dir/"+dir)
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

// copyExtensions puts the project's extensions.yaml into assert-dir at the
// place it holds in the tree.
func copyExtensions() error {
	content, err := os.ReadFile("test-dir/AgnosConfig/extensions.yaml")
	if err != nil {
		return err
	}
	if err := os.MkdirAll("assert-dir/AgnosConfig", 0o755); err != nil {
		return err
	}
	return os.WriteFile("assert-dir/AgnosConfig/extensions.yaml", content, 0o644)
}
