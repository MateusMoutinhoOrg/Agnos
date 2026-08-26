package api

type StartProps struct {
	Path        string
	ProjectName string
	Module      *string
	Force       bool
}

type BuildProps struct {
	Path    string
	Project string
	Force   bool
}

type CoreApi struct {
	Start func(props StartProps) error
	Build func(props BuildProps) error
}
