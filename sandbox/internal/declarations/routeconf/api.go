package routeconf

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"

// Path is one entry of a route's `paths`: the slice of request segments from
// Start to End, both inclusive, End -1 standing for the last segment. The slice
// is read as "/" + its segments joined by "/", compared against Trigger when
// one is declared, and bound to Input.<Id> either way.
type Path struct {
	Id          string
	Start       int
	End         int
	Type        string // "string" | "integer" | "number" | "uuid"
	Trigger     triggerconf.Trigger
	Description string
}

// Parameter is one entry of a route's `parameters`: one value read off the
// request under Key, from the first of Sources that carries it, and bound to
// Input.<Id>. Key is the external spelling — the query key or the header
// name, a header matched without regard to case — and defaults to Id.
type Parameter struct {
	Id          string
	Key         string
	Type        string   // "string" | "integer" | "number" | "boolean" | "datetime" | "string-array" | "integer-array"
	Sources     []string // "query" | "header" | "cookie", in the order they are read
	Required    bool
	Default     string
	HasDefault  bool
	Trigger     triggerconf.Trigger
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

// Body describes a route's request body. The generated ReadBody applies
// Required, MaxBytes and Schema before the handler runs; a middleware in front
// of the route may still refuse a request before a byte of the body is read.
// Schema is declared under `json-schema` on a json body and under
// `form-schema` on a form one — SchemaKeyOf names which — and SchemaKeys
// records every key it was read under, so verify reports the wrong one.
type Body struct {
	Type        string // "none" | "raw" | "text" | "json" | "form"
	Required    bool
	MaxBytes    int
	ContentType string
	Schema      *Schema
	HasSchema   bool
	SchemaKeys  []string
}

// RouteConf is the parsed form of sandbox/internal/routes/<name>/route.yaml —
// the declarative description of one http route, which `agnos build` turns
// into that route's generated new.go and input.go. It is written by
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
	ResponseType string
	// Segments is how many segments the request path has to have for the
	// route to run; HasSegments is false on a route that takes any count.
	Segments    int
	HasSegments bool
	Paths       []Path
	Parameters  []Parameter
	Category    string
	Summary     string
	Description string
	Examples    []string
	Hidden      bool
	Body        Body
	// Legacy lists every key of a pre-routes declaration found in the
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
