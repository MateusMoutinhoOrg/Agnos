package add_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddParameterInternal parses the target route's route.yaml, inserts the new
// entry into `parameters` and writes the file back.
func AddParameterInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteParameterProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	position, err := utils.CheckRoutePosition(sandbox, "parameter", props.Position, len(conf.Parameters))
	if err != nil {
		return err
	}

	parameter, err := utils.NewRouteParameter(sandbox, props)
	if err != nil {
		return err
	}
	if utils.FindRouteParameter(sandbox, conf.Parameters, parameter.Key) >= 0 || utils.RouteIdTaken(conf, parameter.Id, "") {
		return sandbox.Deps.Std.Errorf("route %q already has an Entries field named %q", props.Route, parameter.Id)
	}

	sandbox.Deps.Std.Log("add-parameter adding %s to %s \n", parameter.Key, utils.RouteConfPath(sandbox, props.Route))

	conf.Parameters = utils.InsertRouteParameter(conf.Parameters, parameter, position)
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
