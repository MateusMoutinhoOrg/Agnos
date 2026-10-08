package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// BackofficeInit installs the layers the backoffice stands on — server, front
// and database, each one only when it is off — then writes the backoffice and
// turns its mechanic on, and runs build as a follow-up step.
//
// The layers' own deps go in first, outside the transaction, the way
// front-init installs the server's: their internal halves write nothing to
// go.mod, and database-init's store is a remote dep fetched as it installs.
func BackofficeInit(sandbox *api.Sandbox, props api.BackofficeInitProps) error {
	probe := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)

	has_server, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionServer)
	if err != nil {
		return err
	}
	if !has_server {
		if err := serverInitAction.InstallDeps(sandbox, props.Path); err != nil {
			return err
		}
	}

	has_front, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionFront)
	if err != nil {
		return err
	}
	if !has_front {
		if err := frontInitAction.InstallDeps(sandbox, props.Path); err != nil {
			return err
		}
	}

	has_database, err := utils.ExtensionEnabled(sandbox, probe, utils.ExtensionDatabase)
	if err != nil {
		return err
	}
	if !has_database {
		if err := databaseInitAction.InstallDeps(sandbox, props.Path); err != nil {
			return err
		}
	}

	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := BackofficeInitInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
