package deps

import (
	argvdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/argvdeps"
	dbdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/dbdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/iodeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/std"
)

type Deps struct {
	Std std.Lib

	EmbedDeps embeddeps.Lib

	IoLib iodeps.Lib

	SerializeLib serializibles.Lib

	ArgvLib argvdeps.Lib

	DatabaseLib dbdeps.Lib

	RequestLib requestdeps.Lib
}
