package projectconf

type ProjectConf struct {
	Name    string
	Version string

	Render func() string
}
