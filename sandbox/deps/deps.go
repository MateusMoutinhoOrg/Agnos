package deps

import (
	argvdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	hashdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/hashdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/iodeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	serverdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serverdeps"
	sortdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/sortdeps"
	std "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/std"
	stringsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stringsdeps"
	templatedeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/templatedeps"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	Argvdeps      argvdeps.Sandbox
	Embeddeps     embeddeps.Sandbox
	Goimportsdeps goimportsdeps.Sandbox
	Hashdeps      hashdeps.Sandbox
	Iodeps        iodeps.Sandbox
	Rundeps       rundeps.Sandbox
	Serializables serializables.Sandbox
	Serverdeps    serverdeps.Sandbox
	Sortdeps      sortdeps.Sandbox
	Std           std.Sandbox
	Stringsdeps   stringsdeps.Sandbox
	Templatedeps  templatedeps.Sandbox
}
