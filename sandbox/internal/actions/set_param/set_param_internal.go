package set_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetParamInternal parses the target route's route.yaml, rewrites the named
// query parameter in place and writes the file back. It is add-param applied
// to a declaration that already exists: the keys given are written over the
// ones there, --clear takes one off, and the result goes through the same
// constructor, so an edited parameter and a declared one are the same bytes.
func SetParamInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldEditProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, props.Name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("set-param needs the name of the query parameter to edit")
	}
	if utils.RouteFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-param: nothing to change (pass --rename, --type, --description, --default, --required, --array, --min, --max, --example or --clear)")
	}

	index := utils.FindRouteField(sandbox, conf.Params, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no query parameter named %q", props.Route, key)
	}

	field, err := utils.RouteFieldEdited(sandbox, conf.Params[index], props, utils.RouteFieldInQuery)
	if err != nil {
		return err
	}
	if field.Key != key && utils.FindRouteField(sandbox, conf.Params, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("route %q already declares a query parameter named %q", props.Route, field.Key)
	}

	sandbox.Deps.Std.Log("set-param updating %s in %s \n", key, utils.RouteConfPath(sandbox, props.Route))

	conf.Params[index] = field
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
