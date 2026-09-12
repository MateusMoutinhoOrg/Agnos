package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// frontDeps are the contracts the front layer calls into on top of the ones
// the server layer already installs: the asset tree a page and its static
// files are read out of, the engine that executes a page, and the digest every
// generated link is stamped with. Everything else pageio touches — the output
// channels, text conversion and sorting — comes with the server.
var frontDeps = []string{"embeddeps", "templatedeps", "hashdeps"}

// FrontInit installs the deps the front layer depends on and renders the
// "front" asset group into the project, then runs build as a follow-up step. A
// page is answered over http, so a project with no server layer is given one
// first — along with the deps that layer needs, which its internal half does
// not install.
func FrontInit(sandbox *api.Sandbox, path string) error {
	if !smartio.New(sandbox, path, config.ProjectName).IsDir(serverDir) {
		if err := serverInitAction.InstallDeps(sandbox, path); err != nil {
			return err
		}
	}

	for _, dep := range frontDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	io := smartio.New(sandbox, path, config.ProjectName)
	if err := FrontInitInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
