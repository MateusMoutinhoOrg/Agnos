package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListRoutesInternal is every declared route as one line, in run order: the
// rung, the methods, the pattern and the name — the chain a request walks
// until one route answers.
func ListRoutesInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS) ([]string, error) {
	if err := utils.RequireProject(sandbox, io); err != nil {
		return nil, err
	}
	if !io.IsDir(utils.RoutesDir) {
		return nil, sandbox.Deps.StdDeps.Errorf("the project has no server layer: run server-init first")
	}

	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{sandbox.Deps.StdDeps.Sprintf("%-4s %-9s %-8s %-36s %s", "#", "priority", "methods", "pattern", "route")}
	index := 0
	for _, entry := range chain {
		index++
		lines = append(lines, sandbox.Deps.StdDeps.Sprintf("%-4d %-9d %-8s %-36s %s",
			index, entry.Conf.Priority, sandbox.Deps.StringsDeps.Join(entry.Conf.Methods, ","), entry.Conf.Pattern(),
			utils.RouteName(sandbox, entry.Name)))
	}
	return lines, nil
}
