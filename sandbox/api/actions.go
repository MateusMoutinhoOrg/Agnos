package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

type Actions struct {
	Build      func(path string) error
	Start      func(props StartProps) error
	EnableDeps func(path string) error
	RemoveDeps func(path string) error
	DepInstall func(path string, dep string) error
	DepRemove  func(path string, dep string) error
	DepList    func(path string) ([]string, error)
}
