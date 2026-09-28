package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListRoutesInternal is every declared route as one line, in run order: the
// rung, the methods, the pattern and the name — the chain a request walks
// until one route answers.
func ListRoutesInternal(sandbox *api.Sandbox, io *smartio.SmartIO) ([]string, error) {
	if err := utils.RequireProject(sandbox, io); err != nil {
		return nil, err
	}
	if !io.IsDir(utils.RoutesDir) {
		return nil, sandbox.Deps.Std.Errorf("the project has no server layer: run server-init first")
	}

	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{sandbox.Deps.Std.Sprintf("%-4s %-9s %-8s %-36s %s", "#", "priority", "methods", "pattern", "route")}
	index := 0
	for _, entry := range chain {
		index++
		lines = append(lines, sandbox.Deps.Std.Sprintf("%-4d %-9d %-8s %-36s %s",
			index, entry.Conf.Priority, sandbox.Deps.Stringsdeps.Join(entry.Conf.Methods, ","), entry.Conf.Pattern(),
			utils.RouteIdentifier(sandbox, entry.Name)))
	}
	return lines, nil
}
