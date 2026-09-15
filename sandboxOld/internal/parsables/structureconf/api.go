package structureconf

// Item is one described element of the project tree: a file, a directory, or a
// pattern standing for a family of them ("libs/<lib>/<lib>.go"). Name is the
// path of the element relative to its parent, so it may itself hold slashes.
type Item struct {
	Name        string
	Description string

	// Dir marks the element as a directory: the rendered tree writes it with a
	// trailing "/", and `verify` requires a directory rather than a file at
	// its path.
	Dir bool

	// Gen marks an element `build` rewrites and no one edits by hand. The
	// rendered tree prefixes its description with "(gen)".
	Gen bool

	// Order is the position the item takes among its siblings. HasOrder
	// distinguishes "order: 0" from an absent key: an item with no order is
	// listed after every ordered one, alphabetically by name.
	Order    int
	HasOrder bool

	Children []Item
}

// StructureConf is <ProjectName>Config/structure.yaml: the shape of the
// project as its author describes it, one nested Item per element worth
// documenting. docs/Structure is rendered from it, and `verify` rejects an
// item whose path is no longer on disk.
type StructureConf struct {
	Items []Item

	GetItem func(path string) (*Item, error)
	Render  func() string
}
