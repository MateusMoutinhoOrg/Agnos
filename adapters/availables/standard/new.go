package standard

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/goimportsdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/iodeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/serializables"
	std "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/std"
	verb "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/verb"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	embeddeps.Bind(&deps)
	goimportsdeps.Bind(&deps)
	iodeps.Bind(&deps)
	rundeps.Bind(&deps)
	serializables.Bind(&deps)
	std.Bind(&deps)
	verb.Bind(&deps)
	return deps
}
