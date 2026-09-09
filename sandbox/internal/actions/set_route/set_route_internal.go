package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetRouteInternal parses the target route's route.yaml, overwrites every
// route-level key the caller supplied (empty strings are "leave as is";
// --example appends) and writes the file back.
func SetRouteInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RouteProps) error {
	conf, err := utils.LoadRouteConf(deps, io, props.Route)
	if err != nil {
		return err
	}
	if props.Hidden && props.Visible {
		return deps.Std.Errorf("--hidden and --visible are mutually exclusive")
	}

	changed := false
	if raw := deps.Stringsdeps.TrimSpace(props.Method); raw != "" {
		method, err := utils.RouteMethod(deps, raw)
		if err != nil {
			return err
		}
		conf.Method, changed = method, true
	}
	if help := deps.Stringsdeps.TrimSpace(props.Help); help != "" {
		conf.Help, changed = help, true
	}
	if category := deps.Stringsdeps.TrimSpace(props.Category); category != "" {
		conf.Category, changed = category, true
	}
	if long := deps.Stringsdeps.TrimSpace(props.LongDescription); long != "" {
		conf.LongDescription, changed = long, true
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
		return deps.Std.Errorf("set-route: nothing to change (pass --method, --help, --category, --long-description, --hidden, --visible or --example)")
	}

	deps.Std.Log("set-route updating %s \n", utils.RouteConfPath(deps, props.Route))

	return utils.SaveRouteConf(deps, io, props.Route, conf)
}
