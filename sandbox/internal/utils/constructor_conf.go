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
// sandbox/internal/<name>/<name> when the layer splits its package into
// several — the server's route, server and errors — and keeps the one that
// builds the contract under its own name one level down.
func ConstructorSource(io *smartio.SmartIO, name string) string {
	nested := "sandbox/internal/" + name + "/" + name
	if io.IsFile(nested + "/new.go") {
		return nested
	}
	return "sandbox/internal/" + name
}
