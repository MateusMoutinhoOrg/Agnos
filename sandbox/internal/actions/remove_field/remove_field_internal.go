package remove_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveFieldInternal parses the target route's route.yaml, drops the named
// field from the origin --in names and writes the file back. It is the exact
// inverse of add-field, one origin selector and all.
func RemoveFieldInternal(deps *deps.Deps, io *smartio.SmartIO, route string, in string, name string) error {
	origin, err := utils.RouteFieldOrigin(deps, in)
	if err != nil {
		return err
	}

	conf, err := utils.LoadRouteConf(deps, io, route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(deps, name)
	if key == "" {
		return deps.Std.Errorf("remove-field needs the name of the field to drop")
	}

	deps.Std.Log("remove-field removing %s from %s \n", key, utils.RouteConfPath(deps, route))

	switch origin {
	case utils.RouteFieldInPath:
		return removeSegment(deps, io, route, conf, key)
	case utils.RouteFieldInHeader:
		return removeHeader(deps, io, route, conf, key)
	case utils.RouteFieldInQuery:
		return removeParam(deps, io, route, conf, key)
	}
	return removeBodyProperty(deps, io, route, conf, key)
}

// removeSegment drops one segment of `paths`: a captured one by its name, a
// trigger by the identifier it spells.
func removeSegment(deps *deps.Deps, io *smartio.SmartIO, route string, conf *routeconf.RouteConf, key string) error {
	index := utils.FindRouteSegment(deps, conf.Paths, key)
	if index < 0 {
		identifier, err := utils.RouteIdentifierSegment(deps, key)
		if err != nil {
			return deps.Std.Errorf("route %q declares no path segment named %q", route, key)
		}
		for i, segment := range conf.Paths {
			if segment.Field == nil && segment.Identifier == identifier {
				index = i
				break
			}
		}
	}
	if index < 0 {
		return deps.Std.Errorf("route %q declares no path segment named %q", route, key)
	}

	conf.Paths = utils.RemoveRouteSegment(conf.Paths, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}

func removeHeader(deps *deps.Deps, io *smartio.SmartIO, route string, conf *routeconf.RouteConf, key string) error {
	index := utils.FindRouteField(deps, conf.Headers, key)
	if index < 0 {
		return deps.Std.Errorf("route %q declares no header named %q", route, key)
	}
	conf.Headers = utils.RemoveRouteField(conf.Headers, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}

func removeParam(deps *deps.Deps, io *smartio.SmartIO, route string, conf *routeconf.RouteConf, key string) error {
	index := utils.FindRouteField(deps, conf.Params, key)
	if index < 0 {
		return deps.Std.Errorf("route %q declares no query parameter named %q", route, key)
	}
	conf.Params = utils.RemoveRouteField(conf.Params, index)
	return utils.SaveRouteConf(deps, io, route, conf)
}

// removeBodyProperty drops one property of the body's json-schema, named by
// the same dotted path add-field declared it with, and unlists it from its
// parent's `required`.
func removeBodyProperty(deps *deps.Deps, io *smartio.SmartIO, route string, conf *routeconf.RouteConf, key string) error {
	if conf.Body.Schema == nil {
		return deps.Std.Errorf("route %q declares no body json-schema", route)
	}

	parts := deps.Stringsdeps.Split(key, ".")
	parent := conf.Body.Schema

	for _, segment := range parts[:len(parts)-1] {
		child := findProperty(parent, segment)
		if child == nil {
			return deps.Std.Errorf("route %q declares no body property named %q", route, key)
		}
		parent = child
	}

	leaf := parts[len(parts)-1]
	index := -1
	for i, property := range parent.Properties {
		if property.Name == leaf {
			index = i
			break
		}
	}
	if index < 0 {
		return deps.Std.Errorf("route %q declares no body property named %q", route, key)
	}

	properties := make([]routeconf.SchemaProperty, 0, len(parent.Properties)-1)
	properties = append(properties, parent.Properties[:index]...)
	properties = append(properties, parent.Properties[index+1:]...)
	parent.Properties = properties

	required := make([]string, 0, len(parent.Required))
	for _, name := range parent.Required {
		if name != leaf {
			required = append(required, name)
		}
	}
	parent.Required = required

	return utils.SaveRouteConf(deps, io, route, conf)
}

// findProperty returns the named property of an object schema, or nil.
func findProperty(schema *routeconf.Schema, name string) *routeconf.Schema {
	for _, property := range schema.Properties {
		if property.Name == name {
			return property.Schema
		}
	}
	return nil
}
