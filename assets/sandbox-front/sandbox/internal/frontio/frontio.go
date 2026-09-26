package frontio

// The project's file layer for the front: it maps a request path onto one file
// of the embedded assets/frontend/ tree and names the media type that file is
// served as. It holds no opinion about what that tree is — hand-written html,
// or the dist/ of any bundler — so whatever lands in assets/frontend/ is
// served as it is.
//
// It is a leaf of sandbox/internal/: it imports the contracts and nothing
// else under internal/, so the frontend route and anything else in the
// sandbox may reach for it. It is also what tells `build` this project carries
// the front layer, the way sandbox/internal/server tells it about the server
// one. The route that serves the tree is the project's; this package, and the
// path check in it, is rewritten by every build.

import (
	"{{.Module}}/sandbox/api"
)

const (
	// Root is the directory of the embedded asset tree every file is read
	// from, as sandbox.Deps.Embeddeps spells a path: slash-separated and
	// relative to the root of the assets package, so "assets/frontend" on disk.
	Root = "frontend"
	// Index is the file a directory is answered with, "/" included.
	Index = "index.html"
	// NotFound is the file a path naming no file is answered with, under a
	// 404. front-init writes it once; it is the project's to restyle.
	NotFound = "404.html"
	// RevalidateCache is the cache header every file is served with: the
	// browser may store the answer but must ask before reusing it, so a new
	// build is picked up on the next load whatever the file is named.
	RevalidateCache = "no-cache"
)

// Resolve reads the file one request path names under Root, trying in order
// the path itself, the path plus ".html", and the path as a directory holding
// Index — so /about answers about.html or about/index.html, and / answers
// index.html. It returns the path relative to Root it read, for ContentTypeOf,
// and false when the path is unsafe or names nothing, which a caller reads as
// "not mine".
func Resolve(sandbox *api.Sandbox, requested string) (string, []byte, bool) {
	relative, ok := SafePath(sandbox, requested)
	if !ok {
		return "", nil, false
	}

	candidates := []string{Index}
	if relative != "" {
		candidates = []string{relative, relative + ".html", relative + "/" + Index}
	}

	for _, candidate := range candidates {
		content, err := sandbox.Deps.Embeddeps.ReadFile(Root + "/" + candidate)
		if err == nil {
			return candidate, content, true
		}
	}
	return "", nil, false
}

// SafePath turns a request path into the path relative to Root it names, and
// reports false on anything that could climb out of it. The path is
// attacker-controlled, and the embed adapter cleans what it is handed, so
// "/../asset.go" would otherwise reach a file outside Root. A segment names
// one entry of the tree and nothing else: "." and ".." are the two spellings
// that move rather than name, and a backslash or a NUL inside one means the
// request was encoded to hide something from the dispatch's split. A leading
// and a trailing slash are dropped; "" names Root itself.
func SafePath(sandbox *api.Sandbox, requested string) (string, bool) {
	trimmed := sandbox.Deps.Stringsdeps.Trim(requested, "/")
	if trimmed == "" {
		return "", true
	}

	segments := sandbox.Deps.Stringsdeps.Split(trimmed, "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", false
		}
		if sandbox.Deps.Stringsdeps.ContainsAny(segment, "\\\x00") {
			return "", false
		}
	}
	return sandbox.Deps.Stringsdeps.Join(segments, "/"), true
}

// ExtensionOf returns the extension of the last segment of a path, dot
// included, or "" when it has none.
func ExtensionOf(sandbox *api.Sandbox, path string) string {
	cut := sandbox.Deps.Stringsdeps.LastIndex(path, ".")
	if cut < 0 || cut < sandbox.Deps.Stringsdeps.LastIndex(path, "/") {
		return ""
	}
	return path[cut:]
}

// ContentTypeOf reads the media type off a file's extension, matched in lower
// case. An extension nothing below claims is served as opaque bytes rather
// than guessed at, which is also what keeps an unknown file from being
// rendered as html by the browser.
func ContentTypeOf(sandbox *api.Sandbox, relative string) string {
	switch sandbox.Deps.Stringsdeps.ToLower(ExtensionOf(sandbox, relative)) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".json", ".map":
		return "application/json"
	case ".webmanifest":
		return "application/manifest+json"
	case ".txt":
		return "text/plain; charset=utf-8"
	case ".xml":
		return "application/xml"
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
	case ".avif":
		return "image/avif"
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
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	}

	return "application/octet-stream"
}
