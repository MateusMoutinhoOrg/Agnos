package pageio

// The project's own render layer, sitting on top of the two 0-opinionated
// contracts it composes: sandbox.Deps.Embeddeps reads the asset, sandbox.Deps.Templatedeps
// executes it. What this package adds is the opinion neither contract may
// hold — a fixed set of native functions every page is rendered with, so a
// template names an asset by path and gets back a link that is checked at
// render time and busted on every content change.
//
// It is a leaf of sandbox/internal/: it imports the contracts and nothing
// else under internal/, so both the routes and anything else in the sandbox
// may reach for it. It is also what tells `build` this project carries the
// front layer, the way sandbox/internal/server tells it about the server one.

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/templatedeps"
)

const (
	// StaticRoot is the directory of the embedded asset tree every static
	// helper reads from, as sandbox.Deps.Embeddeps spells a path: slash-separated
	// and relative to the root of the assets package, so
	// "assets/frontend/static" on disk.
	StaticRoot = "frontend/static"
	// StaticMount is the url prefix the static route answers under. It is
	// rendered from the `identifier` of that route's first segment, so
	// renaming the mount in the declaration moves every link the helpers
	// build instead of breaking them in silence.
	StaticMount = {{printf "%q" .StaticMount}}
	// ShaLength is how much of the hex digest a link carries. Twelve hex
	// characters are 48 bits, far past what a cache key needs to stay
	// unique across one asset tree, and short enough to read in a url.
	ShaLength = 12
	// ImmutableCache is the cache header for a request whose sha matches the
	// asset it names. The url changes whenever the bytes do, so the answer
	// can be kept for a year and never revalidated.
	ImmutableCache = "public, max-age=31536000, immutable"
	// RevalidateCache is the cache header for a request with no sha, or a
	// sha that does not match. The url is not proof of its content, so the
	// browser may store the answer but must ask before reusing it.
	RevalidateCache = "no-cache"
)

// Render reads one template out of the embedded asset tree and executes it
// over vars with every helper of this package registered.
//
// path is spelled as sandbox.Deps.Embeddeps spells one — slash-separated and relative
// to the root of the assets package, so "frontend/pages/home.html". The error
// reports an asset that could not be read, a template that does not parse, or
// a helper that failed: all three are authoring or packaging mistakes, so a
// caller answers 500 rather than serving a half-rendered page.
func Render(sandbox *api.Sandbox, path string, vars any) ([]byte, error) {
	source, err := sandbox.Deps.Embeddeps.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rendered, err := sandbox.Deps.Templatedeps.Render(templatedeps.RenderProps{
		Name:   path,
		Source: string(source),
		Vars:   vars,
		Funcs:  funcMap(sandbox, vars),
	})
	if err != nil {
		return nil, err
	}

	return []byte(rendered), nil
}

// ShortSha is the digest a static link carries, and the one the static route
// compares an incoming `?sha=` against. Both sides truncate here rather than
// each on its own: a mismatch in length would send every request down the
// revalidate branch, which is a silent loss of caching rather than a failure.
func ShortSha(sandbox *api.Sandbox, content []byte) string {
	digest := sandbox.Deps.Hashdeps.Sha256Hex(content)
	if len(digest) < ShaLength {
		return digest
	}
	return digest[:ShaLength]
}

// funcMap is the set of native functions every Render registers. Each one
// returns an error rather than a placeholder, because the template engine
// aborts the whole execution on it: an asset path that no longer exists
// fails loudly at render time instead of shipping a page with a dead link.
//
// vars is threaded through so `include` renders a partial over the same
// values the page itself is rendered over.
func funcMap(sandbox *api.Sandbox, vars any) map[string]any {
	return map[string]any{
		"staticref": func(relative string) (string, error) {
			return staticRef(sandbox, relative)
		},
		"cssref": func(relative string) (string, error) {
			return cssRef(sandbox, relative)
		},
		"jsref": func(relative string) (string, error) {
			return jsRef(sandbox, relative)
		},
		"dirref": func(dir string) (string, error) {
			return dirRef(sandbox, dir)
		},
		"inline": func(relative string) (string, error) {
			return inline(sandbox, relative)
		},
		"include": func(path string) (string, error) {
			content, err := Render(sandbox, path, vars)
			if err != nil {
				return "", err
			}
			return string(content), nil
		},
	}
}

// staticRef is the `staticref` helper: the url of one asset under StaticRoot,
// stamped with the digest of the bytes it currently holds. An asset that
// cannot be read is reported rather than linked, which is what turns a typo in
// a page into a failing render instead of a 404 in the browser.
func staticRef(sandbox *api.Sandbox, relative string) (string, error) {
	content, err := readStatic(sandbox, relative)
	if err != nil {
		return "", err
	}
	return StaticMount + "/" + relative + "?sha=" + ShortSha(sandbox, content), nil
}

// cssRef is the `cssref` helper: the whole stylesheet tag for one asset.
func cssRef(sandbox *api.Sandbox, relative string) (string, error) {
	href, err := staticRef(sandbox, relative)
	if err != nil {
		return "", err
	}
	return `<link rel="stylesheet" href="` + href + `">`, nil
}

// jsRef is the `jsref` helper: the whole script tag for one asset. It is
// deferred, so a page may declare its scripts in the head and still parse
// before any of them runs.
func jsRef(sandbox *api.Sandbox, relative string) (string, error) {
	src, err := staticRef(sandbox, relative)
	if err != nil {
		return "", err
	}
	return `<script src="` + src + `" defer></script>`, nil
}

// dirRef is the `dirref` helper: every stylesheet and script at or below one
// directory of StaticRoot, as tags, in a fixed order — styles first so the
// page never paints unstyled, then scripts, each side sorted by path so the
// same tree always renders the same bytes. An extension neither side claims is
// skipped: an image or a font is referenced by the markup that uses it, not by
// a tag of its own. "" and "." both name StaticRoot itself.
func dirRef(sandbox *api.Sandbox, dir string) (string, error) {
	prefix, err := staticPrefix(sandbox, dir)
	if err != nil {
		return "", err
	}

	listed, err := sandbox.Deps.Embeddeps.ListFilesRecursively(StaticRoot + "/" + prefix)
	if err != nil {
		return "", err
	}
	sandbox.Deps.Sortdeps.Strings(listed)

	styles := []string{}
	scripts := []string{}
	for _, name := range listed {
		relative := prefix + name

		switch sandbox.Deps.Stringsdeps.ToLower(extensionOf(sandbox, name)) {
		case ".css":
			tag, err := cssRef(sandbox, relative)
			if err != nil {
				return "", err
			}
			styles = append(styles, tag)
		case ".js", ".mjs":
			tag, err := jsRef(sandbox, relative)
			if err != nil {
				return "", err
			}
			scripts = append(scripts, tag)
		}
	}

	return sandbox.Deps.Stringsdeps.Join(append(styles, scripts...), "\n"), nil
}

// inline is the `inline` helper: the content of one asset written straight
// into the page. It carries no digest and no tag, because nothing is fetched —
// it is for critical css, an icon, or a script small enough that a second
// request costs more than the bytes.
func inline(sandbox *api.Sandbox, relative string) (string, error) {
	content, err := readStatic(sandbox, relative)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// readStatic reads one asset under StaticRoot, refusing a path that could name
// something outside it.
func readStatic(sandbox *api.Sandbox, relative string) ([]byte, error) {
	if err := checkRelative(sandbox, relative); err != nil {
		return nil, err
	}
	return sandbox.Deps.Embeddeps.ReadFile(StaticRoot + "/" + relative)
}

// staticPrefix resolves the directory one helper was pointed at into what has
// to be put back in front of every name the listing returns, so a listed
// "app.js" becomes the "scripts/app.js" the other helpers take. "" and "."
// both name StaticRoot itself, whose prefix is empty.
func staticPrefix(sandbox *api.Sandbox, dir string) (string, error) {
	if dir == "" || dir == "." {
		return "", nil
	}
	if err := checkRelative(sandbox, dir); err != nil {
		return "", err
	}
	return sandbox.Deps.Stringsdeps.TrimSuffix(dir, "/") + "/", nil
}

// checkRelative refuses a path that does not stay under the directory it is
// resolved against. These paths are written by whoever wrote the template and
// never come from a request, so this is not the static route's defence against
// a caller — it is what keeps a helper unable to reach the rest of the asset
// tree by mistake.
func checkRelative(sandbox *api.Sandbox, relative string) error {
	if relative == "" {
		return sandbox.Deps.Std.Errorf("asset path is empty")
	}
	if sandbox.Deps.Stringsdeps.HasPrefix(relative, "/") {
		return sandbox.Deps.Std.Errorf("asset path %s is absolute", relative)
	}
	if sandbox.Deps.Stringsdeps.Contains(relative, "..") {
		return sandbox.Deps.Std.Errorf("asset path %s climbs out of %s", relative, StaticRoot)
	}
	if sandbox.Deps.Stringsdeps.ContainsAny(relative, "\\\x00") {
		return sandbox.Deps.Std.Errorf("asset path %s is not slash-separated", relative)
	}
	return nil
}

// extensionOf returns the asset's extension, dot included, or "" when it has
// none. It reads the last dot of the whole path, which is safe here because
// every path a helper resolves is slash-separated and already checked.
func extensionOf(sandbox *api.Sandbox, name string) string {
	cut := sandbox.Deps.Stringsdeps.LastIndex(name, ".")
	if cut < 0 {
		return ""
	}
	return name[cut:]
}
