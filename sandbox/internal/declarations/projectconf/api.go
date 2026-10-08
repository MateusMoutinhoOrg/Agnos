package projectconf

type ProjectConf struct {
	ProjectName string
	Version     string

	Render func() string
}
