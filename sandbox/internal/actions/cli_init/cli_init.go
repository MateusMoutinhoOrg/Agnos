package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CliDeps are the contracts the cli layer calls into: the three output
// channels and the argv parser the dispatch is handed, the text conversion the
// help command is written through, and OpinatedAgnosCli — the dispatch
// itself, last because its contract imports the first two. It is exported for
// the layer that composes CliInitInternal into its own transaction —
// server-init does — and still has to install this set first.
var CliDeps = []string{"std", "argvdeps", "stringsdeps", utils.OpinatedAgnosCli}

// CliInit installs the deps the cli layer depends on and renders the "cli"
// asset group into the project, then runs build as a follow-up step.
func CliInit(sandbox *api.Sandbox, path string) error {
	for _, dep := range CliDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := CliInitInternal(sandbox, io, path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
