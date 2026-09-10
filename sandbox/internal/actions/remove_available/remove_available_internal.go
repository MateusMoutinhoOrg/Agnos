package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveAvailableInternal deletes one available, selection and generated New()
// together. The standard one is refused: cmd/main/main.go imports it by name,
// so removing it is removing the entry point's constructor.
func RemoveAvailableInternal(deps *deps.Deps, io *smartio.SmartIO, path string, available string) error {
	deps.Std.Log("remove-available started with path %s available %s \n", path, available)

	if available == utils.StandardAvailable {
		return deps.Std.Errorf("the %s available cannot be removed: cmd/main/main.go imports it", utils.StandardAvailable)
	}

	if !io.IsDir(utils.AvailableDir(available)) {
		return deps.Std.Errorf("available %q does not exist", available)
	}

	utils.RemoveTree(deps, io, []string{utils.AvailableDir(available)})
	return nil
}
