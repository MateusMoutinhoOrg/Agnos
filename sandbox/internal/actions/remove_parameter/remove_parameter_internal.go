package remove_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveParameterInternal parses the target route's route.yaml, drops the
// named entry of `parameters` and writes the file back. It is the exact
// inverse of add-parameter.
func RemoveParameterInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, route string, name string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	index := utils.FindRouteParameter(sandbox, conf.Parameters, name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no parameter named %q", route, name)
	}

	sandbox.Deps.StdDeps.Logf("remove-parameter removing %s from %s \n", conf.Parameters[index].Key, utils.RouteConfPath(sandbox, io, route))

	conf.Parameters = utils.RemoveAt(conf.Parameters, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}
