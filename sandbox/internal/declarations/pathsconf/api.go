package pathreplacerconf

type PathReplacerEntry struct {
	Original    string
	Replacement string
}

type PathReplacerConf struct {
	Entries []PathReplacerEntry

	AddEntry func(original string, replacement string)
	Format   func(path string) string
	Render   func() string
}
