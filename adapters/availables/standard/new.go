package standard

import (
	argvdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/argvdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/goimportsdeps"
	hashdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/hashdeps"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/interviewer"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/iodeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/serializables"
	serverdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/serverdeps"
	sortdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/sortdeps"
	std "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/std"
	stringsdeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/stringsdeps"
	templatedeps "github.com/MateusMoutinhoOrg/Agnos/adapters/libs/templatedeps"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	argvdeps.Bind(&deps)
	embeddeps.Bind(&deps)
	goimportsdeps.Bind(&deps)
	hashdeps.Bind(&deps)
	interviewer.Bind(&deps)
	iodeps.Bind(&deps)
	rundeps.Bind(&deps)
	serializables.Bind(&deps)
	serverdeps.Bind(&deps)
	sortdeps.Bind(&deps)
	std.Bind(&deps)
	stringsdeps.Bind(&deps)
	templatedeps.Bind(&deps)
	return deps
}
