package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos/bootstrap/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos/bootstrap/sandbox/contracts/deps"
	internallib "github.com/MateusMoutinhoOrg/Agnos/bootstrap/sandbox/internal/lib"
)

// New injects a Deps struct into the library and returns the api.Lib entry
// point. It delegates to the internal lib constructor, which stores the deps
// on the struct and runs the factories over it, each of which fills one
// function field with a closure reading those deps.
func New(d deps.Deps) api.Lib {
	return internallib.New(d)
}
