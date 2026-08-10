package embeddeps

type Lib struct {
	ReadFile(path string) ([]byte, error)

	// will list the file of the given dir
	ListFiles(path string) ([]string, error)

	// will list the file of the given dir and subdirs
	ListFilesRecursively(path string) ([]string, error)
}

