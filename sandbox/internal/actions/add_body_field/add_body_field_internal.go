package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddBodyFieldInternal parses the target route's route.yaml, declares one
// property of the body's json-schema at the dotted path props.Name and writes
// the file back. A route that declared no body becomes a json one here:
// declaring a property is what says it takes a body at all.
func AddBodyFieldInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteBodyFieldProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	name := utils.RouteFieldName(sandbox, props.Name)
	if name == "" {
		return sandbox.Deps.Std.Errorf("a body property needs a name")
	}

	if conf.Body.Type == routeconf.BodyNone {
		conf.Body.Type = "json"
		conf.Body.ContentType = routeconf.DefaultJsonContentType
	}
	if conf.Body.Type != "json" {
		return sandbox.Deps.Std.Errorf("route %q declares a %q body, which carries no json-schema", props.Route, conf.Body.Type)
	}
	if conf.Body.Schema == nil {
		conf.Body.Schema = &routeconf.Schema{Type: "object"}
		conf.Body.HasSchema = true
	}

	parts := utils.SplitSchemaPath(sandbox, name)
	parent, err := utils.WalkSchemaTo(sandbox, conf.Body.Schema, parts[:len(parts)-1])
	if err != nil {
		return err
	}

	leaf := parts[len(parts)-1]
	if utils.SchemaPropertyOf(parent, leaf) != nil {
		return sandbox.Deps.Std.Errorf("route %q already declares a body property named %q", props.Route, name)
	}

	schema, err := utils.RouteBodyPropertySchema(sandbox, props)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-body-field adding %s to %s \n", name, utils.RouteConfPath(sandbox, props.Route))

	utils.InsertSchemaProperty(parent, leaf, schema)
	if props.Required {
		parent.Required = utils.AppendUnique(parent.Required, []string{leaf})
		sandbox.Deps.Sortdeps.Strings(parent.Required)
	}

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
