package api

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
)

type Memory struct {
	Database keepdeps.KeepDatabase

	GetProjectName func() string
	GetVersion     func() string
	IsStateActive  func(state string) bool
	ActiveState    func(state string)
	DeactiveState  func(state string)
}

type MemoryApi struct {
	NewMemory func(path string) Memory
}
