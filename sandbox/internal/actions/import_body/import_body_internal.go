package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ImportBodyInternal parses the target route's route.yaml, reads one example
// payload and declares a body property for every key that example carries,
// then writes the file back. A route that declared no body becomes a json one
// here, exactly as one property declared by hand would make it.
//
// It is add-body-field run once per key, which is what the interview could not
// offer: a payload of ten keys is ten questionnaires by hand and one pasted
// document here. What it infers is a starting point and nothing more — a type
// per key, the objects and the lists around them, and, with --infer-format,
// the four formats a string may spell. Every bound after that is
// set-body-field's.
func ImportBodyInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteBodyImportProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	document, err := readExample(sandbox, props)
	if err != nil {
		return err
	}

	if conf.Body.Type == routeconf.BodyNone {
		conf.Body.Type = "json"
		conf.Body.ContentType = routeconf.DefaultJsonContentType
	}
	if conf.Body.Type != "json" {
		return sandbox.Deps.Std.Errorf("route %q declares a %q body, which carries no json-schema", props.Route, conf.Body.Type)
	}
	if conf.Body.Schema == nil || props.Replace {
		conf.Body.Schema = &routeconf.Schema{Type: "object"}
		conf.Body.HasSchema = true
	}

	inferred, err := inferSchema(sandbox, document, props)
	if err != nil {
		return err
	}
	if inferred.Type != "object" {
		return sandbox.Deps.Std.Errorf("an example payload is a json object: a body schema's root is the object its keys are declared in")
	}

	sandbox.Deps.Std.Log("import-body reading %s from the example \n", utils.RouteConfPath(sandbox, props.Route))

	added, skipped := mergeSchema(sandbox, conf.Body.Schema, inferred, "")
	for _, name := range added {
		sandbox.Deps.Std.Log("import-body adding %s \n", name)
	}
	for _, name := range skipped {
		sandbox.Deps.Std.Log("import-body leaving %s as it is: already declared \n", name)
	}
	if len(added) == 0 {
		return sandbox.Deps.Std.Errorf("the example declares no property this route does not already have (--replace starts the schema over)")
	}

	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// readExample is the payload itself: the document typed on the command line,
// or the file named instead of it. The file is read off the host rather than
// through the project's SmartIO — an example payload is something the person
// has lying about, not a file of the project being edited.
func readExample(sandbox *api.Sandbox, props api.RouteBodyImportProps) (string, error) {
	inline := sandbox.Deps.Stringsdeps.TrimSpace(props.Json)
	file := sandbox.Deps.Stringsdeps.TrimSpace(props.File)

	if inline != "" && file != "" {
		return "", sandbox.Deps.Std.Errorf("--json and --file are two ways of giving the same example: pass one of them")
	}
	if inline != "" {
		return inline, nil
	}
	if file == "" {
		return "", sandbox.Deps.Std.Errorf("import-body needs an example payload: --json '{...}' or --file payload.json")
	}

	content, err := sandbox.Deps.Iodeps.ReadFile(file)
	if err != nil {
		return "", sandbox.Deps.Std.Errorf("could not read the example payload at %q", file)
	}
	return string(content), nil
}
