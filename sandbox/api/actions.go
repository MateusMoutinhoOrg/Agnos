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

// AddDepProps describes one dep to install. Dep is either a name of the
// embedded catalog or, when it holds a "/", the module path of another agnos
// repo — the same disambiguation `go get` makes. Adapter picks which
// implementation of a catalog dep fills the contract ("" takes the dep's
// declared default-adapter); As names the copied contract of a remote one (""
// takes the last segment of the module path); RemoteAvailable is the available
// of the remote repo the generated shim builds its sandbox from ("" is
// standard).
type AddDepProps struct {
	Path            string
	Dep             string
	Adapter         string
	As              string
	RemoteAvailable string
}

// SetDepProps describes one remote dep to move to another version of the
// module it was copied from. RemoteAvailable is the available of the remote
// repo the regenerated shim builds its sandbox from ("" is standard).
type SetDepProps struct {
	Path            string
	Dep             string
	Version         string
	RemoteAvailable string
}

// RemoveDepProps describes one dep to uninstall. A dep with adapters
// installed is refused, because which of them was meant is not a question the
// tree can answer; WithAdapters is the caller saying "all of them".
type RemoveDepProps struct {
	Path         string
	Dep          string
	WithAdapters bool
}

// AddAdapterProps describes one further implementation to install for a
// contract the project already has. Available names the available that should
// switch to it; "" installs the package and leaves every selection alone.
type AddAdapterProps struct {
	Path      string
	Adapter   string
	Available string
}

// SetAdapterProps describes one selection to change: which adapter fills a
// dep's field in one available. Available is the standard one when empty.
type SetAdapterProps struct {
	Path      string
	Dep       string
	Adapter   string
	Available string
}

// ExtensionInfo is one row of ListExtensions: one generation mechanic of the
// catalog, what it looks after, and whether this project turned it on.
type ExtensionInfo struct {
	Name    string
	Enabled bool
	Help    string
}

// DepInfo is one row of ListDeps: a dep of the embedded catalog, and what the
// project holds of it.
type DepInfo struct {
	Name           string
	Field          string
	Help           string
	DefaultAdapter string
	Installed      bool
	Adapters       []string
}

// AdapterInfo is one row of ListAdapters: an adapter of the embedded catalog
// or one installed in the project, and which availables bind it. Origin is
// "catalog" for one the catalog installs and "generated" for the shim of a dep
// copied from a remote repo.
type AdapterInfo struct {
	Name       string
	Dep        string
	Help       string
	Module     string
	Origin     string
	Installed  bool
	Availables []string
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
// captured one, and is normalized to start with "/". Array collects a []T
// field: every occurrence of a query key, or — on the last segment of the
// path, and there alone — every segment left in the URL. Default, Min and Max
// are the raw literals typed on the command line ("" means unset); Position is
// the index to insert at (< 0 appends).
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

// RouteFieldEditProps describes the change set-segment, set-header or
// set-param applies to one field a route already declares. Name is the field
// as it is declared now and Rename the spelling it takes on ("" leaves it
// alone); every other key overwrites what is there when it is given, and an
// empty one leaves it as it is. Clear is how a key is taken off again —
// "description", "examples", "default", "required", "array", "min" or "max" —
// because an empty string cannot say "unset this" and "leave it alone" at once.
type RouteFieldEditProps struct {
	Path        string
	Route       string
	Name        string
	Rename      string
	Identifier  string
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
	Array       bool
	Min         string
	Max         string
	Clear       []string
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

// RouteBodyFieldEditProps describes the change set-body-field applies to one
// property a route's body json-schema already declares. Name is the dotted
// path it sits at and Rename the leaf spelling it takes on (it stays in the
// object it is declared in); every other key is one keyword of the supported
// subset, overwriting what is there when it is given. Clear names the keywords
// to take off instead — "required", "array", "min", "max", "exclusive-min",
// "exclusive-max", "format", "pattern", "enum", "const", "nullable",
// "min-items", "max-items", "unique-items" or "additional-properties" — which
// is the one thing an empty value cannot say.
type RouteBodyFieldEditProps struct {
	Path                   string
	Route                  string
	Name                   string
	Rename                 string
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
	Clear                  []string
}

// RouteBodyImportProps describes one example payload to read a route's body
// json-schema off. Json is the document itself and File a path to read it from
// — exactly one of the two — and the inference walks it: an object becomes an
// object property, a list an array of whatever its first item is, and a scalar
// the type it is written as. Required lists every key the example carries in
// its object's required set, InferFormat reads an email, a uuid, a date-time
// or a uri back as the format it spells, and Replace drops the schema that is
// there instead of adding to it.
type RouteBodyImportProps struct {
	Path        string
	Route       string
	Json        string
	File        string
	Required    bool
	Replace     bool
	InferFormat bool
}

// PageProps describes one html page to declare: the project directory, the
// name the page carries (it becomes the route, its Go package and the html
// file), the literal path segment it answers on ("" defaults to /<name>), the
// <title> the scaffolded html carries ("" defaults to the name) and the
// one-line help its route.yaml is declared with ("" derives one from the name).
type PageProps struct {
	Path    string
	Name    string
	Trigger string
	Title   string
	Help    string
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

	// EnableExtension turns one generation mechanic on in the project's
	// extensions.yaml and rebuilds, so what that mechanic owns is rendered
	// from here on.
	EnableExtension func(path string, name string) error

	// DisableExtension turns one generation mechanic off. Nothing is removed:
	// agnos stops rendering what that mechanic owns and the files it wrote
	// become the project's, to keep or to edit by hand. Deleting them is what
	// the matching <x>-purge is for.
	DisableExtension func(path string, name string) error

	// ListExtensions returns one row per generation mechanic of the catalog,
	// saying which ones this project turned on.
	ListExtensions func(path string) ([]ExtensionInfo, error)

	// DepsInit adds the dependency layer (sandbox/deps/ and
	// adapters/availables/standard/) to a project that has none.
	DepsInit func(path string) error

	// DepsPurge removes the dependency layer and every installed dep with it.
	DepsPurge func(path string) error

	// AddDep installs one dep of the built-in list: its contract under
	// sandbox/deps/, one adapter filling it under adapters/libs/ and that
	// adapter's go.mod require.
	AddDep func(props AddDepProps) error

	// RemoveDep uninstalls one installed dep: its adapters, their requires,
	// and then the contract itself. It refuses a dep that still has an
	// adapter installed unless props.WithAdapters says to take those too.
	RemoveDep func(props RemoveDepProps) error

	// ListDeps returns one row per dep of the embedded catalog, saying which
	// the project has installed and which adapters fill each one.
	ListDeps func(path string) ([]DepInfo, error)

	// SetDep re-copies one remote dep at another version of its module and
	// regenerates the shim that converts it.
	SetDep func(props SetDepProps) error

	// AddAdapter installs one further implementation of a contract the
	// project already has, and — when props.Available names one — switches
	// that available to it.
	AddAdapter func(props AddAdapterProps) error

	// RemoveAdapter uninstalls one adapter, its require and its files. It
	// refuses one that an available still binds, and one written by the
	// generator as half of a remote dep.
	RemoveAdapter func(path string, adapter string) error

	// SetAdapter changes which adapter fills one dep's field in one
	// available, the only place that choice is recorded.
	SetAdapter func(props SetAdapterProps) error

	// ListAdapters returns one row per adapter, of the embedded catalog and
	// of the project, with the availables that bind each one.
	ListAdapters func(path string) ([]AdapterInfo, error)

	// AddAvailable creates one further available, seeded with the standard
	// available's selection so it starts filling every field.
	AddAvailable func(path string, available string) error

	// RemoveAvailable deletes one available. The standard one is refused: it
	// is what cmd/main/main.go imports.
	RemoveAvailable func(path string, available string) error

	// CliInit adds the CLI layer (cmd/main, the dispatcher and the help and
	// version commands) to a project that has none.
	CliInit func(path string) error

	// CliPurge removes the CLI layer and every command declared in it.
	CliPurge func(path string) error

	// AddCommand declares a new command: its entries.yaml, its generated
	// new.go and a handler.go to fill in.
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

	// AddRoute declares a new route: its route.yaml, its generated new.go
	// and a handler.go to fill in.
	AddRoute func(path string, name string, method string, trigger string, help string, category string) error

	// RemoveRoute deletes one route and unwires it from the dispatch.
	RemoveRoute func(path string, name string) error

	// SetRoute rewrites the route-level keys of one route's route.yaml.
	SetRoute func(props RouteProps) error

	// AddSegment appends one segment to a route's path: a literal one when
	// props.Identifier is set, a captured one otherwise — and, with
	// props.Array, one taking every segment left in the path.
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

	// SetSegment rewrites one declared segment of a route's path, named
	// either by its capture name or by the identifier it spells.
	SetSegment func(props RouteFieldEditProps) error

	// SetHeader rewrites one declared request header of a route.
	SetHeader func(props RouteFieldEditProps) error

	// SetParam rewrites one declared query parameter of a route.
	SetParam func(props RouteFieldEditProps) error

	// SetBodyField rewrites one property of a route's body json-schema, at
	// the dotted path props.Name.
	SetBodyField func(props RouteBodyFieldEditProps) error

	// ImportBody declares a route's body json-schema from an example
	// payload, inferring one property per key the example carries.
	ImportBody func(props RouteBodyImportProps) error

	// ShowRoute renders one route's whole declaration — its path, its
	// headers, its query parameters and its body schema — as the lines of
	// a tree, ready to print.
	ShowRoute func(path string, route string) ([]string, error)

	// FrontInit adds the html front layer (sandbox/internal/pageio, the
	// route serving assets/frontend/static and that tree's skeleton) to a
	// project that has none, installing the server layer first when it is
	// missing.
	FrontInit func(path string) error

	// FrontPurge removes the front layer, the static route and the route of
	// every declared page, leaving assets/frontend/ untouched.
	FrontPurge func(path string) error

	// AddPage declares a new html page: the route that answers it and the
	// html template under assets/frontend/pages/ that it renders.
	AddPage func(props PageProps) error

	// RemovePage deletes one page, its route and its html template both.
	RemovePage func(path string, name string) error

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

	// Interview runs the interactive session over a project: it asks what is
	// to be done, generates the questions from the declaration of the command
	// that answers it, and runs that command with the answers bound onto it.
	// It writes nothing of its own — every command it dispatches runs the
	// action behind it, which persists and builds for itself.
	Interview func(path string) error
}
