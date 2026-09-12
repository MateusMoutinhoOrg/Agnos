package remove_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveSegmentInternal parses the target route's route.yaml, drops the named
// segment from `paths` and writes the file back. It is the exact inverse of
// add-segment: a capture answers to its name, a literal to the identifier it
// spells.
func RemoveSegmentInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("remove-segment needs the name of the segment to drop")
	}

	index := findSegment(sandbox, conf, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no path segment named %q", route, key)
	}

	sandbox.Deps.Std.Log("remove-segment removing %s from %s \n", key, utils.RouteConfPath(sandbox, route))

	conf.Paths = utils.RemoveRouteSegment(conf.Paths, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}

// findSegment locates the segment key names, looking it up as a capture first
// and as a literal identifier second, or returns -1.
func findSegment(sandbox *api.Sandbox, conf *routeconf.RouteConf, key string) int {
	if index := utils.FindRouteSegment(sandbox, conf.Paths, key); index >= 0 {
		return index
	}

	identifier, err := utils.RouteIdentifierSegment(sandbox, key)
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
