package ignorableconf

type IgnorableConf struct {
	Paths []string

	AddPath     func(path string)
	IsIgnorable func(path string) bool
	Render      func() string
}
