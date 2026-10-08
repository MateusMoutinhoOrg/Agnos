package set_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetParameterInternal parses the target route's route.yaml, rebuilds the
// named entry of `parameters` with the changes applied and writes the file
// back.
func SetParameterInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetParameterProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	index := utils.FindRouteParameter(sandbox, conf.Parameters, props.Name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no parameter named %q", props.Route, props.Name)
	}
	if utils.IsRouteParameterEditEmpty(sandbox, props) {
		return sandbox.Deps.StdDeps.Errorf("set-parameter was given nothing to change")
	}

	current := conf.Parameters[index]
	parameter, err := utils.EditRouteParameter(sandbox, current, props)
	if err != nil {
		return err
	}
	if utils.RouteIdTaken(conf, parameter.Id, current.Id) {
		return sandbox.Deps.StdDeps.Errorf("route %q already has an Input field named %q", props.Route, parameter.Id)
	}

	sandbox.Deps.StdDeps.Logf("set-parameter rewriting %s of %s \n", current.Key, utils.RouteConfPath(sandbox, io, props.Route))

	conf.Parameters[index] = parameter
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
