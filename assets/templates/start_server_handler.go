package start_server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// InternalPureHandler backs `start-server`: it serves sandbox.Server until the
// process is asked to stop, and answers the command line once it has.
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := sandbox.Server.Serve(api.ServeProps{
		Addr:              entries.Addr,
		ReadTimeoutMs:     entries.ReadTimeoutMs,
		WriteTimeoutMs:    entries.WriteTimeoutMs,
		ShutdownTimeoutMs: entries.ShutdownTimeoutMs,
	})
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "server stopped: "+err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
