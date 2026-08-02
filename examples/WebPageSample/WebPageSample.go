package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/api"
)

// This sample caches a web page for 60 seconds: repeated reads inside the
// window are served from the cache, and only after the TTL expires does it
// hit the network again. The target returns a fresh UUID on every request,
// so a cache HIT visibly serves the same body until it expires — swap it
// for any mutable endpoint you like.
const targetURL = "https://httpbin.org/uuid"

// fetchCached serves url from the cache when a live entry exists and only
// performs the HTTP request on a miss, then caches the body for 60 seconds.
func fetchCached(l agnostypes.Lib, url string) (body string, cached bool, err error) {
	if entry, ok := l.Get(url); ok {
		return entry.Value, true, nil // cache HIT — no network call
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	page := string(raw)
	l.Set(url, page, 60) // cache for 60 seconds
	return page, false, nil
}

func main() {
	l := agnoslib.New(agnosadapter.New("webpage.json"))

	// Poll every 25s. Within the 60s TTL the body is served from cache;
	// once it expires, the next poll fetches a fresh (different) page.
	for i := 0; i < 6; i++ {
		body, cached, err := fetchCached(l, targetURL)
		if err != nil {
			fmt.Println("fetch error:", err)
			return
		}

		source := "MISS (fetched)"
		if cached {
			source = "HIT  (cached) "
		}
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), source, strings.TrimSpace(body))

		time.Sleep(25 * time.Second)
	}
}
