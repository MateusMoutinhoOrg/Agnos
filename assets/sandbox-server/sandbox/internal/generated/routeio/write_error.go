package routeio

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// WriteError writes one failure onto the response: the body is always the same
// JSON object — {"error": "...", "field": "..."} — so a client parses one shape
// whatever went wrong, and the failure is logged on the progress channel as it
// is written.
//
// It is the writer, not the policy: what decides *whether* a failure is
// answered this way is the project's own Handle* file, and this is what those
// files call by default. A route raises a failure through Fail instead, which
// is what routes it to the right one.
//
// It returns the message as an error, so a handler answers and reports in one
// line. The status the chain reads is the one written here, never the one
// returned.
func WriteError(sandbox *api.Sandbox, response serverdeps.Response, status int, field string, message string) error {
	sandbox.Deps.Std.Log("route error %d %s %s \n", status, field, message)

	body := sandbox.Deps.Serializables.CreateObject()
	body.AddItemToObject("error", message)
	body.AddItemToObject("field", field)

	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(status)
	response.Write([]byte(sandbox.Deps.Serializables.SerializeToJson(body)))

	return sandbox.Deps.Std.Errorf("%s", message)
}
