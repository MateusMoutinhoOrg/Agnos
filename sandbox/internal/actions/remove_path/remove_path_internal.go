package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemovePathInternal parses the target route's route.yaml, drops the named
// entry of `paths` and writes the file back. It is the exact inverse of
// add-path, and it refuses to leave a route with no path at all.
func RemovePathInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string, id string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	index := utils.FindRoutePath(sandbox, conf.Paths, id)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no path with id %q", route, utils.RouteEntryId(sandbox, id))
	}
	if len(conf.Paths) == 1 {
		return sandbox.Deps.Std.Errorf("%q is the last path of route %q: a route declares one at least", conf.Paths[index].Id, route)
	}

	sandbox.Deps.Std.Log("remove-path removing %s from %s \n", conf.Paths[index].Id, utils.RouteConfPath(sandbox, route))

	conf.Paths = utils.RemoveRoutePath(conf.Paths, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}
