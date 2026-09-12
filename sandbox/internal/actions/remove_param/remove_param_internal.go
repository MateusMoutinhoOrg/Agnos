package remove_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveParamInternal parses the target route's route.yaml, drops the named
// query parameter from `params` and writes the file back. It is the exact
// inverse of add-param.
func RemoveParamInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("remove-param needs the name of the query parameter to drop")
	}

	index := utils.FindRouteField(sandbox, conf.Params, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no query parameter named %q", route, key)
	}

	sandbox.Deps.Std.Log("remove-param removing %s from %s \n", key, utils.RouteConfPath(sandbox, route))

	conf.Params = utils.RemoveRouteField(conf.Params, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}
