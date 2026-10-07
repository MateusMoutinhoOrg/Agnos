package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// OpinatedPrefix starts the name of every opinated lib: a catalog dep whose
// contract carries an agnos mechanic itself — its declarations, its dispatch,
// its failures — rather than a library's raw capability. It is the one kind of
// dep that is not 0-opinionated, which is why the name says so.
const OpinatedPrefix = "OpinatedAgnos"

// The four opinated libs, one per mechanic that has code: what used to be
// generated under sandbox/internal/generated/ for that mechanic, outside the
// sandbox, installed by its -init.
const (
	OpinatedAgnosCli      = OpinatedPrefix + "Cli"
	OpinatedAgnosServer   = OpinatedPrefix + "Server"
	OpinatedAgnosFront    = OpinatedPrefix + "Front"
	OpinatedAgnosDatabase = OpinatedPrefix + "Database"
)

// opinatedLibs maps each mechanic to the lib its code lives in.
var opinatedLibs = map[string]string{
	ExtensionSandboxCli:      OpinatedAgnosCli,
	ExtensionSandboxServer:   OpinatedAgnosServer,
	ExtensionSandboxFront:    OpinatedAgnosFront,
	ExtensionSandboxDatabase: OpinatedAgnosDatabase,
}

// OpinatedLib is the lib one mechanic's code lives in, "" for a mechanic with
// none.
func OpinatedLib(extension string) string {
	return opinatedLibs[extension]
}

// IsOpinatedLib reports whether a dep name is an opinated lib's.
func IsOpinatedLib(sandbox *api.Sandbox, name string) bool {
	return sandbox.Deps.Stringsdeps.HasPrefix(name, OpinatedPrefix)
}
