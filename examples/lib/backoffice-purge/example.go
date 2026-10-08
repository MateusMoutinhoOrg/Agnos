package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The backoffice-purge example: remove the admin backoffice backoffice-init
// wrote.
//
// It calls the same actions `agnos backoffice-init` and `agnos
// backoffice-purge` call, and writes only inside test-dir.
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

	if err := lib.Actions.BackofficeInit(api.BackofficeInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.BackofficePurge(api.BackofficePurgeProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/routeprops",
		"sandbox/internal/routes",
		"sandbox/internal/commands",
	} {
		copyTree("test-dir/"+dir, "assert-dir/"+dir)
	}
	for _, file := range []string{
		"AgnosConfig/extensions.yaml",
		"sandbox/api/config.go",
	} {
		copyFile("test-dir/"+file, "assert-dir/"+file)
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
		copyFile(path, filepath.Join(dest, relative))
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// copyFile copies one file to dest, creating the directories it needs.
func copyFile(source string, dest string) {
	content, err := os.ReadFile(source)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(dest, content, 0o644); err != nil {
		panic(err)
	}
}
