package add_binding

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddBindingInternal declares one further binding: a second answer to
// "which adapter wins for each field", for a build that swaps one
// implementation — a lambda entry point for the http server, say.
//
// It starts as a copy of the standard binding's selection, not empty: every
// binding has to fill every field, so an empty one would be a tree that
// fails `verify` the moment it is written. Point it at another adapter with
// `set-adapter --binding <name>`.
func AddBindingInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string, binding string) error {
	sandbox.Deps.StdDeps.Logf("add-binding started with path %s binding %s \n", path, binding)

	if err := utils.ValidateBindingName(sandbox, binding); err != nil {
		return err
	}

	if io.IsDir(utils.BindingDir(binding)) {
		return sandbox.Deps.StdDeps.Errorf("binding %q already exists", binding)
	}

	conf, err := utils.LoadBindingConf(sandbox, io, utils.StandardBinding)
	if err != nil {
		return sandbox.Deps.StdDeps.Errorf("no %s binding to copy the selection from: run `agnos deps-init` first", utils.StandardBinding)
	}

	return io.CreateFile(utils.BindingConfPath(binding), []byte(conf.Render()))
}
