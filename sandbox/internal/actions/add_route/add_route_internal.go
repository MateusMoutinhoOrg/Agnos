package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddRouteInternal writes the two hand-written files of a new route package.
// It refuses to overwrite an existing route (via io.WriteFile). The trigger is
// normalized to start with "/", so "users" and "/users" produce the same
// route.yaml — and with --starts-with it becomes a prefix instead, which the
// route answers every path beginning with.
func AddRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddRouteProps) error {
	if sandbox.Deps.Stringsdeps.TrimSpace(props.Help) == "" {
		return sandbox.Deps.Std.Errorf("add-route requires --help")
	}
	if sandbox.Deps.Stringsdeps.TrimSpace(props.Category) == "" {
		return sandbox.Deps.Std.Errorf("add-route requires --category")
	}
	if props.Priority < 0 {
		return sandbox.Deps.Std.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
	}

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

	trigger := props.Trigger
	if sandbox.Deps.Stringsdeps.TrimSpace(trigger) == "" {
		trigger = "/" + identifier
	}

	segment, err := utils.RouteIdentifierSegment(sandbox, trigger)
	if props.StartsWith {
		segment, err = utils.RoutePrefixIdentifier(sandbox, trigger)
	}
	if err != nil {
		return err
	}

	verb, err := utils.RouteMethod(sandbox, props.Method)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-route creating %s \n", utils.RouteDir(sandbox, name))

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier": identifier,
		"Package":    pkg,
		"Module":     module_conf.Module,
		"Method":     verb,
		"Trigger":    segment,
		"StartsWith": props.StartsWith,
		"Priority":   props.Priority,
		"Help":       sandbox.Deps.Stringsdeps.TrimSpace(props.Help),
		"Category":   sandbox.Deps.Stringsdeps.TrimSpace(props.Category),
	}

	dir := utils.RouteDir(sandbox, name)

	route, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/route_route.yaml", vars)
	if err != nil {
		return err
	}
	if err := io.WriteFile(dir+"/route.yaml", route); err != nil {
		return err
	}

	handler, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/route_handler.go", vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/handler.go", handler)
}
