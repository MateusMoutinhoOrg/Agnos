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
// The three kinds are edited by the three halves of add-segment. --identifier
// makes the segment a trigger spelling that literal and --starts-with one
// spelling that prefix, whatever either was before; every other key edits a
// capture, and a trigger has none of them to edit.
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
		return sandbox.Deps.Std.Errorf("set-segment: nothing to change (pass --identifier, --starts-with, --rename, --type, --description, --array, --min, --max, --example or --clear)")
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
// --identifier or --starts-with is given, and otherwise the capture that was
// there with the keys written over it.
func editedSegment(sandbox *api.Sandbox, conf *routeconf.RouteConf, index int, props api.RouteFieldEditProps) (routeconf.Segment, error) {
	current := conf.Paths[index]

	identifier := sandbox.Deps.Stringsdeps.TrimSpace(props.Identifier)
	starts_with := sandbox.Deps.Stringsdeps.TrimSpace(props.StartsWith)

	if identifier != "" && starts_with != "" {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf("a literal segment is spelled by --identifier or by --starts-with, never both")
	}

	if starts_with != "" {
		spelling, err := utils.RoutePrefixIdentifier(sandbox, starts_with)
		if err != nil {
			return routeconf.Segment{}, err
		}
		return routeconf.Segment{Identifier: spelling, StartsWith: true}, nil
	}

	if identifier != "" {
		spelling, err := utils.RouteIdentifierSegment(sandbox, identifier)
		if err != nil {
			return routeconf.Segment{}, err
		}
		return routeconf.Segment{Identifier: spelling}, nil
	}

	if current.Field == nil {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf(
			"the segment %q is a literal, which carries no type, no bounds and no description: --identifier and --starts-with are the whole of what it spells",
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
// generated from: whatever takes the rest of the path only reads as a suffix
// while it is the last one, so --array and --starts-with on any other segment
// are refused where add-segment would have refused them too — and a route
// never carries both, since each one already takes everything left.
func checkRestLast(sandbox *api.Sandbox, conf *routeconf.RouteConf) error {
	rest := utils.RouteRestIndex(conf.Paths)
	prefix := -1
	for i, segment := range conf.Paths {
		if segment.Field == nil && segment.StartsWith {
			prefix = i
		}
	}

	if rest >= 0 && prefix >= 0 {
		return sandbox.Deps.Std.Errorf(
			"this route already ends in %q, which takes every segment left in the path: a route has one such segment, not two",
			conf.Paths[prefix].Identifier)
	}

	if rest >= 0 && rest != len(conf.Paths)-1 {
		return sandbox.Deps.Std.Errorf(
			"the captured segment %q takes the rest of the path, so it is always the last one: only the last segment may be an array",
			conf.Paths[rest].Field.Key)
	}

	if prefix >= 0 && prefix != len(conf.Paths)-1 {
		return sandbox.Deps.Std.Errorf(
			"the starts-with segment %q leaves every segment after it unmatched, so it is always the last one",
			conf.Paths[prefix].Identifier)
	}

	return nil
}
