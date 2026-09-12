package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCliExampleInternal deletes the example's directory and everything in
// it — the example.sh, the golden result.yaml and any TestDir or AssertDir left
// behind by the last run.
func RemoveCliExampleInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	return utils.RemoveExample(sandbox, io, utils.ExampleCliSide, name)
}
