package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// OpinionatedPrefix starts the name of every opinionated lib: a catalog dep whose
// contract carries an agnos mechanic itself — its declarations, its dispatch,
// its failures — rather than a library's raw capability. It is the one kind of
// dep that is not 0-opinionated, which is why the name says so.
const OpinionatedPrefix = "OpinionatedAgnos"

// The four opinionated libs, one per mechanic that has code: what used to be
// generated under sandbox/internal/generated/ for that mechanic, outside the
// sandbox, installed by its -init.
const (
	OpinionatedAgnosCli      = OpinionatedPrefix + "Cli"
	OpinionatedAgnosServer   = OpinionatedPrefix + "Server"
	OpinionatedAgnosFront    = OpinionatedPrefix + "Front"
	OpinionatedAgnosDatabase = OpinionatedPrefix + "Database"
)

// opinionatedLibs maps each mechanic to the lib its code lives in.
var opinionatedLibs = map[string]string{
	ExtensionCli:      OpinionatedAgnosCli,
	ExtensionServer:   OpinionatedAgnosServer,
	ExtensionFront:    OpinionatedAgnosFront,
	ExtensionDatabase: OpinionatedAgnosDatabase,
}

// OpinionatedLib is the lib one mechanic's code lives in, "" for a mechanic with
// none.
func OpinionatedLib(extension string) string {
	return opinionatedLibs[extension]
}

// IsOpinionatedLib reports whether a dep name is an opinionated lib's.
func IsOpinionatedLib(sandbox *api.Sandbox, name string) bool {
	return sandbox.Deps.StringsDeps.HasPrefix(name, OpinionatedPrefix)
}
