package start_server

import (
	"{{.Module}}/sandbox/api"
	server "{{.Module}}/sandbox/internal/server"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	err := server.ServerMain(sandbox, api.ServeProps{
		Addr:           command.GetString("addr"),
		ReadTimeoutMs:  command.GetInt("read_timeout_ms"),
		WriteTimeoutMs: command.GetInt("write_timeout_ms"),
	})
	if err != nil {
		sandbox.Deps.Std.Error("server stopped: %s \n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
