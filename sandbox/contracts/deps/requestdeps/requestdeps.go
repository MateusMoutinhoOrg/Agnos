package requestdeps

// This package is the sandbox's *copy* of the api an HTTP client library
// exposes — the same mechanic as verbdeps, keepdeps, embeddeps and iodeps,
// for the same reason: opening a socket is an OS-bound effect, so `net/http`
// may not appear inside the sandbox. The contract is restated here, and the
// adapter — which lives outside the sandbox — is what fills it.
//
// It differs from the other copies in one way: a request is created per call
// rather than injected once, so the entry point is the Deps.NewRequest
// function field rather than a library struct. Everything below is what that
// field hands back.
//
// The tracker in sandbox/ never calls it — nothing it does leaves the
// machine. It is carried as a standing capability of the template, filled by
// the standard adapter over `net/http`. See the Deps.NewRequest field.

// Request is one HTTP request under construction, handed back by
// Deps.NewRequest already bound to a url. The setters mutate the pending
// request and may be called in any order; nothing leaves the machine until
// Fetch is called, and a Request may be sent more than once.
//
// The method defaults to GET and the body to none, so a plain read is
// NewRequest followed by Fetch.
type Request struct {
	// AddHeader sets one header on the pending request, replacing whatever
	// value that key carried before.
	AddHeader func(key string, value string)

	// SetMethod sets the HTTP method the request is sent with — "POST",
	// "PUT", "DELETE". It defaults to "GET".
	SetMethod func(method string)

	// SetBody sets the bytes sent as the request body. It defaults to none.
	SetBody func(body []byte)

	// Fetch sends the request and returns the response. The error reports a
	// request that could not be built or a round trip that failed; an HTTP
	// error status is *not* an error here, it is reported by
	// Response.GetStatusCode. A Response returned without an error holds an
	// open body the caller must Close.
	Fetch func() (Response, error)
}

// Response is one HTTP response, handed back by Request.Fetch with its body
// still open. The caller must Close it, whether or not the body is read.
type Response struct {
	// GetStatusCode returns the response's HTTP status code.
	GetStatusCode func() int

	// GetHeader returns the first value of one response header, or "" when
	// the response carries no such header.
	GetHeader func(key string) string

	// ReadBody reads at most size bytes of the response body, or the whole
	// body when size is -1. A body shorter than size is returned whole
	// rather than reported as an error, so a short read is not a failure.
	ReadBody func(size int) ([]byte, error)

	// Close releases the response body. It must be called for every Response
	// returned without an error, or the underlying connection leaks.
	Close func() error
}
