package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetBodyInternal parses the target route's route.yaml, overwrites every body
// key the caller supplied (an empty string and a negative --max-bytes are
// "leave as is") and writes the file back. The schema is grown property by
// property by add-body-field; here it is only deleted, or carried from a json
// body to a form one and back.
func SetBodyInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetBodyProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}
	if props.Required && props.Optional {
		return sandbox.Deps.StdDeps.Errorf("--required and --optional are mutually exclusive")
	}

	changed := false

	if props.DropSchema {
		conf.Body.Schema, conf.Body.HasSchema, changed = nil, false, true
	}
	if raw := sandbox.Deps.StringsDeps.TrimSpace(props.Type); raw != "" {
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
			return sandbox.Deps.StdDeps.Errorf("--max-bytes 0 accepts no body at all: declare `--type none` instead")
		}
		conf.Body.MaxBytes, changed = props.MaxBytes, true
	}
	if content := sandbox.Deps.StringsDeps.TrimSpace(props.ContentType); content != "" {
		conf.Body.ContentType, changed = content, true
	}

	if !changed {
		return sandbox.Deps.StdDeps.Errorf("set-body: nothing to change (pass --type, --required, --optional, --max-bytes, --content-type or --drop-schema)")
	}
	if conf.Body.Type == routeconf.BodyNone && conf.Body.Required {
		return sandbox.Deps.StdDeps.Errorf("a `none` body cannot be required: it is never read")
	}

	sandbox.Deps.StdDeps.Logf("set-body updating %s \n", utils.RouteConfPath(sandbox, io, props.Route))

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// setType rewrites how the body is read, carrying the content-type along: a
// json or a form body demands its own, and a route that takes no body demands
// none, so the dispatch stops answering 415 for a request it no longer reads.
//
// The schema follows a json body turned form and back — it is the same subset
// under another key — as long as it is flat enough for a form to carry; any
// other type carries none, and --drop-schema has to say so.
func setType(sandbox *api.Sandbox, conf *routeconf.RouteConf, props api.SetBodyProps, raw string) error {
	kind, err := utils.RouteBodyType(sandbox, raw)
	if err != nil {
		return err
	}
	if conf.Body.HasSchema && routeconf.SchemaKeyOf(kind) == "" {
		return sandbox.Deps.StdDeps.Errorf("route %q declares a body schema, which only a json or a form body carries: pass --drop-schema to delete it", props.Route)
	}
	if conf.Body.HasSchema && kind == "form" {
		if err := utils.CheckFormSchema(sandbox, props.Route, conf.Body.Schema); err != nil {
			return sandbox.Deps.StdDeps.Errorf("%s (pass --drop-schema to start the form over)", err.Error())
		}
	}

	if sandbox.Deps.StringsDeps.TrimSpace(props.ContentType) == "" {
		switch {
		case kind == routeconf.BodyNone:
			conf.Body.ContentType = ""
		case kind == "json" && (conf.Body.ContentType == "" || conf.Body.ContentType == routeconf.DefaultFormContentType):
			conf.Body.ContentType = routeconf.DefaultJsonContentType
		case kind == "form" && (conf.Body.ContentType == "" || conf.Body.ContentType == routeconf.DefaultJsonContentType):
			conf.Body.ContentType = routeconf.DefaultFormContentType
		}
	}

	conf.Body.Type = kind
	return nil
}
