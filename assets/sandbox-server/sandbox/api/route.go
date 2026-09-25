package api

// TriggerType is how a Trigger compares the text it is handed.
type TriggerType int

const (
	// EqualTrigger matches a text that is exactly the trigger's Value.
	EqualTrigger TriggerType = iota
	// PrefixTrigger matches a text that begins with the Value.
	PrefixTrigger
	// SuffixTrigger matches a text that ends with the Value.
	SuffixTrigger
	// RegexTrigger matches a text the Value, a regular expression, matches.
	RegexTrigger
)

// Trigger is the condition a path slice or a parameter value has to meet for
// a route to run at all — the parsed form of one `trigger:` of route.yaml.
// Failing it is a non-match, never a 400: the request is for some other route.
type Trigger struct {
	// Exist tells a declared trigger from none at all; an entry with none
	// matches whatever the request brought.
	Exist bool
	// Type is how Value is compared.
	Type TriggerType
	// Value is what the text is compared against.
	Value string
}

// Path is one entry of `paths` in route.yaml: the slice of request segments
// from Start to End, both inclusive, End -1 standing for the last segment. The
// slice reads as "/" followed by its segments joined by "/", and is bound to
// the Entries field tagged with its Id.
type Path struct {
	// Id is the Entries field the slice is bound to.
	Id string
	// Start is the index of the first segment of the slice.
	Start int
	// End is the index of the last segment of the slice, -1 for the last
	// segment of the request.
	End int
	// Description is the one-line help text.
	Description string
	// Trigger is what the slice has to match for the route to run.
	Trigger Trigger
}

// ParameterFont is one place of the request a Parameter is read from.
type ParameterFont int

const (
	// HeaderParam reads a request header, matched without regard to case.
	HeaderParam ParameterFont = iota
	// QueryParam reads a query-string parameter.
	QueryParam
)

// ParameterType is the type a Parameter is converted to before it reaches
// Entries.
type ParameterType int

const (
	// StringType is bound as a string.
	StringType ParameterType = iota
	// NumberType is bound as a float64.
	NumberType
	// BooleanType is bound as a bool: true/1 or false/0.
	BooleanType
	// DateTimeType is bound as a string that has to read as RFC 3339.
	DateTimeType
	// StringArrayType is bound as a []string: every occurrence of a query
	// key, or a header's comma-separated values.
	StringArrayType
)

// Parameter is one entry of `parameters` in route.yaml: one value read off the
// request under Key, from the first of Fonts that carries it, and bound to the
// Entries field tagged with its Id.
type Parameter struct {
	// Id is the Entries field the value is bound to.
	Id string
	// Key is the query key or the header name the value is read under.
	Key string
	// Fonts are the places the value is read from, in order: the first one
	// that brings a value wins.
	Fonts []ParameterFont
	// Required reports that the request is answered 400 without it.
	Required bool
	// Type is what the value is converted to.
	Type ParameterType
	// Default is the value bound when the request brings none, spelled as
	// route.yaml writes it; HasDefault tells an empty default from none.
	Default    string
	HasDefault bool
	// Trigger is what the value has to match for the route to run.
	Trigger Trigger
	// Description is the one-line help text.
	Description string
}

// RouteBody is the request body a route declares — the parsed form of `body:`
// in its route.yaml. The dispatch enforces ContentType and MaxBytes before the
// handler runs; reading the body itself is the handler's to ask for, through
// the ReadBody its own package generates.
type RouteBody struct {
	// Type is the declared type: "none", "raw", "text" or "json".
	Type string
	// Required reports that an absent or empty body is rejected.
	Required bool
	// MaxBytes is the largest body the route reads; a longer one is 413.
	MaxBytes int
	// ContentType is the media type the route accepts; a divergent one is
	// 415. It is empty when the route accepts any.
	ContentType string
	// Schema is the declared json-schema in canonical form, "" when the
	// route declares none.
	Schema string
}

// RouteFailure is one way a request did not get answered: what the dispatch
// would have written, and why. Every failure the server layer raises — a
// header that will not bind, a body the schema rejected, a handler that
// returned an error, a path nothing matched — arrives at one of the project's
// own Handle* files as this, read off the Failure of the route it is handed.
type RouteFailure struct {
	// Status is the http status the failure carries: one of the Status*
	// constants of server.go.
	Status int
	// Field is the header, query parameter, path segment or json path that
	// failed, "" when the failure names no single value.
	Field string
	// Message is the one-line reason, written for whoever called the route.
	Message string
	// Cause is what went wrong underneath — an error's text, or the value a
	// handler panicked with — "" when there is nothing below the message. It
	// is for the log, not for the caller.
	Cause string
}

// Route is one http route of the project, as the sandbox offers it: the whole
// of what its route.yaml declares, plus the handler behind it. Server.Routes
// holds one per sandbox/internal/routeslist/<name>/, each built by that
// package's generated NewRoute, in run order — lowest Priority first.
//
// What Server.Routes holds is the declaration alone: nothing is ever bound
// onto it. The dispatch copies it with BindRoute, puts the request and the
// response on the copy, and hands the copy to IsActionable and RequestHandler.
type Route struct {
	// Name is the package directory of the route, snake_case.
	Name string
	// AcceptMethods are the http methods it answers to ("GET", "POST").
	AcceptMethods []string
	// Priority is the rung this route runs on when several match one
	// request: the dispatch runs them from the lowest upwards and stops at
	// the first handler that sets a status.
	Priority int
	// ResponseType is the Content-Type set on the response before the
	// handler runs; the handler may set another.
	ResponseType string
	// Pattern is its path as it reads in docs and messages.
	Pattern string
	// Category groups it on the generated Routes page.
	Category string
	// Help is the one-line description.
	Help string
	// LongDescription is the paragraph the Routes page prints.
	LongDescription string
	// Examples are whole requests the Routes page prints.
	Examples []string
	// Hidden keeps it off the Routes page without disabling it.
	Hidden bool
	// Paths are the slices of the request path it reads, in declaration
	// order.
	Paths []Path
	// Parameters are the headers and query parameters it reads, in
	// declaration order.
	Parameters []Parameter
	// Body is the request body declaration.
	Body RouteBody

	// InternalPurehandler is the route package's own InternalPureHandler,
	// closed over the sandbox: a func(route *Route, entries *Entries,
	// response *serverdeps.Response) error whose Entries is that package's
	// generated struct. It is held as any because every route's Entries is a
	// type of its own; RequestHandler builds and fills one through
	// Deps.Reflectdeps and calls it.
	InternalPurehandler any

	// Request is the http request this copy was bound from and Response the
	// one being written. Both are handed over as any: sandbox/api may name no
	// type of sandbox/deps, so the server layer reads them back through
	// routeio.RequestOf and routeio.ResponseOf.
	Request  any
	Response any

	// Failure is why this route is being handed to one of the project's
	// Handle* files, nil on a normal run. It is set by routeio.Fail, which
	// is the one way any part of the server layer raises a failure.
	Failure *RouteFailure

	// IsActionable reports whether one bound copy — its Request set — is for
	// this route: the method is accepted, and every path slice and every
	// parameter declaring a trigger matches it.
	IsActionable func(bound *Route) bool
	// MatchesPath reports whether one bound copy's request path is for this
	// route whatever its method — what tells a 405 from a 404.
	MatchesPath func(bound *Route) bool
	// RequestHandler binds one bound copy's request onto a fresh Entries and
	// runs InternalPurehandler with it. It returns the failure the handler
	// did not answer itself, nil otherwise; what it answered with is the
	// status it wrote, and a handler writing none hands the request to the
	// next route of the chain.
	RequestHandler func(bound *Route) error
}

// NewRoute returns an empty Route with every slice open. The generic
// sandbox/internal/server/route.NewRoute fills the matcher and the handler on
// top of it, and a route's generated NewRoute its declaration.
func NewRoute() *Route {
	return &Route{
		AcceptMethods: []string{},
		Examples:      []string{},
		Paths:         []Path{},
		Parameters:    []Parameter{},
	}
}

// BindRoute copies one declaration into the route a single request runs on:
// the same declared fields — the slices are read-only and shared — with no
// request, response or failure yet. The dispatch calls it once per request, so
// two requests in flight never share a bound value.
func BindRoute(route *Route) *Route {
	bound := *route
	bound.Request = nil
	bound.Response = nil
	bound.Failure = nil
	return &bound
}
