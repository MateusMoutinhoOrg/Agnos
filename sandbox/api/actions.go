package api

// The runtimes `build` can hand a rendered project to. RuntimeGo resolves the
// module graph and compiles every package, so a build that reports success is
// a build the Go toolchain accepted; RuntimeNone renders only, which is what
// the removal commands use — dropping a command or a flag may leave
// hand-written code referring to what is gone.
const (
	// RuntimeGo resolves the module graph and compiles every package after
	// the render.
	RuntimeGo = "go"
	// RuntimeNone renders only, leaving the result unchecked.
	RuntimeNone = "none"
)

// BuildProps describes one (re)render of a project: the directory holding it
// and the runtime that then checks the result.
type BuildProps struct {
	Path    string
	Runtime string
}

// CompileProps describes one cross-compile run: the directory holding the
// project and the target names to build. Each name is one of the keys
// `agnos compile` accepts (linux86, linuxarm64, linuxi32, mac86, macarm64,
// windows86, windowsi32) or "all", which expands to every target.
type CompileProps struct {
	Path    string
	Targets []string
}

// StartProps describes one project to scaffold: the directory to write it
// into, the name it carries in <Name>Config/project.yaml, the module path for
// go.mod (nil derives it from the name) and whether an existing directory may
// be written over.
type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

// ExecTestProps describes one run of the project's example suite: the
// directory holding the project, the single example name to run (empty runs
// every one, both sides) and whether the goldens are rewritten with what the
// run produced instead of compared against it.
type ExecTestProps struct {
	Path   string
	Only   string
	Update bool
}

// FieldProps describes one flag or positional arg to add to a command's
// entries.yaml. Default, Min and Max are the raw literals typed on the
// command line ("" means unset) so the action can tell "not given" from a
// zero value; Position is the index to insert at (< 0 appends).
type FieldProps struct {
	Path        string
	Command     string
	Name        string
	Identifiers []string
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
	Array       bool
	Min         string
	Max         string
	Position    int
}

// CommandProps carries the command-level keys of entries.yaml that
// set-command may rewrite. Empty strings leave the current value alone;
// Identifiers / Examples are appended (deduplicated).
type CommandProps struct {
	Path            string
	Command         string
	Help            string
	Category        string
	LongDescription string
	Hidden          bool
	Visible         bool
	Identifiers     []string
	Examples        []string
}

// RouteProps carries the route-level keys of route.yaml that set-route may
// rewrite. Empty strings leave the current value alone; Examples are appended
// (deduplicated), and Hidden / Visible are the two sides of one switch.
type RouteProps struct {
	Path            string
	Route           string
	Method          string
	Help            string
	Category        string
	LongDescription string
	Hidden          bool
	Visible         bool
	Examples        []string
}

// RouteFieldProps describes one field to add to a route's route.yaml. It
// covers the three origins that read a value off the request line — a captured
// path segment, a header and a query parameter — which differ only in where
// the entry lands. Identifier declares a literal path segment instead of a
// captured one, and is normalized to start with "/". Default, Min and Max are
// the raw literals typed on the command line ("" means unset); Position is the
// index to insert at (< 0 appends).
type RouteFieldProps struct {
	Path        string
	Route       string
	Name        string
	Identifier  string
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
	Array       bool
	Min         string
	Max         string
	Position    int
}

// RouteBodyProps describes the body envelope of one route — everything about
// the request body but the json-schema, which is grown property by property
// with AddBodyField. Type is "none", "raw", "text" or "json"; Required and
// Optional are the two sides of one switch, as are the empty strings and
// MaxBytes < 0 that mean "leave as is". DropSchema deletes the declared
// json-schema.
type RouteBodyProps struct {
	Path        string
	Route       string
	Type        string
	Required    bool
	Optional    bool
	MaxBytes    int
	ContentType string
	DropSchema  bool
}

// RouteBodyFieldProps describes one property of a route's body json-schema.
// Name is the dotted path it sits at ("address.city"), and every other field
// is one keyword of the supported subset: the raw literals typed on the
// command line, where "" means unset. Type is "string", "boolean", "int",
// "float" or "object", and Array wraps the whole of it in an array schema.
// AdditionalProperties and NoAdditionalProperties are the two sides of one
// switch.
type RouteBodyFieldProps struct {
	Path                   string
	Route                  string
	Name                   string
	Type                   string
	Required               bool
	Array                  bool
	Min                    string
	Max                    string
	ExclusiveMin           string
	ExclusiveMax           string
	Format                 string
	Pattern                string
	Enum                   []string
	Const                  string
	Nullable               bool
	MinItems               string
	MaxItems               string
	UniqueItems            bool
	AdditionalProperties   bool
	NoAdditionalProperties bool
}

// DocProps describes one doc to create under docs/. Name is the doc's
// directory, optionally nested under its parent ("PublicApi/api.Actions").
// Themes are the theme ids of <ProjectName>Config/themes.yaml the doc belongs
// to: required on a first-level doc, forbidden on a sub-doc.
type DocProps struct {
	Path        string
	Name        string
	Description string
	Themes      []string
}

// Actions is the whole set of operations agnos performs on a project. Every
// field takes the project directory as its first input (`path`, or the Path of
// a props struct) and reports failure as an error; the ones that change the
// tree re-render it before returning, so a project is always left in a built
// state.
type Actions struct {
	// Build re-renders every generated file of the project and hands the
	// result to the runtime named by the props.
	Build func(props BuildProps) error

	// Compile cross-compiles the project's cmd/ binaries into release/, one
	// file per named target.
	Compile func(props CompileProps) error

	// Verify checks the project against the schema every generator assumes
	// and writes nothing; it reports every violation at once.
	Verify func(path string) error

	// Start scaffolds a new project: the config directory, go.mod, the
	// sandbox skeleton and a first build.
	Start func(props StartProps) error

	// DepsInit adds the dependency layer (sandbox/deps/ and
	// adapters/availables/standard/) to a project that has none.
	DepsInit func(path string) error

	// DepsPurge removes the dependency layer and every installed dep with it.
	DepsPurge func(path string) error

	// DepInstall installs one dep of the built-in list: its contract under
	// sandbox/deps/, its adapter under adapters/libs/ and its go.mod require.
	DepInstall func(path string, dep string) error

	// DepRemove uninstalls one installed dep, contract, adapter and require.
	DepRemove func(path string, dep string) error

	// DepList returns the names of the deps installed in the project.
	DepList func(path string) ([]string, error)

	// CliInit adds the CLI layer (cmd/main, the dispatcher and the help and
	// version commands) to a project that has none.
	CliInit func(path string) error

	// CliPurge removes the CLI layer and every command declared in it.
	CliPurge func(path string) error

	// AddCommand declares a new command: its entries.yaml, its generated
	// entries.go and a handler.go to fill in.
	AddCommand func(path string, name string, help string, category string) error

	// RemoveCommand deletes one command and unwires it from the dispatcher.
	RemoveCommand func(path string, name string) error

	// SetCommand rewrites the command-level keys of one command's
	// entries.yaml.
	SetCommand func(props CommandProps) error

	// AddFlag declares one flag on a command.
	AddFlag func(props FieldProps) error

	// RemoveFlag deletes one declared flag from a command.
	RemoveFlag func(path string, command string, name string) error

	// AddArg declares one positional argument on a command.
	AddArg func(props FieldProps) error

	// RemoveArg deletes one declared positional argument from a command.
	RemoveArg func(path string, command string, name string) error

	// ServerInit adds the http server layer (sandbox/internal/server, the
	// routeio package, the health route and the start-server command) to a
	// project that has none, installing the CLI layer first when it is
	// missing.
	ServerInit func(path string) error

	// ServerPurge removes the server layer and every route declared in it.
	ServerPurge func(path string) error

	// AddRoute declares a new route: its route.yaml, its generated
	// entries.go and a handler.go to fill in.
	AddRoute func(path string, name string, method string, trigger string, help string, category string) error

	// RemoveRoute deletes one route and unwires it from the dispatch.
	RemoveRoute func(path string, name string) error

	// SetRoute rewrites the route-level keys of one route's route.yaml.
	SetRoute func(props RouteProps) error

	// AddSegment appends one segment to a route's path: a literal one when
	// props.Identifier is set, a captured one otherwise.
	AddSegment func(props RouteFieldProps) error

	// RemoveSegment deletes one segment from a route's path, named either
	// by its capture name or by the identifier it spells.
	RemoveSegment func(path string, route string, name string) error

	// AddHeader declares one request header on a route.
	AddHeader func(props RouteFieldProps) error

	// RemoveHeader deletes one declared header from a route.
	RemoveHeader func(path string, route string, name string) error

	// AddParam declares one query-string parameter on a route.
	AddParam func(props RouteFieldProps) error

	// RemoveParam deletes one declared query parameter from a route.
	RemoveParam func(path string, route string, name string) error

	// SetBody rewrites the body keys of one route's route.yaml.
	SetBody func(props RouteBodyProps) error

	// AddBodyField declares one property of a route's body json-schema, at
	// the dotted path props.Name.
	AddBodyField func(props RouteBodyFieldProps) error

	// RemoveBodyField deletes one property from a route's body
	// json-schema.
	RemoveBodyField func(path string, route string, name string) error

	// AddDoc creates one doc directory under docs/, with its props.yaml and
	// a doc.md to fill in.
	AddDoc func(props DocProps) error

	// RemoveDoc deletes one doc directory and everything under it.
	RemoveDoc func(path string, name string) error

	// AddCliExample creates one example under examples/cli/, with an
	// example.sh stub that already runs.
	AddCliExample func(path string, name string) error

	// RemoveCliExample deletes one example of examples/cli/ whole.
	RemoveCliExample func(path string, name string) error

	// AddLibExample creates one example under examples/lib/, with an
	// example.go stub that already runs.
	AddLibExample func(path string, name string) error

	// RemoveLibExample deletes one example of examples/lib/ whole.
	RemoveLibExample func(path string, name string) error

	// ExecTest runs the project's examples and checks each one against its
	// golden result.yaml, reporting every example that diverged.
	ExecTest func(props ExecTestProps) error

	// UpdateTest runs one example by name, both sides, and rewrites its
	// golden result.yaml with what the run produced, printing the changes.
	UpdateTest func(path string, name string) error
}
