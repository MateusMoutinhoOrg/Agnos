package requestdeps

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
)

// requestTimeout bounds one whole round trip — connect, send, and read the
// response. Without it net/http waits forever, so a hung server would hang
// the calling program with no way for the sandbox to intervene: the contract
// exposes no cancellation, which makes the bound this adapter's job.
const requestTimeout = 30 * time.Second

// Bind fills deps.Deps.RequestLib.NewRequest with the HTTP request
// implementation built on the standard library's net/http package.
func Bind(deps *deps.Deps) {
	deps.Requestdeps.NewRequest = newRequest
}

// newRequest builds one HTTP request bound to url.
func newRequest(url string) requestdeps.Request {
	headers := make(map[string]string)
	method := "GET"
	var reqBody []byte

	return requestdeps.Request{
		AddHeader: func(key string, value string) {
			headers[key] = value
		},
		SetMethod: func(m string) {
			method = m
		},
		SetBody: func(body []byte) {
			reqBody = body
		},
		Fetch: func() (requestdeps.Response, error) {
			var bodyReader io.Reader
			if reqBody != nil {
				bodyReader = bytes.NewReader(reqBody)
			}

			req, err := http.NewRequest(method, url, bodyReader)
			if err != nil {
				return requestdeps.Response{}, err
			}

			for k, v := range headers {
				req.Header.Add(k, v)
			}

			client := &http.Client{Timeout: requestTimeout}
			resp, err := client.Do(req)
			if err != nil {
				return requestdeps.Response{}, err
			}

			return requestdeps.Response{
				GetStatusCode: func() int {
					return resp.StatusCode
				},
				GetHeader: func(key string) string {
					return resp.Header.Get(key)
				},
				ReadBody: func(size int) ([]byte, error) {
					if size == -1 {
						return io.ReadAll(resp.Body)
					}
					buf := make([]byte, size)
					n, err := io.ReadFull(resp.Body, buf)
					// A body shorter than size is the contract's normal
					// case, not a failure: io.ReadFull reports it as EOF
					// or ErrUnexpectedEOF, and both mean "that was all
					// of it".
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return buf[:n], nil
					}
					return buf[:n], err
				},
				Close: func() error {
					return resp.Body.Close()
				},
			}, nil
		},
	}
}
