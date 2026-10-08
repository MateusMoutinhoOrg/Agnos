package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RouteChainEntry is one declared route as the chain lays it down: its package
// name, the directory it sits in and its parsed declaration.
type RouteChainEntry struct {
	Name string
	Dir  string
	Conf *routeconf.RouteConf
}

// LoadRouteChain reads every route.yaml under sandbox/internal/routes and
// returns them in the order the dispatch runs them: by priority, lowest first,
// then by name — the order the build collector lays Server.Routes down in.
func LoadRouteChain(sandbox *api.Sandbox, io *stagedfs.StagedFS) ([]RouteChainEntry, error) {
	chain := []RouteChainEntry{}

	for _, unit := range RouteDirs(sandbox, io) {
		content, err := io.ReadFile(unit.Dir + "/" + RouteConfFile)
		if err != nil {
			continue
		}
		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("%s/%s: %w", unit.Dir, RouteConfFile, err)
		}
		chain = append(chain, RouteChainEntry{Name: unit.Name, Dir: unit.Dir, Conf: conf})
	}

	SortRouteChain(sandbox, chain)
	return chain, nil
}

// SortRouteChain puts a chain in run order, in place.
func SortRouteChain(sandbox *api.Sandbox, chain []RouteChainEntry) {
	sandbox.Deps.SortDeps.SliceStable(chain, func(i int, j int) bool {
		left, right := chain[i].Conf, chain[j].Conf
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		return chain[i].Name < chain[j].Name
	})
}

// RouteRelativePriority is the rung one rung below (before) or above (after)
// the route named, for --before and --after. At most one of the two is given;
// it reports false when neither is.
func RouteRelativePriority(sandbox *api.Sandbox, io *stagedfs.StagedFS, before string, after string) (int, bool, error) {
	before = sandbox.Deps.StringsDeps.TrimSpace(before)
	after = sandbox.Deps.StringsDeps.TrimSpace(after)

	if before != "" && after != "" {
		return 0, false, sandbox.Deps.StdDeps.Errorf("--before and --after exclude each other")
	}
	if before == "" && after == "" {
		return 0, false, nil
	}

	if before != "" {
		other, err := LoadRouteConf(sandbox, io, before)
		if err != nil {
			return 0, false, err
		}
		if other.Priority == 0 {
			return 0, false, sandbox.Deps.StdDeps.Errorf(
				"--before %s: it runs on rung 0, and nothing runs below it — spread the chain with rebalance-routes first", RouteName(sandbox, before))
		}
		return other.Priority - 1, true, nil
	}

	other, err := LoadRouteConf(sandbox, io, after)
	if err != nil {
		return 0, false, err
	}
	return other.Priority + 1, true, nil
}
