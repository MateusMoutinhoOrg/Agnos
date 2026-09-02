package goimportsdeps

// This package is the sandbox's *copy* of the api a Go-source import parser
// exposes — the same mechanic as argvdeps, dbdeps, embeddeps, iodeps,
// requestdeps and std, for the same reason: parsing Go source pulls in
// `go/parser` / `go/token`, which are OS-adjacent enough to stay outside the
// sandbox. The contract is restated here, and the adapter — which lives
// outside the sandbox — is what fills it.

// Lib is the Go-source import parser injected whole as the Deps.Goimportsdeps
// field.
type Lib struct {
	// GetImports parses the given Go source file content and returns every
	// import path it declares, in source order. It errors when the content
	// is not parsable Go.
	GetImports func(content string) ([]string, error)
}
