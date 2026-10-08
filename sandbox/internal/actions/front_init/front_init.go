package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// frontDeps are the contracts the front layer calls into on top of the ones
// the server layer already installs: the asset tree every file is read out
// of, and OpinionatedAgnosFront — the file layer itself, last because its
// contract imports the first.
var frontDeps = []string{"embeddeps", utils.OpinionatedAgnosFront}

// FrontInit installs the deps the front layer depends on and turns the front
// mechanic on, then runs build as a follow-up step, which renders the group.
// The front is answered over http, so a project with no server layer is given one
// first — along with the deps that layer needs, which its internal half does
// not install.
func FrontInit(sandbox *api.Sandbox, props api.FrontInitProps) error {
	has_server, err := utils.ExtensionEnabled(sandbox, stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName), utils.ExtensionServer)
	if err != nil {
		return err
	}
	if !has_server {
		if err := serverInitAction.InstallDeps(sandbox, props.Path); err != nil {
			return err
		}
	}

	if err := InstallDeps(sandbox, props.Path); err != nil {
		return err
	}

	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := FrontInitInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}

// InstallDeps installs the contracts the front layer calls into on top of the
// server's. It is exported for the same reason server-init exports its own: a
// layer that composes FrontInitInternal into its own transaction — the
// backoffice does — still has to install this set first.
func InstallDeps(sandbox *api.Sandbox, path string) error {
	for _, dep := range frontDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}
	return nil
}
