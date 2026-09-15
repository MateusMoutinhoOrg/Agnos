package add_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddParamInternal parses the target route's route.yaml, adds the new query
// parameter to `params` (refusing a duplicate) and writes the file back. The
// query string is the one origin an array field may be read from.
func AddParamInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	field, err := utils.NewRouteField(sandbox, props, utils.RouteFieldInQuery)
	if err != nil {
		return err
	}
	if utils.FindRouteField(sandbox, conf.Params, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("route %q already declares a query parameter named %q", props.Route, field.Key)
	}

	position, err := utils.CheckRoutePosition(sandbox, "query parameter", props.Position, len(conf.Params))
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-param adding %s to %s \n", field.Key, utils.RouteConfPath(sandbox, props.Route))

	conf.Params = utils.InsertRouteField(conf.Params, field, position)
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
