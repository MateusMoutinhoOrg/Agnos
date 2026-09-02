package standard

import (
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/embeddeps"
	iolib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	dblib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/keep"
	requestlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/requestdeps"
	serializelib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/serializebles"
	stdlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	argvlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/verb"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	stdlib.BindStdLib(&deps)
	argvlib.BindArgvLib(&deps)
	dblib.BindDatabaseLib(&deps)
	embeddeps.BindEmbedDeps(&deps)
	iolib.BindIoLib(&deps)
	requestlib.BindRequestLib(&deps)
	serializelib.BindSerializeLib(&deps)
	return deps
}
