package remove_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveHeaderInternal parses the target route's route.yaml, drops the named
// header from `headers` and writes the file back. It is the exact inverse of
// add-header.
func RemoveHeaderInternal(deps *deps.Deps, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(deps, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(deps, name)
	if key == "" {
		return deps.Std.Errorf("remove-header needs the name of the header to drop")
	}

	index := utils.FindRouteField(deps, conf.Headers, key)
	if index < 0 {
		return deps.Std.Errorf("route %q declares no header named %q", route, key)
	}

	deps.Std.Log("remove-header removing %s from %s \n", key, utils.RouteConfPath(deps, route))

	conf.Headers = utils.RemoveRouteField(conf.Headers, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}
