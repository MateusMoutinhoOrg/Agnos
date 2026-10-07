package opinatedagnosserver

import (
	"fmt"

	opinatedagnosserver "{{.Module}}/sandbox/deps/OpinatedAgnosServer"
	serializables "{{.Module}}/sandbox/deps/serializables"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// The writers below answer a request in one call: the Content-Type, the
// status and the body, in the order the http library needs them. Each answers
// the request, so a handler that calls one has ended the chain.

// writeError writes one failure onto the response: the body is always the same
// JSON object — {"error": "...", "field": "..."} — so a client parses one shape
// whatever went wrong. It returns the message as an error, so a Handle* file
// answers and reports in one line. The status the chain reads is the one
// written here, never the one returned.
func writeError(serializer serializables.Sandbox, response serverdeps.Response, status int, field string, message string) error {
	body := serializer.CreateObject()
	body.AddItemToObject("error", message)
	body.AddItemToObject("field", field)

	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(status)
	response.Write([]byte(serializer.SerializeToJson(body)))

	return fmt.Errorf("%s", message)
}

// writeJSON answers with one document serialized as JSON.
func writeJSON(serializer serializables.Sandbox, response serverdeps.Response, status int, document *serializables.SerializibleObject) error {
	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(status)
	return response.Write([]byte(serializer.SerializeToJson(document)))
}

// writeText answers with one text as text/plain.
func writeText(response serverdeps.Response, status int, text string) error {
	response.SetHeader("Content-Type", "text/plain; charset=utf-8")
	response.SetStatus(status)
	return response.Write([]byte(text))
}

// redirect answers by sending the caller to location, with one of the
// redirect statuses.
func redirect(response serverdeps.Response, status int, location string) error {
	response.SetHeader("Location", location)
	response.SetStatus(status)
	return nil
}

// tracked wraps one response so the dispatch can tell whether a handler
// answered the request. It returns the wrapper to hand to the handlers and a
// reader of the status the request was answered with, 0 while none was.
//
// Answering is what ends a chain, and there are two ways to answer: setting a
// status, or writing a byte — the http library sends a 200 ahead of the first
// byte of a body, so a Write before any SetStatus is that 200, said out loud.
// SetHeader alone answers nothing, which is how a middleware adds a header to
// whatever answers after it. That rule lives here rather than in serverdeps
// because a raw dep states what a library can do and never what a project
// does with it — the contract is a struct of function fields precisely so the
// lib can wrap it like this.
func tracked(response serverdeps.Response) (serverdeps.Response, func() int) {
	status := 0

	wrapper := response
	wrapper.SetStatus = func(code int) {
		if status != 0 {
			return
		}
		status = code
		response.SetStatus(code)
	}
	wrapper.Write = func(body []byte) error {
		if status == 0 {
			wrapper.SetStatus(opinatedagnosserver.StatusOk)
		}
		return response.Write(body)
	}

	return wrapper, func() int {
		return status
	}
}
