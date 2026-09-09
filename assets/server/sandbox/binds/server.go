package binds

import (
	api "{{.Module}}/sandbox/api"
	deps "{{.Module}}/sandbox/deps"
	server "{{.Module}}/sandbox/internal/server"
)

// ServerBind wires the generated server.ServerMain (built by `agnos build`
// from every route's route.yaml) onto the sandbox.
func ServerBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Server.Serve = func(props api.ServeProps) error {
		return server.ServerMain(deps, props)
	}
}
