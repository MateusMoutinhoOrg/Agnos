package set_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetHeaderInternal parses the target route's route.yaml, rewrites the named
// request header in place and writes the file back. It is add-header applied
// to a declaration that already exists: the keys given are written over the
// ones there, --clear takes one off, and the result goes through the same
// constructor, so an edited header and a declared one are the same bytes.
func SetHeaderInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldEditProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, props.Name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("set-header needs the name of the request header to edit")
	}
	if utils.RouteFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-header: nothing to change (pass --rename, --type, --description, --default, --required, --min, --max, --example or --clear)")
	}

	index := utils.FindRouteField(sandbox, conf.Headers, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no request header named %q", props.Route, key)
	}

	field, err := utils.RouteFieldEdited(sandbox, conf.Headers[index], props, utils.RouteFieldInHeader)
	if err != nil {
		return err
	}
	if field.Key != key && utils.FindRouteField(sandbox, conf.Headers, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("route %q already declares a request header named %q", props.Route, field.Key)
	}

	sandbox.Deps.Std.Log("set-header updating %s in %s \n", key, utils.RouteConfPath(sandbox, props.Route))

	conf.Headers[index] = field
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
