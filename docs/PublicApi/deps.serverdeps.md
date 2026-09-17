# `deps.Serverdeps`

`sandbox/deps/serverdeps`

## `Sandbox`

Sandbox is the http-server library injected whole as the Deps.Serverdeps field. A server is bound to one address and one handler, so it is created per call rather than injected once: what the sandbox holds is this one-field struct.

| Field | Type | Description |
| --- | --- | --- |
| `NewServer` | `func(props ServerProps) Server` | NewServer builds a server over the given props. It binds nothing until Server.Listen is called. |

## `ServerProps`

ServerProps is everything one server needs to run: where to listen, how long a request and a response may take, and the single function every request is handed to.

| Field | Type | Description |
| --- | --- | --- |
| `Addr` | `string` | Addr is the address to listen on, in the host:port spelling (":8080", "127.0.0.1:3000"). |
| `ReadTimeoutMs` | `int` | ReadTimeoutMs is how long a request has to arrive, in milliseconds. Zero means no timeout. |
| `WriteTimeoutMs` | `int` | WriteTimeoutMs is how long a response has to be written, in milliseconds. Zero means no timeout. |
| `Handler` | `func(request Request, response Response)` | Handler is called once per request, whatever the method or the path. Everything the request needs to be answered is reachable from the two arguments; the handler returns once the response is written. |

## `Server`

Server is one built-but-not-yet-listening http server.

| Field | Type | Description |
| --- | --- | --- |
| `Listen` | `func() error` | Listen binds the address and serves until Shutdown is called or the server fails. It blocks. |
| `Shutdown` | `func() error` | Shutdown stops the server, letting the requests in flight finish. |

## `Request`

Request is one incoming http request, read through function fields only: the sandbox never holds the library's own request type.

| Field | Type | Description |
| --- | --- | --- |
| `GetMethod` | `func() string` | GetMethod returns the http method in upper case ("GET", "POST"). |
| `GetPath` | `func() string` | GetPath returns the raw request path, query string excluded ("/users/acme/create"). Slicing it into segments is the sandbox's business. |
| `GetHeader` | `func(key string) string` | GetHeader returns the first value of the named header, matched without regard to case, or "" when it is absent. |
| `GetQueryParam` | `func(name string) string` | GetQueryParam returns the first value of the named query parameter, or "" when it is absent. |
| `GetQueryAll` | `func(name string) []string` | GetQueryAll returns every value of the named query parameter, in the order they appear, for a field declared `array: true`. |
| `ReadBody` | `func(limit int) ([]byte, error)` | ReadBody reads at most limit bytes of the request body (-1 reads it whole) and reports an error when the body is longer than limit or cannot be read. The body is read once: a second call returns what the first one read. |
| `GetRemoteAddr` | `func() string` | GetRemoteAddr returns the address the request came from, in the host:port spelling. |

## `Response`

Response is the one http response being written, through function fields only. Headers and status are set before the first Write.

| Field | Type | Description |
| --- | --- | --- |
| `SetHeader` | `func(key string, value string)` | SetHeader sets one response header, replacing whatever value it had. |
| `SetStatus` | `func(code int)` | SetStatus writes the status line. It is called at most once, before any Write; without it the status is 200. |
| `Write` | `func(body []byte) error` | Write appends bytes to the response body. |

[every contract](doc.md)
