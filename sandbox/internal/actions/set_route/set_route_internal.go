package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetRouteInternal parses the target route's route.yaml, overwrites every
// route-level key the caller supplied (empty strings are "leave as is";
// --method replaces the whole list; --example appends) and writes the file
// back.
func SetRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}
	if props.Hidden && props.Visible {
		return sandbox.Deps.Std.Errorf("--hidden and --visible are mutually exclusive")
	}

	changed := false
	if len(props.Methods) > 0 {
		methods, err := utils.RouteMethodList(sandbox, props.Methods)
		if err != nil {
			return err
		}
		conf.Methods, changed = methods, true
	}
	if response_type := sandbox.Deps.Stringsdeps.TrimSpace(props.ResponseType); response_type != "" {
		conf.ResponseType, changed = response_type, true
	}
	if help := sandbox.Deps.Stringsdeps.TrimSpace(props.Help); help != "" {
		conf.Help, changed = help, true
	}
	if category := sandbox.Deps.Stringsdeps.TrimSpace(props.Category); category != "" {
		conf.Category, changed = category, true
	}
	if long := sandbox.Deps.Stringsdeps.TrimSpace(props.LongDescription); long != "" {
		conf.LongDescription, changed = long, true
	}
	if props.HasPriority {
		if props.Priority < 0 {
			return sandbox.Deps.Std.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
		}
		conf.Priority, changed = props.Priority, true
	}
	if props.Hidden {
		conf.Hidden, changed = true, true
	}
	if props.Visible {
		conf.Hidden, changed = false, true
	}
	if len(props.Examples) > 0 {
		conf.Examples, changed = utils.AppendUnique(conf.Examples, props.Examples), true
	}
	if !changed {
		return sandbox.Deps.Std.Errorf("set-route: nothing to change (pass --method, --response-type, --priority, --help, --category, --long-description, --hidden, --visible or --example)")
	}

	sandbox.Deps.Std.Log("set-route updating %s \n", utils.RouteConfPath(sandbox, props.Route))

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
