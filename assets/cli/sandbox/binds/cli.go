package binds

import (
	api "{{.Module}}/sandbox/api"
	deps "{{.Module}}/sandbox/deps"
	cli "{{.Module}}/sandbox/internal/cli"
{{range .Commands}}
	{{.Name}} "{{$.Module}}/sandbox/internal/commands/{{.Name}}"{{end}}
)

func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.Commands = append(sandbox.Cli.Commands,
{{- range .Commands}}
		{{.Name}}.NewCommand(deps, sandbox),
{{- end}}
	)

	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, sandbox, args)
	}
}
