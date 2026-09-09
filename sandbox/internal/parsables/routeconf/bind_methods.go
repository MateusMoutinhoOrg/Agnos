package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func BindMethods(deps *deps.Deps, conf *RouteConf) {
	conf.Render = func() string {
		return Render(deps, conf)
	}
	conf.Pattern = func() string {
		return Pattern(conf)
	}
	conf.IdentifierCount = func() int {
		return IdentifierCount(conf)
	}
	conf.IdentifierLen = func() int {
		return IdentifierLen(conf)
	}
	conf.SchemaJson = func() string {
		return SchemaJson(deps, conf)
	}
}

// Pattern is the route's path as docs and messages spell it: every segment in
// order, a trigger bringing its own leading slash and a capture entering as
// "/{name}". A route with no segment at all reads as "/".
func Pattern(conf *RouteConf) string {
	pattern := ""
	for _, segment := range conf.Paths {
		if segment.Field != nil {
			pattern += "/{" + segment.Field.Key + "}"
			continue
		}
		pattern += segment.Identifier
	}
	if pattern == "" {
		return "/"
	}
	return pattern
}

// IdentifierCount is how many trigger segments the route fixes.
func IdentifierCount(conf *RouteConf) int {
	count := 0
	for _, segment := range conf.Paths {
		if segment.Field == nil {
			count++
		}
	}
	return count
}

// IdentifierLen is the total number of characters the route's triggers spell.
func IdentifierLen(conf *RouteConf) int {
	length := 0
	for _, segment := range conf.Paths {
		if segment.Field == nil {
			length += len(segment.Identifier)
		}
	}
	return length
}
