package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, explain_error := explainRouteAction.ExplainRoute(sandbox, api.ExplainRouteProps{
		Path:        props.Path,
		Method:      entries.Method,
		RequestPath: entries.RequestPath,
		Headers:     entries.Header,
		Cookies:     entries.Cookie,
	})
	if explain_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", explain_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
