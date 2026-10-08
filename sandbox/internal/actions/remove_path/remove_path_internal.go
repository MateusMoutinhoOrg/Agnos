package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemovePathInternal parses the target route's route.yaml, drops the named
// entry of `paths` and writes the file back. It is the exact inverse of
// add-path, and it refuses to leave a route with no path at all.
func RemovePathInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, route string, id string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	index := utils.FindRoutePath(sandbox, conf.Paths, id)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no path with id %q", route, utils.GoIdentifier(sandbox, id))
	}
	if len(conf.Paths) == 1 {
		return sandbox.Deps.StdDeps.Errorf("%q is the last path of route %q: a route declares one at least", conf.Paths[index].Id, route)
	}

	sandbox.Deps.StdDeps.Logf("remove-path removing %s from %s \n", conf.Paths[index].Id, utils.RouteConfPath(sandbox, io, route))

	conf.Paths = utils.RemoveAt(conf.Paths, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}
