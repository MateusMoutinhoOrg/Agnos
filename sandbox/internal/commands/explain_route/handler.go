package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_route"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	lines, explain_error := explainRouteAction.ExplainRoute(sandbox, api.ExplainRouteProps{
		Path:        command.GetString("path"),
		Method:      command.GetString("method"),
		RequestPath: command.GetString("request-path"),
		Headers:     command.GetStrings("header"),
		Cookies:     command.GetStrings("cookie"),
	})
	if explain_error != nil {
		sandbox.Deps.Std.Error("%s\n", explain_error.Error())
		return api.ExitFailure
	}

	for _, line := range lines {
		sandbox.Deps.Std.Printf("%s\n", line)
	}
	return api.ExitOk
}
