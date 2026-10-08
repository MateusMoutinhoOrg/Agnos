package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The front-init example: add the html front layer to a project that has none.
//
// It calls the same action `agnos front-init` calls, and writes only inside
// test-dir.
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

	// What result.yaml records: the same set the cli side copies — the file
	// layer's contract, the route serving assets/front, and that tree's
	// index.html.
	for _, dir := range []string{
		"sandbox/deps/OpinionatedAgnosFront",
		"sandbox/internal/routes/front",
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
