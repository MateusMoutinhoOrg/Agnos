package standard

import (
	argvlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/argvdeps"
	dblib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/dbdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/embeddeps"
	iolib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/iodeps"
	requestlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/requestdeps"
	serializelib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/serializebles"
	stdlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/libs/std"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	deps.Std = stdlib.NewStdLib()
	deps.ArgvLib.New = argvlib.NewArgvLib
	deps.DatabaseLib = dblib.NewDatabaseLib()
	deps.EmbedDeps = embeddeps.NewEmbedDeps()
	deps.IoLib = iolib.NewIoLib()
	deps.RequestLib.NewRequest = requestlib.NewRequestFactory()
	deps.SerializeLib = serializelib.NewSerializeLib()
	return deps
}
