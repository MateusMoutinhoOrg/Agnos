package standard

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	keep "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/keep"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/requestdeps"
	serializes "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/serializebles"
	std "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	verb "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/verb"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	std.BindStdLib(&deps)
	verb.BindArgvLib(&deps)
	keep.BindDatabaseLib(&deps)
	embeddeps.BindEmbedDeps(&deps)
	iodeps.BindIoLib(&deps)
	requestdeps.BindRequestLib(&deps)
	serializes.BindSerializeLib(&deps)
	return deps
}
