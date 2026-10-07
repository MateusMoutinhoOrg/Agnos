package api

import (
	opinatedagnosserver "{{.Module}}/sandbox/deps/OpinatedAgnosServer"
)

// Every type here is the OpinatedAgnosServer contract's own, aliased so a
// handler reads it as api.<Name>: sandbox/deps/OpinatedAgnosServer holds the
// whole doc of each.

// PathType is what one segment a Path reads has to convert to.
type PathType = opinatedagnosserver.PathType

// StringPath takes any slice, bound as a string.
const StringPath = opinatedagnosserver.StringPath

// IntegerPath takes one segment reading as a whole number, bound as an int.
const IntegerPath = opinatedagnosserver.IntegerPath

// NumberPath takes one segment reading as a number, bound as a float64.
const NumberPath = opinatedagnosserver.NumberPath

// UuidPath takes one segment reading as a canonical uuid, bound as a string.
const UuidPath = opinatedagnosserver.UuidPath

// Path is one entry of `paths` in route.yaml.
type Path = opinatedagnosserver.Path

// ParameterFont is one place of the request a Parameter is read from.
type ParameterFont = opinatedagnosserver.ParameterFont

// HeaderParam reads a request header, matched without regard to case.
const HeaderParam = opinatedagnosserver.HeaderParam

// QueryParam reads a query-string parameter.
const QueryParam = opinatedagnosserver.QueryParam

// CookieParam reads a request cookie.
const CookieParam = opinatedagnosserver.CookieParam

// ParameterType is the type a Parameter is converted to before it reaches
// Entries.
type ParameterType = opinatedagnosserver.ParameterType

// StringType is bound as a string.
const StringType = opinatedagnosserver.StringType

// NumberType is bound as a float64.
const NumberType = opinatedagnosserver.NumberType

// BooleanType is bound as a bool: true/1 or false/0.
const BooleanType = opinatedagnosserver.BooleanType

// DateTimeType is bound as a string that has to read as RFC 3339.
const DateTimeType = opinatedagnosserver.DateTimeType

// StringArrayType is bound as a []string.
const StringArrayType = opinatedagnosserver.StringArrayType

// IntegerType is bound as an int.
const IntegerType = opinatedagnosserver.IntegerType

// IntegerArrayType is bound as a []int.
const IntegerArrayType = opinatedagnosserver.IntegerArrayType

// AnyMethod is the one entry of AcceptMethods that accepts every http method.
const AnyMethod = opinatedagnosserver.AnyMethod

// Parameter is one entry of `parameters` in route.yaml.
type Parameter = opinatedagnosserver.Parameter

// RouteBody is the request body a route declares — the parsed form of `body:`
// in its route.yaml.
type RouteBody = opinatedagnosserver.RouteBody

// RouteFailure is one way a request did not get answered; an
// InternalPureHandler refuses a request by returning one, built by
// Deps.OpinatedAgnosServer.Fail.
type RouteFailure = opinatedagnosserver.RouteFailure

// Route is one http route of the project: the whole of what its route.yaml
// declares, plus the handler behind it.
type Route = opinatedagnosserver.Route
