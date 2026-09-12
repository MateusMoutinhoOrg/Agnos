package binds

import (
	api "{{.Module}}/sandbox/api"
	server "{{.Module}}/sandbox/internal/server"
)

// ServerBind wires the generated server.ServerMain (built by `agnos build`
// from every route's route.yaml) onto the sandbox.
func ServerBind(sandbox *api.Sandbox) {
	sandbox.Server.Serve = func(props api.ServeProps) error {
		return server.ServerMain(sandbox, props)
	}
}
