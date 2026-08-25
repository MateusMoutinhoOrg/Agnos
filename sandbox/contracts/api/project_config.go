package api

type ProjectConfig struct {
	GetProjectName    func() string
	GetProjectVersion func() string
	SetProjectName    func(string) error
	SetProjectVersion func(string) error
}
type ProjectConfigApi struct {
	NewProjectConfig func(path string) ProjectConfig
}
