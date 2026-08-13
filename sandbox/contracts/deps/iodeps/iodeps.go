package iodeps

// Lib is the embedded Io library api, defining all filesystem operations
// the sandbox requires.
type Lib struct {
	ReadFile             func(path string) ([]byte, error)
	WriteFile            func(path string, content []byte) error
	IsDir                func(path string) bool
	IsFile               func(path string) bool
	Exist                func(path string) bool
	CreateDir            func(path string)
	ListDirs             func(path string) []string
	ListFiles            func(path string) []string
	ListAll              func(path string) []string
	ListDirsRecursively  func(path string) []string
	ListFilesRecursively func(path string) []string
	ListAllRecursively   func(path string) []string
}
