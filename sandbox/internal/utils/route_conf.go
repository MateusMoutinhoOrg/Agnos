package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// The route helpers below mirror the command ones of command_conf.go, one for
// one: a route is declared, edited and looked up exactly the way a command is,
// only against route.yaml instead of command.yaml.

// RouteName normalizes a user-typed route name into its canonical
// spelling: lowercased, spaces and underscores turned into dashes
// ("Create User" -> "create-user").
func RouteName(sandbox *api.Sandbox, name string) string {
	return CommandName(sandbox, name)
}

// ValidateRouteName reports whether a user-typed route name normalizes to a
// usable identifier. It becomes a directory name and a Go package clause, so
// the same alphabet a command name is held to applies.
func ValidateRouteName(sandbox *api.Sandbox, name string) error {
	identifier := RouteName(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.StdDeps.Errorf("a route needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.StdDeps.Errorf("invalid route name %q: a route name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.StdDeps.Errorf(
				"invalid route name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory %s of sandbox/internal/routes and a Go package name)",
				name, RoutePackage(sandbox, name))
		}
	}
	return nil
}

// RoutePackage is the Go package / directory name for a route: the identifier
// with dashes turned into underscores ("create-user" -> "create_user").
func RoutePackage(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.StringsDeps.ReplaceAll(RouteName(sandbox, name), "-", "_")
}

// RoutesDir is the tree the routes are declared in, the server layer's mirror
// of sandbox/internal/commands: every directory at any depth holding a
// route.yaml is a route, and every other one a folder grouping them.
const RoutesDir = "sandbox/internal/routes"

// RouteConfFile is the declaration of one route, and what makes its
// directory a route.
const RouteConfFile = "route.yaml"

// GeneratedRoutes is every route the server group renders itself, each at the
// top of RoutesDir: health, which answers that the server is up, and openapi,
// which answers the OpenAPI document of every route. No verb declares, removes,
// renames or renumbers one of them.
func GeneratedRoutes() []string {
	return []string{"health", "openapi"}
}

// IsGeneratedRoute reports whether a user-typed route name spells one of
// GeneratedRoutes.
func IsGeneratedRoute(sandbox *api.Sandbox, name string) bool {
	return contains(GeneratedRoutes(), RoutePackage(sandbox, name))
}

// RouteDirs is every route declared under RoutesDir, at any depth. A
// generated route whose route.yaml is still pending in the transaction — the
// first build after the server group gained it — is listed too, so the build
// that renders it also renders its generated.new.go and generated.input.go.
func RouteDirs(sandbox *api.Sandbox, io *stagedfs.StagedFS) []UnitDir {
	units := FindUnitDirs(sandbox, io, RoutesDir, RouteConfFile)
	for _, name := range GeneratedRoutes() {
		if _, found := FindUnitDirIn(units, name); found {
			continue
		}
		dir := RoutesDir + "/" + name
		if _, err := io.ReadFile(dir + "/" + RouteConfFile); err == nil {
			units = append(units, UnitDir{Name: name, Dir: dir})
		}
	}
	return units
}

// RouteDir is the project-relative directory holding the route named, in
// whatever folder it sits; one no route.yaml declares yet lands at the top of
// RoutesDir.
func RouteDir(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) string {
	pkg := RoutePackage(sandbox, name)
	if dir, found := FindUnitDir(sandbox, io, RoutesDir, RouteConfFile, pkg); found {
		return dir
	}
	return RoutesDir + "/" + pkg
}

// RouteConfPath is the project-relative path of a route's route.yaml.
func RouteConfPath(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) string {
	return RouteDir(sandbox, io, name) + "/" + RouteConfFile
}

// LoadRouteConf reads and parses the route.yaml of the route named, wherever
// under sandbox/internal/routes it sits.
func LoadRouteConf(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) (*routeconf.RouteConf, error) {
	if err := ValidateRouteName(sandbox, name); err != nil {
		return nil, err
	}
	path := RouteConfPath(sandbox, io, name)
	content, err := io.ReadFile(path)
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("route %q not found in %s", RouteName(sandbox, name), RoutesDir)
	}
	conf, err := routeconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.StdDeps.Errorf("%s: %w", path, err)
	}
	return conf, nil
}

// SaveRouteConf renders conf back over the route.yaml of the route named.
func SaveRouteConf(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, conf *routeconf.RouteConf) error {
	return io.WriteFile(RouteConfPath(sandbox, io, name), []byte(conf.Render()))
}

// GoIdentifier is the exported Go name a kebab-case, snake_case or spaced name
// is spelled as: every word title-cased and joined — "out-file" is OutFile,
// "token_sha256" TokenSha256, "api products" ApiProducts. It is the one rule
// every generated identifier follows: the Input field an arg, a flag, a path
// or a parameter binds to, a database record, a table's methods.
func GoIdentifier(sandbox *api.Sandbox, raw string) string {
	parts := sandbox.Deps.StringsDeps.FieldsFunc(sandbox.Deps.StringsDeps.TrimSpace(raw), func(r rune) bool {
		return r == '-' || r == '_' || r == ' ' || r == '.'
	})
	id := ""
	for _, part := range parts {
		id += sandbox.Deps.StringsDeps.ToUpper(part[:1]) + part[1:]
	}
	return id
}

// ValidateEntryId reports whether id — the Input field the name raw typed on
// the command line becomes — is one the generated.input.go can spell: an
// ASCII letter first, then ASCII letters and digits. It runs before anything is
// written, so a name the Go compiler would refuse never reaches a declaration.
func ValidateEntryId(sandbox *api.Sandbox, kind string, raw string, id string) error {
	if id == "" {
		return sandbox.Deps.StdDeps.Errorf("a %s needs a name", kind)
	}
	first := id[0]
	if !((first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z')) {
		return sandbox.Deps.StdDeps.Errorf("invalid %s name %q: it must start with an ASCII letter (it becomes the Go field Input.%s)", kind, raw, id)
	}
	for i := 0; i < len(id); i++ {
		letter := id[i]
		valid := (letter >= 'A' && letter <= 'Z') ||
			(letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9')
		if !valid {
			return sandbox.Deps.StdDeps.Errorf("invalid %s name %q: only ASCII letters, digits, dashes, underscores, dots and spaces are allowed (it becomes the Go field Input.%s)", kind, raw, id)
		}
	}
	return nil
}

// RouteReservedIds are the Input fields the generated.input.go spells
// itself — FullRoute on every route, Body on one declaring a body — so no
// path or parameter may take them.
var RouteReservedIds = []string{"FullRoute", "Body"}

// RouteSegmentIndex reads a --start or --end typed on the command line, "" as
// the fallback given.
func RouteSegmentIndex(sandbox *api.Sandbox, label string, raw string, fallback int) (int, error) {
	raw = sandbox.Deps.StringsDeps.TrimSpace(raw)
	if raw == "" {
		return fallback, nil
	}
	value, err := sandbox.Deps.StringsDeps.Atoi(raw)
	if err != nil {
		return 0, sandbox.Deps.StdDeps.Errorf("--%s %q is not a whole number", label, raw)
	}
	return value, nil
}

// NewRoutePath builds a routeconf.Path from the raw values typed on the
// command line, holding it to the rules verify checks: an id, a start that is
// not negative and an end that is -1 or not before it.
func NewRoutePath(sandbox *api.Sandbox, props api.AddPathProps) (routeconf.Path, error) {
	path := routeconf.Path{
		Id:          GoIdentifier(sandbox, props.Name),
		Description: sandbox.Deps.StringsDeps.TrimSpace(props.Description),
	}
	if err := ValidateEntryId(sandbox, "path", props.Name, path.Id); err != nil {
		return path, err
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
		return path, sandbox.Deps.StdDeps.Errorf("--start %d is negative: a slice starts at segment 0 or later", start)
	}
	if end != routeconf.LastSegment && end < start {
		return path, sandbox.Deps.StdDeps.Errorf("--end %d is before --start %d: use -1 for the last segment", end, start)
	}
	path.Start, path.End = start, end

	kind := sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(props.Type))
	if kind == "" {
		kind = routeconf.DefaultPathType
	}
	if !contains(routeconf.PathTypes, kind) {
		return path, sandbox.Deps.StdDeps.Errorf("unknown path type %q (use one of %s)", props.Type, sandbox.Deps.StringsDeps.Join(routeconf.PathTypes, ", "))
	}
	if kind != routeconf.DefaultPathType && start != end {
		return path, sandbox.Deps.StdDeps.Errorf("a path of type %s reads one segment: give it the same --start and --end", kind)
	}
	path.Type = kind

	trigger, err := NewTrigger(sandbox, TriggerProps{
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
// the command line: a known type, known sources (the query string when none is
// given), and a default that parses as the type and never sits beside
// `required`.
func NewRouteParameter(sandbox *api.Sandbox, props api.AddParameterProps) (routeconf.Parameter, error) {
	parameter := routeconf.Parameter{
		Key:         sandbox.Deps.StringsDeps.TrimSpace(props.Name),
		Description: sandbox.Deps.StringsDeps.TrimSpace(props.Description),
		Examples:    props.Examples,
		Required:    props.Required,
		Sources:     []string{},
	}
	parameter.Id = GoIdentifier(sandbox, parameter.Key)
	if err := ValidateEntryId(sandbox, "parameter", parameter.Key, parameter.Id); err != nil {
		return parameter, err
	}

	kind := sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(props.Type))
	if kind == "" {
		kind = "string"
	}
	if !contains(routeconf.ParameterTypes, kind) {
		return parameter, sandbox.Deps.StdDeps.Errorf("unknown parameter type %q (use one of %s)", props.Type, sandbox.Deps.StringsDeps.Join(routeconf.ParameterTypes, ", "))
	}
	parameter.Type = kind

	for _, raw := range props.Sources {
		source := sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(raw))
		if !contains(routeconf.ParameterSources, source) {
			return parameter, sandbox.Deps.StdDeps.Errorf("unknown source %q (use one of %s)", raw, sandbox.Deps.StringsDeps.Join(routeconf.ParameterSources, ", "))
		}
		if !contains(parameter.Sources, source) {
			parameter.Sources = append(parameter.Sources, source)
		}
	}
	if len(parameter.Sources) == 0 {
		parameter.Sources = []string{"query"}
	}

	if parameter.Required && kind == "boolean" {
		return parameter, sandbox.Deps.StdDeps.Errorf("a boolean parameter cannot be required (its absence already means false)")
	}

	if value := sandbox.Deps.StringsDeps.TrimSpace(props.Default); value != "" {
		if parameter.Required {
			return parameter, sandbox.Deps.StdDeps.Errorf("a parameter cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := CheckRouteParameterLiteral(sandbox, kind, value); err != nil {
			return parameter, err
		}
		parameter.Default, parameter.HasDefault = value, true
	}

	trigger, err := NewTrigger(sandbox, TriggerProps{
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

// CheckRouteParameterLiteral reports whether a raw default parses as the
// parameter type it is declared for.
func CheckRouteParameterLiteral(sandbox *api.Sandbox, kind string, raw string) error {
	switch kind {
	case "integer", "integer-array":
		if _, err := sandbox.Deps.StringsDeps.Atoi(raw); err != nil {
			return sandbox.Deps.StdDeps.Errorf("default %q is not a whole number", raw)
		}
	case "number":
		if _, err := sandbox.Deps.StringsDeps.ParseFloat(raw, 64); err != nil {
			return sandbox.Deps.StdDeps.Errorf("default %q is not a number", raw)
		}
	case "boolean":
		if raw != "true" && raw != "false" {
			return sandbox.Deps.StdDeps.Errorf("default %q is not true or false", raw)
		}
	}
	return nil
}

// FindRoutePath returns the index of the path whose id is the one typed, or
// -1.
func FindRoutePath(sandbox *api.Sandbox, paths []routeconf.Path, id string) int {
	wanted := GoIdentifier(sandbox, id)
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
	key := sandbox.Deps.StringsDeps.TrimSpace(name)
	id := GoIdentifier(sandbox, key)
	for i, parameter := range parameters {
		if parameter.Key == key || parameter.Id == id {
			return i
		}
	}
	return -1
}

// RouteIdTaken reports whether id is already an Input field of the route —
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
// inserted into, the same way CheckCommandPosition does for a command's fields.
func CheckRoutePosition(sandbox *api.Sandbox, kind string, position int, size int) (int, error) {
	if position == AppendPosition {
		return size, nil
	}
	if position < 0 {
		return 0, sandbox.Deps.StdDeps.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, size)
	}
	if position > size {
		return 0, sandbox.Deps.StdDeps.Errorf("--position %d is out of range: this route has %d %s(s), so the accepted range is 0 to %d", position, size, kind, size)
	}
	return position, nil
}

// RouteMethod normalizes an http method, refusing one no route may answer.
// ANY — or * — accepts every method.
func RouteMethod(sandbox *api.Sandbox, raw string) (string, error) {
	method := sandbox.Deps.StringsDeps.ToUpper(sandbox.Deps.StringsDeps.TrimSpace(raw))
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
	return "", sandbox.Deps.StdDeps.Errorf("unknown method %q (use one of %s, or ANY)", raw, sandbox.Deps.StringsDeps.Join(RouteMethods, ", "))
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
		return nil, sandbox.Deps.StdDeps.Errorf("--method ANY already accepts every method: pass it alone")
	}
	return methods, nil
}

// ValidateMediaType reports whether raw is a media type a Content-Type header
// can carry: a type and a subtype of token characters ("application/json"),
// optionally followed by parameters ("text/plain; charset=utf-8").
func ValidateMediaType(sandbox *api.Sandbox, raw string) error {
	strs := sandbox.Deps.StringsDeps
	essence := strs.TrimSpace(strs.Split(raw, ";")[0])
	parts := strs.Split(essence, "/")
	if len(parts) != 2 || !isMediaToken(parts[0]) || !isMediaToken(parts[1]) {
		return sandbox.Deps.StdDeps.Errorf("invalid response type %q: it is sent as the Content-Type header, so it must be a media type such as application/json or text/plain", raw)
	}
	return nil
}

// isMediaToken reports whether word is a non-empty run of the characters a
// media type's type or subtype may hold.
func isMediaToken(word string) bool {
	if word == "" {
		return false
	}
	for i := 0; i < len(word); i++ {
		letter := word[i]
		valid := (letter >= 'a' && letter <= 'z') || (letter >= 'A' && letter <= 'Z') ||
			(letter >= '0' && letter <= '9') || containsByte("!#$&-^_.+", letter)
		if !valid {
			return false
		}
	}
	return true
}

// RequireBodyMethod refuses to declare a body on a route answering only
// methods that carry none — GET and HEAD — since no request it answers would
// ever send one. ANY, or any other method beside them, may.
func RequireBodyMethod(sandbox *api.Sandbox, conf *routeconf.RouteConf, route string) error {
	for _, method := range conf.Methods {
		if method != "GET" && method != "HEAD" {
			return nil
		}
	}
	return sandbox.Deps.StdDeps.Errorf("route %s answers only %s, which carry no body: add a method that does first (set-route %s --method POST)",
		RouteName(sandbox, route), sandbox.Deps.StringsDeps.Join(conf.Methods, ", "), RouteName(sandbox, route))
}

// RouteMethods is every http method a route may declare.
var RouteMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// RouteFieldName normalizes a name typed on the command line — a body
// property, say — trimming only the surrounding space: the punctuation of an
// external spelling is its own to keep.
func RouteFieldName(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.StringsDeps.TrimSpace(name)
}

// CheckRouteLiteral reports whether a raw command-line literal parses as the
// given type.
func CheckRouteLiteral(sandbox *api.Sandbox, kind string, label string, raw string) error {
	return checkLiteral(sandbox, kind, label, raw)
}

// ParseRouteBound reads a min/max literal for a numeric body property.
func ParseRouteBound(sandbox *api.Sandbox, kind string, label string, raw string) (float64, error) {
	return parseBound(sandbox, kind, label, raw)
}

// RouteFieldType maps the type spellings accepted on the command line onto the
// canonical route.yaml set; "" defaults to string.
func RouteFieldType(sandbox *api.Sandbox, raw string) (string, bool) {
	switch sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(raw)) {
	case "", "string", "str":
		return "string", true
	case "bool", "boolean":
		return "boolean", true
	case "int", "integer":
		return "int", true
	case "float", "double", "number":
		return "float", true
	default:
		return "", false
	}
}

func checkLiteral(sandbox *api.Sandbox, kind string, label string, raw string) error {
	switch kind {
	case "boolean":
		if raw != "true" && raw != "false" {
			return sandbox.Deps.StdDeps.Errorf("%s for a boolean must be true or false, got %q", label, raw)
		}
	case "int":
		if _, err := sandbox.Deps.StringsDeps.ParseInt(raw, 10, 64); err != nil {
			return sandbox.Deps.StdDeps.Errorf("%s must be an int, got %q", label, raw)
		}
	case "float":
		if _, err := sandbox.Deps.StringsDeps.ParseFloat(raw, 64); err != nil {
			return sandbox.Deps.StdDeps.Errorf("%s must be a float, got %q", label, raw)
		}
	}
	return nil
}

func parseBound(sandbox *api.Sandbox, kind string, label string, raw string) (float64, error) {
	if kind != "int" && kind != "float" {
		return 0, sandbox.Deps.StdDeps.Errorf("%s only applies to int/float fields", label)
	}
	if err := checkLiteral(sandbox, kind, label, raw); err != nil {
		return 0, err
	}
	value, _ := sandbox.Deps.StringsDeps.ParseFloat(raw, 64)
	return value, nil
}
