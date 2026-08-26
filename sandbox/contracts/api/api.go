package api

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

type SandBox struct {
	CliApi //embedding:
	Config
	CoreApi
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps
}
