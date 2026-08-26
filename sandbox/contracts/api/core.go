package api

type CoreApi struct {
	Start func(path string, project_name string, module string) error
	Build func(path string) error
}
