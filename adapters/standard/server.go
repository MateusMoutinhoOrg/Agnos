package standard

import (
	"bytes"
	"io"
	"net/http"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serverdeps"
)

// NewRequestFactory returns the value that fills deps.Deps.NewRequest: the implementation
// of the HTTP request dependency using the standard library's net/http package.
func NewRequestFactory(s *StandardAdapter) func(url string) serverdeps.Request {
	return func(url string) serverdeps.Request {
		headers := make(map[string]string)
		method := "GET"
		var reqBody []byte

		return serverdeps.Request{
			AddHeader: func(key string, value string) {
				headers[key] = value
			},
			SetMethod: func(m string) {
				method = m
			},
			SetBody: func(body []byte) {
				reqBody = body
			},
			Fetch: func() (serverdeps.Response, error) {
				var bodyReader io.Reader
				if reqBody != nil {
					bodyReader = bytes.NewReader(reqBody)
				}

				req, err := http.NewRequest(method, url, bodyReader)
				if err != nil {
					return serverdeps.Response{}, err
				}

				for k, v := range headers {
					req.Header.Add(k, v)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return serverdeps.Response{}, err
				}

				return serverdeps.Response{
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
						if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
							return buf[:n], err
						}
						if err == io.ErrUnexpectedEOF || err == io.EOF {
							return buf[:n], nil
						}
						return buf[:n], nil
					},
					Close: func() error {
						return resp.Body.Close()
					},
				}, nil
			},
		}
	}
}
