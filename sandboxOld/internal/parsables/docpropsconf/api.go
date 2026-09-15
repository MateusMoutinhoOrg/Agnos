package docpropsconf

type DocPropsConf struct {
	Name        string
	Description string
	Themes      []string

	// Order is the position the doc takes in the index that lists it. HasOrder
	// distinguishes "order: 0" from an absent key: a doc with no order is
	// listed after every ordered one, alphabetically by name.
	Order    int
	HasOrder bool

	AddTheme func(id string)
	Render   func() string
}
