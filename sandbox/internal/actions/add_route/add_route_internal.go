package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddRouteInternal writes the two hand-written files of a new route package.
// It refuses to overwrite an existing route (via io.WriteFile). The trigger is
// normalized to start with "/", so "users" and "/users" produce the same
// route.yaml.
func AddRouteInternal(deps *deps.Deps, io *smartio.SmartIO, name string, method string, trigger string, help string, category string) error {
	if deps.Stringsdeps.TrimSpace(help) == "" {
		return deps.Std.Errorf("add-route requires --help")
	}
	if deps.Stringsdeps.TrimSpace(category) == "" {
		return deps.Std.Errorf("add-route requires --category")
	}

	if err := utils.ValidateRouteName(deps, name); err != nil {
		return err
	}
	utils.NoteNormalizedCommandName(deps, name)

	identifier := utils.RouteIdentifier(deps, name)
	pkg := utils.RoutePackage(deps, name)

	if pkg == "health" {
		return deps.Std.Errorf("the health route is generated and cannot be declared")
	}

	if deps.Stringsdeps.TrimSpace(trigger) == "" {
		trigger = "/" + identifier
	}
	segment, err := utils.RouteIdentifierSegment(deps, trigger)
	if err != nil {
		return err
	}

	verb, err := utils.RouteMethod(deps, method)
	if err != nil {
		return err
	}

	deps.Std.Log("add-route creating %s \n", utils.RouteDir(deps, name))

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier": identifier,
		"Package":    pkg,
		"Module":     module_conf.Module,
		"Method":     verb,
		"Trigger":    segment,
		"Help":       deps.Stringsdeps.TrimSpace(help),
		"Category":   deps.Stringsdeps.TrimSpace(category),
	}

	dir := utils.RouteDir(deps, name)

	route, err := deps.Embeddeps.RenderTemplate("templates/route_route.yaml", vars)
	if err != nil {
		return err
	}
	if err := io.WriteFile(dir+"/route.yaml", route); err != nil {
		return err
	}

	handler, err := deps.Embeddeps.RenderTemplate("templates/route_handler.go", vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/handler.go", handler)
}
