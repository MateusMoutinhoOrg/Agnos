package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// frontDeps are the contracts the front layer calls into on top of the ones
// the server layer already installs: the asset tree every file is read out
// of. Everything else frontio touches — text conversion above all — comes
// with the server.
var frontDeps = []string{"embeddeps"}

// FrontInit installs the deps the front layer depends on and turns the front
// mechanic on, then runs build as a follow-up step, which renders the group.
// The front is answered over http, so a project with no server layer is given one
// first — along with the deps that layer needs, which its internal half does
// not install.
func FrontInit(sandbox *api.Sandbox, path string) error {
	has_server, err := utils.ExtensionEnabled(sandbox, smartio.New(sandbox, path, sandbox.Config.ProjectName), utils.ExtensionSandboxServer)
	if err != nil {
		return err
	}
	if !has_server {
		if err := serverInitAction.InstallDeps(sandbox, path); err != nil {
			return err
		}
	}

	for _, dep := range frontDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := FrontInitInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
