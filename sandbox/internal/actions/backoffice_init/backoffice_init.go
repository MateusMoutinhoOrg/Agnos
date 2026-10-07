package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// BackofficeInit installs the layers the backoffice stands on — server, front
// and database, each one only when it is off — then writes the backoffice and
// turns its mechanic on, and runs build as a follow-up step.
//
// The layers' own deps go in first, outside the transaction, the way
// front-init installs the server's: their internal halves write nothing to
// go.mod, and database-init's store is a remote dep fetched as it installs.
func BackofficeInit(sandbox *api.Sandbox, path string) error {
	probe := smartio.New(sandbox, path, sandbox.Config.ProjectName)

	has_server, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionSandboxServer)
	if err != nil {
		return err
	}
	if !has_server {
		if err := serverInitAction.InstallDeps(sandbox, path); err != nil {
			return err
		}
	}

	has_front, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionSandboxFront)
	if err != nil {
		return err
	}
	if !has_front {
		if err := frontInitAction.InstallDeps(sandbox, path); err != nil {
			return err
		}
	}

	has_database, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionSandboxDatabase)
	if err != nil {
		return err
	}
	if !has_database {
		if err := databaseInitAction.InstallDeps(sandbox, path); err != nil {
			return err
		}
	}

	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := BackofficeInitInternal(sandbox, io, path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
