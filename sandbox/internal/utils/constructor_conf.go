package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ConstructorsDir is the project-relative directory holding the sandbox's
// constructor packages. Every directory under it is one package whose
// Constructor(sandbox) fills a field of the Sandbox, and sandbox/generated.new.go is
// nothing but a call to each in turn — which is what lets a project add a
// constructor of its own without a generated file standing in the way.
const ConstructorsDir = "sandbox/constructors"

// ConstructorFile is the file a constructor package declares its Constructor
// in, the one name sandbox/generated.new.go's import of that package resolves against.
const ConstructorFile = "constructor.go"

// RoutePropsDir is the package declaring routeprops.RouteProps, the struct one
// request's chain of routes shares. It sits under sandbox/internal, not in
// sandbox/api, so its fields may name any type of the project — a record of
// one of its databases, say — which sandbox/api, imported by all of them,
// never could.
const RoutePropsDir = "sandbox/internal/routeprops"

// RoutePropsFile is the file of RoutePropsDir declaring RouteProps, before
// GeneratedPrefix: the build rewrites generated.routeprops.go every time, as
// the embedding of every other file's struct — the project's own project.go
// among them. routeprops.go is also the name it had in sandbox/api/, the old
// home the build moves it out of.
const RoutePropsFile = "routeprops.go"

// CommandPropsDir is the package declaring commandprops.CommandProps, the
// struct one command line's chain of commands shares — under sandbox/internal
// for the reason RoutePropsDir is.
const CommandPropsDir = "sandbox/internal/commandprops"

// CommandPropsFile is the file of CommandPropsDir declaring CommandProps, before
// GeneratedPrefix: the build rewrites generated.commandprops.go every time, the
// way it does RoutePropsFile. commandprops.go is also the name it had in
// sandbox/api/, the old home the build moves it out of.
const CommandPropsFile = "commandprops.go"

// ProjectSandboxFile is the sandbox/api/ file declaring api.ProjectSandbox, the
// part of the Sandbox the project types itself. start writes it once and no
// build rewrites it; api.Sandbox embeds it.
const ProjectSandboxFile = "projectsandbox.go"

// ProjectConfigFile is the sandbox/api/ file declaring api.ProjectConfig, the part
// of the Config the project types itself. start writes it once and no build
// rewrites it; api.Config embeds it.
const ProjectConfigFile = "projectconfig.go"

// SandboxPartSuffix ends the name of every sandbox/api/ file whose structs
// api.Sandbox embeds: projectsandbox.go, the project's, and one more per mechanic
// that adds a part of its own (clisandbox.go, serversandbox.go), so none edits
// a file of another. sandbox.go, the aggregate itself, is not one. A name is
// read with GeneratedPrefix taken off: generated.clisandbox.go is a part.
const SandboxPartSuffix = "sandbox.go"

// ConfigPartSuffix is SandboxPartSuffix for api.Config: projectconfig.go, and
// <x>config.go per mechanic — backofficeconfig.go, say. config.go, the
// aggregate itself, is not one.
const ConfigPartSuffix = "config.go"

// ConstructorExempt are the sandbox/api/ files that declare no field of the
// sandbox: sandbox.go is the struct itself, command.go and route.go are the
// shape of one command and of one route, each owned by the contract whose
// New<Name> builds the slice of them, trigger.go is the condition both of
// them match on, and routeprops.go and commandprops.go are the old home of
// what a route's and a command's handler are handed per run — a tree the
// build has not moved them out of yet, or whose layer is off.
var ConstructorExempt = []string{"sandbox.go", "command.go", "route.go", "trigger.go", RoutePropsFile, CommandPropsFile}

// IsConstructorExempt reports a sandbox/api/ file that declares no field of
// the sandbox: one of ConstructorExempt, or a part api.Sandbox or api.Config
// embeds (IsSandboxPart, IsConfigPart). A name is read with GeneratedPrefix
// taken off.
func IsConstructorExempt(sandbox *api.Sandbox, name string) bool {
	name = SourceName(sandbox, name)
	for _, exempt := range ConstructorExempt {
		if exempt == name {
			return true
		}
	}
	return IsSandboxPart(sandbox, name) || IsConfigPart(sandbox, name)
}

// IsSandboxPart reports a sandbox/api/ file whose structs api.Sandbox embeds:
// its name ends with SandboxPartSuffix and is not sandbox.go itself.
func IsSandboxPart(sandbox *api.Sandbox, name string) bool {
	name = SourceName(sandbox, name)
	return name != SandboxPartSuffix && sandbox.Deps.StringsDeps.HasSuffix(name, SandboxPartSuffix)
}

// IsConfigPart reports a sandbox/api/ file whose structs api.Config embeds:
// its name ends with ConfigPartSuffix and is not config.go itself, which is the
// Config contract and a field of the Sandbox.
func IsConfigPart(sandbox *api.Sandbox, name string) bool {
	name = SourceName(sandbox, name)
	return name != ConfigPartSuffix && sandbox.Deps.StringsDeps.HasSuffix(name, ConfigPartSuffix)
}

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
// contract of sandbox/api/: sandbox/internal/<name>, or
// sandbox/internal/<name>/<name> when the package that builds the contract
// sits one level down. A generated layer (the cli and server registries,
// config) and a contract the project writes itself live there alike: the
// first is a generated.new.go, the second a new.go (ConstructorNewFile).
func ConstructorSource(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string) string {
	root := "sandbox/internal"
	nested := root + "/" + name + "/" + name
	if ConstructorNewFile(sandbox, io, nested) != "" {
		return nested
	}
	return root + "/" + name
}

// ConstructorNewFile is the file of a constructor source declaring its
// New<Name>: generated.new.go for a layer the build renders, new.go for one
// the project writes, and "" when the package has neither.
func ConstructorNewFile(sandbox *api.Sandbox, io *stagedfs.StagedFS, source string) string {
	return GeneratedPath(sandbox, io, source, "new.go")
}
