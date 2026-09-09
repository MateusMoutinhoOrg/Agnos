package routeconf

// Field is one header, one query parameter or one captured path segment
// declared in a route's route.yaml. Key is the external spelling — the header
// name (matched without regard to case) or the query key — and the generated
// Entries field is derived from it.
type Field struct {
	Key         string
	Description string
	Examples    []string
	Type        string // "string" | "boolean" | "int" | "float"
	Default     string
	HasDefault  bool
	Required    bool
	Array       bool
	Min         float64
	HasMin      bool
	Max         float64
	HasMax      bool
}

// Segment is one item of a route's `paths` sequence: either a trigger, whose
// Identifier is the literal it matches in the URL (always starting with "/"),
// or a capture, whose Field names the segment and types the value. Exactly one
// of the two is filled.
type Segment struct {
	Identifier string
	Field      *Field
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
// dispatch does not read: the generated Entries.ReadBody applies Required,
// MaxBytes and Schema on demand, so a handler can refuse a request before a
// byte of the body is read.
type Body struct {
	Type        string // "none" | "raw" | "text" | "json"
	Required    bool
	MaxBytes    int
	ContentType string
	Schema      *Schema
	HasSchema   bool
}

// RouteConf is the parsed form of sandbox/internal/routes/<name>/route.yaml —
// the declarative description of one http route, which `agnos build` turns
// into an entries.go struct plus a match/handle pair in
// sandbox/internal/server/servermain.go. It is written by `add-route` and
// rewritten by `add-field` / `remove-field` / `set-route`, never by hand.
type RouteConf struct {
	Method          string
	Paths           []Segment
	Category        string
	Help            string
	LongDescription string
	Examples        []string
	Hidden          bool
	Headers         []Field
	Params          []Field
	Body            Body

	// Render serializes the declaration back to the route.yaml shape.
	Render func() string
	// Pattern is the route's path as it reads in docs and messages: every
	// segment in order, a trigger bringing its own leading slash and a
	// capture entering as "/{name}". The dispatch matches segment by
	// segment, never this text.
	Pattern func() string
	// IdentifierCount is how many trigger segments the route fixes — the
	// first key the match order sorts by, most specific first.
	IdentifierCount func() int
	// IdentifierLen is the total number of characters the route's triggers
	// spell, the tie-break between two routes fixing as many segments.
	IdentifierLen func() int
	// SchemaJson is the declared json-schema as canonical JSON — the text
	// baked into the generated EntriesSchema constant — or "" when the body
	// declares none.
	SchemaJson func() string
}
