package serverdeps

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	serverdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serverdeps"
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

	// listener is what Bind opened, kept for Listen to serve: binding and
	// serving are two steps so the sandbox can tell an address it cannot
	// open from a server that failed while running.
	var listener net.Listener
	bind := func() error {
		opened, err := net.Listen("tcp", props.Addr)
		if err != nil {
			return err
		}
		listener = opened
		return nil
	}

	return serverdeps.Server{
		Bind: bind,
		Listen: func() error {
			if listener == nil {
				if err := bind(); err != nil {
					return err
				}
			}
			// A server stopped through Shutdown is a clean stop, not a
			// failure: net/http reports it as ErrServerClosed, which the
			// contract turns back into a nil error.
			err := inner.Serve(listener)
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		},
		Shutdown: func() error {
			if props.ShutdownTimeoutMs <= 0 {
				return inner.Shutdown(context.Background())
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(props.ShutdownTimeoutMs)*time.Millisecond)
			defer cancel()
			return inner.Shutdown(ctx)
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

	read_body := func(limit int) ([]byte, error) {
		if read {
			return body, body_err
		}
		read = true
		body, body_err = readBody(request, limit)
		return body, body_err
	}

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
		GetHeaders: func() map[string][]string {
			headers := map[string][]string{}
			for key, values := range request.Header {
				headers[key] = append([]string{}, values...)
			}
			return headers
		},
		GetHost: func() string {
			return request.Host
		},
		GetCookie: func(name string) string {
			cookie, err := request.Cookie(name)
			if err != nil {
				return ""
			}
			return cookie.Value
		},
		GetQueryParam: func(name string) string {
			return request.URL.Query().Get(name)
		},
		GetQueryAll: func(name string) []string {
			return request.URL.Query()[name]
		},
		ReadBody: read_body,
		ReadForm: func(limit int) (map[string][]string, error) {
			raw, err := read_body(limit)
			if err != nil {
				return map[string][]string{}, err
			}
			values, err := url.ParseQuery(string(raw))
			if err != nil {
				return map[string][]string{}, err
			}
			return values, nil
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
		AddHeader: func(key string, value string) {
			writer.Header().Add(key, value)
		},
		GetHeader: func(key string) string {
			return writer.Header().Get(key)
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
