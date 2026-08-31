package themesconf

type Theme struct {
	Name        string
	Id          string
	Description string
}

type ThemesConf struct {
	Themes []Theme

	AddTheme func(name string, id string, description string) error
	GetTheme func(name string) (*Theme, error)
	Render   func() string
}
