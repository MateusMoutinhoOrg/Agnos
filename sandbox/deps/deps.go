package deps

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/iodeps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/std"
)

type Deps struct {
	Std std.Lib

	EmbedDeps embeddeps.Lib

	IoLib iodeps.Lib

	SerializeLib serializibles.Lib
}
