package api

import (
	opinionatedagnosserver "{{.Module}}/sandbox/deps/OpinionatedAgnosServer"
)

// Every type here is the OpinionatedAgnosServer contract's own, aliased so a
// handler reads it as api.<Name>: sandbox/deps/OpinionatedAgnosServer holds the
// whole doc of each.

// PathType is what one segment a Path reads has to convert to.
type PathType = opinionatedagnosserver.PathType

// PathString takes any slice, bound as a string.
const PathString = opinionatedagnosserver.PathString

// PathInteger takes one segment reading as a whole number, bound as an int.
const PathInteger = opinionatedagnosserver.PathInteger

// PathNumber takes one segment reading as a number, bound as a float64.
const PathNumber = opinionatedagnosserver.PathNumber

// PathUuid takes one segment reading as a canonical uuid, bound as a string.
const PathUuid = opinionatedagnosserver.PathUuid

// Path is one entry of `paths` in route.yaml.
type Path = opinionatedagnosserver.Path

// ParameterSource is one place of the request a Parameter is read from.
type ParameterSource = opinionatedagnosserver.ParameterSource

// SourceHeader reads a request header, matched without regard to case.
const SourceHeader = opinionatedagnosserver.SourceHeader

// SourceQuery reads a query-string parameter.
const SourceQuery = opinionatedagnosserver.SourceQuery

// SourceCookie reads a request cookie.
const SourceCookie = opinionatedagnosserver.SourceCookie

// ParameterType is the type a Parameter is converted to before it reaches
// Input.
type ParameterType = opinionatedagnosserver.ParameterType

// ParameterString is bound as a string.
const ParameterString = opinionatedagnosserver.ParameterString

// ParameterNumber is bound as a float64.
const ParameterNumber = opinionatedagnosserver.ParameterNumber

// ParameterBoolean is bound as a bool: true/1 or false/0.
const ParameterBoolean = opinionatedagnosserver.ParameterBoolean

// ParameterDateTime is bound as a string that has to read as RFC 3339.
const ParameterDateTime = opinionatedagnosserver.ParameterDateTime

// ParameterStringArray is bound as a []string.
const ParameterStringArray = opinionatedagnosserver.ParameterStringArray

// ParameterInteger is bound as an int.
const ParameterInteger = opinionatedagnosserver.ParameterInteger

// ParameterIntegerArray is bound as a []int.
const ParameterIntegerArray = opinionatedagnosserver.ParameterIntegerArray

// AnyMethod is the one entry of Methods that accepts every http method.
const AnyMethod = opinionatedagnosserver.AnyMethod

// Parameter is one entry of `parameters` in route.yaml.
type Parameter = opinionatedagnosserver.Parameter

// RouteBody is the request body a route declares — the parsed form of `body:`
// in its route.yaml.
type RouteBody = opinionatedagnosserver.RouteBody

// RouteFailure is one way a request did not get answered; a
// Handle refuses a request by returning one, built by
// Deps.OpinionatedAgnosServer.Fail.
type RouteFailure = opinionatedagnosserver.RouteFailure

// Route is one http route of the project: the whole of what its route.yaml
// declares, plus the handler behind it.
type Route = opinionatedagnosserver.Route
