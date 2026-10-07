package database_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DatabaseInitInternal turns the database mechanic on in the project's
// declaration. The code every generated methods.go shares is the
// OpinatedAgnosDatabase lib, installed with the store, and the databases
// themselves are rendered by the follow-up build.
//
// It scaffolds no database of its own: a database is a declaration, and which
// tables a project wants is not something an init can guess. "Declare its
// first one" is the step that follows.
func DatabaseInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("database-init started with path %s \n", path)

	if !io.IsDir(utils.ContractsDir + "/" + DatabaseDep) {
		return sandbox.Deps.Std.Errorf(
			"the database layer needs the %s contract: install it with `agnos add-dep %s --as %s`",
			DatabaseDep, KeepModule, DatabaseDep)
	}

	io.CreateDir(utils.DatabasesDir)

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxDatabase, true)
}
