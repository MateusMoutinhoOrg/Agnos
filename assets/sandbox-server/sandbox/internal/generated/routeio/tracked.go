package routeio

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// Tracked wraps one response so the dispatch can tell whether a handler
// answered the request. It returns the wrapper to hand to the handlers and a
// reader of the status the request was answered with, 0 while none was.
//
// Answering is what ends a chain, and there are two ways to answer: setting a
// status, or writing a byte — the http library sends a 200 ahead of the first
// byte of a body, so a Write before any SetStatus is that 200, said out loud.
// SetHeader alone answers nothing, which is how a middleware adds a header to
// whatever answers after it. That rule lives here rather than in serverdeps
// because a dep states what a library can do and never what this project does
// with it — the contract is a struct of function fields precisely so the
// sandbox can wrap it like this.
func Tracked(response serverdeps.Response) (serverdeps.Response, func() int) {
	status := 0

	tracked := response
	tracked.SetStatus = func(code int) {
		if status != 0 {
			return
		}
		status = code
		response.SetStatus(code)
	}
	tracked.Write = func(body []byte) error {
		if status == 0 {
			tracked.SetStatus(api.StatusOk)
		}
		return response.Write(body)
	}

	return tracked, func() int {
		return status
	}
}
