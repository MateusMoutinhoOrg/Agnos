package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// forceFieldId is the answer that turns a scaffold from something that creates
// a project into something that writes over one.
const forceFieldId = "force"

// removePrefix is how every command that drops one unit of a layer is spelled,
// which is why they are matched by their shape rather than listed one by one.
const removePrefix = "remove-"

// destructiveVerbs are the commands that take something away and are not
// spelled remove-<unit>. They are listed because nothing about their name says
// so: disable-extension stops a mechanic generating, and publish reaches a
// release out of this machine and cannot be taken back.
var destructiveVerbs = map[string]bool{
	"disable-extension": true,
	"publish":           true,
}

// purgeUnits is the layer each purge takes with it, said as the unit a person
// counts it in. It is what lets the confirm screen name what is about to go
// instead of asking to confirm a command line.
var purgeUnits = map[string]string{
	"cli-purge":    "command",
	"server-purge": "route",
	"front-purge":  "page",
	"deps-purge":   "dep",
}

// DestructiveVerb reports a command that removes or overwrites something the
// project already has, from its name alone. A menu is built before any question
// is answered, so this is the most a menu can know.
func DestructiveVerb(sandbox *api.Sandbox, verb string) bool {
	if _, purge := purgeUnits[verb]; purge {
		return true
	}
	if sandbox.Deps.Stringsdeps.HasPrefix(verb, removePrefix) {
		return true
	}
	return destructiveVerbs[verb]
}

// Destructive is DestructiveVerb plus the one command whose answers decide it:
// a scaffold creates a project, and only --force makes it write over one.
func Destructive(sandbox *api.Sandbox, command api.Command, values map[string][]any) bool {
	if verbOf(command) == scaffoldVerb {
		return answeredYes(values, forceFieldId)
	}
	return DestructiveVerb(sandbox, verbOf(command))
}

// Losses is what a purge is about to take away: the unit it is counted in and
// the names themselves, read off the project with the same lists the questions
// are built from. Every other command answers with no unit, and its confirm
// screen says only that it removes something.
func Losses(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command) (string, []string) {
	unit, purge := purgeUnits[verbOf(command)]
	if !purge {
		return "", nil
	}

	names := []string{}
	for _, option := range unitOptions(sandbox, io, unit) {
		names = append(names, option.Id)
	}

	return unit, names
}

// unitOptions reads the units of one layer, by the unit's own word for itself.
func unitOptions(sandbox *api.Sandbox, io *smartio.SmartIO, unit string) []interviewer.AlternativeOption {
	switch unit {
	case "command":
		return commandOptions(sandbox, io)
	case "route":
		return dirOptions(sandbox, io, routesDir)
	case "page":
		return pageOptions(sandbox, io)
	case "dep":
		return dirOptions(sandbox, io, utils.ContractsDir)
	}
	return []interviewer.AlternativeOption{}
}

// SafeFirst orders a menu so that pressing enter without reading it never runs
// something that takes a layer away or installs one. The inert row — the way
// back, or the way out — leads the menu when the row that would otherwise lead
// it is one of those, and trails it when the menu is harmless.
//
// It is the answer to the one thing a guided screen must never do: a person who
// has not read the rows yet is told the default is safe, and the default has to
// be.
func SafeFirst(rows []interviewer.AlternativeOption, inert interviewer.AlternativeOption, dangerous bool) []interviewer.AlternativeOption {
	if !dangerous {
		return append(rows, inert)
	}
	return append([]interviewer.AlternativeOption{inert}, rows...)
}
