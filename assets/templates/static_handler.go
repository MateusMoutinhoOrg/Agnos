package static

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/pageio"
	"{{.Module}}/sandbox/internal/routeio"
)

// RouteHandler answers GET /static/<item...> with one file of pageio.StaticRoot,
// served straight out of the binary through sandbox.Deps.Embeddeps.
//
// The captured segments are attacker-controlled and the path they build is
// cleaned by the embed adapter before the read, so "/static/../asset.go" would
// otherwise reach a file outside pageio.StaticRoot. safeSegments is what forbids
// it: nothing leaves this route that the caller did not name segment by segment
// underneath that directory. This file is yours to edit, and that check is the
// one thing to keep whatever else changes.
//
// `?sha=` is what the pageio helpers stamp on every link they build. It is a
// cache key and never an input: the file that is served is the one the path
// names, whatever the sha says. All it decides is cacheHeader below.
func RouteHandler(sandbox *api.Sandbox, entries *Entries, response serverdeps.Response) int {
	relative, ok := safeSegments(sandbox, entries.Item)
	if !ok {
		return routeio.WriteError(sandbox, response, api.StatusBadRequest, "item",
			"invalid static asset path")
	}

	content, err := sandbox.Deps.Embeddeps.ReadFile(pageio.StaticRoot + "/" + relative)
	if err != nil {
		return routeio.WriteError(sandbox, response, api.StatusNotFound, "item",
			"static asset not found")
	}

	response.SetHeader("Content-Type", contentTypeOf(sandbox, relative))
	response.SetHeader("Cache-Control", cacheHeader(sandbox, entries.Sha, content))
	response.SetStatus(api.StatusOk)
	response.Write(content)

	return api.StatusOk
}

// cacheHeader reports how long the answer may be reused. A request carrying the
// digest of the very bytes being sent came from a link this build generated, so
// that url can never mean anything else and the answer is immutable. A request
// with no sha, or with one from an older build, gets an answer the browser has
// to revalidate — it is still served the current file, so a stale link is slow
// rather than broken.
func cacheHeader(sandbox *api.Sandbox, sha string, content []byte) string {
	if sha != "" && sha == pageio.ShortSha(sandbox, content) {
		return pageio.ImmutableCache
	}
	return pageio.RevalidateCache
}

// safeSegments joins the captured segments into the path relative to
// pageio.StaticRoot, and reports false on anything that could climb out of it. A
// segment names one entry of the tree and nothing else: "." and ".." are the
// two spellings that move rather than name, and a separator inside a segment
// means the request path was percent-encoded to hide one of them from the
// dispatch's split. An empty segment cannot reach here — the dispatch drops
// those — and is refused anyway rather than trusted.
func safeSegments(sandbox *api.Sandbox, segments []string) (string, bool) {
	if len(segments) == 0 {
		return "", false
	}
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", false
		}
		if sandbox.Deps.Stringsdeps.ContainsAny(segment, "/\\\x00") {
			return "", false
		}
	}
	return sandbox.Deps.Stringsdeps.Join(segments, "/"), true
}

// contentTypeOf reads the media type off the asset's extension, matched in
// lower case. An extension nothing below claims is served as opaque bytes
// rather than guessed at, which is also what keeps an unknown asset from being
// rendered as html by the browser.
func contentTypeOf(sandbox *api.Sandbox, relative string) string {
	cut := sandbox.Deps.Stringsdeps.LastIndex(relative, ".")
	if cut < 0 {
		return "application/octet-stream"
	}

	switch sandbox.Deps.Stringsdeps.ToLower(relative[cut:]) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".json", ".map":
		return "application/json"
	case ".txt":
		return "text/plain; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".otf":
		return "font/otf"
	case ".wasm":
		return "application/wasm"
	case ".pdf":
		return "application/pdf"
	}

	return "application/octet-stream"
}
