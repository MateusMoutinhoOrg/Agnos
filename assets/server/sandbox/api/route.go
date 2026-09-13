package api

// RouteField is one value a route binds off a request — the parsed form of one
// captured segment of `paths`, or of one entry under `headers` or `params` in
// that route's route.yaml. It is matched by Id: the header name, the query key
// or the name the captured segment is declared under, all read back with the
// Get* readers of Route.
type RouteField struct {
	// Type is the declared type: "string", "boolean", "int" or "float".
	Type string
	// Id is the name the field is declared and read back under.
	Id string
	// Required reports that the request is rejected without it.
	Required bool
	// Array reports that every occurrence is kept, not just the first; on
	// the last captured segment of `paths` it takes every segment left in
	// the path.
	Array bool
	// Description is the one-line help text.
	Description string
	// Examples are whole requests the docs print.
	Examples []string
	// Default is the value bound when the field is absent, spelled as it
	// is written in route.yaml.
	Default string
	// HasDefault tells a declared empty default from no default at all.
	HasDefault bool
	// Min and Max bound a numeric value; HasMin and HasMax tell a bound of
	// zero from no bound.
	Min    float64
	HasMin bool
	Max    float64
	HasMax bool
}

// RoutePath is one segment of the route's url, in declaration order. Exactly
// one of the two is filled: a literal the request path has to spell, or a
// captured field that matches whatever sits in its place.
type RoutePath struct {
	// Identifier is the literal segment, leading slash included
	// ("/users"); "/" alone fixes the empty path and matches no segment of
	// its own. It is empty on a capture.
	Identifier string
	// Field is the captured segment's declaration, nil on a literal. Only
	// the last entry of Paths may declare an Array, and it then takes every
	// segment left in the path.
	Field *RouteField
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

// Route is one http route of the project, as the sandbox offers it: the whole
// of what its route.yaml declares, plus the handler behind it. Sandbox.Routes
// holds one per sandbox/internal/routes/<name>/, each built by that package's
// generated NewRoute, in match order. The server dispatch mints one instance
// per request through New, fills Items from the path, the headers and the
// query string, and then calls Handler; a caller holding the sandbox can do
// the same.
type Route struct {
	// Name is the package directory of the route, snake_case.
	Name string
	// Method is the http method it answers to ("GET", "POST").
	Method string
	// Pattern is its path as it reads in docs and messages
	// ("/users/{tenant}/create").
	Pattern string
	// Paths are the url's segments, in declaration order.
	Paths []RoutePath
	// Headers are the request headers it binds, in declaration order.
	Headers []RouteField
	// Params are the query parameters it binds, in declaration order.
	Params []RouteField
	// Body is the request body declaration.
	Body RouteBody
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

	// Items holds the values bound to the declaration: one entry per field
	// Id, in request order for an array, one element long for a scalar,
	// absent when nothing was bound. The Get* fields read it.
	Items map[string][]any

	// Request is the http request this instance was bound from and
	// Response the one being written. Both are handed over as any:
	// sandbox/api may name no type of sandbox/deps, so the server layer
	// reads them back through routeio.RequestOf and routeio.ResponseOf.
	Request  any
	Response any

	// GetItem returns every value bound under one Id, nil when none was.
	GetItem func(id string) []any
	// GetString returns the first string bound under one Id, "" when none
	// was.
	GetString func(id string) string
	// GetBool returns the first boolean bound under one Id, false when none
	// was.
	GetBool func(id string) bool
	// GetInt returns the first int bound under one Id, 0 when none was.
	GetInt func(id string) int
	// GetFloat returns the first float bound under one Id, 0 when none was.
	GetFloat func(id string) float64
	// GetStrings returns every string bound under one Id, in order.
	GetStrings func(id string) []string
	// GetInts returns every int bound under one Id, in order.
	GetInts func(id string) []int
	// GetFloats returns every float bound under one Id, in order.
	GetFloats func(id string) []float64

	// New mints a fresh instance of this route: the same declaration, with
	// Items empty and Handler closed over the new instance. The dispatch
	// calls it once per request, so two requests in flight never share
	// bound values.
	New func() *Route

	// Handler runs the route against the values in Items and returns the
	// status it answered with. It is the route package's own RouteHandler,
	// closed over the sandbox and over this instance.
	Handler func() int
}

// NewRoute returns an empty Route with Items open and every Get* reader bound
// to it. A generated NewRoute fills the declaration, New and the handler on top
// of what this returns, so every route reads its values the same way.
func NewRoute() *Route {
	route := &Route{
		Examples: []string{},
		Paths:    []RoutePath{},
		Headers:  []RouteField{},
		Params:   []RouteField{},
		Items:    map[string][]any{},
	}

	route.GetItem = func(id string) []any {
		return route.Items[id]
	}
	route.GetString = func(id string) string {
		value, _ := firstRouteItem(route, id).(string)
		return value
	}
	route.GetBool = func(id string) bool {
		value, _ := firstRouteItem(route, id).(bool)
		return value
	}
	route.GetInt = func(id string) int {
		value, _ := firstRouteItem(route, id).(int)
		return value
	}
	route.GetFloat = func(id string) float64 {
		value, _ := firstRouteItem(route, id).(float64)
		return value
	}
	route.GetStrings = func(id string) []string {
		values := []string{}
		for _, item := range route.Items[id] {
			if value, ok := item.(string); ok {
				values = append(values, value)
			}
		}
		return values
	}
	route.GetInts = func(id string) []int {
		values := []int{}
		for _, item := range route.Items[id] {
			if value, ok := item.(int); ok {
				values = append(values, value)
			}
		}
		return values
	}
	route.GetFloats = func(id string) []float64 {
		values := []float64{}
		for _, item := range route.Items[id] {
			if value, ok := item.(float64); ok {
				values = append(values, value)
			}
		}
		return values
	}

	return route
}

// firstRouteItem is the one value of an Id a scalar reader wants, nil when
// nothing was bound under it.
func firstRouteItem(route *Route, id string) any {
	items := route.Items[id]
	if len(items) == 0 {
		return nil
	}
	return items[0]
}
