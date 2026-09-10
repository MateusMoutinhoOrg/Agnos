package serverdeps

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"{{.Module}}/sandbox/deps"
	serverdeps "{{.Module}}/sandbox/deps/serverdeps"
)

// Bind fills deps.Deps.Serverdeps.NewServer with the http-server
// implementation built on the standard library's net/http package.
func Bind(deps *deps.Deps) {
	deps.Serverdeps.NewServer = newServer
}

// newServer builds one net/http server over props. Every request goes to the
// single props.Handler: no ServeMux, no PathValue, no routing of any kind —
// that is the sandbox's business, and the contract keeps this adapter free of
// it.
func newServer(props serverdeps.ServerProps) serverdeps.Server {
	inner := &http.Server{
		Addr:         props.Addr,
		ReadTimeout:  time.Duration(props.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(props.WriteTimeoutMs) * time.Millisecond,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			props.Handler(newRequest(request), newResponse(writer))
		}),
	}

	return serverdeps.Server{
		Listen: func() error {
			// A server stopped through Shutdown is a clean stop, not a
			// failure: net/http reports it as ErrServerClosed, which the
			// contract turns back into a nil error.
			err := inner.ListenAndServe()
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		},
		Shutdown: func() error {
			return inner.Shutdown(context.Background())
		},
	}
}

// newRequest copies one net/http request onto the sandbox's local
// serverdeps.Request. The body is read at most once and cached, so a second
// ReadBody returns what the first one read rather than an empty slice.
func newRequest(request *http.Request) serverdeps.Request {
	var body []byte
	var body_err error
	read := false

	return serverdeps.Request{
		GetMethod: func() string {
			return request.Method
		},
		GetPath: func() string {
			return request.URL.Path
		},
		GetHeader: func(key string) string {
			return request.Header.Get(key)
		},
		GetQueryParam: func(name string) string {
			return request.URL.Query().Get(name)
		},
		GetQueryAll: func(name string) []string {
			return request.URL.Query()[name]
		},
		ReadBody: func(limit int) ([]byte, error) {
			if read {
				return body, body_err
			}
			read = true
			body, body_err = readBody(request, limit)
			return body, body_err
		},
		GetRemoteAddr: func() string {
			return request.RemoteAddr
		},
	}
}

// readBody drains the request body, refusing one longer than limit. A limit
// of -1 reads it whole. One byte past the limit is read on purpose: that is
// what tells a body exactly at the limit from one over it.
func readBody(request *http.Request, limit int) ([]byte, error) {
	if request.Body == nil {
		return []byte{}, nil
	}
	defer request.Body.Close()

	if limit < 0 {
		return io.ReadAll(request.Body)
	}

	body, err := io.ReadAll(io.LimitReader(request.Body, int64(limit)+1))
	if err != nil {
		return body, err
	}
	if len(body) > limit {
		return body[:limit], errors.New("request body is larger than the declared limit")
	}
	return body, nil
}

// newResponse copies one net/http response writer onto the sandbox's local
// serverdeps.Response.
func newResponse(writer http.ResponseWriter) serverdeps.Response {
	return serverdeps.Response{
		SetHeader: func(key string, value string) {
			writer.Header().Set(key, value)
		},
		SetStatus: func(code int) {
			writer.WriteHeader(code)
		},
		Write: func(body []byte) error {
			_, err := writer.Write(body)
			return err
		},
	}
}
