package routeio

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	serializables "{{.Module}}/sandbox/deps/serializables"
)

// The writers below answer a request in one call: the Content-Type, the
// status and the body, in the order the http library needs them. Each answers
// the request, so a handler that calls one has ended the chain.

// WriteJSON answers with one document serialized as JSON.
func WriteJSON(sandbox *api.Sandbox, response serverdeps.Response, status int, document *serializables.SerializibleObject) error {
	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(status)
	return response.Write([]byte(sandbox.Deps.Serializables.SerializeToJson(document)))
}

// WriteText answers with one text as text/plain.
func WriteText(response serverdeps.Response, status int, text string) error {
	response.SetHeader("Content-Type", "text/plain; charset=utf-8")
	response.SetStatus(status)
	return response.Write([]byte(text))
}

// Redirect answers by sending the caller to location, with one of the
// redirect statuses: api.StatusFound, api.StatusSeeOther,
// api.StatusMovedPermanently, api.StatusTemporaryRedirect or
// api.StatusPermanentRedirect.
func Redirect(response serverdeps.Response, status int, location string) error {
	response.SetHeader("Location", location)
	response.SetStatus(status)
	return nil
}
