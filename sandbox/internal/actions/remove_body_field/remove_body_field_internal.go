package remove_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveBodyFieldInternal parses the target route's route.yaml, drops the
// property named by the same dotted path add-body-field declared it with, and
// writes the file back. It is the exact inverse of add-body-field, its
// parent's `required` entry included.
func RemoveBodyFieldInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, route string, name string) error {
	conf, err := utils.LoadRouteConf(sandbox, io, route)
	if err != nil {
		return err
	}
	if conf.Body.Schema == nil {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no body schema", route)
	}

	key := utils.RouteFieldName(sandbox, name)
	if key == "" {
		return sandbox.Deps.StdDeps.Errorf("remove-body-field needs the dotted path of the property to drop")
	}

	parts := utils.SplitSchemaPath(sandbox, key)
	parent, err := walkTo(sandbox, conf.Body.Schema, parts[:len(parts)-1], route, key)
	if err != nil {
		return err
	}

	if !utils.DropSchemaProperty(parent, parts[len(parts)-1]) {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no body property named %q", route, key)
	}

	sandbox.Deps.StdDeps.Logf("remove-body-field removing %s from %s \n", key, utils.RouteConfPath(sandbox, io, route))

	return utils.SaveRouteConf(sandbox, io, route, conf)
}

// walkTo descends the dotted path to the object the leaf is declared in,
// entering an array of objects through its `items` and reporting the first
// segment that was never declared.
func walkTo(sandbox *api.Sandbox, root *routeconf.Schema, segments []string, route string, key string) (*routeconf.Schema, error) {
	parent := root
	for _, segment := range segments {
		child := utils.SchemaObjectOf(utils.SchemaPropertyOf(parent, segment))
		if child == nil {
			return nil, sandbox.Deps.StdDeps.Errorf("route %q declares no body property named %q", route, key)
		}
		parent = child
	}
	return parent, nil
}
