package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// InternalPureHandlerFile is the one hand-written file of a route package.
const InternalPureHandlerFile = "InternalPureHandler.go"

// AddRouteInternal writes the two hand-written files of a new route package.
// It refuses to overwrite an existing route (via io.WriteFile). The route
// starts with one path, `Route`, reading the whole request path and matching
// the trigger — "/" followed by the name unless --trigger says otherwise, an
// equal comparison unless --trigger-type does.
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

	path, err := utils.NewRoutePath(sandbox, api.RoutePathProps{
		Id:          "Route",
		TriggerType: props.TriggerType,
		Trigger:     trigger,
	})
	if err != nil {
		return err
	}

	methods, err := utils.RouteMethodList(sandbox, props.Methods)
	if err != nil {
		return err
	}

	response_type := sandbox.Deps.Stringsdeps.TrimSpace(props.ResponseType)
	if response_type == "" {
		response_type = routeconf.DefaultResponseType
	}

	sandbox.Deps.Std.Log("add-route creating %s \n", utils.RouteDir(sandbox, name))

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	conf := routeconf.NewEmpty(sandbox)
	conf.Methods = methods
	conf.Priority = props.Priority
	conf.ResponseType = response_type
	conf.Paths = []routeconf.Path{path}
	conf.Category = sandbox.Deps.Stringsdeps.TrimSpace(props.Category)
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
		"Trigger":    path.Trigger.Value,
	}

	handler, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/route_internal_pure_handler.go", vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/"+InternalPureHandlerFile, handler)
}
