package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// cliDeps are the contracts the cli layer calls into: the three output
// channels, the argv parser, text conversion, and the reflection that fills a
// command's Entries.
var cliDeps = []string{"std", "argvdeps", "stringsdeps", "reflectdeps"}

// CliInit installs the deps the cli layer depends on and renders the "cli"
// asset group into the project, then runs build as a follow-up step.
func CliInit(sandbox *api.Sandbox, path string) error {
	for _, dep := range cliDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := CliInitInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
