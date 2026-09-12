package remove_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveHeaderInternal parses the target route's route.yaml, drops the named
// header from `headers` and writes the file back. It is the exact inverse of
// add-header.
func RemoveHeaderInternal(sandbox *api.Sandbox, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("remove-header needs the name of the header to drop")
	}

	index := utils.FindRouteField(sandbox, conf.Headers, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no header named %q", route, key)
	}

	sandbox.Deps.Std.Log("remove-header removing %s from %s \n", key, utils.RouteConfPath(sandbox, route))

	conf.Headers = utils.RemoveRouteField(conf.Headers, index)
	return utils.SaveRouteConf(sandbox, io, route, conf)
}
