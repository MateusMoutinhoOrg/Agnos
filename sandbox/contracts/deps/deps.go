package deps

import (
	"time"

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
// transaction's timestamp can be fixed in a test) and the schema database
// every category and transaction is persisted in.
//
// VerbLib and KeepLib show what changes when the dependency is another
// library built with this same pattern: the whole library arrives as one
// struct field, with no getter method and no bridging type around it. The
// sandbox never imports Verb or Keep — it declares a copy of each api in
// verbdeps and keepdeps, and the adapter, which lives outside the sandbox,
// fills it.
type Deps struct {
	// Now returns the current time, used to stamp categories and
	// transactions as they are created.
	Now func() time.Time
	// VerbLib is the embedded Verb argv-parser library, already initialized
	// by the adapter over the argument vector that adapter chose.
	VerbLib verbdeps.Lib
	// KeepLib is the embedded Keep schema-database library, already wired
	// by the adapter to the storage backend that adapter chose. Every
	// category and transaction the tracker stores lives in it.
	KeepLib keepdeps.Lib
}
