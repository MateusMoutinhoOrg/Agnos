package deps

import (
	argvdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/argvdeps"
	dbdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/dbdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/goimportsdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/iodeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
	serializebles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/std"
)

type Deps struct {
	Argvdeps argvdeps.Lib
	Dbdeps dbdeps.Lib
	Embeddeps embeddeps.Lib
	Goimportsdeps goimportsdeps.Lib
	Iodeps iodeps.Lib
	Requestdeps requestdeps.Lib
	Serializebles serializebles.Lib
	Std std.Lib
}
