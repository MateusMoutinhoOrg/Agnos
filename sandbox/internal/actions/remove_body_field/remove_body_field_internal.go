package remove_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveBodyFieldInternal parses the target route's route.yaml, drops the
// property named by the same dotted path add-body-field declared it with, and
// writes the file back. It is the exact inverse of add-body-field, its
// parent's `required` entry included.
func RemoveBodyFieldInternal(deps *deps.Deps, io *smartio.SmartIO, route string, name string) error {
	conf, err := utils.LoadRouteConf(deps, io, route)
	if err != nil {
		return err
	}
	if conf.Body.Schema == nil {
		return deps.Std.Errorf("route %q declares no body json-schema", route)
	}

	key := utils.RouteFieldName(deps, name)
	if key == "" {
		return deps.Std.Errorf("remove-body-field needs the dotted path of the property to drop")
	}

	parts := utils.SplitSchemaPath(deps, key)
	parent, err := walkTo(deps, conf.Body.Schema, parts[:len(parts)-1], route, key)
	if err != nil {
		return err
	}

	if !utils.DropSchemaProperty(parent, parts[len(parts)-1]) {
		return deps.Std.Errorf("route %q declares no body property named %q", route, key)
	}

	deps.Std.Log("remove-body-field removing %s from %s \n", key, utils.RouteConfPath(deps, route))

	return utils.SaveRouteConf(deps, io, route, conf)
}

// walkTo descends the dotted path to the object the leaf is declared in,
// entering an array of objects through its `items` and reporting the first
// segment that was never declared.
func walkTo(deps *deps.Deps, root *routeconf.Schema, segments []string, route string, key string) (*routeconf.Schema, error) {
	parent := root
	for _, segment := range segments {
		child := utils.SchemaObjectOf(utils.SchemaPropertyOf(parent, segment))
		if child == nil {
			return nil, deps.Std.Errorf("route %q declares no body property named %q", route, key)
		}
		parent = child
	}
	return parent, nil
}
