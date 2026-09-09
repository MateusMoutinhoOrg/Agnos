package remove_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveParamInternal parses the target route's route.yaml, drops the named
// query parameter from `params` and writes the file back. It is the exact
// inverse of add-param.
func RemoveParamInternal(deps *deps.Deps, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(deps, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(deps, name)
	if key == "" {
		return deps.Std.Errorf("remove-param needs the name of the query parameter to drop")
	}

	index := utils.FindRouteField(deps, conf.Params, key)
	if index < 0 {
		return deps.Std.Errorf("route %q declares no query parameter named %q", route, key)
	}

	deps.Std.Log("remove-param removing %s from %s \n", key, utils.RouteConfPath(deps, route))

	conf.Params = utils.RemoveRouteField(conf.Params, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}
