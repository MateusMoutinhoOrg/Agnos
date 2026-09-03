package embeddeps

// This package is the sandbox's *copy* of the api an embedded-asset library
// exposes — the same mechanic as argvdeps and dbdeps, for the same reason:
// reading a file is an OS-bound effect, and compiling one into a binary needs
// the `embed` directive, so neither may appear inside the sandbox. The
// contract is restated here, and the adapter — which lives outside the
// sandbox — is what fills it.
//
// The library asks for an asset the way it would ask a filesystem, by path;
// where those bytes come from is the adapter's decision. The standard adapter
// serves them out of the assets compiled into the binary, so an installed
// `agnos` carries them wherever it runs; another adapter could read a
// directory on disk, an archive, or a network store without the sandbox
// noticing.
//
// The tracker in sandbox/ never reads an asset — its display text is small
// and fixed, so it lives in sandbox/config as compile-time constants. This
// contract is filled anyway, as a standing capability of the template: see
// the Deps.EmbedDeps field for why.

// Lib is the embedded-asset library injected whole as the Deps.EmbedDeps
// field. It is read-only by design: assets ship with the program, and nothing
// in the library ever writes one back.
//
// Every path is slash-separated and relative to the root of the asset tree
// the adapter serves — "report.tmpl", "templates/invoice.tmpl" — never an
// absolute path and never a path reaching outside that root, so the same call
// means the same asset whatever the adapter is backed by.
type Lib struct {
	// ReadFile returns the whole content of one asset. The error reports an
	// asset that does not exist or could not be read; callers inside the
	// sandbox report it rather than assuming the bytes are there, because a
	// missing asset is a packaging mistake, not a user mistake.
	ReadFile func(path string) ([]byte, error)

	// ListFiles returns the names of the assets directly inside the given
	// directory, in lexical order, relative to that directory. Nested
	// directories are not descended into and are not reported. The root
	// itself is addressed as ".".
	ListFiles func(path string) ([]string, error)

	// ListFilesRecursively returns every asset at or below the given
	// directory, in lexical order, as slash-separated paths relative to that
	// directory — "templates/invoice.tmpl" and not just "invoice.tmpl".
	// Directories are never reported, only the files inside them.
	ListFilesRecursively func(path string) ([]string, error)

	// RenderTemplate reads the template from path, renders it using the given variables,
	// and returns the resulting byte slice.
	RenderTemplate func(path string, vars interface{}) ([]byte, error)
}
