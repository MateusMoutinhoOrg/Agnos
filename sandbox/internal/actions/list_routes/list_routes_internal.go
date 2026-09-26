package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListRoutesInternal is every declared route as one line, in run order: the
// rung, the methods, the pattern and the name. The `before` phase comes first
// — the chain a request walks until one route answers — and the `after` phase
// beneath it, which runs once it has.
func ListRoutesInternal(sandbox *api.Sandbox, io *smartio.SmartIO) ([]string, error) {
	if !io.IsDir(utils.RoutesDir) {
		return nil, sandbox.Deps.Std.Errorf("the project has no server layer: run server-init first")
	}

	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{sandbox.Deps.Std.Sprintf("%-4s %-9s %-8s %-36s %s", "#", "priority", "methods", "pattern", "route")}
	phase := routeconf.PhaseBefore
	index := 0
	for _, entry := range chain {
		if entry.Conf.Phase != phase {
			phase, index = entry.Conf.Phase, 0
			lines = append(lines, "", "after the answer")
		}
		index++
		lines = append(lines, sandbox.Deps.Std.Sprintf("%-4d %-9d %-8s %-36s %s",
			index, entry.Conf.Priority, sandbox.Deps.Stringsdeps.Join(entry.Conf.Methods, ","), entry.Conf.Pattern(),
			utils.RouteIdentifier(sandbox, entry.Name)))
	}
	return lines, nil
}
