package extension

type ExtensionApi struct {
	Name string

	Install     func() error
	IsInstalled func() bool
	IsAvailable func() bool

	Verify func() error
	Build  func() error
}
