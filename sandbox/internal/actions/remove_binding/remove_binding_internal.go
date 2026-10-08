package remove_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveAvailableInternal deletes one available, selection and generated New()
// together. The standard one is refused: cmd/main/main.go imports it by name,
// so removing it is removing the entry point's constructor.
func RemoveAvailableInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string, available string) error {
	sandbox.Deps.Std.Log("remove-available started with path %s available %s \n", path, available)

	if available == utils.StandardAvailable {
		return sandbox.Deps.Std.Errorf("the %s available cannot be removed: cmd/main/main.go imports it", utils.StandardAvailable)
	}

	if !io.IsDir(utils.AvailableDir(available)) {
		return sandbox.Deps.Std.Errorf("available %q does not exist", available)
	}

	utils.RemoveTree(sandbox, io, []string{utils.AvailableDir(available)})
	return nil
}
