package start_server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
	server "{{.Module}}/sandbox/internal/server"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	err := server.ServerMain(deps, api.ServeProps{
		Addr:           entries.Addr,
		ReadTimeoutMs:  entries.ReadTimeoutMs,
		WriteTimeoutMs: entries.WriteTimeoutMs,
	})
	if err != nil {
		deps.Std.Error("server stopped: %s \n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
