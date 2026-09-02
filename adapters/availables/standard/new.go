package standard

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	keep "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/keep"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/requestdeps"
	serializebles "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	verb "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/verb"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	embeddeps.Bind(&deps)
	iodeps.Bind(&deps)
	keep.Bind(&deps)
	requestdeps.Bind(&deps)
	serializebles.Bind(&deps)
	std.Bind(&deps)
	verb.Bind(&deps)
	return deps
}
