package api

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// Sandbox is the whole library: every part declared in a
// sandbox/api/<x>sandbox.go file, embedded, plus the two fields every project
// has. sandbox.New returns it, and nothing callable lives outside of it.
type Sandbox struct {
	// CliSandbox is a part of the Sandbox, declared in sandbox/api/clisandbox.go.
	// Embedded, so each of its fields is read as sandbox.<Field> like any
	// contract. Every struct of a sandbox/api/<x>sandbox.go file is one: each
	// mechanic writes its own, the project writes usersandbox.go.
	CliSandbox

	// UserSandbox is a part of the Sandbox, declared in sandbox/api/usersandbox.go.
	// Embedded, so each of its fields is read as sandbox.<Field> like any
	// contract. Every struct of a sandbox/api/<x>sandbox.go file is one: each
	// mechanic writes its own, the project writes usersandbox.go.
	UserSandbox

	// Deps is every capability the sandbox reaches the outside world
	// through. It rides on the api so that a function handed the Sandbox
	// holds the whole of what it needs, and can call another field of the
	// api besides — which is what makes a field a caller replaced take
	// effect everywhere. It is also the one field that does not cross into
	// a consumer: an installed copy of this contract carries the api, never
	// the wiring behind it.
	Deps *deps.Deps

	// Config is what the project knows about itself, built from
	// <ProjectName>Config/project.yaml (see sandbox/api/config.go).
	Config Config
}
