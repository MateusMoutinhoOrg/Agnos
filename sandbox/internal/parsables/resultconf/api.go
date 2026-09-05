package resultconf

// TreeEntry is one file an example wrote into its TestDir: the path it holds
// relative to that directory, and the sha256 of its content, hex-encoded.
type TreeEntry struct {
	File string
	Sha  string
}

// ResultConf is the parsed form of one example's result.yaml — everything a
// run of that example produced, and everything the next run is compared
// against.
type ResultConf struct {
	// CliOutput is the example's standard output and standard error, merged
	// in the order they were written and normalized (see the exec_test
	// action): no absolute paths, no carriage returns.
	CliOutput string

	// ExitCode is the status the example exited with; 0 is success.
	ExitCode int

	// Tree is every file inside the example's TestDir, ordered by File.
	Tree []TreeEntry

	// AddTreeEntry appends one file of the TestDir to Tree.
	AddTreeEntry func(file string, sha string)

	Render func() string
}
