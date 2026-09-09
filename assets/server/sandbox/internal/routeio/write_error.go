package routeio

import (
	"{{.Module}}/sandbox/deps"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// WriteError is the one way a failure reaches the caller, used both by the
// generated dispatch and by every generated ReadBody. The body is always the
// same JSON object — {"error": "...", "field": "..."} — so a client parses one
// shape whatever went wrong, and the failure is logged on the progress channel
// as it is written.
//
// It returns the status it wrote, so a caller answers and reports in one line.
func WriteError(deps *deps.Deps, response serverdeps.Response, status int, field string, message string) int {
	deps.Std.Log("route error %d %s %s \n", status, field, message)

	body := deps.Serializables.CreateObject()
	body.AddItemToObject("error", message)
	body.AddItemToObject("field", field)

	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(status)
	response.Write([]byte(deps.Serializables.SerializeToJson(body)))

	return status
}
