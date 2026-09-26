package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The two field ids a declared name is read from: the name itself, and the
// spelling a flag answers to when the person gave one instead of letting the
// name decide it.
const (
	nameFieldId       = "name"
	identifierFieldId = "identifier"
)

// namedUnits is every command that takes a name the person types and writes a
// different one down. utils.CommandIdentifier is what rewrites it — lowercased,
// spaces and underscores turned into dashes — and every one of these commands
// funnels through it, whether it is spelled CommandIdentifier, FieldName,
// RouteIdentifier, RouteEntryId or RouteFieldName.
//
// The commands that name something that already exists are not here: they are
// answered from a menu of names read off disk, so there is nothing to rewrite.
var namedUnits = map[string]string{
	"add-command":     "command",
	"add-route":       "route",
	"rename-route":    "route",
	"add-page":        "page",
	"add-flag":        "flag",
	"add-arg":         "argument",
	"add-parameter":   "parameter",
	"add-body-field":  "body field",
	"add-path":        "path",
	"add-database":    "database",
	"add-table":       "table",
	"add-table-field": "field",
}

// NormalizedNotes is what the confirm screen has to say beyond the command
// line: the line promises "you could have typed it yourself", and a name that
// gets rewritten on the way in is the one case where running it does something
// else than the line says. Saying so is the difference between a screen that
// teaches the cli and one that misreports it.
//
// It answers nothing when the name is written down exactly as it was typed,
// which is the common case.
func NormalizedNotes(sandbox *api.Sandbox, command api.Command, values map[string][]any) []string {
	noun, named := namedUnits[verbOf(command)]
	if !named {
		return []string{}
	}

	typed := answeredText(sandbox, values, nameFieldId)
	written := utils.CommandIdentifier(sandbox, typed)
	if typed == "" || written == typed {
		return []string{}
	}

	return []string{sandbox.Deps.Std.Sprintf(
		"%q is written down as the %s %s", typed, noun, spelling(command, values, written))}
}

// spelling is how the rewritten name is met again. A flag given no --identifier
// of its own answers to the name with two dashes in front, so that — and not
// the bare name — is what the person will type at it.
func spelling(command api.Command, values map[string][]any, written string) string {
	if verbOf(command) != "add-flag" || len(values[identifierFieldId]) > 0 {
		return written
	}
	return "--" + written
}
