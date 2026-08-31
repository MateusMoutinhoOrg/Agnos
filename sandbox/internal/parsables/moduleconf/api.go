package moduleconf

type ModuleConf struct {
	Module    string
	GoVersion string
	Requires  []string

	Render func() string
}
