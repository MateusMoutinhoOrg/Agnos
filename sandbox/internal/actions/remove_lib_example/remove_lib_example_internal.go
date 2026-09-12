package remove_lib_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveLibExampleInternal deletes the example's directory and everything in
// it — the example.go, the golden result.yaml and any TestDir or AssertDir left
// behind by the last run.
func RemoveLibExampleInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	return utils.RemoveExample(sandbox, io, utils.ExampleLibSide, name)
}
