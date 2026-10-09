package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// retiredGenerated is, per mechanic, every generated file and package an older
// build wrote under sandbox/internal/generated/ that nothing renders any more:
// the dispatch, the binders, the matcher and the shared io packages the
// OpinionatedAgnos libs replaced, and the registries' old cli/cli and
// server/server packages, now one level up, and help, version and help_flag
// from before every command sat in the folder of its category. Each one names a symbol the aliased api dropped
// (api.NewCommand, api.BindRoute, a Request held as any), so a tree still
// carrying it would not compile.
var retiredGenerated = map[string][]string{
	ExtensionCli: {
		LegacyGeneratedDir + "/cli/cli",
		CommandsDir + "/help",
		CommandsDir + "/version",
		CommandsDir + "/help_flag",
		LegacyGeneratedDir + "/cli/command",
		LegacyGeneratedDir + "/cliio",
		LegacyGeneratedDir + "/trigger",
	},
	ExtensionServer: {
		LegacyGeneratedDir + "/server/server",
		LegacyGeneratedDir + "/server/route",
		LegacyGeneratedDir + "/routeio",
		LegacyGeneratedDir + "/trigger",
	},
	ExtensionFront: {
		LegacyGeneratedDir + "/frontio",
	},
	ExtensionDatabase: {
		LegacyGeneratedDir + "/databaseio",
	},
}

// retiredReplacements names, for each retired generated package, what took its
// place — the message verify answers an import of it with, which is the whole
// of the migration a project's hand-written files need.
var retiredReplacements = map[string]string{
	LegacyGeneratedDir + "/cliio":        "sandbox.Deps.OpinionatedAgnosCli: Fail, FailWithCause, FailureOf",
	LegacyGeneratedDir + "/cli/command":  "sandbox.Deps.OpinionatedAgnosCli.NewCommand()",
	LegacyGeneratedDir + "/trigger":      "sandbox.Deps.OpinionatedAgnosCli.MatchTrigger",
	LegacyGeneratedDir + "/routeio":      "sandbox.Deps.OpinionatedAgnosServer: Fail, FailWithCause, FailureOf, WriteError, WriteJSON, WriteText, Redirect, ValidateSchema, ValidateForm and the Read*/Item* readers; route.Request and route.Response for RequestOf and ResponseOf",
	LegacyGeneratedDir + "/server/route": "sandbox.Deps.OpinionatedAgnosServer.NewRoute()",
	LegacyGeneratedDir + "/frontio":      "sandbox.Deps.OpinionatedAgnosFront: Resolve, SafePath, ExtensionOf, ContentTypeOf, and the constants of sandbox/deps/OpinionatedAgnosFront",
	LegacyGeneratedDir + "/databaseio":   "sandbox.Deps.OpinionatedAgnosDatabase: Fail, Collection, ReadString, ReadInt, ReadFloat, TextMatches, IntInRange, FloatInRange",
}

// RetiredReplacement is what replaced one retired generated package, named by
// its project-relative path, and whether the path is one.
func RetiredReplacement(pkg string) (string, bool) {
	replacement, retired := retiredReplacements[pkg]
	return replacement, retired
}

// RetiredGenerated is every retired generated path of one mechanic.
func RetiredGenerated(extension string) []string {
	return retiredGenerated[extension]
}

// RemoveRetiredGenerated removes, from the transaction, every retired
// generated path of one mechanic the tree still carries. The build runs it
// for each mechanic that is on, so a project's first build on the libs drops
// what they replaced; an <x>-purge runs it too, so a purge leaves nothing of
// an older build behind.
func RemoveRetiredGenerated(sandbox *api.Sandbox, io *stagedfs.StagedFS, extension string) {
	for _, path := range RetiredGenerated(extension) {
		if io.IsFile(path) {
			io.RemoveDir(path)
			continue
		}
		if !io.IsDir(path) {
			continue
		}
		for _, entry := range io.ListAllRecursively(path) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(path)
	}
}
