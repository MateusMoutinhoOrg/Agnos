package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CliDeps are the contracts the cli layer calls into: the three output
// channels and the argv parser the dispatch is handed, the text conversion the
// help command is written through, and OpinionatedAgnosCli — the dispatch
// itself, last because its contract imports the first two. It is exported for
// the layer that composes CliInitInternal into its own transaction —
// server-init does — and still has to install this set first.
var CliDeps = []string{"stddeps", "argvdeps", "stringsdeps", utils.OpinionatedAgnosCli}

// CliInit installs the deps the cli layer depends on and renders the "cli"
// asset group into the project, then runs build as a follow-up step.
func CliInit(sandbox *api.Sandbox, props api.CliInitProps) error {
	for _, dep := range CliDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: props.Path, Dep: dep}); err != nil {
			return err
		}
	}
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := CliInitInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
