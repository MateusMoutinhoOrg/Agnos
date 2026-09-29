package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ConstructorsDir is the project-relative directory holding the sandbox's
// constructor packages. Every directory under it is one package whose
// Constructor(sandbox) fills a field of the Sandbox, and sandbox/new.go is
// nothing but a call to each in turn — which is what lets a project add a
// constructor of its own without a generated file standing in the way.
const ConstructorsDir = "sandbox/constructors"

// ConstructorFile is the file a constructor package declares its Constructor
// in, the one name sandbox/new.go's import of that package resolves against.
const ConstructorFile = "constructor.go"

// GeneratedDir is the project-relative directory holding every package the
// build rewrites whole — cli, config, routeio, frontio, databaseio and the
// server's route and server. Nothing under it is the project's to edit; the
// packages that mix a generated file with a hand-written one (a command, a
// route, a database) stay beside it under sandbox/internal/.
const GeneratedDir = "sandbox/internal/generated"

// RoutePropsFile is the sandbox/api/ file declaring api.RouteProps, the struct
// one request's chain of routes shares. The build writes it once and it is the
// project's from then on, like a route's InternalPureHandler.go.
const RoutePropsFile = "routeprops.go"

// CommandPropsFile is the sandbox/api/ file declaring api.CommandProps, the
// struct one command line's chain of commands shares. The build writes it once
// and it is the project's from then on.
const CommandPropsFile = "commandprops.go"

// UserSandboxFile is the sandbox/api/ file declaring api.UserSandbox, the
// part of the Sandbox the project types itself. start writes it once and no
// build rewrites it; api.Sandbox embeds it.
const UserSandboxFile = "usersandbox.go"

// UserConfigFile is the sandbox/api/ file declaring api.UserConfig, the part
// of the Config the project types itself. start writes it once and no build
// rewrites it; api.Config embeds it.
const UserConfigFile = "userconfig.go"

// ConstructorExempt are the sandbox/api/ files that declare no field of the
// sandbox: sandbox.go is the struct itself, command.go and route.go are the
// shape of one command and of one route, each owned by the contract whose
// New<Name> builds the slice of them, trigger.go is the condition both of
// them match on, routeprops.go and commandprops.go are what a route's and a
// command's handler are handed per run, and usersandbox.go and userconfig.go
// are embedded in the Sandbox and the Config rather than fields of their own.
var ConstructorExempt = []string{"sandbox.go", "command.go", "route.go", "trigger.go", RoutePropsFile, CommandPropsFile, UserSandboxFile, UserConfigFile}

// ConstructorDir is the project-relative directory of one constructor package.
func ConstructorDir(name string) string {
	return ConstructorsDir + "/" + name
}

// ConstructorPath is the project-relative path of one constructor package's
// constructor.go.
func ConstructorPath(name string) string {
	return ConstructorDir(name) + "/" + ConstructorFile
}

// ConstructorSource is the project-relative package whose New<Name> builds one
// contract of sandbox/api/. A generated layer lives under GeneratedDir —
// sandbox/internal/generated/<name>, or sandbox/internal/generated/<name>/<name>
// when the layer splits its package into several (the server's route and
// server) and keeps the one that builds the contract under its own name one
// level down. A contract the project writes itself lives at
// sandbox/internal/<name>, with the same one-level-down fallback.
func ConstructorSource(io *smartio.SmartIO, name string) string {
	for _, root := range []string{GeneratedDir, "sandbox/internal"} {
		nested := root + "/" + name + "/" + name
		if io.IsFile(nested + "/new.go") {
			return nested
		}
		if io.IsFile(root + "/" + name + "/new.go") {
			return root + "/" + name
		}
	}
	return "sandbox/internal/" + name
}
