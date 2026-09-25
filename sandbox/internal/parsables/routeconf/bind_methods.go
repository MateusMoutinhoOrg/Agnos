package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *RouteConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
	conf.Pattern = func() string {
		return Pattern(conf)
	}
	conf.SchemaJson = func() string {
		return SchemaJson(sandbox, conf)
	}
}

// Pattern is the route's path as docs and messages spell it: one piece per
// entry of `paths`, in order. An `equal` trigger is its value, a `prefix` its
// value followed by "*", a `suffix` "*" followed by its value, a `regex` its
// value between "~(" and ")", and a plain capture "{Id}". A route with no
// path at all reads as "/".
func Pattern(conf *RouteConf) string {
	pattern := ""
	for _, path := range conf.Paths {
		if !path.Trigger.Exists {
			pattern += "{" + path.Id + "}"
			continue
		}
		switch path.Trigger.Type {
		case "prefix":
			pattern += path.Trigger.Value + "*"
		case "suffix":
			pattern += "*" + path.Trigger.Value
		case "regex":
			pattern += "~(" + path.Trigger.Value + ")"
		default:
			pattern += path.Trigger.Value
		}
	}
	if pattern == "" {
		return "/"
	}
	return pattern
}
