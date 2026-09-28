package explain_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	explainRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/explain_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, explain_error := explainRouteAction.ExplainRoute(sandbox, api.ExplainRouteProps{
		Path:        props.Path,
		Method:      entries.Method,
		RequestPath: entries.RequestPath,
		Headers:     entries.Header,
		Cookies:     entries.Cookie,
	})
	if explain_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", explain_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
