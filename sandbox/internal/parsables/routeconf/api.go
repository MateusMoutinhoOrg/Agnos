package routeconf

// Trigger is the condition a path slice or a parameter value has to meet for
// the route to join the run list: the text compared, and how. Exists is false
// on an entry that declares none — a plain capture, or a parameter that is
// bound and never matched on.
type Trigger struct {
	Exists bool
	Type   string // "equal" | "prefix" | "suffix" | "regex"
	Value  string
}

// Path is one entry of a route's `paths`: the slice of request segments from
// Start to End, both inclusive, End -1 standing for the last segment. The slice
// is read as "/" + its segments joined by "/", compared against Trigger when
// one is declared, and bound to Entries.<Id> either way.
type Path struct {
	Id          string
	Start       int
	End         int
	Trigger     Trigger
	Description string
}

// Parameter is one entry of a route's `parameters`: one value read off the
// request under Key, from the first of Fonts that carries it, and bound to
// Entries.<Id>. Key is the external spelling — the query key or the header
// name, a header matched without regard to case — and defaults to Id.
type Parameter struct {
	Id          string
	Key         string
	Type        string   // "string" | "number" | "boolean" | "datetime" | "string-array"
	Fonts       []string // "query" | "header", in the order they are read
	Required    bool
	Default     string
	HasDefault  bool
	Trigger     Trigger
	Description string
	Examples    []string
}

// SchemaProperty is one named property of an object Schema, kept as an ordered
// slice so the generated struct fields and the canonical schema JSON come out
// the same on every build.
type SchemaProperty struct {
	Name   string
	Schema *Schema
}

// Schema is the subset of JSON Schema a route's body may declare, as a tree.
// Every bound carries its own Has… companion, so an unset bound is told from
// one declared as zero — the same pair Field.Min/HasMin uses. Unknown holds
// every key outside the subset, which build and verify reject rather than
// ignore.
type Schema struct {
	Type                    string // "object"|"array"|"string"|"integer"|"number"|"boolean"|"null"
	Nullable                bool
	Format                  string // "email" | "uuid" | "date-time" | "uri"
	Pattern                 string
	Properties              []SchemaProperty
	Required                []string
	AdditionalProperties    bool
	HasAdditionalProperties bool
	Items                   *Schema
	Enum                    []string
	Const                   string
	HasConst                bool
	Minimum                 float64
	HasMinimum              bool
	Maximum                 float64
	HasMaximum              bool
	ExclusiveMinimum        float64
	HasExclusiveMinimum     bool
	ExclusiveMaximum        float64
	HasExclusiveMaximum     bool
	MinLength               int
	HasMinLength            bool
	MaxLength               int
	HasMaxLength            bool
	MinItems                int
	HasMinItems             bool
	MaxItems                int
	HasMaxItems             bool
	UniqueItems             bool
	Unknown                 []string
}

// Body describes a route's request body. It is the one part of a request the
// dispatch does not read: the generated ReadBody applies Required, MaxBytes and
// Schema on demand, so a handler can refuse a request before a byte of the body
// is read.
type Body struct {
	Type        string // "none" | "raw" | "text" | "json"
	Required    bool
	MaxBytes    int
	ContentType string
	Schema      *Schema
	HasSchema   bool
}

// RouteConf is the parsed form of sandbox/internal/routeslist/<name>/route.yaml —
// the declarative description of one http route, which `agnos build` turns
// into that route's generated new.go and entries.go. It is written by
// `add-route` and rewritten by the path, parameter and body editors, never by
// hand.
type RouteConf struct {
	Methods []string
	// Priority is the rung this route runs on when several match one
	// request: lowest first. It is required: HasPriority tells a declared 0
	// from none at all.
	Priority    int
	HasPriority bool
	// ResponseType is the Content-Type the dispatch sets before the handler
	// runs. It is required.
	ResponseType    string
	Paths           []Path
	Parameters      []Parameter
	Category        string
	Help            string
	LongDescription string
	Examples        []string
	Hidden          bool
	Body            Body
	// Legacy lists every key of a pre-routeslist declaration found in the
	// file (`method`, `headers`, `params`), which verify reports by name.
	Legacy []string

	// Render serializes the declaration back to the route.yaml shape.
	Render func() string
	// Pattern is the route's path as it reads in docs and messages.
	Pattern func() string
	// SchemaJson is the declared json-schema as canonical JSON — the text
	// baked into the generated BodySchema constant — or "" when the body
	// declares none.
	SchemaJson func() string
}
