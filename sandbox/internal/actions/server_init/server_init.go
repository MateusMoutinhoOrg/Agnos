package server_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_install"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// serverDeps are the contracts the server layer calls into: the socket itself,
// the three output channels, text conversion, sorting and the JSON codec the
// error body and the schema validator are written through.
var serverDeps = []string{"std", "stringsdeps", "sortdeps", "serializables", "serverdeps"}

// cliDep is the one further contract the implicit cli-init needs, installed
// only when this project has no cli layer yet.
const cliDep = "argvdeps"

// InstallDeps installs the contracts the server layer calls into, plus the one
// the implicit cli-init needs when the project has no cli layer yet.
//
// It is exported because a layer that composes ServerInitInternal into its own
// transaction — front-init does — still has to install this set first: the
// internal half renders assets and writes nothing to go.mod.
func InstallDeps(deps *deps.Deps, path string) error {
	install := serverDeps
	if !smartio.New(deps, path, config.ProjectName).IsDir(cliDir) {
		install = append(install, cliDep)
	}

	for _, dep := range install {
		if err := depInstallAction.DepInstall(deps, api.DepInstallProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	return nil
}

// ServerInit installs the deps the server layer depends on and renders the
// "server" asset group into the project, then runs build as a follow-up step.
// A server needs an entry point that starts it and that entry point is a
// command, so a project with no cli layer gets one first.
func ServerInit(deps *deps.Deps, path string) error {
	if err := InstallDeps(deps, path); err != nil {
		return err
	}

	io := smartio.New(deps, path, config.ProjectName)
	if err := ServerInitInternal(deps, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
