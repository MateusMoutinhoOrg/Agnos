package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/interview"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

// CommandHandler runs the `interview` command: it hands the session the
// directory to work on and lets it drive the rest of the command surface. The
// commands the session runs report their own results, so there is nothing to
// print here beyond a failure to start.
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if err := interviewAction.Interview(sandbox, props.Path); err != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
