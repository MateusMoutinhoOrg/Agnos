package server_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// serverDeps are the contracts the server layer calls into: the socket itself,
// the three output channels, text conversion, sorting, the JSON codec the
// error body and the schema validator are written through, the reflection
// that fills a route's Entries, and the signal a graceful shutdown waits on.
var serverDeps = []string{"std", "stringsdeps", "sortdeps", "serializables", "serverdeps", "reflectdeps", "signaldeps"}

// cliDep is the one further contract the implicit cli-init needs, installed
// only when this project has no cli layer yet.
const cliDep = "argvdeps"

// InstallDeps installs the contracts the server layer calls into, plus the one
// the implicit cli-init needs when the project has no cli layer yet.
//
// It is exported because a layer that composes ServerInitInternal into its own
// transaction — front-init does — still has to install this set first: the
// internal half renders assets and writes nothing to go.mod.
func InstallDeps(sandbox *api.Sandbox, path string) error {
	has_cli, err := utils.ExtensionEnabled(sandbox, smartio.New(sandbox, path, sandbox.Config.ProjectName), utils.ExtensionSandboxCli)
	if err != nil {
		return err
	}

	install := serverDeps
	if !has_cli {
		install = append(install, cliDep)
	}

	for _, dep := range install {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	return nil
}

// ServerInit installs the deps the server layer depends on and turns the server
// mechanic on, then runs build as a follow-up step, which renders the group.
// A server needs an entry point that starts it and that entry point is a
// command, so a project with no cli layer gets one first.
func ServerInit(sandbox *api.Sandbox, path string) error {
	if err := InstallDeps(sandbox, path); err != nil {
		return err
	}

	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := ServerInitInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
