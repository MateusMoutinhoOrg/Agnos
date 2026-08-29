package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

type SandBox struct {
	Cli    CliApi
	Config Config
	Core   CoreApi
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps
}
