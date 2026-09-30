package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// InternalPureHandlerFile is the one hand-written file of a route package.
const InternalPureHandlerFile = "InternalPureHandler.go"

// MiddlewareHandlerTemplate is the stub `add-route --middleware` writes: a
// handler that declines, so the chain goes on.
const MiddlewareHandlerTemplate = "templates/route_middleware_handler.go"

// RouteHandlerTemplate is the stub every other route starts with: a handler
// that answers.
const RouteHandlerTemplate = "templates/route_internal_pure_handler.go"

// defaultCategory is the heading a route declared without --category is
// listed under in docs/Routes, and defaultMiddlewareCategory the one of a
// middleware.
const defaultCategory = "Routes"
const defaultMiddlewareCategory = "Middleware"

// middlewareResponseType is the response-type a middleware declares when it
// names none: it answers nothing on the happy path, so the type is the one a
// refusal written by hand would carry.
const middlewareResponseType = "text/plain"

// AddRouteInternal writes the two hand-written files of a new route package,
// in the folder props.Dir names under sandbox/internal/routeslist. It refuses
// a name another route already carries, in whatever folder.
//
// The paths come from one of two places. --pattern compiles a url shape into
// one path per literal run and per capture, fixing the segment count unless
// it ends on {*rest}. Otherwise the route starts with one path, `Route`,
// reading the whole request path and matching the trigger — "/" followed by
// the name unless --trigger says otherwise ("/" for a middleware), an equal
// comparison unless --trigger-type does (prefix for a middleware).
//
// The rung is --priority when it is given, one below --before or one above
// --after the route it names, and DefaultRoutePriority — or
// DefaultMiddlewarePriority — otherwise.
func AddRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddRouteProps) error {
	name := props.Name
	if err := utils.ValidateRouteName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedName(sandbox, "route", name)

	identifier := utils.RouteIdentifier(sandbox, name)
	pkg := utils.RoutePackage(sandbox, name)

	if pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be declared")
	}

	group, err := utils.UnitGroup(sandbox, props.Dir)
	if err != nil {
		return err
	}
	if existing, found := utils.FindUnitDir(sandbox, io, utils.RoutesDir, utils.RouteConfFile, pkg); found {
		return sandbox.Deps.Std.Errorf("route %q already exists in %s: a route name is unique across every folder", identifier, existing)
	}
	dir := utils.UnitDirIn(utils.RoutesDir, group, pkg)

	conf := routeconf.NewEmpty(sandbox)

	pattern := sandbox.Deps.Stringsdeps.TrimSpace(props.Pattern)
	if pattern != "" {
		if sandbox.Deps.Stringsdeps.TrimSpace(props.Trigger+props.TriggerType) != "" || props.TriggerNegate || props.TriggerIgnoreCase {
			return sandbox.Deps.Std.Errorf("--pattern declares the paths itself: it excludes --trigger, --trigger-type, --trigger-negate and --trigger-ignore-case")
		}
		compiled, err := utils.CompileRoutePattern(sandbox, pattern)
		if err != nil {
			return err
		}
		conf.Paths = compiled.Paths
		conf.Segments, conf.HasSegments = compiled.Segments, compiled.HasSegments
	} else {
		trigger := props.Trigger
		trigger_type := props.TriggerType
		if sandbox.Deps.Stringsdeps.TrimSpace(trigger) == "" {
			trigger = "/" + identifier
			if props.Middleware {
				trigger = "/"
			}
		}
		if sandbox.Deps.Stringsdeps.TrimSpace(trigger_type) == "" && props.Middleware {
			trigger_type = "prefix"
		}

		path, err := utils.NewRoutePath(sandbox, api.RoutePathProps{
			Id:                "Route",
			TriggerType:       trigger_type,
			Trigger:           trigger,
			TriggerNegate:     props.TriggerNegate,
			TriggerIgnoreCase: props.TriggerIgnoreCase,
		})
		if err != nil {
			return err
		}
		conf.Paths = []routeconf.Path{path}
	}

	raw_methods := props.Methods
	if len(raw_methods) == 0 && props.Middleware {
		raw_methods = []string{routeconf.AnyMethod}
	}
	methods, err := utils.RouteMethodList(sandbox, raw_methods)
	if err != nil {
		return err
	}

	priority, err := addRoutePriority(sandbox, io, props)
	if err != nil {
		return err
	}

	response_type := sandbox.Deps.Stringsdeps.TrimSpace(props.ResponseType)
	if response_type == "" {
		response_type = routeconf.DefaultResponseType
		if props.Middleware {
			response_type = middlewareResponseType
		}
	}
	if err := utils.ValidateMediaType(sandbox, response_type); err != nil {
		return err
	}

	category := sandbox.Deps.Stringsdeps.TrimSpace(props.Category)
	if category == "" {
		category = defaultCategory
		if props.Middleware {
			category = defaultMiddlewareCategory
		}
	}

	sandbox.Deps.Std.Log("add-route creating %s \n", dir)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	conf.Methods = methods
	conf.Priority = priority
	conf.ResponseType = response_type
	conf.Category = category
	conf.Help = sandbox.Deps.Stringsdeps.TrimSpace(props.Help)

	if !props.Middleware {
		warnShadowedRoute(sandbox, io, identifier, conf)
	}

	if err := io.WriteFile(dir+"/"+utils.RouteConfFile, []byte(conf.Render())); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier": identifier,
		"Package":    pkg,
		"Module":     module_conf.Module,
		"Methods":    sandbox.Deps.Stringsdeps.Join(methods, ", "),
		"Trigger":    conf.Pattern(),
	}

	template := RouteHandlerTemplate
	if props.Middleware {
		template = MiddlewareHandlerTemplate
	}
	handler, err := sandbox.Deps.Embeddeps.RenderTemplate(template, vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/"+InternalPureHandlerFile, handler)
}

// addRoutePriority is the rung a new route lands on: --priority when it is
// given, one rung from --before or --after, the default of its kind
// otherwise. It is never negative.
func addRoutePriority(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddRouteProps) (int, error) {
	relative, has_relative, err := utils.RouteRelativePriority(sandbox, io, props.Before, props.After)
	if err != nil {
		return 0, err
	}
	if has_relative && props.HasPriority {
		return 0, sandbox.Deps.Std.Errorf("--priority excludes --before and --after: name the rung, or the route it sits next to")
	}

	priority := api.DefaultRoutePriority
	if props.Middleware {
		priority = api.DefaultMiddlewarePriority
	}
	if props.HasPriority {
		priority = props.Priority
	}
	if has_relative {
		priority = relative
	}

	if priority < 0 {
		return 0, sandbox.Deps.Std.Errorf("--priority %d is negative: the chain runs from zero upwards", priority)
	}
	return priority, nil
}

// warnShadowedRoute says so when another route already matches exactly the
// requests conf does — the same paths, segment count and an overlapping
// method. Both run on every such request, lowest priority first, so whichever
// answers first shadows the other: legitimate for a middleware, a mistake
// almost always for two routes, and never something to find out by accident.
func warnShadowedRoute(sandbox *api.Sandbox, io *smartio.SmartIO, identifier string, conf *routeconf.RouteConf) {
	signature := routeMatchSignature(sandbox, conf)
	for _, unit := range utils.RouteDirs(sandbox, io) {
		other_name := unit.Name
		other, err := utils.LoadRouteConf(sandbox, io, other_name)
		if err != nil || !methodsOverlap(other.Methods, conf.Methods) {
			continue
		}
		if routeMatchSignature(sandbox, other) == signature {
			sandbox.Deps.Std.Error("warning: route %s matches exactly the requests %s does (%s): the one on the lower priority answers first and shadows the other\n",
				identifier, utils.RouteIdentifier(sandbox, other_name), conf.Pattern())
		}
	}
}

// routeMatchSignature is what a route matches on, spelled as one string.
func routeMatchSignature(sandbox *api.Sandbox, conf *routeconf.RouteConf) string {
	signature := sandbox.Deps.Std.Sprintf("segments=%v:%d", conf.HasSegments, conf.Segments)
	for _, path := range conf.Paths {
		signature += sandbox.Deps.Std.Sprintf("|%d:%d:%s:%+v", path.Start, path.End, path.Type, path.Trigger)
	}
	return signature
}

// methodsOverlap reports whether one method answers both lists.
func methodsOverlap(first []string, second []string) bool {
	for _, one := range first {
		for _, two := range second {
			if one == two || one == routeconf.AnyMethod || two == routeconf.AnyMethod {
				return true
			}
		}
	}
	return false
}
