package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RouteClearKeys is every key --clear may take off a route.
var RouteClearKeys = []string{"segments"}

// SetRouteInternal parses the target route's route.yaml, overwrites every
// route-level key the caller supplied (empty strings are "leave as is";
// --method replaces the whole list; --example appends; --before and --after
// land it one rung from another route) and writes the file back.
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
	cleared, err := utils.RouteClearSet(sandbox, props.Clear, RouteClearKeys)
	if err != nil {
		return err
	}
	if cleared["segments"] {
		conf.Segments, conf.HasSegments, changed = 0, false, true
	}

	relative, has_relative, err := utils.RouteRelativePriority(sandbox, io, props.Before, props.After)
	if err != nil {
		return err
	}
	if has_relative && props.HasPriority {
		return sandbox.Deps.Std.Errorf("--priority excludes --before and --after: name the rung, or the route it sits next to")
	}
	if has_relative {
		props.Priority, props.HasPriority = relative, true
	}
	if props.HasPriority {
		if props.Priority < 0 {
			return sandbox.Deps.Std.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
		}
		conf.Priority, changed = props.Priority, true
	}
	if props.HasSegments {
		if props.Segments < 1 {
			return sandbox.Deps.Std.Errorf("--segments %d is below 1: clear it with --clear segments for a route that takes any count", props.Segments)
		}
		conf.Segments, conf.HasSegments, changed = props.Segments, true, true
	}
	if sandbox.Deps.Stringsdeps.TrimSpace(props.Phase) != "" {
		phase, err := utils.RoutePhase(sandbox, props.Phase)
		if err != nil {
			return err
		}
		if phase == routeconf.PhaseAfter && conf.Body.Type != routeconf.BodyNone {
			return sandbox.Deps.Std.Errorf("route %q reads a body, which nothing reads in the after phase: drop it with set-body --type none first", props.Route)
		}
		conf.Phase, changed = phase, true
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
		return sandbox.Deps.Std.Errorf("set-route: nothing to change (pass --method, --response-type, --priority, --before, --after, --segments, --phase, --clear, --help, --category, --long-description, --hidden, --visible or --example)")
	}

	sandbox.Deps.Std.Log("set-route updating %s \n", utils.RouteConfPath(sandbox, props.Route))

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
