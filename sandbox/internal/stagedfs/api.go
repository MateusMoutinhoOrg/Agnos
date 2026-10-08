package stagedfs

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/pathsconf"
)

type StagedFS struct {
	// Root is the target project directory every path is resolved against —
	// the value of the command's --path flag. It is normalized so "", "." and
	// "./" all mean "" (the current working directory, no prefix). Every path
	// handed to a StagedFS method is project-relative; StagedFS joins Root on
	// only at the boundary where it touches the real filesystem, so no
	// generation ever escapes Root.
	Root string

	// sandbox is the api the helpers resolve paths through, deps included.
	// It is held on the struct rather than passed to every helper because a
	// StagedFS is built once, by New, and every method closure it carries is
	// bound then.
	sandbox *api.Sandbox

	Replacers    *pathsconf.PathsConf
	Transactions map[string][]byte

	PendingCreateDirs []string
	PendingRemoveDirs []string

	// Journal is what the last Persist changed on disk, recorded before it
	// changed it, so Undo can put every path back the way it was.
	Journal *Journal

	ReadFile             func(path string) ([]byte, error)
	CreateFile           func(path string, content []byte) error
	WriteFile            func(path string, content []byte) error
	Persist              func() error
	Undo                 func() error
	IsDir                func(path string) bool
	IsFile               func(path string) bool
	Exists               func(path string) bool
	CreateDir            func(path string)
	RemoveDir            func(path string)
	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}

// Journal is the state of every path a Persist touched, as it was before: the
// files it wrote (Missing when a file was not there), the directories it
// created that were not there, and the whole contents of the directories it
// removed. All paths are rooted — the ones the filesystem was handed.
type Journal struct {
	Files        map[string][]byte
	Missing      map[string]bool
	Order        []string
	CreatedDirs  []string
	RemovedDirs  []string
	RemovedFiles map[string][]byte
}
