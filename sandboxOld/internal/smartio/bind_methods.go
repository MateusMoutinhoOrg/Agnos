package smartio

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"

func BindMethods(sandbox *api.Sandbox, io *SmartIO) {
	io.ReadFile = func(path string) ([]byte, error) { return ReadFile(sandbox, io, path) }
	io.WriteFile = func(path string, content []byte) error { return WriteFile(sandbox, io, path, content) }
	io.WriteFileOverwrite = func(path string, content []byte) error { return WriteFileOverwrite(sandbox, io, path, content) }
	io.Persist = func() error { return Persist(sandbox, io) }
	io.IsDir = func(path string) bool { return IsDir(sandbox, io, path) }
	io.IsFile = func(path string) bool { return IsFile(sandbox, io, path) }
	io.Exist = func(path string) bool { return Exist(sandbox, io, path) }
	io.CreateDir = func(path string) { CreateDir(sandbox, io, path) }
	io.RemoveDir = func(path string) { RemoveDir(sandbox, io, path) }
	io.ListDirs = func(path string) []string { return ListDirs(sandbox, io, path) }
	io.ListFiles = func(path string) []string { return ListFiles(sandbox, io, path) }
	io.ListAll = func(path string) []string { return ListAll(sandbox, io, path) }
	io.ListDirsRecursively = func(path string) []string { return ListDirsRecursively(sandbox, io, path) }
	io.ListFilesRecursively = func(path string) []string { return ListFilesRecursively(sandbox, io, path) }
	io.ListAllRecursively = func(path string) []string { return ListAllRecursively(sandbox, io, path) }
}
