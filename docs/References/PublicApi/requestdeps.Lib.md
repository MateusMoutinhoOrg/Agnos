# `requestdeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	NewRequest func(url string) Request
}

type Request struct {
	AddHeader func(key string, value string)
	SetMethod func(method string)
	SetBody   func(body []byte)
	Fetch     func() (Response, error)
}

type Response struct {
	GetStatusCode func() int
	GetHeader     func(key string) string
	ReadBody      func(size int) ([]byte, error)
	Close         func() error
}
```

## Description

The HTTP-request constructor injected whole as `Deps.Requestdeps` — the sandbox's copy of what `net/http` provides. A request is created per call rather than injected once, so the sandbox holds the one-field `Lib`, and everything below it is what `NewRequest` hands back. The setters mutate the pending request in any order; nothing leaves the machine until `Fetch`, and a `Request` may be sent more than once. The method defaults to `GET` and the body to none, so a plain read is `NewRequest` followed by `Fetch`. An HTTP error status is **not** an error from `Fetch`; it is reported by `Response.GetStatusCode`. A `Response` returned without an error holds an open body the caller must `Close`.

The standard adapter bounds every round trip with a 30-second timeout, because the contract exposes no cancellation and the sandbox could not intervene otherwise. Wired by the `requestdeps` dep; **no action uses it yet** — a standing capability.

## Fields

| Field | Description |
| :--- | :--- |
| `NewRequest(url)` | A request bound to `url`. |
| `Request.AddHeader(key, value)` | Sets one header, replacing a previous value. |
| `Request.SetMethod(method)` | `"POST"`, `"PUT"`, … Defaults to `GET`. |
| `Request.SetBody(body)` | The request body. Defaults to none. |
| `Request.Fetch()` | Sends; errors only when the request could not be built or the round trip failed. |
| `Response.GetStatusCode()`, `GetHeader(key)` | Status and the first value of one header (`""` when absent). |
| `Response.ReadBody(size)` | At most `size` bytes, or the whole body for `-1`. A short read is not a failure. |
| `Response.Close()` | Releases the body. Required for every `Response` returned without error. |

## Examples

```go
request := deps.Requestdeps.NewRequest("https://example.com/releases/latest")
request.AddHeader("Accept", "application/json")

response, err := request.Fetch()
if err != nil {
	return err
}
defer response.Close()

if response.GetStatusCode() != 200 {
	return deps.Std.Errorf("unexpected status %d", response.GetStatusCode())
}
body, err := response.ReadBody(-1)
```
