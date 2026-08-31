package projectconf

type ProjectConf struct {
	Name        string
	Version     string
	Description string

	Render func() string
}
