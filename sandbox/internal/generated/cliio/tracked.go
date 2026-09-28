package cliio

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// Tracked builds the one response of a command line, so the dispatch can tell
// whether a command answered it. It returns the response to hand to the
// commands and a reader of the status it was answered with, and whether it
// was answered at all.
//
// Answering is what ends a chain, and there are two ways to answer: setting a
// status, or printing to stdout — which answers ExitOk, the status a command
// that printed its result ends with. Error and Log answer nothing, which is how
// a middleware says something and hands the command line on. Every print goes
// through sandbox.Deps.Std at the moment it is made, so a middleware silencing
// Std.Log silences Log here too.
//
// The first status set is the one the line exits with. A print only implies
// ExitOk, so a status set after it — a failure the same command raises once it
// has printed part of its output — still replaces it: a command that printed
// and then failed never exits 0.
func Tracked(sandbox *api.Sandbox) (*api.CommandResponse, func() (int, bool)) {
	status := 0
	answered := false
	explicit := false

	response := &api.CommandResponse{}
	response.SetStatus = func(code int) {
		if explicit {
			return
		}
		status, answered, explicit = code, true, true
	}
	response.Printf = func(format string, a ...any) (int, error) {
		if !answered {
			status, answered = api.ExitOk, true
		}
		return sandbox.Deps.Std.Printf(format, a...)
	}
	response.Error = func(format string, a ...any) (int, error) {
		return sandbox.Deps.Std.Error(format, a...)
	}
	response.Log = func(format string, a ...any) (int, error) {
		return sandbox.Deps.Std.Log(format, a...)
	}

	return response, func() (int, bool) {
		return status, answered
	}
}
