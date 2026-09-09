package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// The route helpers below mirror the command ones of command_conf.go, one for
// one: a route is declared, edited and looked up exactly the way a command is,
// only against route.yaml instead of entries.yaml.

// RouteIdentifier normalizes a user-typed route name into its canonical
// spelling: lowercased, spaces and underscores turned into dashes
// ("Create User" -> "create-user").
func RouteIdentifier(deps *deps.Deps, name string) string {
	return CommandIdentifier(deps, name)
}

// ValidateRouteName reports whether a user-typed route name normalizes to a
// usable identifier. It becomes a directory name and a Go package clause, so
// the same alphabet a command name is held to applies.
func ValidateRouteName(deps *deps.Deps, name string) error {
	identifier := RouteIdentifier(deps, name)
	if identifier == "" {
		return deps.Std.Errorf("a route needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return deps.Std.Errorf("invalid route name %q: a route name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return deps.Std.Errorf(
				"invalid route name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory sandbox/internal/routes/%s and a Go package name)",
				name, RoutePackage(deps, name))
		}
	}
	return nil
}

// RoutePackage is the Go package / directory name for a route: the identifier
// with dashes turned into underscores ("create-user" -> "create_user").
func RoutePackage(deps *deps.Deps, name string) string {
	return deps.Stringsdeps.ReplaceAll(RouteIdentifier(deps, name), "-", "_")
}

// RouteDir is the project-relative directory holding a route package.
func RouteDir(deps *deps.Deps, name string) string {
	return "sandbox/internal/routes/" + RoutePackage(deps, name)
}

// RouteConfPath is the project-relative path of a route's route.yaml.
func RouteConfPath(deps *deps.Deps, name string) string {
	return RouteDir(deps, name) + "/route.yaml"
}

// LoadRouteConf reads and parses sandbox/internal/routes/<name>/route.yaml.
func LoadRouteConf(deps *deps.Deps, io *smartio.SmartIO, name string) (*routeconf.RouteConf, error) {
	if err := ValidateRouteName(deps, name); err != nil {
		return nil, err
	}
	content, err := io.ReadFile(RouteConfPath(deps, name))
	if err != nil {
		return nil, deps.Std.Errorf("route %q not found in %s", RouteIdentifier(deps, name), RouteDir(deps, name))
	}
	conf, err := routeconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("routes/%s/route.yaml: %w", RoutePackage(deps, name), err)
	}
	return conf, nil
}

// SaveRouteConf renders conf back over sandbox/internal/routes/<name>/route.yaml.
func SaveRouteConf(deps *deps.Deps, io *smartio.SmartIO, name string, conf *routeconf.RouteConf) error {
	return io.WriteFileOverwrite(RouteConfPath(deps, name), []byte(conf.Render()))
}

// RouteIdentifierSegment normalizes a trigger segment to the one spelling a
// declaration carries: always starting with "/", never holding another one.
// "users" and "/users" are the same request; "/a/b" and "/users/" are refused,
// because a trigger is one segment and the leading slash is what makes it read
// as a path everywhere it is printed.
func RouteIdentifierSegment(deps *deps.Deps, raw string) (string, error) {
	trimmed := deps.Stringsdeps.TrimSpace(raw)
	if trimmed == "" {
		return "", deps.Std.Errorf("a path identifier cannot be empty")
	}

	inner := deps.Stringsdeps.Trim(trimmed, "/")
	if inner == "" {
		return "/", nil
	}
	if deps.Stringsdeps.Contains(inner, "/") {
		return "", deps.Std.Errorf("invalid path identifier %q: an identifier is one segment, so it holds no inner or trailing slash", raw)
	}

	return "/" + inner, nil
}

// RouteFieldIn names the three origins that read a value off the request line.
// They share one Field shape and one set of editors, and differ only in the
// rules NewRouteField holds each of them to.
const (
	// RouteFieldInPath is a captured segment of the route's `paths`.
	RouteFieldInPath = "path"
	// RouteFieldInHeader is a declared request header.
	RouteFieldInHeader = "header"
	// RouteFieldInQuery is a declared query-string parameter.
	RouteFieldInQuery = "query"
)

// NewRouteField builds a routeconf.Field from the raw values typed on the
// command line, holding to the rules of the origin it is declared in: a
// captured segment is always required and never defaults, and a header never
// repeats. An array is a query parameter collecting every occurrence of its
// key, or the last path segment taking every segment left in the path.
func NewRouteField(deps *deps.Deps, props api.RouteFieldProps, in string) (routeconf.Field, error) {
	field := routeconf.Field{
		Key:         RouteFieldName(deps, props.Name),
		Description: deps.Stringsdeps.TrimSpace(props.Description),
		Examples:    props.Examples,
		Required:    props.Required,
		Array:       props.Array,
	}
	if field.Key == "" {
		return field, deps.Std.Errorf("a field needs a name")
	}

	kind, ok := FieldType(deps, props.Type)
	if !ok {
		return field, deps.Std.Errorf("unknown type %q (use string, boolean, int or float)", props.Type)
	}
	field.Type = kind

	if field.Array && in == RouteFieldInHeader {
		return field, deps.Std.Errorf("only a query parameter or the last path segment may be an array")
	}
	if field.Required && kind == "boolean" {
		return field, deps.Std.Errorf("a boolean field cannot be required (its absence already means false)")
	}

	if in == RouteFieldInPath {
		if props.Default != "" {
			return field, deps.Std.Errorf("a captured path segment cannot carry a default: it is always present when the route matches")
		}
		field.Required = true
	}

	if props.Default != "" {
		if field.Required {
			return field, deps.Std.Errorf("a field cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := RouteCheckLiteral(deps, kind, "default", props.Default); err != nil {
			return field, err
		}
		field.HasDefault = true
		field.Default = props.Default
	}

	if props.Min != "" {
		value, err := RouteParseBound(deps, kind, "min", props.Min)
		if err != nil {
			return field, err
		}
		field.Min, field.HasMin = value, true
	}
	if props.Max != "" {
		value, err := RouteParseBound(deps, kind, "max", props.Max)
		if err != nil {
			return field, err
		}
		field.Max, field.HasMax = value, true
	}
	if field.HasMin && field.HasMax && field.Min > field.Max {
		return field, deps.Std.Errorf("min (%s) is greater than max (%s)", props.Min, props.Max)
	}

	return field, nil
}

// RouteFieldName normalizes a field name. A header name and a query key are
// external spellings, so only the surrounding space is trimmed: their own
// punctuation is theirs to keep.
func RouteFieldName(deps *deps.Deps, name string) string {
	return deps.Stringsdeps.TrimSpace(name)
}

// FindRouteField returns the index of the field named name in fields, or -1.
func FindRouteField(deps *deps.Deps, fields []routeconf.Field, name string) int {
	key := RouteFieldName(deps, name)
	for i, field := range fields {
		if field.Key == key {
			return i
		}
	}
	return -1
}

// FindRouteSegment returns the index of the captured segment named name in
// segments, or -1. A trigger segment carries no name and never matches.
func FindRouteSegment(deps *deps.Deps, segments []routeconf.Segment, name string) int {
	key := RouteFieldName(deps, name)
	for i, segment := range segments {
		if segment.Field != nil && segment.Field.Key == key {
			return i
		}
	}
	return -1
}

// RouteRestIndex returns the index of the segment that takes the rest of the
// path — the one capture declared `array: true` — or -1 when the route fixes
// its length. It is the one reading of that rule, shared by the editor, the
// collector and verify.
func RouteRestIndex(segments []routeconf.Segment) int {
	for i, segment := range segments {
		if segment.Field != nil && segment.Field.Array {
			return i
		}
	}
	return -1
}

// CheckRoutePosition validates a --position against the list it will be
// inserted into, the same way CheckPosition does for a command's fields.
func CheckRoutePosition(deps *deps.Deps, kind string, position int, size int) (int, error) {
	if position == AppendPosition {
		return size, nil
	}
	if position < 0 {
		return 0, deps.Std.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, size)
	}
	if position > size {
		return 0, deps.Std.Errorf("--position %d is out of range: this route has %d %s(s), so the accepted range is 0 to %d", position, size, kind, size)
	}
	return position, nil
}

// InsertRouteField places field at position inside fields (appending when
// position is negative or past the end).
func InsertRouteField(fields []routeconf.Field, field routeconf.Field, position int) []routeconf.Field {
	if position < 0 || position >= len(fields) {
		return append(fields, field)
	}
	out := make([]routeconf.Field, 0, len(fields)+1)
	out = append(out, fields[:position]...)
	out = append(out, field)
	out = append(out, fields[position:]...)
	return out
}

// RemoveRouteField drops the field at index from fields.
func RemoveRouteField(fields []routeconf.Field, index int) []routeconf.Field {
	out := make([]routeconf.Field, 0, len(fields)-1)
	out = append(out, fields[:index]...)
	out = append(out, fields[index+1:]...)
	return out
}

// InsertRouteSegment places segment at position inside segments (appending
// when position is negative or past the end).
func InsertRouteSegment(segments []routeconf.Segment, segment routeconf.Segment, position int) []routeconf.Segment {
	if position < 0 || position >= len(segments) {
		return append(segments, segment)
	}
	out := make([]routeconf.Segment, 0, len(segments)+1)
	out = append(out, segments[:position]...)
	out = append(out, segment)
	out = append(out, segments[position:]...)
	return out
}

// RemoveRouteSegment drops the segment at index from segments.
func RemoveRouteSegment(segments []routeconf.Segment, index int) []routeconf.Segment {
	out := make([]routeconf.Segment, 0, len(segments)-1)
	out = append(out, segments[:index]...)
	out = append(out, segments[index+1:]...)
	return out
}

// RouteMethod normalizes an http method, refusing one no route may answer.
func RouteMethod(deps *deps.Deps, raw string) (string, error) {
	method := deps.Stringsdeps.ToUpper(deps.Stringsdeps.TrimSpace(raw))
	if method == "" {
		return routeconf.DefaultMethod, nil
	}
	for _, known := range RouteMethods {
		if known == method {
			return method, nil
		}
	}
	return "", deps.Std.Errorf("unknown method %q (use one of %s)", raw, deps.Stringsdeps.Join(RouteMethods, ", "))
}

// RouteMethods is every http method a route may declare.
var RouteMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// RouteCheckLiteral reports whether a raw command-line literal parses as the
// field's type.
func RouteCheckLiteral(deps *deps.Deps, kind string, label string, raw string) error {
	return checkLiteral(deps, kind, label, raw)
}

// RouteParseBound reads a min/max literal for a numeric field.
func RouteParseBound(deps *deps.Deps, kind string, label string, raw string) (float64, error) {
	return parseBound(deps, kind, label, raw)
}
