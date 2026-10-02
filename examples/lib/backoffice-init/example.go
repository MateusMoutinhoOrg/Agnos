package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The backoffice-init example: add the admin backoffice to a project that has
// none of the layers it stands on.
//
// It calls the same action `agnos backoffice-init` calls, and writes only
// inside TestDir.
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

	if err := lib.Actions.BackofficeInit("TestDir"); err != nil {
		panic(err)
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/routeprops",
		"sandbox/internal/commands/backoffice",
		"sandbox/internal/databases/backofficedb",
		"sandbox/internal/server/backoffice/backofficeauth",
		"docs/Backoffice",
	} {
		copyTree("TestDir/"+dir, "AssertDir/"+dir)
	}
	for _, file := range []string{
		".gitignore",
		"AgnosConfig/extensions.yaml",
		"sandbox/api/config.go",
		"sandbox/api/backofficeconfig.go",
	} {
		copyFile("TestDir/"+file, "AssertDir/"+file)
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
