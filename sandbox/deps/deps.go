package deps

import (
	OpinionatedAgnosCli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinionatedAgnosCli"
	argvdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
	embeddeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/embeddeps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	hashdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/hashdeps"
	interviewdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewdeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/iodeps"
	reflectdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/reflectdeps"
	rundeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	serializabledeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializabledeps"
	serverdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serverdeps"
	sortdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/sortdeps"
	stddeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stddeps"
	stringsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stringsdeps"
	templatedeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/templatedeps"
)

// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
	OpinionatedAgnosCli OpinionatedAgnosCli.Contract
	ArgvDeps            argvdeps.Contract
	EmbedDeps           embeddeps.Contract
	GoimportsDeps       goimportsdeps.Contract
	HashDeps            hashdeps.Contract
	InterviewDeps       interviewdeps.Contract
	IoDeps              iodeps.Contract
	ReflectDeps         reflectdeps.Contract
	RunDeps             rundeps.Contract
	SerializableDeps    serializabledeps.Contract
	ServerDeps          serverdeps.Contract
	SortDeps            sortdeps.Contract
	StdDeps             stddeps.Contract
	StringsDeps         stringsdeps.Contract
	TemplateDeps        templatedeps.Contract
}
