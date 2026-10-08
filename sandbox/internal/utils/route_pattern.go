package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
)

// RoutePattern is what one --pattern compiles to: the entries of `paths` it
// stands for, and the segment count a request has to have — HasSegments false
// when the pattern ends on a {*rest}, which takes any count from there on.
type RoutePattern struct {
	Paths       []routeconf.Path
	Segments    int
	HasSegments bool
}

// CompileRoutePattern turns one url shape typed on the command line into the
// paths route.yaml declares. The yaml never holds the pattern itself: it holds
// what the pattern means, and conf.Pattern() draws it back.
//
//	/get-article                 {GetArticle 0..0 equal /get-article}
//	/users/{user}                + {User 1..1}
//	/users/{id:integer}          + {Id 1..1 type integer}
//	/static/{*file}              + {File 1..-1}
//
// Literal segments in a row are one path, named after them. Without a {*…}
// the pattern fixes the segment count, so /get-article/{article} does not run
// for /get-article/42/extra. "/" alone is the root.
func CompileRoutePattern(sandbox *api.Sandbox, raw string) (RoutePattern, error) {
	text := sandbox.Deps.StringsDeps.TrimSpace(raw)
	if !sandbox.Deps.StringsDeps.HasPrefix(text, "/") {
		return RoutePattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q has to start with /", raw)
	}

	segments := []string{}
	for _, segment := range sandbox.Deps.StringsDeps.Split(text, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	if len(segments) == 0 {
		return RoutePattern{Paths: []routeconf.Path{{
			Id:      RoutePathIdOf(sandbox, "/"),
			Start:   0,
			End:     routeconf.LastSegment,
			Type:    routeconf.DefaultPathType,
			Trigger: triggerconf.Trigger{Set: true, Type: "equal", Value: "/"},
		}}}, nil
	}

	compiled := RoutePattern{Paths: []routeconf.Path{}, Segments: len(segments), HasSegments: true}
	taken := map[string]bool{}
	for _, reserved := range RouteReservedIds {
		taken[reserved] = true
	}
	claim := func(id string, index int) (string, error) {
		if id == "" {
			return "", sandbox.Deps.StdDeps.Errorf("--pattern %q names an empty capture at segment %d", raw, index)
		}
		if taken[id] {
			return "", sandbox.Deps.StdDeps.Errorf("--pattern %q binds %s twice, or one Input already carries", raw, id)
		}
		taken[id] = true
		return id, nil
	}

	literal_start := -1
	literal := []string{}
	flush := func(end int) error {
		if len(literal) == 0 {
			return nil
		}
		value := "/" + sandbox.Deps.StringsDeps.Join(literal, "/")
		id := RoutePathIdOf(sandbox, value)
		if id == "" || id[0] < 'A' || id[0] > 'Z' || taken[id] {
			id = "Seg" + sandbox.Deps.StringsDeps.FormatInt(int64(literal_start), 10)
		}
		id, err := claim(id, literal_start)
		if err != nil {
			return err
		}
		compiled.Paths = append(compiled.Paths, routeconf.Path{
			Id:      id,
			Start:   literal_start,
			End:     end,
			Type:    routeconf.DefaultPathType,
			Trigger: triggerconf.Trigger{Set: true, Type: "equal", Value: value},
		})
		literal, literal_start = []string{}, -1
		return nil
	}

	for index, segment := range segments {
		if !sandbox.Deps.StringsDeps.HasPrefix(segment, "{") {
			if sandbox.Deps.StringsDeps.Contains(segment, "{") || sandbox.Deps.StringsDeps.Contains(segment, "}") {
				return RoutePattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q mixes text and a capture in the segment %q: a capture is a whole segment", raw, segment)
			}
			if literal_start < 0 {
				literal_start = index
			}
			literal = append(literal, segment)
			continue
		}

		if err := flush(index - 1); err != nil {
			return RoutePattern{}, err
		}
		if !sandbox.Deps.StringsDeps.HasSuffix(segment, "}") {
			return RoutePattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q opens a capture it does not close: %q", raw, segment)
		}
		inner := segment[1 : len(segment)-1]

		if sandbox.Deps.StringsDeps.HasPrefix(inner, "*") {
			if index != len(segments)-1 {
				return RoutePattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q puts %s before the end: a {*…} capture takes the rest of the path", raw, segment)
			}
			id, err := claim(GoIdentifier(sandbox, inner[1:]), index)
			if err != nil {
				return RoutePattern{}, err
			}
			compiled.Paths = append(compiled.Paths, routeconf.Path{
				Id: id, Start: index, End: routeconf.LastSegment, Type: routeconf.DefaultPathType,
			})
			compiled.Segments, compiled.HasSegments = 0, false
			continue
		}

		name, kind := inner, routeconf.DefaultPathType
		if parts := sandbox.Deps.StringsDeps.Split(inner, ":"); len(parts) == 2 {
			name, kind = parts[0], sandbox.Deps.StringsDeps.ToLower(parts[1])
		}
		if !contains(routeconf.PathTypes, kind) {
			return RoutePattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q declares the unknown type %q (use one of %s)",
				raw, kind, sandbox.Deps.StringsDeps.Join(routeconf.PathTypes, ", "))
		}
		id, err := claim(GoIdentifier(sandbox, name), index)
		if err != nil {
			return RoutePattern{}, err
		}
		compiled.Paths = append(compiled.Paths, routeconf.Path{Id: id, Start: index, End: index, Type: kind})
	}

	if err := flush(len(segments) - 1); err != nil {
		return RoutePattern{}, err
	}
	return compiled, nil
}

// RootPathId is the id of a path that reads the whole request path and names
// no word of it: a trigger of "/" alone.
const RootPathId = "Root"

// RoutePathIdOf is the Input field a path whose trigger is text binds to: the
// words of the trigger joined into one exported Go name — /api/products binds
// ApiProducts, ^/(api/)?admin(/|$) binds ApiAdmin — and Root when the trigger
// names no word. It is the one rule add-route --trigger and a literal run of
// --pattern both spell a path's id with, so the id says what the path reads.
func RoutePathIdOf(sandbox *api.Sandbox, text string) string {
	words := sandbox.Deps.StringsDeps.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
	})
	id := GoIdentifier(sandbox, sandbox.Deps.StringsDeps.Join(words, "-"))
	if id == "" {
		return RootPathId
	}
	if id[0] < 'A' || id[0] > 'Z' {
		return "Path" + id
	}
	for _, reserved := range RouteReservedIds {
		if id == reserved {
			return id + "Path"
		}
	}
	return id
}
