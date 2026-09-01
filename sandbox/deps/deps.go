package deps

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/iodeps"
	keepdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/keepdeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/std"
	verbdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/verbdeps"
)

type Deps struct {
	Std std.Lib

	EmbedDeps embeddeps.Lib

	IoLib iodeps.Lib

	SerializeLib serializibles.Lib

	VerbLib verbdeps.Lib

	KeepLib keepdeps.Lib

	RequestLib requestdeps.Lib
}
