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

// AddRouteInternal writes the two hand-written files of a new route package.
// It refuses to overwrite an existing route (via io.WriteFile).
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
	utils.NoteNormalizedCommandName(sandbox, name)

	identifier := utils.RouteIdentifier(sandbox, name)
	pkg := utils.RoutePackage(sandbox, name)

	if pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be declared")
	}

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

	phase, err := utils.RoutePhase(sandbox, props.Phase)
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

	category := sandbox.Deps.Stringsdeps.TrimSpace(props.Category)
	if category == "" {
		category = defaultCategory
		if props.Middleware {
			category = defaultMiddlewareCategory
		}
	}

	sandbox.Deps.Std.Log("add-route creating %s \n", utils.RouteDir(sandbox, name))

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	conf.Methods = methods
	conf.Priority = priority
	conf.Phase = phase
	conf.ResponseType = response_type
	conf.Category = category
	conf.Help = sandbox.Deps.Stringsdeps.TrimSpace(props.Help)

	dir := utils.RouteDir(sandbox, name)
	if err := io.WriteFile(dir+"/route.yaml", []byte(conf.Render())); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier": identifier,
		"Package":    pkg,
		"Module":     module_conf.Module,
		"Methods":    sandbox.Deps.Stringsdeps.Join(methods, ", "),
		"Trigger":    conf.Pattern(),
		"After":      phase == routeconf.PhaseAfter,
	}

	template := RouteHandlerTemplate
	if props.Middleware || phase == routeconf.PhaseAfter {
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
