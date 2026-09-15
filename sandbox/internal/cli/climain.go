package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// Exit codes. Kept here (not in sandbox/api) so the cli layer has no
// dependency on the contract package. They mirror sandbox/api: 0 success,
// 1 a well-formed command that failed, 2 a command line that was wrong.
const (
	ExitOk      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// The declared types a flag or an arg may carry, as entries.yaml spells them.
const (
	typeBoolean = "boolean"
	typeInt     = "int"
	typeFloat   = "float"
)

// subjectFlag and subjectArg are how a usage error names what it is about.
const (
	subjectFlag = "flag"
	subjectArg  = "arg"
)

// quietId is the flag every command that wants a silent run declares. The
// dispatch does the silencing once, so no handler has to.
const quietId = "quiet"

// helpId is the command the empty command line falls back to.
const helpId = "help"

// CliMain reads the verb, finds the command of Cli.Commands that answers to
// it, binds the rest of the command line to a copy of that command's declared
// flags and args, and calls its handler. Nothing here is generated per command:
// every command is one declaration built by its own NewCommand and collected by
// sandbox/internal/cli/new.go, so this dispatch is the same file in every
// project.
// `help` is reached through that same path — it is a declared command whose
// files `agnos build` happens to write itself — and directly only
// for the empty command line below.
func CliMain(sandbox *api.Sandbox, args []string) int {

	//verb := sandbox.Deps.Argvdeps.New(args)

	//action, err := verb.GetNextStringArg()

	return 0
}
