package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetBodyInternal parses the target route's route.yaml, overwrites every body
// key the caller supplied (an empty string and a negative --max-bytes are
// "leave as is") and writes the file back. The json-schema is grown property
// by property by add-body-field; here it is only deleted.
func SetBodyInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteBodyProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}
	if props.Required && props.Optional {
		return sandbox.Deps.Std.Errorf("--required and --optional are mutually exclusive")
	}

	changed := false

	if props.DropSchema {
		conf.Body.Schema, conf.Body.HasSchema, changed = nil, false, true
	}
	if raw := sandbox.Deps.Stringsdeps.TrimSpace(props.Type); raw != "" {
		if err := setType(sandbox, conf, props, raw); err != nil {
			return err
		}
		changed = true
	}
	if props.Required {
		conf.Body.Required, changed = true, true
	}
	if props.Optional {
		conf.Body.Required, changed = false, true
	}
	if props.MaxBytes >= 0 {
		if props.MaxBytes == 0 {
			return sandbox.Deps.Std.Errorf("--max-bytes 0 accepts no body at all: declare `--type none` instead")
		}
		conf.Body.MaxBytes, changed = props.MaxBytes, true
	}
	if content := sandbox.Deps.Stringsdeps.TrimSpace(props.ContentType); content != "" {
		conf.Body.ContentType, changed = content, true
	}

	if !changed {
		return sandbox.Deps.Std.Errorf("set-body: nothing to change (pass --type, --required, --optional, --max-bytes, --content-type or --drop-schema)")
	}
	if conf.Body.Type == routeconf.BodyNone && conf.Body.Required {
		return sandbox.Deps.Std.Errorf("a `none` body cannot be required: it is never read")
	}

	sandbox.Deps.Std.Log("set-body updating %s \n", utils.RouteConfPath(sandbox, props.Route))

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// setType rewrites how the body is read, carrying the content-type along: a
// json body demands one, and a route that takes no body demands none, so the
// dispatch stops answering 415 for a request it no longer reads.
func setType(sandbox *api.Sandbox, conf *routeconf.RouteConf, props api.RouteBodyProps, raw string) error {
	kind, err := utils.RouteBodyType(sandbox, raw)
	if err != nil {
		return err
	}
	if kind != "json" && conf.Body.HasSchema {
		return sandbox.Deps.Std.Errorf("route %q declares a json-schema, which only a json body carries: pass --drop-schema to delete it", props.Route)
	}

	if sandbox.Deps.Stringsdeps.TrimSpace(props.ContentType) == "" {
		switch {
		case kind == routeconf.BodyNone:
			conf.Body.ContentType = ""
		case kind == "json" && conf.Body.ContentType == "":
			conf.Body.ContentType = routeconf.DefaultJsonContentType
		}
	}

	conf.Body.Type = kind
	return nil
}
