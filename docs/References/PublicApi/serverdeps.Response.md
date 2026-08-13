# `serverdeps.Response`

**Type:** Struct

## Definition

```go
type Response struct {
	GetStatusCode func() int
	GetHeader     func(key string) string
	ReadBody      func(size int) ([]byte, error)
	Close         func() error
}
```

## Description

One HTTP response, declared in `sandbox/contracts/deps/serverdeps/` and handed back by [`serverdeps.Request.Fetch`](/docs/References/PublicApi/serverdeps.Request.md) with its body still open. Like the request it came from, it is the sandbox's **copy** of an HTTP client api: the sandbox may not import `net/http`, so the adapter fills these four fields — see [`standard.New`](/docs/References/PublicApi/standard.New.md).

An HTTP error **status** is not a transport error. `Fetch` returns no error for a `404` or a `500`; the status arrives through `GetStatusCode`, and the caller decides what it means.

The caller **must** `Close` every response returned without an error, whether or not the body is read, or the underlying connection leaks. The financial tracker never creates one — see [`serverdeps.Request`](/docs/References/PublicApi/serverdeps.Request.md) for why the contract ships anyway.

## Fields

| Field | Description |
| :--- | :--- |
| `GetStatusCode func() int` | Returns the response's HTTP status code. |
| `GetHeader func(key string) string` | Returns the first value of one response header, or `""` when the response carries no such header. |
| `ReadBody func(size int) ([]byte, error)` | Reads at most `size` bytes of the body, or the whole body when `size` is `-1`. A body shorter than `size` is returned whole rather than reported as an error, so a short read is not a failure. |
| `Close func() error` | Releases the response body. Must be called for every `Response` returned without an error. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	d := agnosadapter.New("trackerdata")

	response, err := d.NewRequest("https://example.com/rates").Fetch()
	if err != nil {
		fmt.Println("round trip failed:", err)
		return
	}
	defer response.Close()

	// A 404 is not an error from Fetch — it arrives as a status code.
	if response.GetStatusCode() != 200 {
		fmt.Println("unexpected status:", response.GetStatusCode())
		return
	}
	fmt.Println(response.GetHeader("Content-Type")) // application/json

	// A prefix read: 512 bytes, or the whole body when it is shorter.
	head, err := response.ReadBody(512)
	if err != nil {
		fmt.Println("could not read body:", err)
		return
	}
	fmt.Println(string(head))
}
```
