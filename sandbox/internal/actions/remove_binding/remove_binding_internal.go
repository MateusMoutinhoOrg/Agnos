package remove_binding

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveBindingInternal deletes one binding, selection and generated New()
// together. The standard one is refused: cmd/main/main.go imports it by name,
// so removing it is removing the entry point's constructor.
func RemoveBindingInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string, binding string) error {
	sandbox.Deps.StdDeps.Logf("remove-binding started with path %s binding %s \n", path, binding)

	if binding == utils.StandardBinding {
		return sandbox.Deps.StdDeps.Errorf("the %s binding cannot be removed: cmd/main/main.go imports it", utils.StandardBinding)
	}

	if !io.IsDir(utils.BindingDir(binding)) {
		return sandbox.Deps.StdDeps.Errorf("binding %q does not exist", binding)
	}

	utils.RemoveTree(sandbox, io, []string{utils.BindingDir(binding)})
	return nil
}
