package routeio

import (
	"{{.Module}}/sandbox/deps/serverdeps"
)

// Tracked wraps one response so the dispatch can tell whether a handler
// answered the request. It returns the wrapper to hand to the handlers and a
// reader that reports whether a status has been set on it.
//
// Setting the status is the whole of what ends a chain: a handler that writes
// bytes without one has not answered, and the next route of the chain runs.
// That rule lives here rather than in serverdeps because a dep states what a
// library can do and never what this project does with it — the contract is a
// struct of function fields precisely so the sandbox can wrap it like this.
func Tracked(response serverdeps.Response) (serverdeps.Response, func() bool) {
	answered := false

	tracked := serverdeps.Response{
		SetHeader: response.SetHeader,
		Write:     response.Write,
		SetStatus: func(code int) {
			answered = true
			response.SetStatus(code)
		},
	}

	return tracked, func() bool {
		return answered
	}
}
