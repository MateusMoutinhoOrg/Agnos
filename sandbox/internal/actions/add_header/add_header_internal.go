package add_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddHeaderInternal parses the target route's route.yaml, adds the new header
// to `headers` (refusing a duplicate) and writes the file back. The name is
// the external spelling, which the dispatch matches without regard to case.
func AddHeaderInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RouteFieldProps) error {
	conf, err := utils.LoadRouteConf(deps, io, props.Route)
	if err != nil {
		return err
	}

	field, err := utils.NewRouteField(deps, props, utils.RouteFieldInHeader)
	if err != nil {
		return err
	}
	if utils.FindRouteField(deps, conf.Headers, field.Key) >= 0 {
		return deps.Std.Errorf("route %q already declares a header named %q", props.Route, field.Key)
	}

	position, err := utils.CheckRoutePosition(deps, "header", props.Position, len(conf.Headers))
	if err != nil {
		return err
	}

	deps.Std.Log("add-header adding %s to %s \n", field.Key, utils.RouteConfPath(deps, props.Route))

	conf.Headers = utils.InsertRouteField(conf.Headers, field, position)
	return utils.SaveRouteConf(deps, io, props.Route, conf)
}
