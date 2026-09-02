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
	deps.ArgvLib.New = argvlib.NewArgvLib
	deps.DatabaseLib = dblib.NewDatabaseLib()
	deps.EmbedDeps = embeddeps.NewEmbedDeps()
	deps.IoLib = iolib.NewIoLib()
	deps.RequestLib.NewRequest = requestlib.NewRequestFactory()
	deps.SerializeLib = serializelib.NewSerializeLib()
	return deps
}
