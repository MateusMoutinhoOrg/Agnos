package deps

import (
	argvdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/iodeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/requestdeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	std "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/std"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	Argvdeps argvdeps.Lib
	Embeddeps embeddeps.Lib
	Goimportsdeps goimportsdeps.Lib
	Iodeps iodeps.Lib
	Requestdeps requestdeps.Lib
	Rundeps rundeps.Lib
	Serializables serializables.Lib
	Std std.Lib
}
