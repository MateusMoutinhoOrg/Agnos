package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// The route helpers below mirror the command ones of command_conf.go, one for
// one: a route is declared, edited and looked up exactly the way a command is,
// only against route.yaml instead of entries.yaml.

// RouteIdentifier normalizes a user-typed route name into its canonical
// spelling: lowercased, spaces and underscores turned into dashes
// ("Create User" -> "create-user").
func RouteIdentifier(sandbox *api.Sandbox, name string) string {
	return CommandIdentifier(sandbox, name)
}

// ValidateRouteName reports whether a user-typed route name normalizes to a
// usable identifier. It becomes a directory name and a Go package clause, so
// the same alphabet a command name is held to applies.
func ValidateRouteName(sandbox *api.Sandbox, name string) error {
	identifier := RouteIdentifier(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.Std.Errorf("a route needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid route name %q: a route name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.Std.Errorf(
				"invalid route name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory sandbox/internal/routeslist/%s and a Go package name)",
				name, RoutePackage(sandbox, name))
		}
	}
	return nil
}

// RoutePackage is the Go package / directory name for a route: the identifier
// with dashes turned into underscores ("create-user" -> "create_user").
func RoutePackage(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.Stringsdeps.ReplaceAll(RouteIdentifier(sandbox, name), "-", "_")
}

// RoutesDir holds one declared route per sub-directory, the server layer's
// mirror of sandbox/internal/commands.
const RoutesDir = "sandbox/internal/routeslist"

// RouteDir is the project-relative directory holding a route package.
func RouteDir(sandbox *api.Sandbox, name string) string {
	return RoutesDir + "/" + RoutePackage(sandbox, name)
}

// RouteConfPath is the project-relative path of a route's route.yaml.
func RouteConfPath(sandbox *api.Sandbox, name string) string {
	return RouteDir(sandbox, name) + "/route.yaml"
}

// LoadRouteConf reads and parses sandbox/internal/routeslist/<name>/route.yaml.
func LoadRouteConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string) (*routeconf.RouteConf, error) {
	if err := ValidateRouteName(sandbox, name); err != nil {
		return nil, err
	}
	content, err := io.ReadFile(RouteConfPath(sandbox, name))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("route %q not found in %s", RouteIdentifier(sandbox, name), RouteDir(sandbox, name))
	}
	conf, err := routeconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("routeslist/%s/route.yaml: %w", RoutePackage(sandbox, name), err)
	}
	return conf, nil
}

// SaveRouteConf renders conf back over sandbox/internal/routeslist/<name>/route.yaml.
func SaveRouteConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *routeconf.RouteConf) error {
	return io.WriteFileOverwrite(RouteConfPath(sandbox, name), []byte(conf.Render()))
}

// RouteTriggerAliases maps the spellings a trigger type may be typed in onto
// the one route.yaml carries; the canonical names map onto themselves.
var RouteTriggerAliases = map[string]string{
	"starts-with": "prefix",
	"ends-with":   "suffix",
	"exact":       "equal",
	"equals":      "equal",
	"matches":     "regex",
}

// RouteTriggerType normalizes a trigger type typed on the command line: an
// alias becomes the type it stands for, and anything that is neither is
// refused with the list of both.
func RouteTriggerType(sandbox *api.Sandbox, raw string) (string, error) {
	kind := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if canonical, is := RouteTriggerAliases[kind]; is {
		kind = canonical
	}
	if !contains(routeconf.TriggerTypes, kind) {
		return "", sandbox.Deps.Std.Errorf("unknown trigger type %q (use one of %s, or starts-with, ends-with, exact, matches)",
			raw, sandbox.Deps.Stringsdeps.Join(routeconf.TriggerTypes, ", "))
	}
	return kind, nil
}

// RouteTriggerProps is one trigger as it is typed on the command line: the
// value, how it is compared ("" is equal), and the two switches on it. OnPath
// reports a trigger compared against a path slice, which reads with the
// leading slash every slice carries — "users" and "/users" are the same
// request.
type RouteTriggerProps struct {
	Type       string
	Value      string
	Negate     bool
	IgnoreCase bool
	OnPath     bool
}

// RouteTrigger normalizes a trigger typed on the command line: its type
// defaults to equal and takes the aliases of RouteTriggerAliases, and every
// type but suffix and regex compared against a path slice gets its leading
// slash. A regex is taken verbatim, and has to compile.
func RouteTrigger(sandbox *api.Sandbox, props RouteTriggerProps) (routeconf.Trigger, error) {
	value := sandbox.Deps.Stringsdeps.TrimSpace(props.Value)
	raw_kind := sandbox.Deps.Stringsdeps.TrimSpace(props.Type)

	if value == "" {
		if raw_kind != "" {
			return routeconf.Trigger{}, sandbox.Deps.Std.Errorf("--trigger-type %q needs a --trigger to compare against", raw_kind)
		}
		if props.Negate || props.IgnoreCase {
			return routeconf.Trigger{}, sandbox.Deps.Std.Errorf("--trigger-negate and --trigger-ignore-case need a --trigger to apply to")
		}
		return routeconf.Trigger{}, nil
	}
	if raw_kind == "" {
		raw_kind = "equal"
	}
	kind, err := RouteTriggerType(sandbox, raw_kind)
	if err != nil {
		return routeconf.Trigger{}, err
	}

	if kind == "regex" {
		if _, err := sandbox.Deps.Stringsdeps.MatchPattern(value, ""); err != nil {
			return routeconf.Trigger{}, sandbox.Deps.Std.Errorf("invalid regex trigger %q: %s", value, err.Error())
		}
	} else if props.OnPath && kind != "suffix" {
		value = "/" + sandbox.Deps.Stringsdeps.TrimLeft(value, "/")
	}

	return routeconf.Trigger{
		Exists:     true,
		Type:       kind,
		Value:      value,
		Negate:     props.Negate,
		IgnoreCase: props.IgnoreCase,
	}, nil
}

// RouteEntryId turns a path id or a parameter key typed on the command line
// into the exported Go name its Entries field carries: "user-id" -> "UserId",
// "item" -> "Item".
func RouteEntryId(sandbox *api.Sandbox, raw string) string {
	parts := sandbox.Deps.Stringsdeps.FieldsFunc(sandbox.Deps.Stringsdeps.TrimSpace(raw), func(r rune) bool {
		return r == '-' || r == '_' || r == ' ' || r == '.'
	})
	id := ""
	for _, part := range parts {
		id += sandbox.Deps.Stringsdeps.ToUpper(part[:1]) + part[1:]
	}
	return id
}

// RouteReservedIds are the Entries fields every route carries whatever it
// declares, so no path or parameter may take them.
var RouteReservedIds = []string{"FullRoute", "Body"}

// RouteSegmentIndex reads a --start or --end typed on the command line, "" as
// the fallback given.
func RouteSegmentIndex(sandbox *api.Sandbox, label string, raw string, fallback int) (int, error) {
	raw = sandbox.Deps.Stringsdeps.TrimSpace(raw)
	if raw == "" {
		return fallback, nil
	}
	value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
	if err != nil {
		return 0, sandbox.Deps.Std.Errorf("--%s %q is not a whole number", label, raw)
	}
	return value, nil
}

// NewRoutePath builds a routeconf.Path from the raw values typed on the
// command line, holding it to the rules verify checks: an id, a start that is
// not negative and an end that is -1 or not before it.
func NewRoutePath(sandbox *api.Sandbox, props api.RoutePathProps) (routeconf.Path, error) {
	path := routeconf.Path{
		Id:          RouteEntryId(sandbox, props.Id),
		Description: sandbox.Deps.Stringsdeps.TrimSpace(props.Description),
	}
	if path.Id == "" {
		return path, sandbox.Deps.Std.Errorf("a path needs an id")
	}

	start, err := RouteSegmentIndex(sandbox, "start", props.Start, 0)
	if err != nil {
		return path, err
	}
	end, err := RouteSegmentIndex(sandbox, "end", props.End, routeconf.LastSegment)
	if err != nil {
		return path, err
	}
	if start < 0 {
		return path, sandbox.Deps.Std.Errorf("--start %d is negative: a slice starts at segment 0 or later", start)
	}
	if end != routeconf.LastSegment && end < start {
		return path, sandbox.Deps.Std.Errorf("--end %d is before --start %d: use -1 for the last segment", end, start)
	}
	path.Start, path.End = start, end

	kind := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(props.Type))
	if kind == "" {
		kind = routeconf.DefaultPathType
	}
	if !contains(routeconf.PathTypes, kind) {
		return path, sandbox.Deps.Std.Errorf("unknown path type %q (use one of %s)", props.Type, sandbox.Deps.Stringsdeps.Join(routeconf.PathTypes, ", "))
	}
	if kind != routeconf.DefaultPathType && start != end {
		return path, sandbox.Deps.Std.Errorf("a path of type %s reads one segment: give it the same --start and --end", kind)
	}
	path.Type = kind

	trigger, err := RouteTrigger(sandbox, RouteTriggerProps{
		Type:       props.TriggerType,
		Value:      props.Trigger,
		Negate:     props.TriggerNegate,
		IgnoreCase: props.TriggerIgnoreCase,
		OnPath:     true,
	})
	if err != nil {
		return path, err
	}
	path.Trigger = trigger

	return path, nil
}

// NewRouteParameter builds a routeconf.Parameter from the raw values typed on
// the command line: a known type, known fonts (the query string when none is
// given), and a default that parses as the type and never sits beside
// `required`.
func NewRouteParameter(sandbox *api.Sandbox, props api.RouteParameterProps) (routeconf.Parameter, error) {
	parameter := routeconf.Parameter{
		Key:         sandbox.Deps.Stringsdeps.TrimSpace(props.Name),
		Description: sandbox.Deps.Stringsdeps.TrimSpace(props.Description),
		Examples:    props.Examples,
		Required:    props.Required,
		Fonts:       []string{},
	}
	parameter.Id = RouteEntryId(sandbox, parameter.Key)
	if parameter.Id == "" {
		return parameter, sandbox.Deps.Std.Errorf("a parameter needs a name")
	}

	kind := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(props.Type))
	if kind == "" {
		kind = "string"
	}
	if !contains(routeconf.ParameterTypes, kind) {
		return parameter, sandbox.Deps.Std.Errorf("unknown parameter type %q (use one of %s)", props.Type, sandbox.Deps.Stringsdeps.Join(routeconf.ParameterTypes, ", "))
	}
	parameter.Type = kind

	for _, raw := range props.Fonts {
		font := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
		if !contains(routeconf.ParameterFonts, font) {
			return parameter, sandbox.Deps.Std.Errorf("unknown font %q (use one of %s)", raw, sandbox.Deps.Stringsdeps.Join(routeconf.ParameterFonts, ", "))
		}
		if !contains(parameter.Fonts, font) {
			parameter.Fonts = append(parameter.Fonts, font)
		}
	}
	if len(parameter.Fonts) == 0 {
		parameter.Fonts = []string{"query"}
	}

	if parameter.Required && kind == "boolean" {
		return parameter, sandbox.Deps.Std.Errorf("a boolean parameter cannot be required (its absence already means false)")
	}

	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Default); value != "" {
		if parameter.Required {
			return parameter, sandbox.Deps.Std.Errorf("a parameter cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := RouteCheckParameterLiteral(sandbox, kind, value); err != nil {
			return parameter, err
		}
		parameter.Default, parameter.HasDefault = value, true
	}

	trigger, err := RouteTrigger(sandbox, RouteTriggerProps{
		Type:       props.TriggerType,
		Value:      props.Trigger,
		Negate:     props.TriggerNegate,
		IgnoreCase: props.TriggerIgnoreCase,
	})
	if err != nil {
		return parameter, err
	}
	parameter.Trigger = trigger

	return parameter, nil
}

// RouteCheckParameterLiteral reports whether a raw default parses as the
// parameter type it is declared for.
func RouteCheckParameterLiteral(sandbox *api.Sandbox, kind string, raw string) error {
	switch kind {
	case "integer", "integer-array":
		if _, err := sandbox.Deps.Stringsdeps.Atoi(raw); err != nil {
			return sandbox.Deps.Std.Errorf("default %q is not a whole number", raw)
		}
	case "number":
		if _, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64); err != nil {
			return sandbox.Deps.Std.Errorf("default %q is not a number", raw)
		}
	case "boolean":
		if raw != "true" && raw != "false" {
			return sandbox.Deps.Std.Errorf("default %q is not true or false", raw)
		}
	}
	return nil
}

// FindRoutePath returns the index of the path whose id is the one typed, or
// -1.
func FindRoutePath(sandbox *api.Sandbox, paths []routeconf.Path, id string) int {
	wanted := RouteEntryId(sandbox, id)
	for i, path := range paths {
		if path.Id == wanted {
			return i
		}
	}
	return -1
}

// FindRouteParameter returns the index of the parameter read under the key
// typed — or whose id is its exported spelling — or -1.
func FindRouteParameter(sandbox *api.Sandbox, parameters []routeconf.Parameter, name string) int {
	key := sandbox.Deps.Stringsdeps.TrimSpace(name)
	id := RouteEntryId(sandbox, key)
	for i, parameter := range parameters {
		if parameter.Key == key || parameter.Id == id {
			return i
		}
	}
	return -1
}

// RouteIdTaken reports whether id is already an Entries field of the route —
// reserved, a path or a parameter — skipping the entry being edited, named by
// its current id ("" skips nothing).
func RouteIdTaken(conf *routeconf.RouteConf, id string, editing string) bool {
	if contains(RouteReservedIds, id) {
		return true
	}
	for _, path := range conf.Paths {
		if path.Id == id && path.Id != editing {
			return true
		}
	}
	for _, parameter := range conf.Parameters {
		if parameter.Id == id && parameter.Id != editing {
			return true
		}
	}
	return false
}

// CheckRoutePosition validates a --position against the list it will be
// inserted into, the same way CheckPosition does for a command's fields.
func CheckRoutePosition(sandbox *api.Sandbox, kind string, position int, size int) (int, error) {
	if position == AppendPosition {
		return size, nil
	}
	if position < 0 {
		return 0, sandbox.Deps.Std.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, size)
	}
	if position > size {
		return 0, sandbox.Deps.Std.Errorf("--position %d is out of range: this route has %d %s(s), so the accepted range is 0 to %d", position, size, kind, size)
	}
	return position, nil
}

// InsertRoutePath places path at position inside paths (appending when
// position is negative or past the end).
func InsertRoutePath(paths []routeconf.Path, path routeconf.Path, position int) []routeconf.Path {
	if position < 0 || position >= len(paths) {
		return append(paths, path)
	}
	out := make([]routeconf.Path, 0, len(paths)+1)
	out = append(out, paths[:position]...)
	out = append(out, path)
	out = append(out, paths[position:]...)
	return out
}

// RemoveRoutePath drops the path at index from paths.
func RemoveRoutePath(paths []routeconf.Path, index int) []routeconf.Path {
	out := make([]routeconf.Path, 0, len(paths)-1)
	out = append(out, paths[:index]...)
	out = append(out, paths[index+1:]...)
	return out
}

// InsertRouteParameter places parameter at position inside parameters
// (appending when position is negative or past the end).
func InsertRouteParameter(parameters []routeconf.Parameter, parameter routeconf.Parameter, position int) []routeconf.Parameter {
	if position < 0 || position >= len(parameters) {
		return append(parameters, parameter)
	}
	out := make([]routeconf.Parameter, 0, len(parameters)+1)
	out = append(out, parameters[:position]...)
	out = append(out, parameter)
	out = append(out, parameters[position:]...)
	return out
}

// RemoveRouteParameter drops the parameter at index from parameters.
func RemoveRouteParameter(parameters []routeconf.Parameter, index int) []routeconf.Parameter {
	out := make([]routeconf.Parameter, 0, len(parameters)-1)
	out = append(out, parameters[:index]...)
	out = append(out, parameters[index+1:]...)
	return out
}

// RouteMethod normalizes an http method, refusing one no route may answer.
// ANY — or * — accepts every method.
func RouteMethod(sandbox *api.Sandbox, raw string) (string, error) {
	method := sandbox.Deps.Stringsdeps.ToUpper(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if method == "" {
		return routeconf.DefaultMethod, nil
	}
	if method == "*" || method == routeconf.AnyMethod {
		return routeconf.AnyMethod, nil
	}
	for _, known := range RouteMethods {
		if known == method {
			return method, nil
		}
	}
	return "", sandbox.Deps.Std.Errorf("unknown method %q (use one of %s, or ANY)", raw, sandbox.Deps.Stringsdeps.Join(RouteMethods, ", "))
}

// RouteMethodList normalizes a repeated --method into the list route.yaml
// declares, deduplicated in the order typed; none at all is GET. ANY stands
// alone, since it already accepts every method.
func RouteMethodList(sandbox *api.Sandbox, raws []string) ([]string, error) {
	methods := []string{}
	for _, raw := range raws {
		method, err := RouteMethod(sandbox, raw)
		if err != nil {
			return nil, err
		}
		if !contains(methods, method) {
			methods = append(methods, method)
		}
	}
	if len(methods) == 0 {
		methods = []string{routeconf.DefaultMethod}
	}
	if contains(methods, routeconf.AnyMethod) && len(methods) > 1 {
		return nil, sandbox.Deps.Std.Errorf("--method ANY already accepts every method: pass it alone")
	}
	return methods, nil
}

// RouteMethods is every http method a route may declare.
var RouteMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// contains reports whether list holds value.
func contains(list []string, value string) bool {
	for _, one := range list {
		if one == value {
			return true
		}
	}
	return false
}

// RouteFieldName normalizes a name typed on the command line — a body
// property, say — trimming only the surrounding space: the punctuation of an
// external spelling is its own to keep.
func RouteFieldName(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.Stringsdeps.TrimSpace(name)
}

// RouteCheckLiteral reports whether a raw command-line literal parses as the
// given type.
func RouteCheckLiteral(sandbox *api.Sandbox, kind string, label string, raw string) error {
	return checkLiteral(sandbox, kind, label, raw)
}

// RouteParseBound reads a min/max literal for a numeric body property.
func RouteParseBound(sandbox *api.Sandbox, kind string, label string, raw string) (float64, error) {
	return parseBound(sandbox, kind, label, raw)
}
