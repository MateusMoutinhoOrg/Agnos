package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func BindMethods(deps *deps.Deps, io *SmartIO) {
	io.ReadFile = func(path string) ([]byte, error) { return ReadFile(deps, io, path) }
	io.WriteFile = func(path string, content []byte) error { return WriteFile(deps, io, path, content) }
	io.WriteFileOverwrite = func(path string, content []byte) error { return WriteFileOverwrite(deps, io, path, content) }
	io.Persist = func() error { return Persist(deps, io) }
	io.IsDir = func(path string) bool { return IsDir(deps, io, path) }
	io.IsFile = func(path string) bool { return IsFile(deps, io, path) }
	io.Exist = func(path string) bool { return Exist(deps, io, path) }
	io.CreateDir = func(path string) { CreateDir(deps, io, path) }
	io.RemoveDir = func(path string) { RemoveDir(deps, io, path) }
	io.ListDirs = func(path string) []string { return ListDirs(deps, io, path) }
	io.ListFiles = func(path string) []string { return ListFiles(deps, io, path) }
	io.ListAll = func(path string) []string { return ListAll(deps, io, path) }
	io.ListDirsRecursively = func(path string) []string { return ListDirsRecursively(deps, io, path) }
	io.ListFilesRecursively = func(path string) []string { return ListFilesRecursively(deps, io, path) }
	io.ListAllRecursively = func(path string) []string { return ListAllRecursively(deps, io, path) }
}
