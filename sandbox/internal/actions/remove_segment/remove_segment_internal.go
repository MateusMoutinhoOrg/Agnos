package remove_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveSegmentInternal parses the target route's route.yaml, drops the named
// segment from `paths` and writes the file back. It is the exact inverse of
// add-segment: a capture answers to its name, a literal to the identifier it
// spells.
func RemoveSegmentInternal(deps *deps.Deps, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(deps, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(deps, name)
	if key == "" {
		return deps.Std.Errorf("remove-segment needs the name of the segment to drop")
	}

	index := findSegment(deps, conf, key)
	if index < 0 {
		return deps.Std.Errorf("route %q declares no path segment named %q", route, key)
	}

	deps.Std.Log("remove-segment removing %s from %s \n", key, utils.RouteConfPath(deps, route))

	conf.Paths = utils.RemoveRouteSegment(conf.Paths, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}

// findSegment locates the segment key names, looking it up as a capture first
// and as a literal identifier second, or returns -1.
func findSegment(deps *deps.Deps, conf *routeconf.RouteConf, key string) int {
	if index := utils.FindRouteSegment(deps, conf.Paths, key); index >= 0 {
		return index
	}

	identifier, err := utils.RouteIdentifierSegment(deps, key)
	if err != nil {
		return -1
	}
	for i, segment := range conf.Paths {
		if segment.Field == nil && segment.Identifier == identifier {
			return i
		}
	}
	return -1
}
