package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

type Actions struct {
	Build func(path string) error
	Start func(props StartProps) error
}
