package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/interview"
)

// CommandHandler runs the `interview` command: it hands the session the
// directory to work on and lets it drive the rest of the command surface. The
// commands the session runs report their own results, so there is nothing to
// print here beyond a failure to start.
func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	if err := interviewAction.Interview(sandbox, command.GetString("path")); err != nil {
		sandbox.Deps.Std.Error("%s\n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
