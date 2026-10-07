package opinatedagnosfront

import (
	"strings"

	"{{.Module}}/sandbox/deps"
	opinatedagnosfront "{{.Module}}/sandbox/deps/OpinatedAgnosFront"
	"{{.Module}}/sandbox/deps/embeddeps"
)

// Bind fills deps.Deps.OpinatedAgnosFront with the agnos front file layer.
// Nothing here holds a dep: Resolve is handed the embedded tree it reads.
func Bind(deps *deps.Deps) {
	deps.OpinatedAgnosFront = opinatedagnosfront.Sandbox{
		Resolve:       resolve,
		SafePath:      safePath,
		ExtensionOf:   extensionOf,
		ContentTypeOf: contentTypeOf,
	}
}

// resolve reads the file one request path names under Root, trying in order
// the path itself, the path plus ".html", and the path as a directory holding
// Index — so /about answers about.html or about/index.html, and / answers
// index.html. It returns the path relative to Root it read, for contentTypeOf,
// and false when the path is unsafe or names nothing, which a caller reads as
// "not mine".
func resolve(embedded embeddeps.Sandbox, requested string) (string, []byte, bool) {
	relative, ok := safePath(requested)
	if !ok {
		return "", nil, false
	}

	candidates := []string{opinatedagnosfront.Index}
	if relative != "" {
		candidates = []string{relative, relative + ".html", relative + "/" + opinatedagnosfront.Index}
	}

	for _, candidate := range candidates {
		content, err := embedded.ReadFile(opinatedagnosfront.Root + "/" + candidate)
		if err == nil {
			return candidate, content, true
		}
	}
	return "", nil, false
}

// safePath turns a request path into the path relative to Root it names, and
// reports false on anything that could climb out of it. The path is
// attacker-controlled, and the embed adapter cleans what it is handed, so
// "/../asset.go" would otherwise reach a file outside Root. A segment names
// one entry of the tree and nothing else: "." and ".." are the two spellings
// that move rather than name, and a backslash or a NUL inside one means the
// request was encoded to hide something from the dispatch's split. A leading
// and a trailing slash are dropped; "" names Root itself.
func safePath(requested string) (string, bool) {
	trimmed := strings.Trim(requested, "/")
	if trimmed == "" {
		return "", true
	}

	segments := strings.Split(trimmed, "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", false
		}
		if strings.ContainsAny(segment, "\\\x00") {
			return "", false
		}
	}
	return strings.Join(segments, "/"), true
}

// extensionOf returns the extension of the last segment of a path, dot
// included, or "" when it has none.
func extensionOf(path string) string {
	cut := strings.LastIndex(path, ".")
	if cut < 0 || cut < strings.LastIndex(path, "/") {
		return ""
	}
	return path[cut:]
}

// contentTypeOf reads the media type off a file's extension, matched in lower
// case. An extension nothing below claims is served as opaque bytes rather
// than guessed at, which is also what keeps an unknown file from being
// rendered as html by the browser.
func contentTypeOf(relative string) string {
	switch strings.ToLower(extensionOf(relative)) {
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
