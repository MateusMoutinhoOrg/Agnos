package pathsconf

type PathReplacerEntry struct {
	Original    string
	Replacement string
}

type PathsConf struct {
	Entries []PathReplacerEntry

	AddEntry func(original string, replacement string)
	Format   func(path string) string
	Render   func() string
}
