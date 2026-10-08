package opinatedagnosserver

import (
	opinatedagnoscli "{{.Module}}/sandbox/deps/OpinatedAgnosCli"
	serializables "{{.Module}}/sandbox/deps/serializables"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/deps/signaldeps"
	"{{.Module}}/sandbox/deps/std"
)

// This package is the contract of an *opinated* lib: unlike every other dep,
// which restates a library's raw capability and nothing more, it carries the
// agnos server mechanic itself — the declaration a route.yaml becomes, the
// chain that runs a request against those declarations, how a value and a
// body are bound and validated, how a failure is raised and written. The
// project's sandbox/api aliases its types, so a handler still reads them as
// api.Route, api.StatusOk.
//
// What stays in the project is everything it declares or edits: each route's
// route.yaml and InternalPureHandler.go, the generated declarations built from
// them, routeprops, and the Handle* files of sandbox/internal/server/errors/.

// PathType is what one segment a Path reads has to convert to. A segment
// that will not is a non-match: the url is for some other route.
type PathType int

const (
	// StringPath takes any slice, bound as a string.
	StringPath PathType = iota
	// IntegerPath takes one segment reading as a whole number, bound as an
	// int.
	IntegerPath
	// NumberPath takes one segment reading as a number, bound as a float64.
	NumberPath
	// UuidPath takes one segment reading as a canonical uuid, bound as a
	// string.
	UuidPath
)

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
	// Type is what the slice converts to; anything but StringPath reads one
	// segment alone.
	Type PathType
	// Description is the one-line help text.
	Description string
	// Trigger is what the slice has to match for the route to run.
	Trigger opinatedagnoscli.Trigger
}

// ParameterFont is one place of the request a Parameter is read from.
type ParameterFont int

const (
	// HeaderParam reads a request header, matched without regard to case.
	HeaderParam ParameterFont = iota
	// QueryParam reads a query-string parameter.
	QueryParam
	// CookieParam reads a request cookie.
	CookieParam
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
	// IntegerType is bound as an int.
	IntegerType
	// IntegerArrayType is bound as a []int, read the way a StringArrayType
	// is.
	IntegerArrayType
)

// AnyMethod is the one entry of AcceptMethods that accepts every http method.
const AnyMethod = "ANY"

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
	Trigger opinatedagnoscli.Trigger
	// Description is the one-line help text.
	Description string
}

// RouteBody is the request body a route declares — the parsed form of `body:`
// in its route.yaml. The dispatch enforces ContentType and MaxBytes before the
// handler runs; reading the body itself is the handler's to ask for, through
// the ReadBody its own package generates.
type RouteBody struct {
	// Type is the declared type: "none", "raw", "text", "json" or "form".
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
//
// It is an error too: an InternalPureHandler refuses a request by returning
// one, built by Fail, and the dispatch answers it through the Handle* file of
// its Status.
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

// Error is the failure's Message, which is what makes a RouteFailure an error
// a handler can return.
func (failure *RouteFailure) Error() string {
	return failure.Message
}

// Route is one http route of the project, as the sandbox offers it: the whole
// of what its route.yaml declares, plus the handler behind it. Server.Routes
// holds one per directory under sandbox/internal/routeslist holding a
// route.yaml, at any depth, each built by that package's generated NewRoute,
// in run order — lowest Priority first.
//
// What Server.Routes holds is the declaration alone: nothing is ever bound
// onto it. The dispatch copies it with BindRoute, puts the request and the
// response on the copy, and hands the copy to IsActionable and RequestHandler.
type Route struct {
	// Name is the package directory of the route, snake_case.
	Name string
	// AcceptMethods are the http methods it answers to ("GET", "POST"), or
	// AnyMethod alone for every one.
	AcceptMethods []string
	// Priority is the rung this route runs on when several match one
	// request: the dispatch runs them from the lowest upwards and stops at
	// the first handler that sets a status.
	Priority int
	// ResponseType is the Content-Type set on the response before the
	// handler runs; the handler may set another.
	ResponseType string
	// Segments is how many segments the request path has to have for the
	// route to run, 0 for any count.
	Segments int
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

	// ReadBody reads, validates and converts the request body of one bound
	// copy — the route package's own generated ReadBody, closed over the
	// sandbox. RequestHandler calls it before the handler runs and binds
	// what it returns onto Entries.Body; it is nil on a route whose body is
	// `none`. A body that fails comes back as a *RouteFailure, which the
	// dispatch raises.
	ReadBody func(bound *Route) (any, error)

	// InternalPureHandler is the route package's own InternalPureHandler,
	// closed over the sandbox: a func(props *routeprops.RouteProps, entries
	// *Entries, response *serverdeps.Response) error whose Entries is that package's
	// generated struct. It is held as any because every route's Entries is a
	// type of its own; the dispatch builds and fills one by reflection and
	// calls it.
	InternalPureHandler any

	// Request is the http request this copy was bound from and Response the
	// one being written, the response a handler is handed a pointer to.
	Request  serverdeps.Request
	Response serverdeps.Response

	// Props is one request's *routeprops.RouteProps, shared by every route
	// of the chain that runs for it and handed to each InternalPureHandler as
	// its first argument: what a middleware sets on it, the routes after it
	// read. It is held as any because the project types it under
	// sandbox/internal, which no contract may name; a Handle* file reads it
	// back with route.Props.(*routeprops.RouteProps).
	Props any

	// Failure is why this route is being handed to one of the project's
	// Handle* files, nil on a normal run. The dispatch sets it as it raises,
	// which is the one way any part of the server layer raises a failure.
	Failure *RouteFailure

	// IsActionable reports whether one bound copy — its Request set — is for
	// this route: the method is accepted, and every path slice and every
	// parameter declaring a trigger matches it. Nil reads the copy with the
	// lib's own matcher; set it to decide by hand.
	IsActionable func(bound *Route) bool
	// MatchesPath reports whether one bound copy's request path is for this
	// route whatever its method — what tells a 405 from a 404. Nil reads it
	// with the lib's own matcher.
	MatchesPath func(bound *Route) bool
	// RequestHandler binds one bound copy's request onto a fresh Entries and
	// runs InternalPureHandler with it. It returns the failure the handler
	// did not answer itself, nil otherwise; what it answered with is the
	// status it wrote, and a handler writing none hands the request to the
	// next route of the chain. Nil binds with the lib's own binder.
	RequestHandler func(bound *Route) error
}

// Server is the http surface of the sandbox: every route the project declares,
// and the dispatch that reads a request against them. It is built by
// sandbox/internal/generated/server/server/new.go, generated by the build.
type Server struct {
	// Serve opens the server: it hands the props to the lib's ServerMain
	// against these Routes, and blocks until the server stops.
	Serve func(props ServeProps) error
	// Routes is every http route the project declares, in run order —
	// lowest `priority` first, then by name — each built by the generated
	// NewRoute of its own package. The dispatch reads a request against
	// these declarations; a caller holding the sandbox reads the same
	// surface without one.
	Routes []*Route
	// Fail answers one failure with the project's own handler for it: it
	// reads route.Failure and calls the matching Handle* of
	// sandbox/internal/server/errors/. The dispatch raises through it; a
	// handler returns a failure built by Fail rather than calling it
	// directly.
	Fail func(route *Route) error
}

// ServeProps describes one run of the http server: the address to listen on
// — host:port (":8080"), a bare port ("8080"), or a range of ports tried in
// turn, with or without a host ("3000:4000", "127.0.0.1:3000:4000") —
// the two timeouts, in milliseconds, a request and a response are held to, and
// how long the requests in flight get to finish once the process is asked to
// stop (0 waits for them).
type ServeProps struct {
	Addr              string
	ReadTimeoutMs     int
	WriteTimeoutMs    int
	ShutdownTimeoutMs int
}

const (
	// StatusOk reports that the route did what it was asked to do.
	StatusOk = 200
	// StatusCreated reports that the route created what it was asked for.
	StatusCreated = 201
	// StatusNoContent reports success with nothing to send back.
	StatusNoContent = 204
	// StatusMovedPermanently sends the caller to Location for good.
	StatusMovedPermanently = 301
	// StatusFound sends the caller to Location this once.
	StatusFound = 302
	// StatusSeeOther sends the caller to Location with a GET, after a
	// form was handled.
	StatusSeeOther = 303
	// StatusNotModified reports that the caller's cached copy is current.
	StatusNotModified = 304
	// StatusTemporaryRedirect sends the caller to Location this once,
	// keeping its method.
	StatusTemporaryRedirect = 307
	// StatusPermanentRedirect sends the caller to Location for good,
	// keeping its method.
	StatusPermanentRedirect = 308
	// StatusBadRequest reports a request the dispatch could not bind: a
	// missing required field, an unparsable value, one out of range, or a
	// body the declared schema rejects.
	StatusBadRequest = 400
	// StatusUnauthorized reports a request that carries no valid
	// credentials.
	StatusUnauthorized = 401
	// StatusForbidden reports credentials that are valid and not enough.
	StatusForbidden = 403
	// StatusNotFound reports that no declared route matches the path.
	StatusNotFound = 404
	// StatusMethodNotAllowed reports a path a route matches under another
	// method.
	StatusMethodNotAllowed = 405
	// StatusConflict reports a well-formed request the current state refuses.
	StatusConflict = 409
	// StatusPayloadTooLarge reports a body longer than the route's
	// max-bytes.
	StatusPayloadTooLarge = 413
	// StatusUnsupportedMedia reports a content type the route does not
	// declare.
	StatusUnsupportedMedia = 415
	// StatusUnprocessable reports a well-formed request whose content
	// breaks a rule of the domain.
	StatusUnprocessable = 422
	// StatusTooManyRequests reports a caller over its rate.
	StatusTooManyRequests = 429
	// StatusFailure reports a well-formed request the route could not carry
	// out.
	StatusFailure = 500
	// StatusUnavailable reports a server that cannot answer right now.
	StatusUnavailable = 503
)

// MainProps is everything one run of ServerMain reads: the surface, the
// address and timeouts, and the deps the server reaches the outside world
// through. The generated registry builds it per call, so nothing the lib holds
// goes stale.
type MainProps struct {
	// Server is the surface every request is read against. It is a
	// pointer, read as each request runs, so a caller that replaced a field
	// of sandbox.Server is followed.
	Server *Server
	// Serve is the address and timeouts the server runs with.
	Serve ServeProps
	// NewProps builds the one *routeprops.RouteProps every route of one
	// request's chain is handed — a type of the project, which is why the
	// lib is handed a constructor rather than naming it.
	NewProps func() any
	// MatchTrigger is how a path or a parameter is held to its trigger: the
	// OpinatedAgnosCli lib's own, shared with the cli layer.
	MatchTrigger func(trigger opinatedagnoscli.Trigger, text string, segmented bool) bool
	// Std is the channel the server reports on, read at the moment of each
	// print.
	Std *std.Sandbox
	// Serverdeps opens the port; it routes nothing.
	Serverdeps serverdeps.Sandbox
	// Signaldeps hears the request to stop, which shuts the server down.
	Signaldeps signaldeps.Sandbox
}

// Sandbox is the server lib injected whole as the Deps.OpinatedAgnosServer
// field.
type Sandbox struct {
	// ServerMain opens the port through props.Serverdeps — which routes
	// nothing — and runs every request through the chain of
	// props.Server.Routes: each route the request is for runs, lowest
	// `priority` first, until one of them answers — sets a status, or writes
	// a byte, which sends a 200. Every way a request can end without an
	// answer — nothing matched, matched under another method, a value that
	// will not bind, a panic — is raised through props.Server.Fail. Addr may
	// name a range of ports ("3000:4000"), tried in turn; the address it
	// landed on is printed to stdout. The first request to stop the process
	// shuts it down gracefully, after which it returns nil.
	ServerMain func(props MainProps) error

	// NewRoute returns an empty Route with every slice open, the base a
	// route's generated NewRoute fills its declaration on.
	NewRoute func() *Route

	// BindRoute copies one declaration into the route a single request runs
	// on: the same declared fields — the slices are read-only and shared —
	// with no request, response or failure yet.
	BindRoute func(route *Route) *Route

	// Fail is how an InternalPureHandler — or a generated ReadBody — refuses
	// a request: it builds the failure and returns it as an error, which the
	// dispatch raises on the route, reaching the project's own Handle* file
	// for that status:
	//
	//	return sandbox.Deps.OpinatedAgnosServer.Fail(api.StatusUnauthorized, "authorization", "invalid token")
	//
	// A message left empty is filled by that file's own wording.
	Fail func(status int, field string, message string) error

	// FailWithCause is Fail carrying what went wrong underneath — an error's
	// text, meant for the log, never for the caller.
	FailWithCause func(status int, field string, message string, cause string) error

	// FailureOf is the failure one of the project's Handle* files is
	// answering. A file stands for one status and one wording, and passes
	// them as the fallback: a failure that carries its own — a field that
	// would not bind, a body the schema rejected — answers with that, and one
	// that carries none — "nothing matched" is raised with no message at all
	// — answers with the file's.
	FailureOf func(route *Route, status int, message string) RouteFailure

	// WriteError writes one failure onto the response: the body is always
	// the same JSON object — {"error": "...", "field": "..."} — so a client
	// parses one shape whatever went wrong. It is the writer, not the
	// policy: what decides whether a failure is answered this way is the
	// project's own Handle* file. It returns the message as an error, so a
	// Handle* file answers and reports in one line.
	WriteError func(serializables serializables.Sandbox, response serverdeps.Response, status int, field string, message string) error

	// WriteJSON answers with one document serialized as JSON.
	WriteJSON func(serializables serializables.Sandbox, response serverdeps.Response, status int, document *serializables.SerializibleObject) error

	// WriteText answers with one text as text/plain.
	WriteText func(response serverdeps.Response, status int, text string) error

	// Redirect answers by sending the caller to location, with one of the
	// redirect statuses: StatusFound, StatusSeeOther, StatusMovedPermanently,
	// StatusTemporaryRedirect or StatusPermanentRedirect.
	Redirect func(response serverdeps.Response, status int, location string) error

	// ValidateSchema parses body as JSON and checks it against schemaJson,
	// the canonical form of a route's declared json-schema — a subset:
	// type, nullable, const, enum, required, properties,
	// additionalProperties, items, min/maxItems, uniqueItems, min/maxLength,
	// pattern, format (email, uuid, date-time, uri) and the four numeric
	// bounds. It returns the parsed document, the field path of the first
	// violation, that violation's message, and whether the body passed.
	ValidateSchema func(serializables serializables.Sandbox, schemaJson string, body []byte) (*serializables.SerializibleObject, string, string, bool)

	// ValidateForm is ValidateSchema for a body that arrives as `key=value`
	// pairs: each declared key is converted to what its property declares —
	// an integer or a number parsed, a boolean read as true/1/on or
	// false/0/off, an array taking every occurrence — before the schema
	// runs. An empty value counts as absent.
	ValidateForm func(serializables serializables.Sandbox, schemaJson string, form map[string][]string) (*serializables.SerializibleObject, string, string, bool)

	// ReadString returns the named property of a validated document as
	// text; every reader below returns the zero value for a property that is
	// absent or of another kind, since the schema has been enforced by then.
	ReadString func(object *serializables.SerializibleObject, key string) string
	// ReadInt returns the named property as an int.
	ReadInt func(object *serializables.SerializibleObject, key string) int
	// ReadFloat returns the named property as a float64.
	ReadFloat func(object *serializables.SerializibleObject, key string) float64
	// ReadBool returns the named property as a bool.
	ReadBool func(object *serializables.SerializibleObject, key string) bool
	// ReadObject returns the named property as a document to read further,
	// nil when it is absent.
	ReadObject func(object *serializables.SerializibleObject, key string) *serializables.SerializibleObject
	// ReadItems returns the named property's items, in order.
	ReadItems func(object *serializables.SerializibleObject, key string) []*serializables.SerializibleObject
	// ItemString reads one array item as text.
	ItemString func(item *serializables.SerializibleObject) string
	// ItemInt reads one array item as an int.
	ItemInt func(item *serializables.SerializibleObject) int
	// ItemFloat reads one array item as a float64.
	ItemFloat func(item *serializables.SerializibleObject) float64
	// ItemBool reads one array item as a bool.
	ItemBool func(item *serializables.SerializibleObject) bool
}
