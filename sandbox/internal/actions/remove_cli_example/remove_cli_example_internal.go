package remove_cli_example

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCliExampleInternal deletes the example's directory and everything in
// it — the example.sh, the golden result.yaml and any TestDir left behind by
// the last run.
func RemoveCliExampleInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	return utils.RemoveExample(deps, io, utils.ExampleCliSide, name)
}
