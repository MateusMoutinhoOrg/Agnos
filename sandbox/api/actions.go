package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

type Actions struct {
	Build      func(path string) error
	Verify     func(path string) error
	Start      func(props StartProps) error
	DepsInit   func(path string) error
	DepsPurge  func(path string) error
	DepInstall func(path string, dep string) error
	DepRemove  func(path string, dep string) error
	DepList    func(path string) ([]string, error)
	CliInit    func(path string) error
	CliPurge   func(path string) error
}
