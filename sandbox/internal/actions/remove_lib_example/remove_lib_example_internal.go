package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveLibExampleInternal deletes the example's directory and everything in
// it — the example.go, the golden result.yaml and any TestDir left behind by
// the last run.
func RemoveLibExampleInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	return utils.RemoveExample(deps, io, utils.ExampleLibSide, name)
}
