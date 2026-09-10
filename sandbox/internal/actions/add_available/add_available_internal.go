package add_available

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddAvailableInternal declares one further available: a second answer to
// "which adapter wins for each field", for a build that swaps one
// implementation — a lambda entry point for the http server, say.
//
// It starts as a copy of the standard available's selection, not empty: every
// available has to fill every field, so an empty one would be a tree that
// fails `verify` the moment it is written. Point it at another adapter with
// `set-adapter --available <name>`.
func AddAvailableInternal(deps *deps.Deps, io *smartio.SmartIO, path string, available string) error {
	deps.Std.Log("add-available started with path %s available %s \n", path, available)

	if err := utils.ValidateAvailableName(deps, available); err != nil {
		return err
	}

	if io.IsDir(utils.AvailableDir(available)) {
		return deps.Std.Errorf("available %q already exists", available)
	}

	conf, err := utils.LoadAvailableConf(deps, io, utils.StandardAvailable)
	if err != nil {
		return deps.Std.Errorf("no %s available to copy the selection from: run `agnos deps-init` first", utils.StandardAvailable)
	}

	return io.WriteFile(utils.AvailableConfPath(available), []byte(conf.Render()))
}
