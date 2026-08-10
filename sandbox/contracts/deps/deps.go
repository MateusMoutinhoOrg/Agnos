package deps

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/embeddeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"
)

// Deps is the dependency contract every adapter must satisfy. It is a
// struct of function fields, not an interface: an adapter fills every
// field with the behavior it provides, and the library calls those fields
// directly.
//
// Each field is one injectable behavior the library needs. For this
// financial tracker the injectable behaviors are the clock (so a
// transaction's timestamp can be fixed in a test), the formatted writer the
// command-line interface reports through, and the schema database every
// category and transaction is persisted in.
//
// VerbLib, KeepLib and EmbedDeps show what changes when the dependency is
// another library built with this same pattern: the whole library arrives as
// one struct field, with no getter method and no bridging type around it. The
// sandbox never imports Verb, Keep, or the `embed` machinery — it declares a
// copy of each api in verbdeps, keepdeps and embeddeps, and the adapter,
// which lives outside the sandbox, fills it.
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
	// VerbLib is the embedded Verb argv-parser library, already initialized
	// by the adapter over the argument vector that adapter chose.
	VerbLib verbdeps.Lib
	// KeepLib is the embedded Keep schema-database library, already wired
	// by the adapter to the storage backend that adapter chose. Every
	// category and transaction the tracker stores lives in it.
	KeepLib keepdeps.Lib
	// EmbedDeps is the embedded-asset library, already rooted by the adapter
	// at the asset directory that adapter chose. It is where every piece of
	// standing text the library displays comes from — the usage screen, the
	// version, and each message the command-line interface prints — so no
	// display text is written in the sandbox itself. Only Sandboxmain reads
	// it: a program that calls the library functions directly never touches
	// an asset.
	EmbedDeps embeddeps.Lib
}
