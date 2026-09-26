package rebalance_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// healthRoute is the route build renders itself.
const healthRoute = "health"

// RebalanceRoutesInternal gives every route a rung of its own, Step apart and
// in the order the chain runs them now: the first on Step, the next on twice
// Step, and so on — each phase on its own ladder. Two routes that shared a
// rung ran by name, and still do, each on a rung of its own; what changes is
// that --before and --after have room again.
func RebalanceRoutesInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RebalanceRoutesProps) error {
	if props.Step < 1 {
		return sandbox.Deps.Std.Errorf("--step %d is below 1: the rungs have to be apart", props.Step)
	}
	if !io.IsDir(utils.RoutesDir) {
		return sandbox.Deps.Std.Errorf("the project has no server layer: run server-init first")
	}

	chain, err := utils.LoadRouteChain(sandbox, io)
	if err != nil {
		return err
	}

	rung := map[string]int{routeconf.PhaseBefore: 0, routeconf.PhaseAfter: 0}
	for _, entry := range chain {
		if entry.Name == healthRoute {
			continue
		}
		rung[entry.Conf.Phase] += props.Step
		if entry.Conf.Priority == rung[entry.Conf.Phase] {
			continue
		}
		sandbox.Deps.Std.Log("rebalance-routes %s: %d -> %d \n", utils.RouteIdentifier(sandbox, entry.Name), entry.Conf.Priority, rung[entry.Conf.Phase])
		entry.Conf.Priority = rung[entry.Conf.Phase]
		if err := utils.SaveRouteConf(sandbox, io, entry.Name, entry.Conf); err != nil {
			return err
		}
	}
	return nil
}
