package set_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetBodyFieldInternal parses the target route's route.yaml, rewrites the
// property at the dotted path props.Name and writes the file back. It is
// add-body-field applied to a property that already exists: the keywords
// already declared are read back, the ones given are written over them, and
// the whole is built again by the same constructor — so the alternative it
// replaces, remove-body-field followed by add-body-field, is never the way to
// add a bound that was forgotten.
//
// The objects the path passes through are not created here: a property that is
// not declared yet is add-body-field's, and creating one silently would make a
// typo in the path a new property instead of an error.
func SetBodyFieldInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteBodyFieldEditProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	name := utils.RouteFieldName(sandbox, props.Name)
	if name == "" {
		return sandbox.Deps.Std.Errorf("set-body-field needs the dotted path of the property to edit")
	}
	if utils.RouteBodyFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-body-field: nothing to change (pass --rename, --type, --required, --array, one of the schema keywords, or --clear)")
	}
	if conf.Body.Schema == nil {
		return sandbox.Deps.Std.Errorf("route %q declares no body json-schema: add-body-field declares the first property", props.Route)
	}

	parts := utils.SplitSchemaPath(sandbox, name)
	parent, err := parentOf(sandbox, conf.Body.Schema, parts)
	if err != nil {
		return err
	}

	leaf := parts[len(parts)-1]
	current := utils.SchemaPropertyOf(parent, leaf)
	if current == nil {
		return sandbox.Deps.Std.Errorf("route %q declares no body property named %q", props.Route, name)
	}

	edit, err := utils.RouteBodyFieldEdited(sandbox, current, utils.SchemaDemands(parent, leaf), props)
	if err != nil {
		return err
	}

	renamed := leaf
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename); value != "" {
		renamed = value
		if renamed != leaf && utils.SchemaPropertyOf(parent, renamed) != nil {
			return sandbox.Deps.Std.Errorf("route %q already declares a body property named %q", props.Route, renamed)
		}
	}

	sandbox.Deps.Std.Log("set-body-field updating %s in %s \n", name, utils.RouteConfPath(sandbox, props.Route))
	for _, keyword := range edit.Dropped {
		sandbox.Deps.Std.Log("set-body-field dropping %s: the new type carries none \n", keyword)
	}

	utils.DropSchemaProperty(parent, leaf)
	utils.InsertSchemaProperty(parent, renamed, edit.Schema)
	if edit.Required {
		parent.Required = utils.AppendUnique(parent.Required, []string{renamed})
		sandbox.Deps.Sortdeps.Strings(parent.Required)
	}

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// parentOf walks the dotted path to the object the property is declared in,
// refusing a path that names an object nothing declared — the one thing
// add-body-field does and this does not.
func parentOf(sandbox *api.Sandbox, root *routeconf.Schema, parts []string) (*routeconf.Schema, error) {
	parent := root
	for _, segment := range parts[:len(parts)-1] {
		property := utils.SchemaPropertyOf(parent, segment)
		if property == nil {
			return nil, sandbox.Deps.Std.Errorf("no body property named %q is declared, so nothing under it can be edited", segment)
		}

		child := utils.SchemaObjectOf(property)
		if child == nil || child.Type != "object" {
			return nil, sandbox.Deps.Std.Errorf("body property %q is not an object, so nothing under it can be edited", segment)
		}
		parent = child
	}
	return parent, nil
}
