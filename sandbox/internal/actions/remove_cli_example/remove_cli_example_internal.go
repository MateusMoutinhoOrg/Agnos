package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCliExampleInternal deletes the example's directory and everything in
// it — the example.sh, the golden result.yaml and any test-dir or assert-dir left
// behind by the last run.
func RemoveCliExampleInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionExample); err != nil {
		return err
	}

	return utils.RemoveExample(sandbox, io, utils.ExampleCliSide, name)
}
