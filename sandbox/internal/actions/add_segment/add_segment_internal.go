package add_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddSegmentInternal parses the target route's route.yaml, inserts the new
// segment into `paths` and writes the file back. The URL is the concatenation
// of that sequence, so where the segment lands is what --position says.
func AddSegmentInternal(deps *deps.Deps, io *smartio.SmartIO, props api.RouteFieldProps) error {
	conf, err := utils.LoadRouteConf(deps, io, props.Route)
	if err != nil {
		return err
	}

	position, err := utils.CheckRoutePosition(deps, "segment", props.Position, len(conf.Paths))
	if err != nil {
		return err
	}

	segment, err := newSegment(deps, conf, props)
	if err != nil {
		return err
	}

	deps.Std.Log("add-segment adding a segment to %s \n", utils.RouteConfPath(deps, props.Route))

	conf.Paths = utils.InsertRouteSegment(conf.Paths, segment, position)
	if err := checkRestLast(deps, conf, position); err != nil {
		return err
	}
	return utils.SaveRouteConf(deps, io, props.Route, conf)
}

// checkRestLast refuses a `paths` the new segment left in a shape no matcher
// can be generated from: the segment taking the rest of the path only reads as
// a suffix while it is the last one. Which of the two segments moved decides
// what the message asks for.
func checkRestLast(deps *deps.Deps, conf *routeconf.RouteConf, position int) error {
	at := utils.RouteRestIndex(conf.Paths)
	if at < 0 || at == len(conf.Paths)-1 {
		return nil
	}

	if at == position {
		return deps.Std.Errorf(
			"the captured segment %q takes the rest of the path, so it is always the last one: leave --position out to append it",
			conf.Paths[at].Field.Key)
	}
	return deps.Std.Errorf(
		"the captured segment %q takes the rest of the path, so it is always the last one: put the new segment before it with --position %d",
		conf.Paths[at].Field.Key, at)
}

// newSegment builds the segment the flags describe: a trigger, normalized to
// start with "/", when --identifier is given, and a captured field named on
// the command line otherwise. --array makes that capture take every segment
// left in the path, which only the last segment may do.
func newSegment(deps *deps.Deps, conf *routeconf.RouteConf, props api.RouteFieldProps) (routeconf.Segment, error) {
	identifier := deps.Stringsdeps.TrimSpace(props.Identifier)
	named := utils.RouteFieldName(deps, props.Name) != ""
	if identifier == "" && !named {
		return routeconf.Segment{}, deps.Std.Errorf("a segment needs a name, which captures it, or --identifier, which spells it literally")
	}
	if identifier != "" {
		if named {
			return routeconf.Segment{}, deps.Std.Errorf("a segment is either a literal (--identifier) or a capture (a name), never both")
		}
		if props.Array {
			return routeconf.Segment{}, deps.Std.Errorf("--array belongs to a capture, which becomes a []T field: a literal segment spells one segment of the URL")
		}
		spelling, err := utils.RouteIdentifierSegment(deps, identifier)
		if err != nil {
			return routeconf.Segment{}, err
		}
		return routeconf.Segment{Identifier: spelling}, nil
	}

	field, err := utils.NewRouteField(deps, props, utils.RouteFieldInPath)
	if err != nil {
		return routeconf.Segment{}, err
	}
	if utils.FindRouteSegment(deps, conf.Paths, field.Key) >= 0 {
		return routeconf.Segment{}, deps.Std.Errorf("route %q already captures a segment named %q", props.Route, field.Key)
	}
	return routeconf.Segment{Field: &field}, nil
}
