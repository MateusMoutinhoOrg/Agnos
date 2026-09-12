package start_server

import (
	"{{.Module}}/sandbox/api"
	server "{{.Module}}/sandbox/internal/server"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	err := server.ServerMain(sandbox, api.ServeProps{
		Addr:           entries.Addr,
		ReadTimeoutMs:  entries.ReadTimeoutMs,
		WriteTimeoutMs: entries.WriteTimeoutMs,
	})
	if err != nil {
		sandbox.Deps.Std.Error("server stopped: %s \n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
