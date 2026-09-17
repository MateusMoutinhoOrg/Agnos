package set_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetSegmentInternal parses the target route's route.yaml, rewrites the named
// segment of `paths` in place and writes the file back. A segment is named the
// way remove-segment names one: a capture by its name, a trigger by the
// identifier it spells.
//
// The two kinds are edited by the two halves of add-segment. --identifier
// makes the segment a trigger spelling that literal, whatever it was before;
// every other key edits a capture, and a trigger has none of them to edit.
func SetSegmentInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldEditProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	key := utils.RouteFieldName(sandbox, props.Name)
	if key == "" {
		return sandbox.Deps.Std.Errorf("set-segment needs the name of the segment to edit")
	}
	if utils.RouteFieldEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-segment: nothing to change (pass --identifier, --rename, --type, --description, --array, --min, --max, --example or --clear)")
	}

	index := utils.FindRoutePathSegment(sandbox, conf.Paths, key)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no path segment named %q", props.Route, key)
	}

	segment, err := editedSegment(sandbox, conf, index, props)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("set-segment updating %s in %s \n", key, utils.RouteConfPath(sandbox, props.Route))

	conf.Paths[index] = segment
	if err := checkRestLast(sandbox, conf); err != nil {
		return err
	}
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// editedSegment builds the segment the change leaves behind: a trigger when
// --identifier is given, and otherwise the capture that was there with the
// keys written over it.
func editedSegment(sandbox *api.Sandbox, conf *routeconf.RouteConf, index int, props api.RouteFieldEditProps) (routeconf.Segment, error) {
	current := conf.Paths[index]

	if identifier := sandbox.Deps.Stringsdeps.TrimSpace(props.Identifier); identifier != "" {
		spelling, err := utils.RouteIdentifierSegment(sandbox, identifier)
		if err != nil {
			return routeconf.Segment{}, err
		}
		return routeconf.Segment{Identifier: spelling}, nil
	}

	if current.Field == nil {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf(
			"the segment %q is a literal, which carries no type, no bounds and no description: --identifier is the whole of what it spells",
			current.Identifier)
	}

	field, err := utils.RouteFieldEdited(sandbox, *current.Field, props, utils.RouteFieldInPath)
	if err != nil {
		return routeconf.Segment{}, err
	}
	if field.Key != current.Field.Key && utils.FindRouteSegment(sandbox, conf.Paths, field.Key) >= 0 {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf("route %q already captures a segment named %q", props.Route, field.Key)
	}

	return routeconf.Segment{Field: &field}, nil
}

// checkRestLast refuses a `paths` the edit left in a shape no matcher can be
// generated from: the segment taking the rest of the path only reads as a
// suffix while it is the last one, so --array on any other segment is refused
// where add-segment would have refused it too.
func checkRestLast(sandbox *api.Sandbox, conf *routeconf.RouteConf) error {
	at := utils.RouteRestIndex(conf.Paths)
	if at < 0 || at == len(conf.Paths)-1 {
		return nil
	}

	return sandbox.Deps.Std.Errorf(
		"the captured segment %q takes the rest of the path, so it is always the last one: only the last segment may be an array",
		conf.Paths[at].Field.Key)
}
