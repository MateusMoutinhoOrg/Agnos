package deps

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/embeddeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/iodeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/requestdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"
)

// Deps is the dependency contract every adapter must satisfy. It is a
// struct of function fields, not an interface: an adapter fills every
// field with the behavior it provides, and the library calls those fields
// directly.
//
// Each field is one injectable behavior. The tracker's own behavior is built
// from three of them — the clock (so a transaction's timestamp can be fixed
// in a test), the formatted writer the command-line interface reports
// through, and the schema database every category and transaction is
// persisted in.
//
// The contract is wider than this one tracker uses. Because the repository is
// a template, EmbedDeps, IoLib and NewRequest are carried here as standing
// capabilities: a derived library reads embedded assets, touches the
// filesystem, or speaks HTTP without designing a contract for it first, and
// the standard adapter already fills all three. The demonstration library in
// sandbox/ never calls them, which is the point — they are offered, not
// required. Every field is documented in PublicApi.md whether the
// demonstration exercises it or not.
//
// VerbLib, KeepLib, EmbedDeps, IoLib and NewRequest show what changes when
// the dependency is another library built with this same pattern: the whole
// library arrives as one struct field, with no getter method and no bridging
// type around it. The sandbox never imports Verb, Keep, the `embed`
// machinery, os, or net/http — it declares a copy of each api in verbdeps,
// keepdeps, embeddeps, iodeps and requestdeps, and the adapter, which lives
// outside the sandbox, fills it.
type Deps struct {
	// Now returns the current time, used to stamp categories and
	// transactions as they are created.
	Now func() time.Time
	// Printf writes one formatted message to the interface's output,
	// returning the number of bytes written and the write failure, if any.
	// It is the only way the library emits text: api.Lib.Sandboxmain — the
	// command-line interface — reports every result, every error, and its
	// usage screen through it, so the sandbox never touches a stream itself.
	Printf func(format string, a ...any) (n int, err error)
	// Error writes one formatted error message to the interface's error output.
	Error func(format string, a ...any) (n int, err error)
	// VerbLib is the embedded Verb argv-parser library, already initialized
	// by the adapter over the argument vector that adapter chose.
	VerbLib verbdeps.Lib
	// KeepLib is the embedded Keep schema-database library, already wired
	// by the adapter to the storage backend that adapter chose. Every
	// category and transaction the tracker stores lives in it.
	KeepLib keepdeps.Lib
	// EmbedDeps is the embedded-asset library, already rooted by the adapter
	// at the asset directory that adapter chose. It is how the sandbox reads
	// the files shipped under /assets/ — templates, long-form text, images,
	// any payload better kept as a file than as a Go constant — without
	// importing the `embed` machinery itself.
	//
	// The tracker does not use it: its display text is small and fixed, so it
	// lives in sandbox/config as compile-time constants instead. The field is
	// filled anyway, so a library derived from this template has the asset
	// mechanic ready. See docs/Tutorials/HandleAssets.md.
	EmbedDeps embeddeps.Lib
	// IoLib is the filesystem library, reaching whatever root the adapter
	// gave it: reading and writing files, creating directories, testing
	// existence, and listing a tree one level deep or all the way down.
	//
	// The tracker does not use it — every record it keeps is persisted
	// through KeepLib — but a derived library that must touch the filesystem
	// directly finds the contract already declared and already filled.
	IoLib iodeps.Lib
	// NewRequest opens an HTTP request against the given url, handing back
	// the request object the caller configures — method, headers, body — and
	// then sends with Fetch. It is a function field rather than a library
	// struct because a request is created per call, not injected once.
	//
	// The tracker does not use it: nothing it does leaves the machine. A
	// derived library that must speak HTTP has the contract ready, filled by
	// the standard adapter over net/http.
	NewRequest func(url string) requestdeps.Request
}
