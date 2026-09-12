package add_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddHeaderInternal parses the target route's route.yaml, adds the new header
// to `headers` (refusing a duplicate) and writes the file back. The name is
// the external spelling, which the dispatch matches without regard to case.
func AddHeaderInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	field, err := utils.NewRouteField(sandbox, props, utils.RouteFieldInHeader)
	if err != nil {
		return err
	}
	if utils.FindRouteField(sandbox, conf.Headers, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("route %q already declares a header named %q", props.Route, field.Key)
	}

	position, err := utils.CheckRoutePosition(sandbox, "header", props.Position, len(conf.Headers))
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-header adding %s to %s \n", field.Key, utils.RouteConfPath(sandbox, props.Route))

	conf.Headers = utils.InsertRouteField(conf.Headers, field, position)
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
