package add_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddSegmentInternal parses the target route's route.yaml, inserts the new
// segment into `paths` and writes the file back. The URL is the concatenation
// of that sequence, so where the segment lands is what --position says.
func AddSegmentInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RouteFieldProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	position, err := utils.CheckRoutePosition(sandbox, "segment", props.Position, len(conf.Paths))
	if err != nil {
		return err
	}

	segment, err := newSegment(sandbox, conf, props)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-segment adding a segment to %s \n", utils.RouteConfPath(sandbox, props.Route))

	conf.Paths = utils.InsertRouteSegment(conf.Paths, segment, position)
	if err := checkRestLast(sandbox, conf, position); err != nil {
		return err
	}
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}

// checkRestLast refuses a `paths` the new segment left in a shape no matcher
// can be generated from: whatever takes the rest of the path only reads as a
// suffix while it is the last one. Two things do — a capture declared `array`
// and a `starts-with` trigger — and a route declaring both says twice what it
// binds once. Which of the two segments moved decides what the message asks
// for.
func checkRestLast(sandbox *api.Sandbox, conf *routeconf.RouteConf, position int) error {
	rest := utils.RouteRestIndex(conf.Paths)
	prefix := routePrefixIndex(conf)

	if rest >= 0 && prefix >= 0 {
		return sandbox.Deps.Std.Errorf(
			"this route already ends in %q, which takes every segment left in the path: a route has one such segment, not two",
			conf.Paths[prefix].Identifier)
	}

	at := rest
	label := "the captured segment %q"
	if at < 0 {
		at = prefix
		label = "the starts-with segment %q"
	}
	if at < 0 || at == len(conf.Paths)-1 {
		return nil
	}

	spelling := conf.Paths[at].Identifier
	if conf.Paths[at].Field != nil {
		spelling = conf.Paths[at].Field.Key
	}

	if at == position {
		return sandbox.Deps.Std.Errorf(
			label+" takes the rest of the path, so it is always the last one: leave --position out to append it",
			spelling)
	}
	return sandbox.Deps.Std.Errorf(
		label+" takes the rest of the path, so it is always the last one: put the new segment before it with --position %d",
		spelling, at)
}

// routePrefixIndex returns the index of the `starts-with` trigger, or -1. It is
// the trigger's half of RouteRestIndex: both name a segment that leaves
// everything after it unmatched.
func routePrefixIndex(conf *routeconf.RouteConf) int {
	for i, segment := range conf.Paths {
		if segment.Field == nil && segment.StartsWith {
			return i
		}
	}
	return -1
}

// newSegment builds the segment the flags describe, one of three shapes:
// --identifier spells a literal the URL has to match, --starts-with spells a
// prefix it only has to begin with — which leaves every segment after it
// unread, so the route answers a whole subtree — and a name captures whatever
// sits in that place. --array makes a capture take every segment left in the
// path, which only the last segment may do.
func newSegment(sandbox *api.Sandbox, conf *routeconf.RouteConf, props api.RouteFieldProps) (routeconf.Segment, error) {
	identifier := sandbox.Deps.Stringsdeps.TrimSpace(props.Identifier)
	starts_with := sandbox.Deps.Stringsdeps.TrimSpace(props.StartsWith)
	named := utils.RouteFieldName(sandbox, props.Name) != ""

	if identifier != "" && starts_with != "" {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf("a literal segment is spelled by --identifier or by --starts-with, never both")
	}
	if identifier == "" && starts_with == "" && !named {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf("a segment needs a name, which captures it, or --identifier / --starts-with, which spell it literally")
	}

	if identifier != "" || starts_with != "" {
		if named {
			return routeconf.Segment{}, sandbox.Deps.Std.Errorf("a segment is either a literal (--identifier / --starts-with) or a capture (a name), never both")
		}
		if props.Array {
			return routeconf.Segment{}, sandbox.Deps.Std.Errorf("--array belongs to a capture, which becomes a []T field: a literal segment spells the URL rather than binding it")
		}

		if starts_with != "" {
			spelling, err := utils.RoutePrefixIdentifier(sandbox, starts_with)
			if err != nil {
				return routeconf.Segment{}, err
			}
			return routeconf.Segment{Identifier: spelling, StartsWith: true}, nil
		}

		spelling, err := utils.RouteIdentifierSegment(sandbox, identifier)
		if err != nil {
			return routeconf.Segment{}, err
		}
		return routeconf.Segment{Identifier: spelling}, nil
	}

	field, err := utils.NewRouteField(sandbox, props, utils.RouteFieldInPath)
	if err != nil {
		return routeconf.Segment{}, err
	}
	if utils.FindRouteSegment(sandbox, conf.Paths, field.Key) >= 0 {
		return routeconf.Segment{}, sandbox.Deps.Std.Errorf("route %q already captures a segment named %q", props.Route, field.Key)
	}
	return routeconf.Segment{Field: &field}, nil
}
