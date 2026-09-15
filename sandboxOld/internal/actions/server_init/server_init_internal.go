package server_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/cli_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// cliDir is what tells a project that already has a cli layer from one that
// has to be given one first — the same directory `build` reads hasCli from.
const cliDir = "sandbox/internal/cli"

// startServerDir holds the one command the server layer writes: the entry
// point that opens the port.
const startServerDir = "sandbox/internal/commands/start_server"

// ServerInitInternal renders every embedded asset under assets/server into the
// target project at the path it holds inside that group, and writes the
// start-server command beside it.
//
// A project with no cli layer is given one first, on this same open SmartIO:
// actions compose by sharing one transaction, so there is no intermediate
// Persist and no intermediate build between the two halves.
func ServerInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("server-init started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	if !io.IsDir(cliDir) {
		if err := cliInitAction.CliInitInternal(sandbox, io, path); err != nil {
			return err
		}
	}

	if err := utils.RenderGroup(sandbox, io, "server", vars); err != nil {
		return err
	}

	return writeStartServer(sandbox, io, vars)
}

// writeStartServer scaffolds the start-server command, leaving an existing one
// alone: like a route's handler.go, it is written once and then the project's.
func writeStartServer(sandbox *api.Sandbox, io *smartio.SmartIO, vars map[string]interface{}) error {
	if io.IsDir(startServerDir) {
		sandbox.Deps.Std.Log("server-init: %s already exists, keeping it \n", startServerDir)
		return nil
	}

	if err := utils.RenderTemplateToDest(sandbox, io, "templates/start_server_entries.yaml", vars, startServerDir+"/entries.yaml"); err != nil {
		return err
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/start_server_handler.go", vars, startServerDir+"/handler.go")
}
