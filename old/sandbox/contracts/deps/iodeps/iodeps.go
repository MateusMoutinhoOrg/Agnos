package iodeps

// This package is the sandbox's *copy* of the api a filesystem library
// exposes — the same mechanic as verbdeps, keepdeps and embeddeps, for the
// same reason: touching a file is an OS-bound effect, so `os` and
// `path/filepath` may not appear inside the sandbox. The contract is restated
// here, and the adapter — which lives outside the sandbox — is what fills it.
//
// The tracker in sandbox/ never calls it: every record it keeps is persisted
// through Deps.KeepLib. It is carried as a standing capability of the
// template, filled by the standard adapter over `os` and `path/filepath`, so
// a derived library that must touch the filesystem directly finds the
// contract already declared and already wired. See the Deps.IoLib field.

// Lib is the filesystem library injected whole as the Deps.IoLib field.
//
// Paths are whatever the host operating system accepts, resolved by the
// adapter — unlike embeddeps.Lib, which is always slash-separated and rooted
// at an asset tree. The listing functions report paths that already include
// the directory they were given, so a result can be passed straight back in.
//
// The predicates report false rather than an error: a path that cannot be
// stat'd is not a directory and is not a file, which is the answer the caller
// wanted either way.
type Lib struct {
	// ReadFile returns the whole content of the file at path. The error
	// reports a file that does not exist or could not be read.
	ReadFile func(path string) ([]byte, error)

	// WriteFile writes content to path, creating any missing parent
	// directory first and truncating an existing file. The error reports a
	// directory or a file that could not be written.
	WriteFile func(path string, content []byte) error

	// IsDir reports whether path exists and is a directory.
	IsDir func(path string) bool

	// IsFile reports whether path exists and is not a directory.
	IsFile func(path string) bool

	// Exist reports whether anything exists at path, directory or file.
	Exist func(path string) bool

	// CreateDir creates the directory at path together with any missing
	// parent. It reports nothing: a directory that already exists and a
	// directory just created are the same outcome to the caller.
	CreateDir func(path string)

	// ListDirs returns the directories directly inside path. Nested
	// directories are not descended into.
	ListDirs func(path string) []string

	// ListFiles returns the files directly inside path. Directories are not
	// reported.
	ListFiles func(path string) []string

	// ListAll returns every entry directly inside path, directories and
	// files alike.
	ListAll func(path string) []string

	// ListDirsRecursively returns every directory at or below path,
	// excluding path itself.
	ListDirsRecursively func(path string) []string

	// ListFilesRecursively returns every file at or below path, at any
	// depth. Directories are never reported.
	ListFilesRecursively func(path string) []string

	// ListAllRecursively returns every entry at or below path, directories
	// and files alike, excluding path itself.
	ListAllRecursively func(path string) []string
}
