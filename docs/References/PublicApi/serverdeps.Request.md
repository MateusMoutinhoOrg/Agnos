# `serverdeps.Request`

**Type:** Struct

## Definition

```go
type Request struct {
	AddHeader func(key string, value string)
	SetMethod func(method string)
	SetBody   func(body []byte)
	Fetch     func() (Response, error)
}
```

## Description

One HTTP request under construction, declared in `sandbox/contracts/deps/serverdeps/` and handed back by the [`deps.Deps.NewRequest`](/docs/References/PublicApi/deps.Deps.md) field already bound to a url. It is the sandbox's **copy** of the api an HTTP client library exposes, for the same reason [`iodeps.Lib`](/docs/References/PublicApi/iodeps.Lib.md) exists: opening a socket is an OS-bound effect, so `net/http` may not appear inside the sandbox. The adapter, which lives outside it, fills every field — see [`standard.New`](/docs/References/PublicApi/standard.New.md).

It differs from the other injected libraries in one way: a request is created **per call** rather than injected once, so the entry point is a function field on `Deps` rather than a library struct. Everything here is what that function hands back.

The financial tracker **does not call it** — nothing it does leaves the machine. It is carried as a standing capability of the template, so a derived library that must speak HTTP finds the contract already declared and already wired.

The setters mutate the pending request and may be called in any order; nothing leaves the machine until `Fetch` is called, and a `Request` may be sent more than once. The method defaults to `GET` and the body to none, so a plain read is `NewRequest` followed by `Fetch`.

## Fields

| Field | Description |
| :--- | :--- |
| `AddHeader func(key string, value string)` | Sets one header on the pending request, replacing whatever value that key carried before. |
| `SetMethod func(method string)` | Sets the HTTP method the request is sent with — `"POST"`, `"PUT"`, `"DELETE"`. Defaults to `"GET"`. |
| `SetBody func(body []byte)` | Sets the bytes sent as the request body. Defaults to none. |
| `Fetch func() (Response, error)` | Sends the request and returns the [`Response`](/docs/References/PublicApi/serverdeps.Response.md). The error reports a request that could not be built or a round trip that failed; an HTTP error **status** is not an error here. A `Response` returned without an error holds an open body the caller must `Close`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	// The adapter fills NewRequest with net/http calls, bounding every round
	// trip with a timeout the sandbox cannot set for itself.
	d := agnosadapter.New("trackerdata")

	request := d.NewRequest("https://example.com/rates")
	request.SetMethod("POST")
	request.AddHeader("Content-Type", "application/json")
	request.SetBody([]byte(`{"base":"BRL"}`))

	response, err := request.Fetch()
	if err != nil {
		fmt.Println("round trip failed:", err)
		return
	}
	// Close is mandatory for every Response returned without an error.
	defer response.Close()

	fmt.Println(response.GetStatusCode()) // 200

	body, err := response.ReadBody(-1) // -1 reads the whole body
	if err != nil {
		fmt.Println("could not read body:", err)
		return
	}
	fmt.Println(len(body))
}
```
