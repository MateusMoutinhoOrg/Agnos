package deps

import "github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps/agnosdeps"

// Deps is the dependency contract every bootstrap adapter must satisfy.
// Like every contract in this template it is a struct of fields, not an
// interface.
//
// TrackerLib shows what changes when the dependency is another library
// built with this same pattern: the whole library arrives as one struct
// field, with no getter method and no bridging type around it.
type Deps struct {
	// Println writes a line to the library's output.
	Println func(a ...any)
	// TrackerLib is the embedded Agnos-Cli financial-tracker library, already
	// initialized by the adapter.
	TrackerLib agnosdeps.Lib
}
