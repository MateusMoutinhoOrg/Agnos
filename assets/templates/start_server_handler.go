package start_server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// Handle backs `start-server`: it serves sandbox.Server until the
// process is asked to stop, and answers the command line once it has.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := sandbox.Server.Serve(api.ServeProps{
		Addr:              input.Addr,
		ReadTimeoutMs:     input.ReadTimeoutMs,
		WriteTimeoutMs:    input.WriteTimeoutMs,
		ShutdownTimeoutMs: input.ShutdownTimeoutMs,
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "server stopped: "+err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
