package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveLibExampleInternal deletes the example's directory and everything in
// it — the example.go, the golden result.yaml and any test-dir or assert-dir left
// behind by the last run.
func RemoveLibExampleInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionExample); err != nil {
		return err
	}

	return utils.RemoveExample(sandbox, io, utils.ExampleLibSide, name)
}
