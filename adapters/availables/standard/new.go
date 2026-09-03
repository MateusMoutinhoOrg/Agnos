package standard

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/goimportsdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/requestdeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/serializables"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	verb "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/verb"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	embeddeps.Bind(&deps)
	goimportsdeps.Bind(&deps)
	iodeps.Bind(&deps)
	requestdeps.Bind(&deps)
	rundeps.Bind(&deps)
	serializables.Bind(&deps)
	std.Bind(&deps)
	verb.Bind(&deps)
	return deps
}
