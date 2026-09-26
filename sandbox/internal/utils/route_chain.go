package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RouteChainEntry is one declared route as the chain lays it down: its package
// name and its parsed declaration.
type RouteChainEntry struct {
	Name string
	Conf *routeconf.RouteConf
}

// LoadRouteChain reads every sandbox/internal/routeslist/<name>/route.yaml and
// returns them in the order the dispatch runs them: the `before` phase by
// priority, lowest first, then by name — the order the build collector lays
// Server.Routes down in — and the `after` phase the same way behind it.
func LoadRouteChain(sandbox *api.Sandbox, io *smartio.SmartIO) ([]RouteChainEntry, error) {
	chain := []RouteChainEntry{}

	for _, dir := range io.ListDirs(RoutesDir) {
		parts := sandbox.Deps.Stringsdeps.Split(dir, "/")
		name := parts[len(parts)-1]
		if name == "" {
			continue
		}
		content, err := io.ReadFile(RoutesDir + "/" + name + "/route.yaml")
		if err != nil {
			continue
		}
		conf, err := routeconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: %w", name, err)
		}
		chain = append(chain, RouteChainEntry{Name: name, Conf: conf})
	}

	SortRouteChain(sandbox, chain)
	return chain, nil
}

// SortRouteChain puts a chain in run order, in place.
func SortRouteChain(sandbox *api.Sandbox, chain []RouteChainEntry) {
	sandbox.Deps.Sortdeps.SliceStable(chain, func(i int, j int) bool {
		left, right := chain[i].Conf, chain[j].Conf
		left_after, right_after := left.Phase == routeconf.PhaseAfter, right.Phase == routeconf.PhaseAfter
		if left_after != right_after {
			return right_after
		}
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		return chain[i].Name < chain[j].Name
	})
}

// RouteRelativePriority is the rung one rung below (before) or above (after)
// the route named, for --before and --after. At most one of the two is given;
// it reports false when neither is.
func RouteRelativePriority(sandbox *api.Sandbox, io *smartio.SmartIO, before string, after string) (int, bool, error) {
	before = sandbox.Deps.Stringsdeps.TrimSpace(before)
	after = sandbox.Deps.Stringsdeps.TrimSpace(after)

	if before != "" && after != "" {
		return 0, false, sandbox.Deps.Std.Errorf("--before and --after exclude each other")
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
			return 0, false, sandbox.Deps.Std.Errorf(
				"--before %s: it runs on rung 0, and nothing runs below it — spread the chain with rebalance-routes first", RouteIdentifier(sandbox, before))
		}
		return other.Priority - 1, true, nil
	}

	other, err := LoadRouteConf(sandbox, io, after)
	if err != nil {
		return 0, false, err
	}
	return other.Priority + 1, true, nil
}

// RoutePhase normalizes a --phase: "" is before.
func RoutePhase(sandbox *api.Sandbox, raw string) (string, error) {
	phase := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if phase == "" {
		return routeconf.PhaseBefore, nil
	}
	if !contains(routeconf.Phases, phase) {
		return "", sandbox.Deps.Std.Errorf("unknown phase %q (use one of %s)", raw, sandbox.Deps.Stringsdeps.Join(routeconf.Phases, ", "))
	}
	return phase, nil
}
