package smartio

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/pathreplacerconf"
)

type SmartIO struct {
	// Root is the target project directory every path is resolved against —
	// the value of the command's --path flag. It is normalized so "", "." and
	// "./" all mean "" (the current working directory, no prefix). Every path
	// handed to a SmartIO method is project-relative; SmartIO joins Root on
	// only at the boundary where it touches the real filesystem, so no
	// generation ever escapes Root.
	Root string

	Ignore       *ignorableconf.IgnorableConf
	Replacers    *pathreplacerconf.PathReplacerConf
	Transactions map[string][]byte

	PendingCreateDirs []string
	PendingRemoveDirs []string

	ReadFile             func(path string) ([]byte, error)
	WriteFile            func(path string, content []byte) error
	WriteFileOverwrite   func(path string, content []byte) error
	Persist              func() error
	IsDir                func(path string) bool
	IsFile               func(path string) bool
	Exist                func(path string) bool
	CreateDir            func(path string)
	RemoveDir            func(path string)
	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}
